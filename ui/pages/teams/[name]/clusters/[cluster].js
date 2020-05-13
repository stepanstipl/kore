import React from 'react'
import PropTypes from 'prop-types'
import axios from 'axios'
import moment from 'moment'
import { Typography, Collapse, Row, Col, List, Button, Form, Card, Badge, message, Drawer } from 'antd'
const { Text } = Typography
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()
import KoreApi from '../../../../lib/kore-api'
import Breadcrumb from '../../../../lib/components/layout/Breadcrumb'
import PlanOptionsForm from '../../../../lib/components/plans/PlanOptionsForm'
import ComponentStatusTree from '../../../../lib/components/common/ComponentStatusTree'
import ResourceStatusTag from '../../../../lib/components/resources/ResourceStatusTag'
import { clusterProviderIconSrcMap } from '../../../../lib/utils/ui-helpers'
import copy from '../../../../lib/utils/object-copy'
import FormErrorMessage from '../../../../lib/components/forms/FormErrorMessage'
import { inProgressStatusList } from '../../../../lib/utils/ui-helpers'
import apiPaths from '../../../../lib/utils/api-paths'
import ServiceCredential from '../../../../lib/components/teams/service/ServiceCredential'
import ServiceCredentialForm from '../../../../lib/components/teams/service/ServiceCredentialForm'


class ClusterPage extends React.Component {
  static propTypes = {
    team: PropTypes.object.isRequired,
    user: PropTypes.object.isRequired,
    cluster: PropTypes.object.isRequired,
    serviceCredentials: PropTypes.object.isRequired,
  }

  constructor(props) {
    super(props)
    this.state = {
      cluster: props.cluster,
      components: {},
      editMode: false,
      clusterParams: props.cluster.spec.configuration,
      formErrorMessage: null,
      validationErrors: null,
      serviceCredentials: props.serviceCredentials,
      createServiceCredential: false
    }
  }

