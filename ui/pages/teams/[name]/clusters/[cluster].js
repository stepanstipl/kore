import React from 'react'
import PropTypes from 'prop-types'
import moment from 'moment'
import { Divider, Typography, Collapse, Icon, Row, Col, List, Button, Form, Card, Badge, message, Drawer, Tooltip } from 'antd'
const { Paragraph, Text } = Typography

import KoreApi from '../../../../lib/kore-api'
import Breadcrumb from '../../../../lib/components/layout/Breadcrumb'
import UsePlanForm from '../../../../lib/components/plans/UsePlanForm'
import ComponentStatusTree from '../../../../lib/components/common/ComponentStatusTree'
import ResourceStatusTag from '../../../../lib/components/resources/ResourceStatusTag'
import { clusterProviderIconSrcMap } from '../../../../lib/utils/ui-helpers'
import copy from '../../../../lib/utils/object-copy'
import { featureEnabled, KoreFeatures } from '../../../../lib/utils/features'
import FormErrorMessage from '../../../../lib/components/forms/FormErrorMessage'
import { inProgressStatusList } from '../../../../lib/utils/ui-helpers'
import apiPaths from '../../../../lib/utils/api-paths'
import ServiceCredential from '../../../../lib/components/teams/service/ServiceCredential'
import ServiceCredentialForm from '../../../../lib/components/teams/service/ServiceCredentialForm'
import NamespaceClaim from '../../../../lib/components/teams/namespace/NamespaceClaim'
import NamespaceClaimForm from '../../../../lib/components/teams/namespace/NamespaceClaimForm'
import ClusterAccessInfo from '../../../../lib/components/teams/cluster/ClusterAccessInfo'
import { isReadOnlyCRD } from '../../../../lib/utils/crd-helpers'
import ClusterApplicationServiceForm from '../../../../lib/components/teams/cluster/applications/ClusterApplicationServiceForm'
import Service from '../../../../lib/components/teams/service/Service'

class ClusterPage extends React.Component {
  static propTypes = {
    team: PropTypes.object.isRequired,
    user: PropTypes.object.isRequired,
    cluster: PropTypes.object.isRequired
  }

  constructor(props) {
    super(props)
    this.state = {
      cluster: props.cluster,
      components: {},
      namespaceClaims: false,
      editMode: false,
      clusterParams: props.cluster.spec.configuration,
      formErrorMessage: null,
      validationErrors: null,
      createNamespace: false,
      serviceCredentials: false,
      serviceKinds: false,
      createServiceCredential: false,
      createApplicationService: false,
      services: false,
      applicationServices: false,
      revealBindings: {}
    }
  }

  static getInitialProps = async (ctx) => {
    const api = await KoreApi.client(ctx)
    let [ team, cluster ] = await Promise.all([
      api.GetTeam(ctx.query.name),
      api.GetCluster(ctx.query.name, ctx.query.cluster)
    ])
    if ((!cluster || !team) && ctx.res) {
      /* eslint-disable-next-line require-atomic-updates */
      ctx.res.statusCode = 404
    }
    return { team, cluster }
  }

  fetchCommonData = async () => {
    if (featureEnabled(KoreFeatures.SERVICES)) {
      const serviceKinds = await (await KoreApi.client()).ListServiceKinds()
      return { serviceKinds: serviceKinds.items }
    }
    return {}
  }

  fetchNamespacesData = async () => {
    const team = this.props.team.metadata.name
    const api = await KoreApi.client()
    let [ namespaceClaims, serviceCredentials ] = await Promise.all([
      api.ListNamespaces(team),
      featureEnabled(KoreFeatures.SERVICES) ? api.ListServiceCredentials(team, this.state.cluster.metadata.name) : Promise.resolve({ items: [] }),
    ])
    namespaceClaims = namespaceClaims.items.filter(ns => ns.spec.cluster.name === this.props.cluster.metadata.name)
    serviceCredentials = serviceCredentials.items

    const revealBindings = {}
    featureEnabled(KoreFeatures.SERVICES) && namespaceClaims.filter(nc => serviceCredentials.filter(sc => sc.spec.clusterNamespace === nc.spec.name).length > 0).forEach(nc => revealBindings[nc.spec.name] = true)

    return { namespaceClaims, serviceCredentials, revealBindings }
  }

