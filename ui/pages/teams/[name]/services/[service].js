import React from 'react'
import PropTypes from 'prop-types'
import axios from 'axios'
import moment from 'moment'
import { Typography, Collapse, Row, Col, List, Button, Form, Avatar, Card, Badge, message, Drawer } from 'antd'
const { Text } = Typography
import KoreApi from '../../../../lib/kore-api'
import Breadcrumb from '../../../../lib/components/layout/Breadcrumb'
import UsePlanForm from '../../../../lib/components/plans/UsePlanForm'
import ComponentStatusTree from '../../../../lib/components/common/ComponentStatusTree'
import ResourceStatusTag from '../../../../lib/components/resources/ResourceStatusTag'
import copy from '../../../../lib/utils/object-copy'
import FormErrorMessage from '../../../../lib/components/forms/FormErrorMessage'
import { inProgressStatusList } from '../../../../lib/utils/ui-helpers'
import apiPaths from '../../../../lib/utils/api-paths'
import ServiceCredential from '../../../../lib/components/teams/service/ServiceCredential'
import ServiceCredentialForm from '../../../../lib/components/teams/service/ServiceCredentialForm'

class ServicePage extends React.Component {
  static propTypes = {
    team: PropTypes.object.isRequired,
    user: PropTypes.object.isRequired,
    service: PropTypes.object.isRequired,
    serviceCredentials: PropTypes.object.isRequired,
  }

  constructor(props) {
    super(props)
    this.state = {
      service: props.service,
      components: {},
      editMode: false,
      serviceParams: props.service.spec.configuration,
      formErrorMessage: null,
      validationErrors: null,
      serviceCredentials: props.serviceCredentials,
      createServiceCredential: false
    }
  }

  static getInitialProps = async ctx => {
    const api = await KoreApi.client(ctx)
    const { team, service, serviceCredentials } = await (axios.all([
      api.GetTeam(ctx.query.name),
      api.GetService(ctx.query.name, ctx.query.service),
      api.ListServiceCredentials(ctx.query.name, '', ctx.query.service)
    ]).then(axios.spread((team, service, serviceCredentials) => {
      return { team, service, serviceCredentials }
    })))

    if ((!service || !team) && ctx.res) {
      /* eslint-disable-next-line require-atomic-updates */
      ctx.res.statusCode = 404
    }
    return { team, service, serviceCredentials }
  }

  interval = null
  api = null
  startRefreshing = async () => {
    this.api = await KoreApi.client()
    this.interval = setInterval(async () => {
      await this.refreshService()
    }, 5000)
  }

