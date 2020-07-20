import PropTypes from 'prop-types'
import { Card, Alert, Form, Select, Tag, Tooltip, Typography } from 'antd'
const { Option } = Select
const { Paragraph } = Typography
import { pluralize, titleize } from 'inflect'

import KoreApi from '../../kore-api'
import copy from '../../utils/object-copy'
import ManagePlanForm from './ManagePlanForm'
import { filterCloudAccountList, getProviderCloudInfo } from '../../utils/cloud'

/**
 * ManageClusterPlanForm is for *managing* a cluster plan.
 */
class ManageClusterPlanForm extends ManagePlanForm {
  static propTypes = {
    displayUnassociatedPlanWarning: PropTypes.bool
  }

  resourceType = () => 'cluster'
  estimateSupported = () => true

  cloudInfo = getProviderCloudInfo(this.props.kind)

  async fetchComponentData() {
    const { kind } = this.props
    const api = await KoreApi.client()
    const [ schema, accountManagementList, planList ] = await Promise.all([
      api.GetPlanSchema(kind),
      api.ListAccounts(),
      api.ListPlans(kind)
    ])
    const accountManagement = accountManagementList.items.find(a => a.spec.provider === this.props.kind)
    if (accountManagement) {
      accountManagement.spec.rules = filterCloudAccountList(accountManagement.spec.rules, planList.items)
    }
    this.setState({
      schema,
      accountManagement,
      dataLoading: false
    })
  }

  process = async (err, values) => {
    if (err) {
      this.setFormSubmitting(false, 'Validation failed')
      return
    }
    try {
      const api = await KoreApi.client()
      values.name = this.getMetadataName(values)
      const planResource = KoreApi.resources().generatePlanResource(this.props.kind, { ...values, configuration: { ...this.state.planValues } })
      const planResult = await api.UpdatePlan(values.name, planResource)

      if (this.accountManagementRulesEnabled()) {
        const accountMgtResource = copy(this.state.accountManagement)
        const currentRule = this.props.data && this.props.data.automatedCloudAccount ? accountMgtResource.spec.rules.find(r => r.name === this.props.data.automatedCloudAccount.name) : null
        if (values.automatedCloudAccount) {
          // add to the new rule
          const addedToRule = accountMgtResource.spec.rules.find(r => r.name === values.automatedCloudAccount)
          if (addedToRule) {
            addedToRule.plans.push(values.name)
            // remove from the existing rule if it's been changed
            if (currentRule && currentRule.name !== values.automatedCloudAccount) {
              currentRule.plans = currentRule.plans.filter(p => p !== values.name)
            }
            await api.UpdateAccount(accountMgtResource.metadata.name, accountMgtResource)
            planResult.amend = { automatedCloudAccount: addedToRule }
          } else {
            console.error(`Error occurred setting automated ${this.cloudInfo.accountNoun}, could not find rule with name`, values.automatedCloudAccount)
          }
        } else {
          // remove from the existing rule, if one exists
          if (currentRule) {
            currentRule.plans = currentRule.plans.filter(p => p !== values.name)
            await api.UpdateAccount(accountMgtResource.metadata.name, accountMgtResource)
            planResult.amend = { automatedCloudAccount: false }
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
    if (!this.props.data || !this.props.data.automatedCloudAccount) {
      return true
    }
    const planRule = this.state.accountManagement.spec.rules.find(r => r.name === this.props.data.automatedCloudAccount.name)
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
          message={`Associate with Kore managed ${pluralize(this.cloudInfo.accountNoun)}`}
          description={`Make this plan available to teams using Kore managed ${pluralize(this.cloudInfo.accountNoun)}.`}
          type="info"
          style={{ marginBottom: '20px' }}
        />
        <Form.Item label={`${this.cloudInfo.cloud} automated ${this.cloudInfo.accountNoun}`} validateStatus={this.fieldError('automatedCloudAccount') ? 'error' : ''} help={this.fieldError('automatedCloudAccount') || `Which ${this.cloudInfo.cloud} automated ${this.cloudInfo.accountNoun} this plan is associated with`}>
          {form.getFieldDecorator('automatedCloudAccount', {
            initialValue: data && data.automatedCloudAccount && data.automatedCloudAccount.name
          })(
            <Select placeholder={`${this.cloudInfo.cloud} automated ${this.cloudInfo.accountNoun}`} allowClear={this.allowAutomatedProjectSelectionClear()} disabled={this.disableAutomatedProjectSelection()}>
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
        {data && data.automatedCloudAccount && (
          <Paragraph>{this.cloudInfo.cloud} {this.cloudInfo.accountNoun} automation: <Tooltip overlay={`When using Kore managed ${this.cloudInfo.cloud} ${pluralize(this.cloudInfo.accountNoun)}, clusters using this plan will provisioned inside this ${this.cloudInfo.accountNoun} type.`}><Tag style={{ marginLeft: '10px' }}>{data.automatedCloudAccount.name}</Tag></Tooltip></Paragraph>
        )}
        {displayUnassociatedPlanWarning && (
          <Alert
            message={`This plan not associated with any ${this.cloudInfo.cloud} automated ${pluralize(this.cloudInfo.accountNoun)} and will not be available for teams to use. Set this below or go to ${titleize(this.cloudInfo.accountNoun)} automation settings to review this.`}
            type="warning"
            showIcon
            style={{ marginBottom: '20px' }}
          />
        )}

        {this.defaultFormHeader(formErrorMessage, mode, data)}

        {this.associateWithAccountManagement()}
      </>      
    )
  }
}

const WrappedManageClusterPlanForm = Form.create({ name: 'plan' })(ManageClusterPlanForm)

export default WrappedManageClusterPlanForm