  fetchApplicationServicesData = async () => {
    const team = this.props.team.metadata.name
    let services = await (await KoreApi.client()).ListServices(team)
    services = services.items.filter(s => s.spec.cluster && s.spec.cluster.name && s.spec.kind !== 'app')
    const applicationServices = services.filter(s => s.spec.cluster.namespace === this.props.cluster.metadata.namespace && s.spec.cluster.name === this.props.cluster.metadata.name)

    return { services, applicationServices }
  }

  componentDidMount = () => {
    this.startRefreshing()
    this.fetchCommonData().then(data => {
      this.setState({ ...data })
      this.fetchNamespacesData().then(data => this.setState({ ...data }))
      if (featureEnabled(KoreFeatures.APPLICATION_SERVICES)) {
        this.fetchApplicationServicesData().then(data => this.setState({ ...data }))
      }
    })
  }

  componentWillUnmount = () => {
    if (this.interval) {
      clearInterval(this.interval)
    }
  }

  interval = null
  api = null
  startRefreshing = async () => {
    this.api = await KoreApi.client()
    this.interval = setInterval(async () => {
      await this.refreshCluster()
    }, 5000)
  }

  refreshCluster = async () => {
    const cluster = await this.api.GetCluster(this.props.team.metadata.name, this.state.cluster.metadata.name)
    if (cluster) {
      this.setState({
        cluster: cluster,
        // Keep the params up to date with the cluster, unless we're in edit mode.
        clusterParams: this.state.editMode ? this.state.clusterParams : copy(cluster.spec.configuration)
      })
    } else {
      this.setState({ cluster: { ...this.state.cluster, deleted: true } })
    }
  }

  handleResourceUpdated = resourceType => {
    return (updatedResource, done) => {
      this.setState((state) => {
        return {
          [resourceType]: state[resourceType].map(r => r.metadata.name !== updatedResource.metadata.name ? r : { ...r, status: updatedResource.status })
        }
      }, done)
    }
  }

  handleResourceDeleted = resourceType => {
    return (name, done) => {
      this.setState((state) => {
        const revealBindings = copy(state.revealBindings)
        if (resourceType === 'serviceCredentials') {
          const serviceCred = state.serviceCredentials.find(sc => sc.metadata.name === name)
          revealBindings[serviceCred.spec.clusterNamespace] = Boolean(state.serviceCredentials.filter(sc => sc.metadata.name !== name && !sc.deleted && sc.spec.clusterNamespace === serviceCred.spec.clusterNamespace).length)
        }

        return {
          [resourceType]: state[resourceType].map(r => r.metadata.name !== name ? r : { ...r, deleted: true }),
          revealBindings
        }
      }, done)
    }
  }

  createNamespace = value => () => this.setState({ createNamespace: value })

  handleNamespaceCreated = namespaceClaim => {
    this.setState({
      namespaceClaims: this.state.namespaceClaims.concat([namespaceClaim]),
      createNamespace: false
    })
    message.loading(`Namespace "${namespaceClaim.spec.name}" requested on cluster "${namespaceClaim.spec.cluster.name}"`)
  }

  deleteNamespace = async (name, done) => {
    const team = this.props.team.metadata.name
    try {
      const namespaceClaims = copy(this.state.namespaceClaims)
      const namespaceClaim = namespaceClaims.find(nc => nc.metadata.name === name)
      await (await KoreApi.client()).RemoveNamespace(team, namespaceClaim.metadata.name)
      namespaceClaim.status.status = 'Deleting'
      namespaceClaim.metadata.deletionTimestamp = new Date()
      this.setState({ namespaceClaims }, done)
      message.loading(`Namespace deletion requested: ${namespaceClaim.spec.name}`)
    } catch (err) {
      console.error('Error deleting namespace', err)
      message.error('Error deleting namespace, please try again.')
    }
  }

  deleteServiceCredential = async (name, done) => {
    const team = this.props.team.metadata.name
    try {
      await (await KoreApi.client()).DeleteServiceCredentials(team, name)
      this.setState((state) => {
        return {
          serviceCredentials: state.serviceCredentials.map(r => r.metadata.name !== name ? r : {
            ...r,
            status: { ...r.status, status: 'Deleting' },
            metadata: {
              ...r.metadata,
              deletionTimestamp: new Date()
            }
          })
        }
      }, done)

      message.loading('Deletion of service access requested')
    } catch (err) {
      console.error('Error deleting service access', err)
      message.error('Error deleting service access, please try again.')
    }
  }

  createServiceCredential = (value) => () => {
    this.setState({ createServiceCredential: value })
  }

