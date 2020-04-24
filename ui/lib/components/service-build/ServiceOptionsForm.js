import * as React from 'react'
import PropTypes from 'prop-types'

import { Typography, Form, Card, Radio, Modal, Input, Collapse } from 'antd'
const { Text, Title } = Typography

import PlanViewer from '../configure/PlanViewer'
import PlanOptionsForm from '../plans/PlanOptionsForm'
import KoreApi from '../../kore-api'

class ServiceOptionsForm extends React.Component {
  static propTypes = {
    form: PropTypes.any.isRequired,
    team: PropTypes.object.isRequired,
    selectedServiceKind: PropTypes.string.isRequired,
    servicePlans: PropTypes.array.isRequired,
    teamServices: PropTypes.array.isRequired,
    onServicePlanOverridden: PropTypes.func,
    validationErrors: PropTypes.array
  }

  onServicePlanChange = e => {
    if (this.props.form.getFieldValue('serviceName')) {
      this.props.form.setFieldsValue({ 'serviceName': this.generateServiceName(e.target.value) })
    }
  }

  generateServiceName = selectedServicePlan => {
    let serviceName = `${this.props.team.metadata.name}-${selectedServicePlan}`
    const matchingServices = this.props.teamServices.filter(tc => tc.metadata.name.indexOf(serviceName) === 0)
    if (matchingServices.length) {
      serviceName = `${serviceName}-${matchingServices.length + 1}`
    }
    return serviceName
  }

  showServicePlanDetails = servicePlanName => {
    return () => {
      const selectedServicePlan = this.props.servicePlans.find(p => p.metadata.name === servicePlanName)
      Modal.info({
        title: (<><Title level={4}>{selectedServicePlan.spec.description}</Title><Text>{selectedServicePlan.spec.summary}</Text></>),
        content: <PlanViewer
          plan={selectedServicePlan}
          getPlanSchema={async () => await (await KoreApi.client()).GetServicePlanSchema(selectedServicePlan.spec.kind)}
        />,
        width: 700,
        onOk() {}
      })
    }
  }

  onServicePlanOverridden = paramValues => {
    if (this.props.onServicePlanOverridden) {
      this.props.onServicePlanOverridden(paramValues)
    }
  }

  render() {
    const { getFieldDecorator, getFieldValue } = this.props.form
    const { servicePlans } = this.props
    const selectedServicePlan = getFieldValue('servicePlan')

    const checkForDuplicateName = (rule, value) => {
      const matchingService = this.props.teamServices.find(tc => tc.metadata.name === value)
      if (!matchingService) {
        return Promise.resolve()
      }
      return Promise.reject('This name is already used!')
    }

    return (
      <Card title="Service options">
        <Form.Item label="ServicePlan">
          {getFieldDecorator('servicePlan', {
            rules: [{ required: true, message: 'Please select your service plan!' }],
          })(
            <Radio.Group onChange={this.onServicePlanChange}>
              {servicePlans.map((p, idx) => (
                <Radio.Button key={idx} value={p.metadata.name}>{p.spec.description}</Radio.Button>
              ))}
            </Radio.Group>
          )}
          {selectedServicePlan ?
            <a style={{ marginLeft: '20px' }} onClick={this.showServicePlanDetails(selectedServicePlan)}>View service plan details</a> :
            null
          }
        </Form.Item>
        {selectedServicePlan ? (
          <Form.Item label="Service name">
            {getFieldDecorator('serviceName', {
              rules: [
                { required: true, message: 'Please enter service name!' },
                { pattern: '^[a-z][a-z0-9-]{0,38}[a-z0-9]$', message: 'Name must consist of lower case alphanumeric characters or "-", it must start with a letter and end with an alphanumeric and must be no longer than 40 characters' },
                { validator: checkForDuplicateName }
              ],
              initialValue: this.generateServiceName(selectedServicePlan)
            })(
              <Input />
            )}
          </Form.Item>
        ) : null}
        {selectedServicePlan ? (
          <Collapse>
            <Collapse.Panel header="Customize service parameters">
              <PlanOptionsForm
                team={this.props.team}
                plan={selectedServicePlan}
                getPlanDetails={async (team, plan) => await (await KoreApi.client()).GetTeamServicePlanDetails(team, plan)}
                getPlanConfiguration={(planDetails) => planDetails.servicePlan.configuration}
                validationErrors={this.props.validationErrors}
                onPlanChange={this.onServicePlanOverridden}
                mode="create"
              />
            </Collapse.Panel>
          </Collapse>
        ) : null}
      </Card>
    )
  }
}

const WrappedServiceOptionsForm = Form.create({ name: 'service_options' })(ServiceOptionsForm)

export default WrappedServiceOptionsForm
