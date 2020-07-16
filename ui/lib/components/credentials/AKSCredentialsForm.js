import { Form, Input, Alert, Card, Checkbox } from 'antd'

import VerifiedAllocatedResourceForm from '../resources/VerifiedAllocatedResourceForm'
import KoreApi from '../../kore-api'
import AllocationHelpers from '../../utils/allocation-helpers'

class AKSCredentialsForm extends VerifiedAllocatedResourceForm {

  getResource = async (metadataName) => {
    const api = await KoreApi.client()
    const aksCredentialsResult = await api.GetAKSCredentials(this.props.team, metadataName)
    aksCredentialsResult.allocation = await AllocationHelpers.getAllocationForResource(aksCredentialsResult)
    return aksCredentialsResult
  }

  putResource = async (values) => {
    const api = await KoreApi.client()
    const resourceName = this.getMetadataName(values)
    const secretName = resourceName
    const teamResources = KoreApi.resources().team(this.props.team)
    if (!this.props.data || this.state.replaceKey) {
      const secretData = {
        subscription_id: btoa(values.subscriptionID),
        tenant_id: btoa(values.tenantID),
        client_id: btoa(values.clientID),
        client_secret: btoa(values.clientSecret)
      }
      const secretResource = teamResources.generateSecretResource(secretName, 'azure-credentials', `Azure subscription ${values.subscriptionID} credential`, secretData)
      await api.UpdateTeamSecret(this.props.team, secretName, secretResource)
    }
    const aksCredResource = teamResources.generateAKSCredentialsResource(resourceName, values, secretName)
    const aksResult = await api.UpdateAKSCredentials(this.props.team, values.name, aksCredResource)
    aksResult.allocation = await this.storeAllocation(aksCredResource, values)
    return aksResult
  }

  allocationFormFieldsInfo = {
    allocationMissing: {
      infoMessage: 'This subscription credential is not allocated to any teams',
      infoDescription: 'Give the subscription credential a Name and Description below and enter Allocated team(s) as appropriate. Once complete, click Save to allocate it.'
    },
    nameSection: {
      infoMessage: 'Help Kore teams understand this subscription credential',
      infoDescription: 'Give this subscription credential a name and description to help teams choose the correct one.',
      nameHelp: 'The name for the subscription credential eg. MyOrg subscription-one',
      descriptionHelp: 'A description of the subscription credential to help when choosing between them'
    },
    allocationSection: {
      infoMessage: 'Make this subscription credential available to teams in Kore',
      infoDescription: 'This will give teams the ability to create clusters within the subscription.',
      allTeamsWarningMessage: 'This subscription credential will be available to all teams',
      allTeamsWarningDescription: 'No teams exist in Kore yet, therefore currently this subscription credential will be available to all teams created in the future. If you wish to restrict this please return here and allocate to teams once they have been created.',
      allocateExtra: 'If nothing selected then this subscription credential will be available to ALL teams'
    }
  }

  resourceFormFields = () => {
    const { form, data } = this.props
    const { replaceKey } = this.state
    return (
      <Card style={{ marginBottom: '20px' }}>
        <Alert
          message="Subscription"
          description="Retrieve these values from Azure. Providing these gives Kore the ability to create clusters within the subscription."
          type="info"
          style={{ marginBottom: '20px' }}
        />
        <Form.Item label="Subscription ID" validateStatus={this.fieldError('subscriptionID') ? 'error' : ''} help={this.fieldError('subscriptionID') || 'The Azure subscription that Kore will be able to build clusters within.'}>
          {form.getFieldDecorator('subscriptionID', {
            rules: [{ required: true, message: 'Please enter your subscription ID!' }],
            initialValue: data && data.spec.subscriptionID
          })(
            <Input placeholder="Subscription ID" />,
          )}
        </Form.Item>
        <Form.Item label="Tenant" validateStatus={this.fieldError('tenantID') ? 'error' : ''} help={this.fieldError('tenantID') || 'The tenant for the service principal scoped to the subscription.'}>
          {form.getFieldDecorator('tenantID', {
            rules: [{ required: true, message: 'Please enter your tenant!' }],
            initialValue: data && data.spec.tenantID
          })(
            <Input placeholder="Tenant" />,
          )}
        </Form.Item>
        <Form.Item label="App ID" validateStatus={this.fieldError('clientID') ? 'error' : ''} help={this.fieldError('clientID') || 'The App ID for the service principal scoped to the subscription.'}>
          {form.getFieldDecorator('clientID', {
            rules: [{ required: true, message: 'Please enter your App ID!' }],
            initialValue: data && data.spec.clientID
          })(
            <Input placeholder="App ID" />,
          )}
        </Form.Item>

        {data ? (
          <>
            <Alert
              message="For security reasons, the password cannot be retrieved after creation. If you need to replace the password, tick the box below."
              type="warning"
              style={{ marginTop: '10px' }}
            />
            <Form.Item label="Replace password">
              <Checkbox id="aks_credentials_replace_key" onChange={(e) => this.setState({ replaceKey: e.target.checked })}></Checkbox>
            </Form.Item>
          </>
        ) : null}

        {!data || replaceKey ? (
          <>
            <Form.Item label="Password" validateStatus={this.fieldError('clientSecret') ? 'error' : ''} help={this.fieldError('clientSecret') || 'The password for the service principal scoped to the subscription.'}>
              {form.getFieldDecorator('clientSecret', {
                rules: [{ required: true, message: 'Please enter your password!' }],
                initialValue: data && data.spec.accessKeyID
              })(
                <Input placeholder="Password" type="password" />,
              )}
            </Form.Item>
          </>
        ) : null}

      </Card>
    )
  }
}

const WrappedAKSCredentialsForm = Form.create({ name: 'aks_credentials' })(AKSCredentialsForm)

export default WrappedAKSCredentialsForm