  handleServiceCredentialCreated = serviceCredential => {
    this.setState((state) => {
      const revealBindings = copy(state.revealBindings)
      revealBindings[serviceCredential.spec.clusterNamespace] = true
      return {
        createServiceCredential: false,
        serviceCredentials: [ ...state.serviceCredentials, serviceCredential ],
        revealBindings
      }
    })
    message.loading(`Service access with secret name "${serviceCredential.spec.secretName}" requested`)
  }

  handleApplicationServiceCreated = async (applicationService) => {
    this.setState((state) => {
      return {
        createApplicationService: false,
        services: [ ...state.services, applicationService ],
        applicationServices: [ ...state.applicationServices, applicationService ],
      }
    })
    const commonData = await this.fetchCommonData()
    const namespaceData = await this.fetchNamespacesData()
    this.setState({ ...commonData, ...namespaceData })
  }

  deleteApplicationService = async (name, done) => {
    const team = this.props.team.metadata.name
    try {
      const applicationServices = copy(this.state.applicationServices)
      const applicationService = applicationServices.find(s => s.metadata.name === name)
      await (await KoreApi.client()).DeleteService(team, applicationService.metadata.name)
      applicationService.status.status = 'Deleting'
      applicationService.metadata.deletionTimestamp = new Date()
      this.setState({ applicationServices }, done)
      message.loading(`Application service deletion requested: ${applicationService.metadata.name}`)
    } catch (err) {
      console.error('Error deleting application service', err)
      message.error('Error deleting application service, please try again.')
    }
  }

  onClusterConfigChanged = (updatedClusterParams) => {
    this.setState({
      clusterParams: updatedClusterParams
    })
  }

  onEditClick = (e) => {
    e.stopPropagation()
    this.setState({ editMode: true })
  }

  onCancelClick = (e) => {
    e.stopPropagation()
    this.setState({
      editMode: false,
      clusterParams: copy(this.state.cluster.spec.configuration)
    })
  }

  onSubmit = async (e) => {
    e.preventDefault()
    this.setState({ saving: true, validationErrors: null, formErrorMessage: null })
    const clusterUpdated = copy(this.state.cluster)
    clusterUpdated.spec.configuration = this.state.clusterParams
    try {
      await this.api.UpdateCluster(this.props.team.metadata.name, this.state.cluster.metadata.name, clusterUpdated)
      this.setState({
        cluster: { ...this.state.cluster, status: { ...this.state.cluster.status, status: 'Pending' } },
        saving: false,
        validationErrors: null,
        formErrorMessage: null,
        editMode: false
      })
    } catch (err) {
      this.setState({
        saving: false,
        formErrorMessage: (err.fieldErrors && err.message) ? err.message : 'An error occurred saving the cluster, please try again',
        validationErrors: err.fieldErrors // This will be undefined on non-validation errors, which is fine.
      })
    }
  }

  revealBindings = (namespaceName) => (key) => {
    const revealBindings = copy(this.state.revealBindings)
    revealBindings[namespaceName] = Boolean(key.length)
    this.setState({ revealBindings })
  }

  getCardTitle = (title, resources) => (
    <span>{title} {resources && <Badge showZero={true} style={{ marginLeft: '10px', backgroundColor: '#1890ff' }} count={resources.filter(s => !s.deleted).length} />}</span>
  )

