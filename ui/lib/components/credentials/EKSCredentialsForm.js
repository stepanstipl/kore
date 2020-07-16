import { Form, Input, Alert, Card, Checkbox } from 'antd'

import VerifiedAllocatedResourceForm from '../resources/VerifiedAllocatedResourceForm'
import KoreApi from '../../kore-api'
import AllocationHelpers from '../../utils/allocation-helpers'

class EKSCredentialsForm extends VerifiedAllocatedResourceForm {

  getResource = async (metadataName) => {
    const api = await KoreApi.client()
    const eksCredentialsResult = await api.GetEKSCredentials(this.props.team, metadataName)
    eksCredentialsResult.allocation = await AllocationHelpers.getAllocationForResource(eksCredentialsResult)
    return eksCredentialsResult
  }

  putResource = async (values) => {
    const api = await KoreApi.client()
    const resourceName = this.getMetadataName(values)
    const secretName = resourceName
    const teamResources = KoreApi.resources().team(this.props.team)
    if (!this.props.data || this.state.replaceKey) {
      const secretData = {
        access_key_id: btoa(values.accessKeyID),
        access_secret_key: btoa(values.secretAccessKey)
      }
      const secretResource = teamResources.generateSecretResource(secretName, 'aws-credentials', `AWS account ${values.accountID} credential`, secretData)
      await api.UpdateTeamSecret(this.props.team, secretName, secretResource)
    }
    const eksCredResource = teamResources.generateEKSCredentialsResource(resourceName, values, secretName)
    const eksResult = await api.UpdateEKSCredentials(this.props.team, values.name, eksCredResource)
    eksResult.allocation = await this.storeAllocation(eksCredResource, values)
    return eksResult
  }

  allocationFormFieldsInfo = {
    allocationMissing: {
      infoMessage: 'This account credential is not allocated to any teams',
      infoDescription: 'Give the account credential a Name and Description below and enter Allocated team(s) as appropriate. Once complete, click Save to allocate it.'
    },
    nameSection: {
      infoMessage: 'Help Kore teams understand this account credential',
      infoDescription: 'Give this account credential a name and description to help teams choose the correct one.',
      nameHelp: 'The name for the account credential eg. MyOrg project-one',
      descriptionHelp: 'A description of the account credential to help when choosing between them'
    },
    allocationSection: {
      infoMessage: 'Make this account credential available to teams in Kore',
      infoDescription: 'This will give teams the ability to create clusters within the account.',
      allTeamsWarningMessage: 'This account credential will be available to all teams',
      allTeamsWarningDescription: 'No teams exist in Kore yet, therefore currently this account credential will be available to all teams created in the future. If you wish to restrict this please return here and allocate to teams once they have been created.',
      allocateExtra: 'If nothing selected then this account credential will be available to ALL teams'
    }
  }

  resourceFormFields = () => {
    const { form, data } = this.props
    const { replaceKey } = this.state
    return (
      <Card style={{ marginBottom: '20px' }}>
        <Alert
          message="Account and access key"
          description="Retrieve these values from your AWS account. Providing these gives Kore the ability to create clusters within the account."
          type="info"
          style={{ marginBottom: '20px' }}
        />
        <Form.Item label="Account ID" validateStatus={this.fieldError('accountID') ? 'error' : ''} help={this.fieldError('accountID') || 'The AWS account that Kore will be able to build clusters within.'}>
          {form.getFieldDecorator('accountID', {
            rules: [{ required: true, message: 'Please enter your account ID!' }],
            initialValue: data && data.spec.accountID
          })(
            <Input placeholder="Account ID" />,
          )}
        </Form.Item>

        {data ? (
          <>
            <Alert
              message="For security reasons, the access key cannot be retrieved after creation. If you need to replace the key, tick the box below."
              type="warning"
              style={{ marginTop: '10px' }}
            />
            <Form.Item label="Replace access key">
              <Checkbox id="eks_credentials_replace_key" onChange={(e) => this.setState({ replaceKey: e.target.checked })}></Checkbox>
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
    )
  }
}

const WrappedEKSCredentialsForm = Form.create({ name: 'eks_credentials' })(EKSCredentialsForm)

export default WrappedEKSCredentialsForm
