import * as React from 'react'
import PropTypes from 'prop-types'
import { Checkbox, Form, Input, Alert, Card, Select, Typography } from 'antd'
const { Option } = Select
const { Paragraph } = Typography

import VerifiedAllocatedResourceForm from '../resources/VerifiedAllocatedResourceForm'
import KoreApi from '../../kore-api'
import AllocationHelpers from '../../utils/allocation-helpers'
import { patterns } from '../../utils/validation'

class AWSOrganizationForm extends VerifiedAllocatedResourceForm {

  static propTypes = {
    user: PropTypes.object.isRequired
  }

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
      const secretResource = teamResources.generateSecretResource(secretName, 'aws-credentials', `AWS creds for control tower in OU ${values.ouName}`, secretData)
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

  checkEmailDomain = (rule, value) => {
    // it's not an email address, let the standard validation pick it up
    const d = (e) => e.substr(e.indexOf('@') + 1)
    if (value.indexOf('@') === -1) {
      return Promise.resolve()
    }
    const domain = d(value)
    const userEmailDomain = d(this.props.user.email)
    if (domain === userEmailDomain) {
      return Promise.resolve()
    }
    return Promise.reject(`This email must have the same domain as the user (${userEmailDomain})`)
  }

  resourceFormFields = () => {
    const { form, data } = this.props
    const { replaceKey } = this.state
    return (
      <>
        <Card style={{ marginBottom: '20px' }}>
          <Alert
            message="AWS Organization Access"
            type="info"
            description={
              <>
                <Paragraph>
                  Kore can create AWS accounts within your AWS Organization as required for teams.
                </Paragraph>
                <Paragraph>
                  See <a target="_blank" rel="noopener noreferrer" style={{ textDecoration:'underline' }} href="https://docs.appvia.io/kore/guide/admin/aws_accounting/#1-pre-requisites">AWS Account Management pre-requisites</a>.
                </Paragraph>
                <Paragraph>
                  Providing these details grants Kore the ability to create Accounts for teams within the AWS organization. Teams will then be able to provision clusters within their Accounts.
                </Paragraph>
                <Paragraph style={{ marginBottom: 0 }}>
                  See <a target="_blank" rel="noopener noreferrer" style={{ textDecoration:'underline' }} href="https://docs.appvia.io/kore/guide/admin/aws_accounting/#aws-account-factory-access">AWS Account Factory Access</a> for how to obtain suitable access details.
                </Paragraph>
              </>
            }
            style={{ marginBottom: '20px' }}
          />
          <Alert
            message={ <> The Control Tower region does <b>not</b> dictate where EKS clusters or applications will run </> }
            type="warning"
            style={{ marginTop: '10px' }}
          />
          <Form.Item label="Region" validateStatus={this.fieldError('region') ? 'error' : ''} help={this.fieldError('region') || 'The region where AWS Control Tower is enabled in the master account'}>
            {form.getFieldDecorator('region', {
              rules: [{ required: true, message: 'Please select the region where Control Tower is enabled' }],
              initialValue: data && data.spec.region
            })(
              <Select>
                <Option value="us-east-1">us-east-1 (Northern Virginia)</Option>
                <Option value="us-east-2">us-east-2 (Ohio)</Option>
                <Option value="us-west-2">us-west-2 (Oregon)</Option>
                <Option value="eu-west-1">eu-west-1 (Ireland)</Option>
                <Option value="ap-southeast-2">ap-southeast-2 (Sydney)</Option>
              </Select>,
            )}
          </Form.Item>
          <Form.Item label="Role ARN" validateStatus={this.fieldError('roleARN') ? 'error' : ''} help={this.fieldError('roleARN') || 'The role to assume when provisioning accounts.'}>
            {form.getFieldDecorator('roleARN', {
              rules: [{ required: true, message: 'Please enter the role ARN!' }, patterns.amazonIamRoleArn],
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
            message="AWS Organization Details"
            description="Kore needs to know where to create AWS Accounts in the AWS Organisation"
            type="info"
            style={{ marginBottom: '20px' }}
          />

          <Form.Item label="Organization Unit name" validateStatus={this.fieldError('ouName') ? 'error' : ''} help={this.fieldError('ouName') || 'The name of the parent Organizational Unit (OU) to use for provisioning accounts'}>
            {form.getFieldDecorator('ouName', {
              rules: [{ required: true, message: 'Please select the Organization Unit name!' }],
              initialValue: (data && data.spec.ouName) || 'Custom'
            })(
              <Select>
                <Option value="Custom">Custom</Option>
              </Select>,
            )}
          </Form.Item>
        </Card>

        <Card style={{ marginBottom: '20px' }}>
          <Alert
            message="SSO User"
            type="info"
            description={
              <>
                <Paragraph>
                  The organization account owner for all Kore provisioned accounts.
                </Paragraph>
                <Paragraph style={{ marginBottom: 0 }}>
                  See <a target="_blank" rel="noopener noreferrer" style={{ textDecoration:'underline' }} href="https://docs.appvia.io/kore/guide/admin/aws_accounting/#sso-user-for-aws-account-administration">SSO User for AWS Account Administration</a> for more details.
                </Paragraph>
              </>
            }
          />
          <Alert
            message="The user will have root access to accounts created and a secure email address owned by the organization is required."
            type="warning"
            showIcon
            style={{ marginBottom: '20px', marginTop: '20px' }}
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
              rules: [
                { required: true, message: 'Please enter the email address!' },
                patterns.email,
                { validator: this.checkEmailDomain }
              ],
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