  static getInitialProps = async ctx => {
    const api = await KoreApi.client(ctx)
    const { team, cluster, serviceCredentials } = await (axios.all([
      api.GetTeam(ctx.query.name),
      api.GetCluster(ctx.query.name, ctx.query.cluster),
      publicRuntimeConfig.featureGates['services'] ? api.ListServiceCredentials(ctx.query.name, ctx.query.cluster, '') : Promise.resolve({ items: [] })
    ]).then(axios.spread((team, cluster, serviceCredentials) => {
      return { team, cluster, serviceCredentials }
    })))

    if ((!cluster || !team) && ctx.res) {
      /* eslint-disable-next-line require-atomic-updates */
      ctx.res.statusCode = 404
    }
    return { team, cluster, serviceCredentials }
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
          [resourceType]: {
            ...state[resourceType],
            items: state[resourceType].items.map(r => r.metadata.name !== updatedResource.metadata.name ? r : { ...r, status: updatedResource.status })
          }
        }
      }, done)
    }
  }

  handleResourceDeleted = resourceType => {
    return (name, done) => {
      this.setState((state) => {
        return {
          [resourceType]: {
            ...state[resourceType],
            items: state[resourceType].items.map(r => r.metadata.name !== name ? r : { ...r, deleted: true })
          }
        }
      }, done)
    }
  }

  deleteServiceCredential = async (name, done) => {
    const team = this.props.team.metadata.name
    try {
      await (await KoreApi.client()).DeleteServiceCredentials(team, name)

      this.setState((state) => {
        return {
          serviceCredentials: {
            ...state.serviceCredentials,
            items: state.serviceCredentials.items.map(r => r.metadata.name !== name ? r : {
              ...r,
              status: { ...r.status, status: 'Deleting' },
              metadata: {
                ...r.metadata,
                deletionTimestamp: new Date()
              }
            })
          }
        }
      }, done)

      message.loading(`Service Credential deletion requested: ${name}`)
    } catch (err) {
      console.error('Error deleting service credential', err)
      message.error('Error deleting service credential, please try again.')
    }
  }

  createServiceCredential = value => {
    return () => {
      this.setState({
        createServiceCredential: value
      })
    }
  }

  handleServiceCredentialCreated = serviceCredential => {
    this.setState((state) => {
      return {
        createServiceCredential: false,
        serviceCredentials: {
          ...state.serviceCredentials,
          items: [ ...state.serviceCredentials.items, serviceCredential]
        }
      }
    })
    message.loading(`Service credential "${serviceCredential.metadata.name}" requested`)
  }

  componentDidMount = () => {
    this.startRefreshing()
  }

  componentWillUnmount = () => {
    if (this.interval) {
      clearInterval(this.interval)
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
      // await this.refreshCluster()
    } catch (err) {
      this.setState({
        saving: false,
        formErrorMessage: (err.fieldErrors && err.message) ? err.message : 'An error occurred saving the cluster, please try again',
        validationErrors: err.fieldErrors // This will be undefined on non-validation errors, which is fine.
      })
    }
  }

  render = () => {
    const { team, user } = this.props
    const { cluster, serviceCredentials, createServiceCredential } = this.state
    const created = moment(cluster.metadata.creationTimestamp).fromNow()
    const deleted = cluster.metadata.deletionTimestamp ? moment(cluster.metadata.deletionTimestamp).fromNow() : false
    const clusterNotEditable = !cluster || !cluster.status || inProgressStatusList.includes(cluster.status.status)
    const editClusterFormConfig = {
      layout: 'horizontal', labelAlign: 'left', hideRequiredMark: true,
      labelCol: { xs: 24, xl: 10 }, wrapperCol: { xs: 24, xl: 14 }
    }

    return (
      <div>
        <Breadcrumb
          items={[
            { text: team.spec.summary, href: '/teams/[name]', link: `/teams/${team.metadata.name}` },
            { text: `Cluster: ${cluster.metadata.name}` }
          ]}
        />

        <Row type="flex" gutter={[16,16]}>
          <Col span={24} xl={12}>
            <List.Item>
              <List.Item.Meta
                avatar={<img src={clusterProviderIconSrcMap[cluster.spec.kind]} height="32px" />}
                title={<Text>{cluster.spec.kind} <Text style={{ fontFamily: 'monospace', marginLeft: '15px' }}>{cluster.metadata.name}</Text></Text>}
                description={
                  <div>
                    <Text type='secondary'>Created {created}</Text>
                    {deleted ? <Text type='secondary'><br/>Deleted {deleted}</Text> : null }
                  </div>
                }
              />
            </List.Item>
          </Col>
          <Col span={24} xl={12}>
            <Collapse style={{ marginTop: '12px' }}>
              <Collapse.Panel header="Detailed Cluster Status" extra={(<ResourceStatusTag resourceStatus={cluster.status} />)}>
                <ComponentStatusTree team={team} user={user} component={cluster} />
              </Collapse.Panel>
            </Collapse>
          </Col>
        </Row>
        <Row type="flex" gutter={[16,16]} style={{ marginBottom: '12px' }}>
          <Col span={24} xl={24}>
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
                  <PlanOptionsForm
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
          </Col>
        </Row>

        {publicRuntimeConfig.featureGates['services'] ? (
          <>
            <Drawer
              title="Create service credential"
              placement="right"
              closable={false}
              onClose={this.createServiceCredential(false)}
              visible={createServiceCredential}
              width={700}
            >
              <ServiceCredentialForm
                team={team}
                clusters={{ items: [cluster] }}
                handleSubmit={this.handleServiceCredentialCreated}
                handleCancel={this.createServiceCredential(false)}
              />
            </Drawer>

            <Row type="flex" gutter={[16,16]}>
              <Col span={24} xl={24}>
                <Card
                  title={<div><Text style={{ marginRight: '10px' }}>Service credentials</Text><Badge style={{ backgroundColor: '#1890ff' }} count={serviceCredentials.items.filter(c => !c.deleted).length} /></div>}
                  style={{ marginBottom: '20px' }}
                  extra={
                    <div>
                      <Button type="primary" onClick={this.createServiceCredential(true)}>+ New</Button>
                    </div>
                  }
                >
                  <List
                    dataSource={serviceCredentials.items}
                    renderItem={serviceCredential => {
                      return (
                        <ServiceCredential
                          team={team.metadata.name}
                          serviceCredential={serviceCredential}
                          deleteServiceCredential={this.deleteServiceCredential}
                          handleUpdate={this.handleResourceUpdated('serviceCredentials')}
                          handleDelete={this.handleResourceDeleted('serviceCredentials')}
                          refreshMs={10000}
                          propsResourceDataKey="serviceCredential"
                          resourceApiPath={`${apiPaths.team(team.metadata.name).serviceCredentials}/${serviceCredential.metadata.name}`}
                        />
                      )
                    }}
                  >
                  </List>
                </Card>
              </Col>
            </Row>
          </>
        ): null}
      </div>
    )
  }
}
export default ClusterPage