  render = () => {
    const { team, user } = this.props
    const { cluster, namespaceClaims, applicationServices, serviceCredentials, serviceKinds, createServiceCredential, createApplicationService } = this.state

    const created = moment(cluster.metadata.creationTimestamp).fromNow()
    const deleted = cluster.metadata.deletionTimestamp ? moment(cluster.metadata.deletionTimestamp).fromNow() : false
    const clusterNotEditable = !cluster || isReadOnlyCRD(cluster) || !cluster.status || inProgressStatusList.includes(cluster.status.status)
    const editClusterFormConfig = {
      layout: 'horizontal', labelAlign: 'left', hideRequiredMark: true,
      labelCol: { xs: 24, xl: 10 }, wrapperCol: { xs: 24, xl: 14 }
    }

    const hasActiveNamespaces = namespaceClaims && Boolean(namespaceClaims.filter(c => !c.deleted).length)
    const hasActiveApplicationServices = applicationServices && Boolean(applicationServices.filter(c => !c.deleted).length)

    return (
      <div>
        <Breadcrumb
          items={[
            { text: team.spec.summary, href: '/teams/[name]', link: `/teams/${team.metadata.name}` },
            { text: 'Clusters', href: '/teams/[name]/[tab]', link: `/teams/${team.metadata.name}/clusters` },
            { text: cluster.metadata.name }
          ]}
        />

        <Row type="flex" gutter={[16,16]}>
          <Col span={24} xl={12}>
            <List.Item>
              <List.Item.Meta
                className="large-list-item"
                avatar={<img src={clusterProviderIconSrcMap[cluster.spec.kind]} />}
                title={<Text style={{ marginTop: '15px', display: 'block' }}>{cluster.metadata.name}</Text>}
                description={
                  <div>
                    <Text type='secondary'>Created {created}</Text>
                    {deleted ? <Text type='secondary'><br/>Deleted {deleted}</Text> : null }
                  </div>
                }
              />
              <div>
                <ClusterAccessInfo team={this.props.team} />
              </div>
            </List.Item>
          </Col>
          <Col span={24} xl={12} style={{ marginTop: '14px' }}>
            <Collapse style={{ marginTop: '12px' }}>
              <Collapse.Panel header="Detailed Cluster Status" extra={(<ResourceStatusTag resourceStatus={cluster.status} />)}>
                <ComponentStatusTree team={team} user={user} component={cluster} />
              </Collapse.Panel>
            </Collapse>
          </Col>
        </Row>

        <Divider />

        <Card
          title={this.getCardTitle('Namespaces', namespaceClaims)}
          extra={<Button type="primary" onClick={this.createNamespace(true)}>+ New namespace</Button>}
          style={{ marginBottom: '20px' }}
        >
          {!namespaceClaims && <Icon type="loading" />}
          {namespaceClaims && !hasActiveNamespaces && <Paragraph style={{ marginBottom: 0 }} type="secondary">No namespaces found for this cluster</Paragraph>}
          {namespaceClaims && namespaceClaims.map((namespaceClaim, idx) => {
            const filteredServiceCredentials = (serviceCredentials || []).filter(sc => sc.spec.clusterNamespace === namespaceClaim.spec.name)
            const activeServiceCredentials = filteredServiceCredentials.filter(nc => !nc.deleted)
            return (
              <React.Fragment key={namespaceClaim.metadata.name}>
                <NamespaceClaim
                  key={namespaceClaim.metadata.name}
                  team={team.metadata.name}
                  namespaceClaim={namespaceClaim}
                  deleteNamespace={this.deleteNamespace}
                  handleUpdate={this.handleResourceUpdated('namespaceClaims')}
                  handleDelete={this.handleResourceDeleted('namespaceClaims')}
                  refreshMs={15000}
                  propsResourceDataKey="namespaceClaim"
                  resourceApiPath={`/teams/${team.metadata.name}/namespaceclaims/${namespaceClaim.metadata.name}`}
                />
                {!namespaceClaim.deleted && featureEnabled(KoreFeatures.SERVICES) && (
                  <>
                    <Collapse onChange={this.revealBindings(namespaceClaim.spec.name)} activeKey={this.state.revealBindings[namespaceClaim.spec.name] ? ['bindings'] : []}>
                      <Collapse.Panel
                        key="bindings"
                        header={<span>Cloud service access <Badge showZero={true} style={{ marginLeft: '10px', backgroundColor: '#1890ff' }} count={activeServiceCredentials.length} /></span>}
                        extra={
                          <Tooltip title="Provide this namespace with access to a cloud service">
                            <Icon
                              type="plus"
                              onClick={e => {
                                e.stopPropagation()
                                this.createServiceCredential({ cluster, namespaceClaim })()
                              }}
                            />
                          </Tooltip>
                        }
                      >
                        <List
                          size="small"
                          locale={{ emptyText: 'No cloud service access found' }}
                          dataSource={filteredServiceCredentials}
                          renderItem={serviceCredential => (
                            <ServiceCredential
                              viewPerspective="cluster"
                              team={team.metadata.name}
                              serviceCredential={serviceCredential}
                              serviceKind={serviceKinds.find(kind => kind.metadata.name === serviceCredential.spec.kind)}
                              deleteServiceCredential={this.deleteServiceCredential}
                              handleUpdate={this.handleResourceUpdated('serviceCredentials')}
                              handleDelete={this.handleResourceDeleted('serviceCredentials')}
                              refreshMs={10000}
                              propsResourceDataKey="serviceCredential"
                              resourceApiPath={`${apiPaths.team(team.metadata.name).serviceCredentials}/${serviceCredential.metadata.name}`}
                            />
                          )}
                        >
                        </List>
                      </Collapse.Panel>
                    </Collapse>

                    {!namespaceClaim.deleted && idx < namespaceClaims.length - 1 && <Divider />}
                  </>
                )}

              </React.Fragment>
            )
          })}

          <Drawer
            title="Create namespace"
            placement="right"
            closable={false}
            onClose={this.createNamespace(false)}
            visible={Boolean(this.state.createNamespace)}
            width={700}
          >
            <NamespaceClaimForm team={team.metadata.name} cluster={cluster} handleSubmit={this.handleNamespaceCreated} handleCancel={this.createNamespace(false)}/>
          </Drawer>

          {featureEnabled(KoreFeatures.SERVICES) && (
            <Drawer
              title="Create service access"
              placement="right"
              closable={false}
              onClose={this.createServiceCredential(false)}
              visible={Boolean(createServiceCredential)}
              width={700}
            >
              {Boolean(createServiceCredential) &&
                <ServiceCredentialForm
                  team={team}
                  creationSource="namespace"
                  clusters={ [createServiceCredential.cluster] }
                  namespaceClaims={ [createServiceCredential.namespaceClaim]}
                  handleSubmit={this.handleServiceCredentialCreated}
                  handleCancel={this.createServiceCredential(false)}
                />
              }
            </Drawer>
          )}

        </Card>

        {featureEnabled(KoreFeatures.APPLICATION_SERVICES) && (
          <>
            <Card
              title={this.getCardTitle('Application services', applicationServices)}
              extra={<Button type="primary" onClick={() => this.setState({ createApplicationService: true })}>+ New application service</Button>}
            >
              {!applicationServices && <Icon type="loading" />}
              {applicationServices && !hasActiveApplicationServices && <Text type="secondary">No service applications found for this cluster</Text>}
              {applicationServices && (
                <List
                  className="hide-empty-text"
                  locale={{ emptyText: <div/> }}
                  dataSource={applicationServices}
                  renderItem={service => (
                    <Service
                      team={team.metadata.name}
                      service={service}
                      serviceKind={serviceKinds.find(sk => sk.metadata.name === service.spec.kind)}
                      deleteService={this.deleteApplicationService}
                      handleUpdate={this.handleResourceUpdated('applicationServices')}
                      handleDelete={this.handleResourceDeleted('applicationServices')}
                      refreshMs={10000}
                      propsResourceDataKey="service"
                      resourceApiPath={`/teams/${team.metadata.name}/services/${service.metadata.name}`}
                    />
                  )}
                />
              )}
            </Card>
            <Drawer
              title="Create cluster application service"
              placement="right"
              closable={false}
              onClose={() => this.setState({ createApplicationService: false })}
              visible={Boolean(createApplicationService)}
              width={700}
            >
              {Boolean(createApplicationService) &&
                <ClusterApplicationServiceForm
                  team={team}
                  cluster={cluster}
                  teamServices={this.state.services}
                  handleSubmit={this.handleApplicationServiceCreated}
                  handleCancel={() => this.setState({ createApplicationService: false })}
                />
              }
            </Drawer>
          </>
        )}

        <Divider />

        <Collapse>
          <Collapse.Panel header="Cluster Parameters">
            <Form {...editClusterFormConfig} onSubmit={(e) => this.onSubmit(e)}>
              <FormErrorMessage message={this.state.formErrorMessage} />
              <Form.Item label="" colon={false}>
                {!this.state.editMode ? (
                  <Button icon="edit" htmlType="button" disabled={clusterNotEditable} onClick={(e) => this.onEditClick(e)}>Edit</Button>
                ) : (
                  <>
                    <Button type="primary" icon="save" htmlType="submit" loading={this.state.saving} disabled={this.state.saving || clusterNotEditable}>Save</Button>
                    &nbsp;
                    <Button icon="stop" htmlType="button" onClick={(e) => this.onCancelClick(e)}>Cancel</Button>
                  </>
                )}
              </Form.Item>
              <UsePlanForm
                team={team}
                resourceType="cluster"
                kind={cluster.spec.kind}
                plan={cluster.spec.plan}
                planValues={this.state.clusterParams}
                mode={this.state.editMode ? 'edit' : 'view'}
                validationErrors={this.state.validationErrors}
                onPlanChange={this.onClusterConfigChanged}
              />
            </Form>
          </Collapse.Panel>
        </Collapse>

      </div>
    )
  }
}
export default ClusterPage
