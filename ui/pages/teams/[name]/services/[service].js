import React from 'react'
import PropTypes from 'prop-types'
import axios from 'axios'
import moment from 'moment'
import { Typography, Collapse, Row, Col, List, Button, Form, Avatar } from 'antd'
const { Text } = Typography

import KoreApi from '../../../../lib/kore-api'
import Breadcrumb from '../../../../lib/components/layout/Breadcrumb'
import PlanOptionsForm from '../../../../lib/components/plans/PlanOptionsForm'
import ComponentStatusTree from '../../../../lib/components/common/ComponentStatusTree'
import ResourceStatusTag from '../../../../lib/components/resources/ResourceStatusTag'
import copy from '../../../../lib/utils/object-copy'
import FormErrorMessage from '../../../../lib/components/forms/FormErrorMessage'
import { inProgressStatusList } from '../../../../lib/utils/ui-helpers'

class ServicePage extends React.Component {
  static propTypes = {
    team: PropTypes.object.isRequired,
    user: PropTypes.object.isRequired,
    service: PropTypes.object.isRequired
  }

  constructor(props) {
    super(props)
    this.state = {
      service: props.service,
      components: {},
      editMode: false,
      serviceParams: props.service.spec.configuration,
      formErrorMessage: null,
      validationErrors: null
    }
  }

  static getInitialProps = async ctx => {
    const api = await KoreApi.client(ctx)
    const { team, service } = await (axios.all([
      api.GetTeam(ctx.query.name), 
      api.GetService(ctx.query.name, ctx.query.service)
    ]).then(axios.spread((team, service) => { 
      return { team, service } 
    })))

    if ((!service || !team) && ctx.res) {
      /* eslint-disable-next-line require-atomic-updates */
      ctx.res.statusCode = 404
    }
    return { team, service }
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
    const { service } = this.state
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

        <List.Item actions={[<ResourceStatusTag key="status" resourceStatus={service.status} />]}>
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

        <Row type="flex" gutter={[16,16]}>
          <Col span={24} xl={12}>
            <Collapse defaultActiveKey={['0']}>
              <Collapse.Panel header="Detailed Service Status" extra={(<ResourceStatusTag resourceStatus={service.status} />)}>
                <ComponentStatusTree team={team} user={user} component={service} />
              </Collapse.Panel>
            </Collapse>
          </Col>
          <Col span={24} xl={12}>
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
                  <PlanOptionsForm
                    team={team}
                    resourceType="service"
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
      </div>
    )
  }
}
export default ServicePage
