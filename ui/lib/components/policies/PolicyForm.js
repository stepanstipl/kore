import * as React from 'react'
import PropTypes from 'prop-types'
import { Card, Alert, Form, Input, Button } from 'antd'

import KoreApi from '../../kore-api'
import FormErrorMessage from '../forms/FormErrorMessage'
import canonical from '../../utils/canonical'
import Policy from './Policy'
import AllocationsFormItem from '../forms/AllocationsFormItem'
import AllocationHelpers  from '../../utils/allocation-helpers'

class PolicyForm extends React.Component {
  static propTypes = {
    form: PropTypes.any.isRequired,
    policy: PropTypes.object,
    allocatedTeams: PropTypes.array,
    kind: PropTypes.string.isRequired,
    handleSubmit: PropTypes.func.isRequired
  }

  state = {
    submitting: false,
    policy: this.props.policy || {},
    allocation: null,
    allocatedTeams: this.props.allocatedTeams || ['*'],
    formErrorMessage: false
  }

  componentDidMountComplete = null
  componentDidMount() {
    // To disabled submit button at the beginning.
    this.props.form.validateFields()
    this.componentDidMountComplete = Promise.resolve().then(async () => {
      const existingAllocation = await AllocationHelpers.getAllocationForResource(this.props.policy)
      this.setState({
        allocatedTeams: existingAllocation ? existingAllocation.spec.teams : ['*']
      })
    })
  }

  onPolicyChange = (updatedPolicy) => {
    this.setState({ policy: updatedPolicy })
  }

  onAllocationChange = (updatedAllocation) => {
    this.setState({ allocatedTeams: updatedAllocation })
  }

  disableButton = fieldsError => {
    if (this.state.submitting) {
      return true
    }
    return Object.keys(fieldsError).some(field => fieldsError[field])
  }

  setFormSubmitting = (submitting = true, formErrorMessage = false) => {
    this.setState({
      submitting,
      formErrorMessage
    })
  }

  getMetadataName = values => {
    const policy = this.props.policy
    return (policy && policy.metadata && policy.metadata.name) || canonical(values.summary)
  }

  _process = async (err, values) => {
    if (err) {
      this.setFormSubmitting(false, 'Validation failed')
      return
    }
    if (!this.state.allocatedTeams || this.state.allocatedTeams.length === 0) {
      this.setFormSubmitting(false, 'You must allocate this policy to all teams or to one or more specific teams.')
      return
    }
    try {
      values.name = this.getMetadataName(values)
      const api = await KoreApi.client()
      const policyResource = KoreApi.resources().generatePolicyResource(this.props.kind, { ...values, properties: this.state.policy.spec.properties })
      const result = await api.UpdatePlanPolicy(values.name, policyResource)
      result.allocation = await AllocationHelpers.storeAllocation({ resourceToAllocate: policyResource, teams: this.state.allocatedTeams })
      return await this.props.handleSubmit(result)
    } catch (err) {
      console.log(err)
      const message = (err.fieldErrors && err.message) ? err.message : 'An error occurred saving the plan, please try again'
      this.setFormSubmitting(false, message)
    }
  }

  handleSubmit = e => {
    e.preventDefault()
    this.setFormSubmitting()
    this.props.form.validateFields(this._process)
  }

  // Only show error after a field is touched.
  fieldError = fieldKey => this.props.form.isFieldTouched(fieldKey) && this.props.form.getFieldError(fieldKey)

  render() {
    const { form } = this.props
    const { submitting, policy, formErrorMessage } = this.state
    const { getFieldDecorator, getFieldsError } = form

    const formConfig = {
      layout: 'horizontal',
      labelAlign: 'left',
      labelCol: {
        sm: { span: 24 },
        md: { span: 8 }
      },
      wrapperCol: {
        sm: { span: 24 },
        md: { span: 16 }
      },
      hideRequiredMark: true
    }

    return (
      <Form {...formConfig} onSubmit={this.handleSubmit}>

        <FormErrorMessage message={formErrorMessage} />

        <Card style={{ marginBottom: '20px' }}>
          <Alert
            message="Help Kore administrators understand this policy"
            description="Give this policy a name and description to help admins understand it"
            type="info"
            style={{ marginBottom: '20px' }}
          />
          <Form.Item label="Name" validateStatus={this.fieldError('summary') ? 'error' : ''} help={this.fieldError('summary') || 'Name of the policy'}>
            {getFieldDecorator('summary', {
              rules: [{ required: true, message: 'Please enter the name!' }],
              initialValue: policy && policy.spec.summary
            })(
              <Input placeholder="Name" />,
            )}
          </Form.Item>
          <Form.Item label="Description" validateStatus={this.fieldError('description') ? 'error' : ''} help={this.fieldError('description') || 'Description of this policy to explain it to admins of Kore'}>
            {getFieldDecorator('description', {
              rules: [{ required: true, message: 'Please enter the description!' }],
              initialValue: policy && policy.spec.description
            })(
              <Input placeholder="Description" />,
            )}
          </Form.Item>
          <AllocationsFormItem allocatedTeams={this.state.allocatedTeams} onAllocationChange={this.onAllocationChange} />
        </Card>

        <Card style={{ marginBottom: '20px' }}>
          <Alert
            message="Set policy rules"
            description="By default, all plan values are read-only for teams. By selecting Allow Update, teams that have this policy applied will be able to make changes to that property when creating or editing clusters. By selecting Disallow Update, teams that have this policy applied will NEVER be able to edit that property, even if another policy would allow it to be edited."
            type="info"
            style={{ marginBottom: '20px' }}
          />
          <Policy policy={policy} mode="edit" onPolicyUpdate={this.onPolicyChange} />
        </Card>

        <Form.Item>
          <Button id="policy_save" type="primary" htmlType="submit" loading={submitting} disabled={this.disableButton(getFieldsError())}>Save</Button>
        </Form.Item>
      </Form>
    )
  }
}

const WrappedPolicyForm = Form.create({ name: 'policy' })(PolicyForm)

export default WrappedPolicyForm
