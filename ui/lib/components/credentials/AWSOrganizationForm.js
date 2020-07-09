import * as React from 'react'
import VerifiedAllocatedResourceForm from '../resources/VerifiedAllocatedResourceForm'
import KoreApi from '../../kore-api'
import { Checkbox, Form, Input, Alert, Card, Select } from 'antd'
const { Option } = Select
import AllocationHelpers from '../../utils/allocation-helpers'

class AWSOrganizationForm extends VerifiedAllocatedResourceForm {

  getResource = async metadataName => {
    const api = await KoreApi.client()
    const awsOrgResult = await api.GetAWSOrganization(this.props.team, metadataName)
    awsOrgResult.allocation = await AllocationHelpers.getAllocationForResource(awsOrgResult)
    return awsOrgResult
  }

  putResource = async values => {
    const api = await KoreApi.client()
    values.name = this.getMetadataName(values)
    const secretName = values.name
    const teamResources = KoreApi.resources().team(this.props.team)
    if (!this.props.data || this.state.replaceKey) {
      const secretData = {
        access_key_id: btoa(values.accessKeyID),
        access_secret_key: btoa(values.secretAccessKey)
      }
      const secretResource = teamResources.generateSecretResource(secretName, 'aws-org', `AWS creds for control tower in OU ${values.ouName}`, secretData)
      await api.UpdateTeamSecret(this.props.team, secretName, secretResource)
    }
    const awsOrgResource = teamResources.generateAWSOrganizationResource(values, secretName)
    const awsOrgResult = await api.UpdateAWSOrganization(this.props.team, values.name, awsOrgResource)
    awsOrgResult.allocation = await this.storeAllocation(awsOrgResource, values)
    return awsOrgResult
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
      infoDescription: 'This will give teams the ability to create AWS Accounts inside the organization. Within these scoped projects the team will have the ability to create clusters.',
      allTeamsWarningMessage: 'This organization credential will be available to all teams',
      allTeamsWarningDescription: 'No teams exist in Kore yet, therefore currently this organization credential will be available to all teams created in the future. If you wish to restrict this please return here and allocate to teams once they have been created.',
      allocateExtra: 'If nothing selected then this organization credential will be available to ALL teams'
    }
  }

  resourceFormFields = () => {
    const { form, data } = this.props
    const { replaceKey } = this.state
    return (
      <>
        <Card style={{ marginBottom: '20px' }}>
          <Alert
            message="AWS Organization details"
            description="Retrieve these values from your AWS organization. Providing these gives Kore the ability to create Accounts for teams within the AWS organization. Teams will then be able to provision clusters within their Accounts."
            type="info"
            style={{ marginBottom: '20px' }}
          />
          <Form.Item label="Organization Unit name" validateStatus={this.fieldError('ouName') ? 'error' : ''} help={this.fieldError('ouName') || 'The name of the parent Organizational Unit (OU) to use for provisioning accounts'}>
            {form.getFieldDecorator('ouName', {
              rules: [{ required: true, message: 'Please enter the Organization Unit name!' }],
              initialValue: (data && data.spec.ouName) || 'Custom'
            })(
              <Select>
                <Option value="Custom">Custom</Option>
              </Select>,
            )}
          </Form.Item>
          <Form.Item label="Region" validateStatus={this.fieldError('region') ? 'error' : ''} help={this.fieldError('region') || 'The region where Control Tower is enabled in the master account'}>
            {form.getFieldDecorator('region', {
              rules: [{ required: true, message: 'Please enter the region!' }],
              initialValue: data && data.spec.region
            })(
              <Input placeholder="Region" />,
            )}
          </Form.Item>
          <Form.Item label="Role ARN" validateStatus={this.fieldError('roleARN') ? 'error' : ''} help={this.fieldError('roleARN') || 'The role to assume when provisioning accounts'}>
            {form.getFieldDecorator('roleARN', {
              rules: [{ required: true, message: 'Please enter the role ARN!' }],
              initialValue: data && data.spec.roleARN
            })(
              <Input placeholder="Role ARN" />,
            )}
          </Form.Item>

          {data ? (
            <>
              <Alert
                message="For security reasons, the access key is not shown after creation of the organization credential"
                type="warning"
                style={{ marginTop: '10px' }}
              />
              <Form.Item label="Replace access key">
                <Checkbox id="awsorg_replace_key" onChange={(e) => this.setState({ replaceKey: e.target.checked })} />
              </Form.Item>
            </>
          ) : null}

          {!data || replaceKey ? (
            <>
              <Form.Item label="Access key ID" validateStatus={this.fieldError('accessKeyID') ? 'error' : ''} help={this.fieldError('accessKeyID') || 'The Access key ID part of the access key'}>
                {form.getFieldDecorator('accessKeyID', {
                  rules: [{ required: true, message: 'Please enter your Access key ID!' }],
                  initialValue: data && data.spec.accessKeyID
                })(
                  <Input placeholder="Access key ID" />,
                )}
              </Form.Item>
              <Form.Item label="Secret access key" validateStatus={this.fieldError('secretAccessKey') ? 'error' : ''} help={this.fieldError('secretAccessKey') || 'The Secret access key part of the access key'}>
                {form.getFieldDecorator('secretAccessKey', {
                  rules: [{ required: true, message: 'Please enter your Secret access key!' }],
                  initialValue: data && data.spec.secretAccessKey
                })(
                  <Input placeholder="Secret access key" type="password" />,
                )}
              </Form.Item>
            </>

          ) : null}

        </Card>


        <Card style={{ marginBottom: '20px' }}>
          <Alert
            message="SSO User"
            description="The user who will be the organisational account owner for all Kore provisioned accounts."
            type="info"
            style={{ marginBottom: '20px' }}
          />

          <Form.Item label="First name" validateStatus={this.fieldError('ssoUserFirstName') ? 'error' : ''} help={this.fieldError('ssoUserFirstName') || ''}>
            {form.getFieldDecorator('ssoUserFirstName', {
              rules: [{ required: true, message: 'Please enter the first name!' }],
              initialValue: data && data.spec.ssoUser.firstName
            })(
              <Input placeholder="First name" />,
            )}
          </Form.Item>
          <Form.Item label="Last name" validateStatus={this.fieldError('ssoUserLastName') ? 'error' : ''} help={this.fieldError('ssoUserLastName') || ''}>
            {form.getFieldDecorator('ssoUserLastName', {
              rules: [{ required: true, message: 'Please enter the last name!' }],
              initialValue: data && data.spec.ssoUser.lastName
            })(
              <Input placeholder="Last name" />,
            )}
          </Form.Item>
          <Form.Item label="Email address" validateStatus={this.fieldError('ssoUserEmailAddress') ? 'error' : ''} help={this.fieldError('ssoUserEmailAddress') || ''}>
            {form.getFieldDecorator('ssoUserEmailAddress', {
              rules: [{ required: true, message: 'Please enter the email address!' }],
              initialValue: data && data.spec.ssoUser.email
            })(
              <Input placeholder="Email address" />,
            )}
          </Form.Item>
        </Card>
      </>
    )
  }
}

const WrappedAWSOrganizationForm = Form.create({ name: 'aws_organization' })(AWSOrganizationForm)

export default WrappedAWSOrganizationForm
