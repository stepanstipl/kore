import * as React from 'react'
import VerifiedAllocatedResourceForm from '../../components/forms/VerifiedAllocatedResourceForm'
import V1alpha1Organization from '../../kore-api/model/V1alpha1Organization'
import V1ObjectMeta from '../../kore-api/model/V1ObjectMeta'
import V1Secret from '../../kore-api/model/V1Secret'
import V1SecretSpec from '../../kore-api/model/V1SecretSpec'
import V1SecretReference from '../../kore-api/model/V1SecretReference'
import V1alpha1OrganizationSpec from '../../kore-api/model/V1alpha1OrganizationSpec'
import KoreApi from '../../kore-api'
import { Form, Input, Alert, Card } from 'antd'
import AllocationHelpers from '../../utils/allocation-helpers'

class GCPOrganizationForm extends VerifiedAllocatedResourceForm {

  generateSecretResource = values => {
    const resource = new V1Secret()
    resource.setApiVersion('config.kore.appvia.io')
    resource.setKind('Secret')

    const meta = new V1ObjectMeta()
    meta.setName(this.getMetadataName(values))
    meta.setNamespace(this.props.team)
    resource.setMetadata(meta)

    const spec = new V1SecretSpec()
    spec.setType('gcp-org')
    spec.setDescription(`GCP admin project Service Account for ${values.parentID}`)
    spec.setData({ key: values.account })
    resource.setSpec(spec)

    return resource
  }

  generateGCPOrganizationResource = values => {
    const name = this.getMetadataName(values)
    const resource = new V1alpha1Organization()
    resource.setApiVersion('gcp.compute.kore.appvia.io/v1alpha1')
    resource.setKind('Organization')

    const meta = new V1ObjectMeta()
    meta.setName(name)
    meta.setNamespace(this.props.team)
    resource.setMetadata(meta)

    const spec = new V1alpha1OrganizationSpec()
    spec.setParentType('organization')
    spec.setParentID(values.parentID)
    spec.setBillingAccount(values.billingAccount)
    spec.setServiceAccount('kore')

    const secret = new V1SecretReference()
    secret.setName(name)
    secret.setNamespace(this.props.team)
    spec.setCredentialsRef(secret)

    resource.setSpec(spec)

    return resource
  }

  getResource = async metadataName => {
    const api = await KoreApi.client()
    const gcpOrgResult = await api.GetGCPOrganization(this.props.team, metadataName)
    gcpOrgResult.allocation = await AllocationHelpers.getAllocationForResource(gcpOrgResult)
    return gcpOrgResult
  }

  putResource = async values => {
    const api = await KoreApi.client()
    const metadataName = this.getMetadataName(values)
    await api.UpdateTeamSecret(this.props.team, metadataName, this.generateSecretResource(values))
    const gcpOrg = this.generateGCPOrganizationResource(values)
    const gcpOrgResult = await api.UpdateGCPOrganization(this.props.team, metadataName, gcpOrg)
    gcpOrgResult.allocation = await this.storeAllocation(gcpOrg, values)
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
        <Form.Item label="Service Account JSON" labelCol={{ span: 24 }} wrapperCol={{ span: 24 }} validateStatus={this.fieldError('account') ? 'error' : ''} help={this.fieldError('account') || Boolean(data) || 'The Service Account key in JSON format, with project creation permissions.'}>
          {!data ? (
            form.getFieldDecorator('account', {
              rules: [{ required: true, message: 'Please enter your Service Account!' }]
            })(
              <Input.TextArea autoSize={{ minRows: 4, maxRows: 10  }} placeholder="Service Account JSON" />,
            )
          ) : (
            <Alert
              message="For security reasons, the Service Account key is not shown after creation of the organization"
              type="warning"
              style={{ marginBottom: '-20px', marginTop: '10px' }}
            />
          )}
        </Form.Item>
      </Card>
    )
  }
}

const WrappedGCPOrganizationForm = Form.create({ name: 'gcp_organization' })(GCPOrganizationForm)

export default WrappedGCPOrganizationForm
