import V1alpha1EKSCredentials from '../../kore-api/model/V1alpha1EKSCredentials'
import V1alpha1EKSCredentialsSpec from '../../kore-api/model/V1alpha1EKSCredentialsSpec'
import V1ObjectMeta from '../../kore-api/model/V1ObjectMeta'
import V1Secret from '../../kore-api/model/V1Secret'
import V1SecretSpec from '../../kore-api/model/V1SecretSpec'
import V1SecretReference from '../../kore-api/model/V1SecretReference'
import VerifiedAllocatedResourceForm from '../resources/VerifiedAllocatedResourceForm'
import KoreApi from '../../kore-api'
import { Form, Input, Alert, Card } from 'antd'
import AllocationHelpers from '../../utils/allocation-helpers'

class EKSCredentialsForm extends VerifiedAllocatedResourceForm {

  generateSecretResource = values => {
    const resource = new V1Secret()
    resource.setApiVersion('config.kore.appvia.io')
    resource.setKind('Secret')

    const meta = new V1ObjectMeta()
    meta.setName(this.getMetadataName(values))
    meta.setNamespace(this.props.team)
    resource.setMetadata(meta)

    const spec = new V1SecretSpec()
    spec.setType('aws-credentials')
    spec.setDescription(`AWS account ${values.accountID} credential`)
    spec.setData({
      access_key_id: btoa(values.accessKeyID),
      access_secret_key: btoa(values.secretAccessKey)
    })
    resource.setSpec(spec)

    return resource
  }

  generateEKSCredentialsResource = values => {
    const name = this.getMetadataName(values)
    const resource = new V1alpha1EKSCredentials()
    resource.setApiVersion('aws.compute.kore.appvia.io/v1alpha1')
    resource.setKind('EKSCredentials')

    const meta = new V1ObjectMeta()
    meta.setName(name)
    meta.setNamespace(this.props.team)
    resource.setMetadata(meta)

    const spec = new V1alpha1EKSCredentialsSpec()
    spec.setAccountID(values.accountID)

    const secret = new V1SecretReference()
    secret.setName(name)
    secret.setNamespace(this.props.team)
    spec.setCredentialsRef(secret)

    resource.setSpec(spec)

    return resource
  }

  getResource = async metadataName => {
    const api = await KoreApi.client()
    const eksCredentialsResult = await api.GetEKSCredentials(this.props.team, metadataName)
    eksCredentialsResult.allocation = await AllocationHelpers.getAllocationForResource(eksCredentialsResult)
    return eksCredentialsResult
  }

  putResource = async values => {
    const api = await KoreApi.client()
    const metadataName = this.getMetadataName(values)
    await api.UpdateTeamSecret(this.props.team, metadataName, this.generateSecretResource(values))
    const eksCredResource = this.generateEKSCredentialsResource(values)
    const eksResult = await api.UpdateEKSCredentials(this.props.team, metadataName, eksCredResource)
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

        {!data ? (
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
        ) : (
          <Alert
            message="For security reasons, the access key is not shown after creating the account credential"
            type="warning"
            style={{ marginTop: '10px' }}
          />
        )}

      </Card>
    )
  }
}

const WrappedEKSCredentialsForm = Form.create({ name: 'eks_credentials' })(EKSCredentialsForm)

export default WrappedEKSCredentialsForm
