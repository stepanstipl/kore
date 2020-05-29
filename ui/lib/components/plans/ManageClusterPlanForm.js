import { Card, Alert, Form, Select, Tag, Tooltip, Typography } from 'antd'
const { Option } = Select
const { Paragraph } = Typography
import PropTypes from 'prop-types'

import KoreApi from '../../kore-api'
import V1Plan from '../../kore-api/model/V1Plan'
import V1PlanSpec from '../../kore-api/model/V1PlanSpec'
import V1ObjectMeta from '../../kore-api/model/V1ObjectMeta'
import copy from '../../utils/object-copy'
import ManagePlanForm from './ManagePlanForm'

/**
 * ManageClusterPlanForm is for *managing* a cluster plan.
 */
class ManageClusterPlanForm extends ManagePlanForm {
  static propTypes = {
    displayUnassociatedPlanWarning: PropTypes.bool
  }

  resourceType = () => 'cluster'

  async fetchComponentData() {
    const { kind } = this.props
    const api = await KoreApi.client()
    const [ schema, accountManagementList ] = await Promise.all([
      api.GetPlanSchema(kind),
      api.ListAccounts()
    ])
    const accountManagement = accountManagementList.items.find(a => a.spec.provider === this.props.kind)
    this.setState({
      schema,
      accountManagement,
      dataLoading: false
    })
  }

  generatePlanResource = (values) => {
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

  process = async (err, values) => {
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
        const currentRule = this.props.data && this.props.data.gcpAutomatedProject ? accountMgtResource.spec.rules.find(r => r.name === this.props.data.gcpAutomatedProject.name) : null
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
      this.setFormSubmitting(false, null, [])
      return await this.props.handleSubmit(planResult)
    } catch (err) {
      console.error('Error submitting form', err)
      const message = (err.fieldErrors && err.message) ? err.message : 'An error occurred saving the plan, please try again'
      this.setFormSubmitting(false, message, err.fieldErrors)
    }
  }

  accountManagementRulesEnabled = () => Boolean(this.state.accountManagement && this.state.accountManagement.spec.rules)

  allowAutomatedProjectSelectionClear = () => {
    // only allow clearing of the automated project if it's a new selection or there's more than one plan in the rule
    // a rule cannot be left with no plans
    if (!this.props.data || !this.props.data.gcpAutomatedProject) {
      return true
    }
    const planRule = this.state.accountManagement.spec.rules.find(r => r.name === this.props.data.gcpAutomatedProject.name)
    if (planRule.plans.length === 1) {
      return false
    }
    return true
  }

  disableAutomatedProjectSelection = () => {
    return !this.allowAutomatedProjectSelectionClear()
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

  formHeader = (formErrorMessage, mode, data) => {
    const { displayUnassociatedPlanWarning } = this.props
    return (
      <>
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

        {this.defaultFormHeader(formErrorMessage, mode, data)}

        <this.associateWithAccountManagement />
      </>      
    )
  }
}

const WrappedManageClusterPlanForm = Form.create({ name: 'plan' })(ManageClusterPlanForm)

export default WrappedManageClusterPlanForm