  refreshService = async () => {
    const service = await this.api.GetService(this.props.team.metadata.name, this.state.service.metadata.name)
    if (service) {
      this.setState({
        service: service,
        // Keep the params up to date with the service, unless we're in edit mode.
        serviceParams: this.state.editMode ? this.state.serviceParams : copy(service.spec.configuration)
      })
    } else {
      this.setState({ service: { ...this.state.service, deleted: true } })
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

  onServiceConfigChanged = (updatedServiceParams) => {
    this.setState({
      serviceParams: updatedServiceParams
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
      serviceParams: copy(this.state.service.spec.configuration)
    })
  }

  onSubmit = async (e) => {
    e.preventDefault()
    this.setState({ saving: true, validationErrors: null, formErrorMessage: null })
    const serviceUpdated = copy(this.state.service)
    serviceUpdated.spec.configuration = this.state.serviceParams
    try {
      await this.api.UpdateService(this.props.team.metadata.name, this.state.service.metadata.name, serviceUpdated)
      this.setState({
        service: { ...this.state.service, status: { ...this.state.service.status, status: 'Pending' } },
        saving: false,
        validationErrors: null,
        formErrorMessage: null,
        editMode: false
      })
      // await this.refreshService()
    } catch (err) {
      this.setState({
        saving: false,
        formErrorMessage: (err.fieldErrors && err.message) ? err.message : 'An error occurred saving the service, please try again',
        validationErrors: err.fieldErrors // This will be undefined on non-validation errors, which is fine.
      })
    }
  }

  render = () => {
    const { team, user } = this.props
    const { service, serviceCredentials, createServiceCredential } = this.state
    const created = moment(service.metadata.creationTimestamp).fromNow()
    const deleted = service.metadata.deletionTimestamp ? moment(service.metadata.deletionTimestamp).fromNow() : false
    const serviceNotEditable = !service || !service.status || inProgressStatusList.includes(service.status.status)
    const editServiceFormConfig = {
      layout: 'horizontal', labelAlign: 'left', hideRequiredMark: true,
      labelCol: { xs: 24, xl: 10 }, wrapperCol: { xs: 24, xl: 14 }
    }

    return (
      <div>
        <Breadcrumb
          items={[
            { text: team.spec.summary, href: '/teams/[name]', link: `/teams/${team.metadata.name}` },
            { text: `Service: ${service.metadata.name}` }
          ]}
        />

        <Row type="flex" gutter={[16,16]}>
          <Col span={24} xl={12}>
            <List.Item>
              <List.Item.Meta
                avatar={<Avatar icon="cloud" />}
                title={<Text>{service.spec.kind} <Text style={{ fontFamily: 'monospace', marginLeft: '15px' }}>{service.metadata.name}</Text></Text>}
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
            <Collapse>
              <Collapse.Panel header="Detailed Service Status" extra={(<ResourceStatusTag resourceStatus={service.status} />)}>
                <ComponentStatusTree team={team} user={user} component={service} />
              </Collapse.Panel>
            </Collapse>
          </Col>
        </Row>

        <Row type="flex" gutter={[16,16]} style={{ marginBottom: '12px' }}>
          <Col span={24} xl={24}>
            <Collapse>
              <Collapse.Panel header="Service Parameters">
                <Form {...editServiceFormConfig} onSubmit={(e) => this.onSubmit(e)}>
                  <FormErrorMessage message={this.state.formErrorMessage} />
                  <Form.Item label="" colon={false}>
                    {!this.state.editMode ? (
                      <Button icon="edit" htmlType="button" disabled={serviceNotEditable} onClick={(e) => this.onEditClick(e)}>Edit</Button>
                    ) : (
                      <>
                        <Button type="primary" icon="save" htmlType="submit" loading={this.state.saving} disabled={this.state.saving || serviceNotEditable}>Save</Button>
                        &nbsp;
                        <Button icon="stop" htmlType="button" onClick={(e) => this.onCancelClick(e)}>Cancel</Button>
                      </>
                    )}
                  </Form.Item>
                  <UsePlanForm
                    team={team}
                    resourceType="service"
                    kind={service.spec.kind}
                    plan={service.spec.plan}
                    planValues={this.state.serviceParams}
                    mode={this.state.editMode ? 'edit' : 'view'}
                    validationErrors={this.state.validationErrors}
                    onPlanChange={this.onServiceConfigChanged}
                  />
                </Form>
              </Collapse.Panel>
            </Collapse>
          </Col>
        </Row>

        <Drawer
          title="Create service binding"
          placement="right"
          closable={false}
          onClose={this.createServiceCredential(false)}
          visible={createServiceCredential}
          width={700}
        >
          <ServiceCredentialForm
            team={team}
            services={{ items: [service] }}
            handleSubmit={this.handleServiceCredentialCreated}
            handleCancel={this.createServiceCredential(false)}
          />
        </Drawer>
        <Row type="flex" gutter={[16,16]}>
          <Col span={24} xl={24}>
            <Card
              title={<div><Text style={{ marginRight: '10px' }}>Service bindings</Text><Badge style={{ backgroundColor: '#1890ff' }} count={serviceCredentials.items.filter(c => !c.deleted).length} /></div>}
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
      </div>
    )
  }
}
export default ServicePage
