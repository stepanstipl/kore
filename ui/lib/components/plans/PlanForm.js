import * as React from 'react'
import { Card, Alert, Form, Input, Button, Select, Tag, Tooltip, Typography } from 'antd'
const { Option } = Select
const { Paragraph } = Typography
import PropTypes from 'prop-types'
import { set } from 'lodash'

import KoreApi from '../../kore-api'
import V1Plan from '../../kore-api/model/V1Plan'
import V1PlanSpec from '../../kore-api/model/V1PlanSpec'
import V1ObjectMeta from '../../kore-api/model/V1ObjectMeta'
import PlanOption from './PlanOption'
import FormErrorMessage from '../forms/FormErrorMessage'
import canonical from '../../utils/canonical'
import copy from '../../utils/object-copy'

class PlanForm extends React.Component {
  static propTypes = {
    form: PropTypes.any.isRequired,
    data: PropTypes.object,
    kind: PropTypes.string.isRequired,
    validationErrors: PropTypes.array,
    handleSubmit: PropTypes.func.isRequired,
    handleValidationErrors: PropTypes.func.isRequired,
    displayUnassociatedPlanWarning: PropTypes.bool
  }

  state = {
    dataLoading: true,
    submitting: false,
    schema: null,
    planValues: this.props.data && this.props.data.spec.configuration || {},
    formErrorMessage: false
  }

  componentDidMountComplete = null
  componentDidMount() {
    this.componentDidMountComplete = this.fetchComponentData().then(() => {
      // To disabled submit button at the beginning.
      this.props.form.validateFields()
    })
  }

  async fetchComponentData() {
    const api = await KoreApi.client()
    const [ schema, accountManagementList ] = await Promise.all([
      api.GetPlanSchema(this.props.kind),
      api.ListAccounts()
    ])

    const accountManagement = accountManagementList.items.find(a => a.spec.provider === this.props.kind)
    this.setState({
      schema,
      accountManagement,
      dataLoading: false
    })
  }

  onValueChange(name, value) {
    // Texture this back into a state update using the nifty lodash set function:
    const newPlanValues = set({ ...this.state.planValues }, name, value)
    this.setState({
      planValues: newPlanValues
    })
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
    const data = this.props.data
    return (data && data.metadata && data.metadata.name) || canonical(values.description)
  }

  generatePlanResource = values => {
    const metadataName = this.getMetadataName(values)

    const planResource = new V1Plan()
    planResource.setApiVersion('config.kore.appvia.io/v1')
    planResource.setKind('Plan')

    const meta = new V1ObjectMeta()
    meta.setName(metadataName)
    planResource.setMetadata(meta)

    const spec = new V1PlanSpec()
    spec.setKind(this.props.kind)
    spec.setDescription(values.description)
    spec.setSummary(values.summary)
    spec.setConfiguration(values.configuration)
    planResource.setSpec(spec)

    return planResource
  }

  generatePlanConfiguration = () => {
    const properties = this.state.schema.properties
    const defaultConfiguration = {}
    Object.keys(properties).forEach(p => properties[p].type === 'boolean' ? defaultConfiguration[p] = false : null)
    return { ...defaultConfiguration, ...this.state.planValues }
  }

  _process = async (err, values) => {
    if (err) {
      this.setFormSubmitting(false, 'Validation failed')
      return
    }
    try {
      const api = await KoreApi.client()
      const metadataName = this.getMetadataName(values)
      const planResult = await api.UpdatePlan(metadataName, this.generatePlanResource({ ...values, configuration: this.generatePlanConfiguration() }))

      if (this.accountManagementRulesEnabled()) {
        const accountMgtResource = copy(this.state.accountManagement)
        const currentRule = this.props.data.gcpAutomatedProject ? accountMgtResource.spec.rules.find(r => r.name === this.props.data.gcpAutomatedProject.name) : null
        if (values.gcpAutomatedProject) {
          // add to the new rule
          const addedToRule = accountMgtResource.spec.rules.find(r => r.name === values.gcpAutomatedProject)
          if (addedToRule) {
            addedToRule.plans.push(metadataName)
            // remove from the existing rule if it's been changed
            if (currentRule && currentRule.name !== values.gcpAutomatedProject) {
              currentRule.plans = currentRule.plans.filter(p => p !== metadataName)
            }
            await api.UpdateAccount(`am-${accountMgtResource.spec.organization.name}`, accountMgtResource)
            planResult.append = { gcpAutomatedProject: addedToRule }
          } else {
            console.error('Error occurred setting automated project, could not find rule with name', values.gcpAutomatedProject)
          }
        } else {
          // remove from the existing rule, if one exists
          if (currentRule) {
            currentRule.plans = currentRule.plans.filter(p => p !== metadataName)
            await api.UpdateAccount(`am-${accountMgtResource.spec.organization.name}`, accountMgtResource)
            planResult.append = { gcpAutomatedProject: false }
          }
        }
      }

      return await this.props.handleSubmit(planResult)
    } catch (err) {
      console.error('Error submitting form', err)

      const message = (err.fieldErrors && err.message) ? err.message : 'An error occurred saving the plan, please try again'
      await this.props.handleValidationErrors(err.fieldErrors)
      this.setFormSubmitting(false, message)
    }
  }

  handleSubmit = e => {
    e.preventDefault()
    this.setFormSubmitting()
    this.props.form.validateFields(this._process)
  }

