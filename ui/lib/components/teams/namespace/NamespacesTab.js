import React from 'react'
import moment from 'moment'
import PropTypes from 'prop-types'
import {
  Button,
  Col,
  Divider,
  Drawer,
  Icon,
  List,
  Modal,
  Row,
  Tooltip,
  Tag,
  Typography,
  Collapse,
  Badge
} from 'antd'
const { Paragraph, Text } = Typography

import KoreApi from '../../../kore-api'
import copy from '../../../utils/object-copy'
import { inProgressStatusList, statusColorMap, statusIconMap } from '../../../utils/ui-helpers'
import { featureEnabled, KoreFeatures } from '../../../utils/features'
import NamespaceClaim from './NamespaceClaim'
import ServiceCredential from '../service/ServiceCredential'
import apiPaths from '../../../utils/api-paths'
import NamespaceClaimForm from './NamespaceClaimForm'
import ServiceCredentialForm from '../../../../lib/components/teams/service/ServiceCredentialForm'
import { loadingMessage, errorMessage } from '../../../utils/message'

class NamespacesTab extends React.Component {

  static propTypes = {
    team: PropTypes.object.isRequired,
    cluster: PropTypes.object.isRequired,
    onNamespaceCountChange: PropTypes.func
  }

  state = {
    dataLoading: true,
    namespaceClaims: [],
    serviceKinds: [],
    serviceCredentials: [],
    revealBindings: {},
    showNamespaceClaimForm: false,
    showServiceCredentialForm: false
  }

  async fetchComponentData() {
    try {
      const { team , cluster, onNamespaceCountChange } = this.props
      const api = await KoreApi.client()
      let [ namespaceClaims, serviceKinds, serviceCredentials ] = await Promise.all([
        api.ListNamespaces(team.metadata.name),
        api.ListServiceKinds(team),
        featureEnabled(KoreFeatures.SERVICES) ? api.ListServiceCredentials(team.metadata.name, cluster.metadata.name) : Promise.resolve({ items: [] }),
      ])
      namespaceClaims = namespaceClaims.items.filter(ns => ns.spec.cluster.name === cluster.metadata.name)
      serviceKinds = serviceKinds.items
      serviceCredentials = serviceCredentials.items

      const revealBindings = {}
      featureEnabled(KoreFeatures.SERVICES) && namespaceClaims.filter(nc => serviceCredentials.filter(sc => sc.spec.clusterNamespace === nc.spec.name).length > 0).forEach(nc => revealBindings[nc.spec.name] = true)

      onNamespaceCountChange && onNamespaceCountChange(namespaceClaims.length)

      return { namespaceClaims, serviceKinds, serviceCredentials, revealBindings }
    } catch (err) {
      console.error('Unable to load data for namespaces tab', err)
      return {}
    }
  }

  componentDidMountComplete = null
  componentDidMount() {
    this.componentDidMountComplete = this.fetchComponentData().then(data => {
      this.setState({ ...data, dataLoading: false })
    })
  }

  componentDidUpdate(prevProps) {
    if (prevProps.team.metadata.name !== this.props.team.metadata.name) {
      this.setState({ dataLoading: true })
      return this.fetchComponentData().then(data => this.setState({ ...data, dataLoading: false }))
    }
  }

  refreshServiceCredentials = async () => {
    let serviceCredentials = []
    try {
      const serviceCredentialsResult = await (await KoreApi.client()).ListServiceCredentials(this.props.team.metadata.name, this.props.cluster.metadata.name)
      serviceCredentials = serviceCredentialsResult.items
    } catch (error) {
      console.error('Failed to get service credentials', error)
    }
    if (serviceCredentials.length > 0) {
      this.setState( (state) => {
        const existingServiceCredentials = copy(state.serviceCredentials)
        serviceCredentials.forEach(sc => {
          const found = existingServiceCredentials.find(esc => esc.metadata.name === sc.metadata.name)
          if (found) {
            found.status = sc.status
          } else {
            existingServiceCredentials.push(sc)
          }
        })
        return { serviceCredentials: existingServiceCredentials }
      })
    }
  }

  handleNamespaceCreated = async (namespaceClaim) => {
    this.setState((state) => ({
      showNamespaceClaimForm: false,
      namespaceClaims: [ ...state.namespaceClaims, namespaceClaim ]
    }), async () => {
      this.props.onNamespaceCountChange && this.props.onNamespaceCountChange(this.state.namespaceClaims.length)
      await this.refreshServiceCredentials()
    })
  }

  handleServiceCredentialCreated = serviceCredential => {
    this.setState((state) => {
      const revealBindings = copy(state.revealBindings)
      revealBindings[serviceCredential.spec.clusterNamespace] = true
      return {
        showServiceCredentialForm: false,
        serviceCredentials: [ ...state.serviceCredentials, serviceCredential ],
        revealBindings
      }
    })
    loadingMessage(`Service access with secret name "${serviceCredential.spec.secretName}" requested`)
  }

  deleteNamespace = async (name, done) => {
    const team = this.props.team.metadata.name
    try {
      loadingMessage(`Namespace deletion requested: ${name}`)
      await (await KoreApi.client()).RemoveNamespace(team, name)
      this.setState((state) => {
        const namespaceClaims = copy(state.namespaceClaims)
        const namespaceClaim = namespaceClaims.find(nc => nc.metadata.name === name)
        namespaceClaim.status.status = 'Deleting'
        namespaceClaim.metadata.deletionTimestamp = new Date()
        return { namespaceClaims }
      }, done)
    } catch (err) {
      console.error('Error deleting namespace', err)
      errorMessage('Error deleting namespace, please try again.')
    }
  }

