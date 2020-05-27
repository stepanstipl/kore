import * as React from 'react'
import PropTypes from 'prop-types'

import { Typography, Form, Modal, Input, Collapse, Select } from 'antd'
const { Text, Title } = Typography
const { Option } = Select

import PlanViewer from '../plans/PlanViewer'
import PlanOptionsForm from '../plans/PlanOptionsForm'

class ServiceOptionsForm extends React.Component {
  static propTypes = {
    form: PropTypes.any.isRequired,
    team: PropTypes.object.isRequired,
    selectedServiceKind: PropTypes.string.isRequired,
    servicePlans: PropTypes.array.isRequired,
    teamServices: PropTypes.array.isRequired,
    onServicePlanSelected: PropTypes.func,
    onServicePlanOverridden: PropTypes.func,
    validationErrors: PropTypes.array
  }

  componentDidUpdate(prevProps) {
    // Reset the selected plan if the credential changes:
    if (this.props.selectedServiceKind !== prevProps.selectedServiceKind) {
      this.props.form.setFieldsValue({ 'servicePlan': null })
    }
  }

  onServicePlanChange = (value) => {
    if (this.props.form.getFieldValue('serviceName')) {
      this.props.form.setFieldsValue({ 'serviceName': this.generateServiceName(value) })
    }
    this.props.onServicePlanSelected && this.props.onServicePlanSelected(value)
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
          resourceType="service"
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
    const { servicePlans, selectedServiceKind } = this.props
    const selectedServicePlan = getFieldValue('servicePlan')

    const checkForDuplicateName = (rule, value) => {
      const matchingService = this.props.teamServices.find(tc => tc.metadata.name === value)
      if (!matchingService) {
        return Promise.resolve()
      }
      return Promise.reject('This name is already used!')
    }

    return (
      <>
        <Form.Item label="Service plan">
          {getFieldDecorator('servicePlan', {
            rules: [{ required: true, message: 'Please select your service plan!' }],
          })(
            <Select onChange={this.onServicePlanChange} placeholder="Choose service plan">
              {servicePlans.map(p => <Option key={p.metadata.name} value={p.metadata.name}>{p.spec.description}</Option>)}
            </Select>
          )}
          {selectedServicePlan && <a onClick={this.showServicePlanDetails(selectedServicePlan)}>View service plan details</a>}
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
          <Collapse defaultActiveKey="plan">
            <Collapse.Panel key="plan" header="Customize service parameters">
              <PlanOptionsForm
                team={this.props.team}
                resourceType="service"
                kind={selectedServiceKind}
                plan={selectedServicePlan}
                validationErrors={this.props.validationErrors}
                onPlanChange={this.onServicePlanOverridden}
                mode="create"
              />
            </Collapse.Panel>
          </Collapse>
        ) : null}
      </>
    )
  }
}

const WrappedServiceOptionsForm = Form.create({ name: 'service_options' })(ServiceOptionsForm)

export default WrappedServiceOptionsForm