  accountManagementRulesEnabled = () => Boolean(this.state.accountManagement && this.state.accountManagement.spec.rules)

  allowAutomatedProjectSelectionClear = () => {
    // only allow clearing of the automated project if it's a new selection or there's more than one plan in the rule
    // a rule cannot be left with no plans
    if (!this.props.data.gcpAutomatedProject) {
      return true
    }
    const planRule = this.state.accountManagement.spec.rules.find(r => r.name === this.props.data.gcpAutomatedProject.name)
    if (planRule.plans.length === 1) {
      return false
    }
    return true
  }

  disableAutomatedProjectSelection = () => {
    if (!this.props.data.gcpAutomatedProject) {
      return false
    }
    const planRule = this.state.accountManagement.spec.rules.find(r => r.name === this.props.data.gcpAutomatedProject.name)
    if (planRule.plans.length === 1) {
      return true
    }
    return false
  }

  associateWithAccountManagement = () => {
    // only give an option to associate if rules exist
    if (!this.accountManagementRulesEnabled()) {
      return null
    }
    const { data, form } = this.props
    return (
      <Card style={{ marginBottom: '20px' }}>
        <Alert
          message="Associate with Kore managed projects"
          description="Make this plan available to teams using Kore managed projects."
          type="info"
          style={{ marginBottom: '20px' }}
        />
        <Form.Item label="GCP automated project" validateStatus={this.fieldError('gcpAutomatedProject') ? 'error' : ''} help={this.fieldError('gcpAutomatedProject') || 'Which GCP automated project this plan is associated with'}>
          {form.getFieldDecorator('gcpAutomatedProject', {
            initialValue: data && data.gcpAutomatedProject && data.gcpAutomatedProject.name
          })(
            <Select placeholder="GCP automated project" allowClear={this.allowAutomatedProjectSelectionClear()} disabled={this.disableAutomatedProjectSelection()}>
              {this.state.accountManagement.spec.rules.map(rule => <Option key={rule.name} value={rule.name}>{rule.name} - {rule.description}</Option>)}
            </Select>
          )}
        </Form.Item>
      </Card>
    )
  }

  // Only show error after a field is touched.
  fieldError = fieldKey => this.props.form.isFieldTouched(fieldKey) && this.props.form.getFieldError(fieldKey)

  render() {
    const { form, kind, data, validationErrors, displayUnassociatedPlanWarning } = this.props
    const { dataLoading, submitting, schema, planValues, formErrorMessage } = this.state
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
      }
    }

    if (dataLoading) {
      return (
        <div>Loading...</div>
      )
    }

    return (
      <Form {...formConfig} onSubmit={this.handleSubmit}>
        {data && data.gcpAutomatedProject && (
          <Paragraph>GCP project automation: <Tooltip overlay="When using Kore managed GCP projects, clusters using this plan will provisioned inside this project type."><Tag style={{ marginLeft: '10px' }}>{data.gcpAutomatedProject.name}</Tag></Tooltip></Paragraph>
        )}
        {displayUnassociatedPlanWarning && (
          <Alert
            message="This plan not associated with any GCP automated projects and will not be available for teams to use. Set this below or go to Project automation settings to review this."
            type="warning"
            showIcon
            style={{ marginBottom: '20px' }}
          />
        )}

        <FormErrorMessage message={formErrorMessage} />

        <Card style={{ marginBottom: '20px' }}>
          <Alert
            message="Help Kore teams understand this plan"
            description="Give this plan a summary and description to help teams choose the correct one."
            type="info"
            style={{ marginBottom: '20px' }}
          />
          <Form.Item label="Summary" validateStatus={this.fieldError('summary') ? 'error' : ''} help={this.fieldError('summary') || 'Summary of the plan'}>
            {getFieldDecorator('summary', {
              rules: [{ required: true, message: 'Please enter the name!' }],
              initialValue: data && data.spec.summary
            })(
              <Input placeholder="Summary" />,
            )}
          </Form.Item>
          <Form.Item label="Description" validateStatus={this.fieldError('description') ? 'error' : ''} help={this.fieldError('description') || 'Description of the plan'}>
            {getFieldDecorator('description', {
              rules: [{ required: true, message: 'Please enter the description!' }],
              initialValue: data && data.spec.description
            })(
              <Input placeholder="Description" />,
            )}
          </Form.Item>
        </Card>

        <this.associateWithAccountManagement />

        <Card style={{ marginBottom: '20px' }}>
          <Alert
            message="Plan configuration"
            description="Set the plan configuration"
            type="info"
            style={{ marginBottom: '20px' }}
          />

          {Object.keys(schema.properties).map(property =>
            <PlanOption
              mode="manage"
              resourceType="cluster"
              kind={kind}
              plan={planValues}
              key={property}
              name={property}
              property={schema.properties[property]}
              value={planValues[property]}
              editable={true}
              hideNonEditable={false}
              onChange={(n, v) => this.onValueChange(n, v)}
              validationErrors={validationErrors}
            />
          )}

        </Card>

        <Form.Item>
          <Button type="primary" htmlType="submit" loading={submitting} disabled={this.disableButton(getFieldsError())}>Save</Button>
        </Form.Item>
      </Form>
    )
  }
}

const WrappedPlanForm = Form.create({ name: 'plan' })(PlanForm)

export default WrappedPlanForm

