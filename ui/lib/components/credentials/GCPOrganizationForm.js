import * as React from 'react'
import VerifiedAllocatedResourceForm from '../resources/VerifiedAllocatedResourceForm'
import KoreApi from '../../kore-api'
import { Checkbox, Form, Input, Alert, Card } from 'antd'
import AllocationHelpers from '../../utils/allocation-helpers'

class GCPOrganizationForm extends VerifiedAllocatedResourceForm {

  getResource = async metadataName => {
    const api = await KoreApi.client()
    const gcpOrgResult = await api.GetGCPOrganization(this.props.team, metadataName)
    gcpOrgResult.allocation = await AllocationHelpers.getAllocationForResource(gcpOrgResult)
    return gcpOrgResult
  }

  putResource = async values => {
    const api = await KoreApi.client()
    values.name = this.getMetadataName(values)
    const secretName = values.name
    const teamResources = KoreApi.resources().team(this.props.team)
    if (!this.props.data || this.state.replaceKey) {
      const secretData = { key: btoa(values.account) }
      const secretResource = teamResources.generateSecretResource(secretName, 'gcp-org', `GCP admin project Service Account for ${values.parentID}`, secretData)
      await api.UpdateTeamSecret(this.props.team, secretName, secretResource)
    }
    const gcpOrgResource = teamResources.generateGCPOrganizationResource(values, secretName)
    const gcpOrgResult = await api.UpdateGCPOrganization(this.props.team, values.name, gcpOrgResource)
    gcpOrgResult.allocation = await this.storeAllocation(gcpOrgResource, values)
    return gcpOrgResult
  }

  allocationFormFieldsInfo = {
    allocationMissing: {
      infoMessage: 'This organization credential is not allocated to any teams',
      infoDescription: 'Give the organization credential a Name and Description below and enter Allocated team(s) as appropriate. Once complete, click Save to allocate it.'
    },
    nameSection: {
      infoMessage: 'Help Kore teams understand this organization credential',
      infoDescription: 'Give this organization credential a name and description to help teams choose the correct one.',
      nameHelp: 'The name for the organization credential eg. MyOrg',
      descriptionHelp: 'A description of the organization credential to help when choosing between them'
    },
    allocationSection: {
      infoMessage: 'Make this organization credential available to teams in Kore',
      infoDescription: 'This will give teams the ability to create GCP projects inside the organization. Within these scoped projects the team will have the ability to create clusters.',
      allTeamsWarningMessage: 'This organization credential will be available to all teams',
      allTeamsWarningDescription: 'No teams exist in Kore yet, therefore currently this organization credential will be available to all teams created in the future. If you wish to restrict this please return here and allocate to teams once they have been created.',
      allocateExtra: 'If nothing selected then this organization credential will be available to ALL teams'
    }
  }

  resourceFormFields = () => {
    const { form, data } = this.props
    const { replaceKey } = this.state
    return (
      <Card style={{ marginBottom: '20px' }}>
        <Alert
          message="GCP Organization details"
          description="Retrieve these values from your GCP organization. Providing these gives Kore the ability to create projects for teams within the GCP organization. Teams will then be able to provision clusters within their projects."
          type="info"
          style={{ marginBottom: '20px' }}
        />
        <Form.Item label="Organization ID" validateStatus={this.fieldError('parentID') ? 'error' : ''} help={this.fieldError('parentID') || 'The GCP organization ID'}>
          {form.getFieldDecorator('parentID', {
            rules: [{ required: true, message: 'Please enter the organization ID!' }],
            initialValue: data && data.spec.parentID
          })(
            <Input placeholder="Organization ID" />,
          )}
        </Form.Item>
        <Form.Item label="Billing account" validateStatus={this.fieldError('billingAccount') ? 'error' : ''} help={this.fieldError('billingAccount') || 'The billing account'}>
          {form.getFieldDecorator('billingAccount', {
            rules: [{ required: true, message: 'Please enter your billing account!' }],
            initialValue: data && data.spec.billingAccount
          })(
            <Input placeholder="Billing account" />,
          )}
        </Form.Item>

        {data ? (
          <>
            <Alert
              message="For security reasons, the Service Account key is not shown after creation of the organization credential"
              type="warning"
              style={{ marginTop: '10px' }}
            />
            <Form.Item label="Replace key">
              <Checkbox id="gcp_org_replace_key" onChange={(e) => this.setState({ replaceKey: e.target.checked })} />
            </Form.Item>
          </>
        ) : null}

        {!data || replaceKey ? (
          <Form.Item label="Service Account JSON" labelCol={{ span: 24 }} wrapperCol={{ span: 24 }} validateStatus={this.fieldError('account') ? 'error' : ''} help={this.fieldError('account') || Boolean(data) || 'The Service Account key in JSON format, with project creation permissions.'}>
            {form.getFieldDecorator('account', {
              rules: [{ required: true, message: 'Please enter your Service Account!' }]
            })(
              <Input.TextArea autoSize={{ minRows: 4, maxRows: 10  }} placeholder="Service Account JSON" />,
            )}
          </Form.Item>
        ) : null}

      </Card>
    )
  }
}

const WrappedGCPOrganizationForm = Form.create({ name: 'gcp_organization' })(GCPOrganizationForm)

export default WrappedGCPOrganizationForm