  deleteNamespaceConfirm = async (name, done) => {
    const namespaceClaim = this.state.namespaceClaims.find(nc => nc.metadata.name === name)
    const serviceCredentials = this.state.serviceCredentials.filter(sc => !sc.deleted && sc.spec.clusterNamespace === namespaceClaim.spec.name)
    if (serviceCredentials.length > 0) {
      return Modal.warning({
        title: 'Warning: namespace cannot be deleted',
        width: 600,
        content: (
          <div>
            <Paragraph strong>All service accesses must be removed before the namespace can be deleted.</Paragraph>
          </div>
        ),
        onOk() {}
      })
    }
    await this.deleteNamespace(name, done)
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

      loadingMessage('Deletion of service access requested')
    } catch (err) {
      console.error('Error deleting service access', err)
      errorMessage('Error deleting service access, please try again.')
    }
  }

  handleResourceUpdated = resourceType => {
    return async (updatedResource, done) => {
      this.setState((state) => {
        return {
          [resourceType]: state[resourceType].map(r => r.metadata.name !== updatedResource.metadata.name ? r : { ...r, status: updatedResource.status })
        }
      }, done)

      await this.refreshServiceCredentials()
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
      }, () => {
        this.props.onNamespaceCountChange && this.props.onNamespaceCountChange(this.state.namespaceClaims.filter(s => !s.deleted).length)
        done
      })
    }
  }

  serviceCredentialList = ({ serviceCredentials }) => {
    return (
      <Row style={{ marginLeft: '50px' }}>
        <Col>
          <Text strong style={{ marginRight: '8px' }}>Access: </Text>
          {serviceCredentials.map(serviceCredential => {
            const status = serviceCredential.status.status || 'Pending'
            const created = moment(serviceCredential.metadata.creationTimestamp).fromNow()
            return (
              <span key={serviceCredential.metadata.name} style={{ marginRight: '5px' }}>
                <Tooltip title={`Created ${created}`}>
                  <Tag color={statusColorMap[status] || 'red'}>{serviceCredential.spec.cluster.name}/{serviceCredential.spec.clusterNamespace} {inProgressStatusList.includes(status) ? <Icon type="loading" /> : <Icon type={statusIconMap[status]} />}</Tag>
                </Tooltip>
              </span>
            )
          })}
        </Col>
      </Row>
    )
  }

  revealBindings = (namespaceName) => (key) => {
    const revealBindings = copy(this.state.revealBindings)
    revealBindings[namespaceName] = Boolean(key.length)
    this.setState({ revealBindings })
  }

  render() {
    const { team, cluster } = this.props
    const { dataLoading, namespaceClaims, serviceKinds, serviceCredentials, showNamespaceClaimForm, showServiceCredentialForm } = this.state

    const hasActiveNamespaces =  Boolean(namespaceClaims.filter(c => !c.deleted).length)

    return (
      <>
        <Button type="primary" onClick={() => this.setState({ showNamespaceClaimForm: true })}>New namespace</Button>

        <Divider />

        {dataLoading ? (
          <Icon type="loading" />
        ) : (
          <>
            {!hasActiveNamespaces && <Paragraph type="secondary">No namespaces found for this team</Paragraph>}

            {namespaceClaims && namespaceClaims.map((namespaceClaim, idx) => {
              const filteredServiceCredentials = (serviceCredentials || []).filter(sc => sc.spec.clusterNamespace === namespaceClaim.spec.name)
              const activeServiceCredentials = filteredServiceCredentials.filter(nc => !nc.deleted)
              return (
                <React.Fragment key={namespaceClaim.metadata.name}>
                  <NamespaceClaim
                    key={namespaceClaim.metadata.name}
                    team={team.metadata.name}
                    namespaceClaim={namespaceClaim}
                    deleteNamespace={this.deleteNamespaceConfirm}
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
                                  this.setState({ showServiceCredentialForm : { cluster, namespaceClaim } })
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
                                cluster={cluster}
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
          </>
        )}

        <Drawer
          title="New namespace"
          visible={showNamespaceClaimForm}
          onClose={() => this.setState({ showNamespaceClaimForm: false })}
          width={900}
        >
          {showNamespaceClaimForm && (
            <NamespaceClaimForm team={team.metadata.name} cluster={cluster} handleSubmit={this.handleNamespaceCreated} handleCancel={() => this.setState({ showNamespaceClaimForm: false })}/>
          )}
        </Drawer>

        {featureEnabled(KoreFeatures.SERVICES) && (
          <Drawer
            title="Create service access"
            placement="right"
            closable={false}
            onClose={ () => this.setState({ showServiceCredentialForm: false }) }
            visible={Boolean(showServiceCredentialForm)}
            width={700}
          >
            {Boolean(showServiceCredentialForm) &&
            <ServiceCredentialForm
              team={team}
              creationSource="namespace"
              clusters={ [cluster] }
              namespaceClaims={ [showServiceCredentialForm.namespaceClaim]}
              handleSubmit={this.handleServiceCredentialCreated}
              handleCancel={ () => this.setState({ showServiceCredentialForm: false }) }
            />
            }
          </Drawer>
        )}


      </>
    )
  }

}

export default NamespacesTab
