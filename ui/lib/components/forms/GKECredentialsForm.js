import V1alpha1GKECredentials from '../../kore-api/model/V1alpha1GKECredentials'
import V1alpha1GKECredentialsSpec from '../../kore-api/model/V1alpha1GKECredentialsSpec'
import V1ObjectMeta from '../../kore-api/model/V1ObjectMeta'
import VerifiedAllocatedResourceForm from '../../components/forms/VerifiedAllocatedResourceForm'
import ResourceVerificationStatus from '../../components/ResourceVerificationStatus'
import FormErrorMessage from '../../components/forms/FormErrorMessage'
import KoreApi from '../../kore-api'
import { Button, Form, Input, Alert, Select, Card } from 'antd'

class GKECredentialsForm extends VerifiedAllocatedResourceForm {

  generateGKECredentialsResource = values => {
    const resource = new V1alpha1GKECredentials()
    resource.setApiVersion('gke.compute.kore.appvia.io/v1alpha1')
    resource.setKind('GKECredentials')

    const meta = new V1ObjectMeta()
    meta.setName(this.getMetadataName(values))
    meta.setNamespace(this.props.team)
    resource.setMetadata(meta)

    const spec = new V1alpha1GKECredentialsSpec()
    spec.setProject(values.project)
    spec.setAccount(values.account)

    resource.setSpec(spec)

    return resource
  }

  getResource = async metadataName => {
    const api = await KoreApi.client()
    const gkeCredentialsResult = await api.GetGKECredential(this.props.team, metadataName)
    gkeCredentialsResult.allocation = await api.GetAllocation(this.props.team, metadataName)
    return gkeCredentialsResult
  }

  putResource = async values => {
    const api = await KoreApi.client()
    const metadataName = this.getMetadataName(values)
    const gkeResult = await api.UpdateGKECredential(this.props.team, metadataName, this.generateGKECredentialsResource(values))
    const allocationResource = this.generateAllocationResource({ group: 'config.kore.appvia.io', version: 'v1', kind: 'GKECredentials' }, values)
    gkeResult.allocation = await api.UpdateAllocation(this.props.team, metadataName, allocationResource)
    return gkeResult
  }

  render() {
    const { form, data, allTeams, saveButtonText } = this.props
    const { getFieldDecorator, getFieldsError, getFieldError, isFieldTouched } = form
    const { formErrorMessage, allocations, submitting, inlineVerificationFailed } = this.state
    const formConfig = {
      layout: 'horizontal',
      labelAlign: 'left',
      hideRequiredMark: true,
      labelCol: {
        sm: { span: 24 },
        md: { span: 7 },
        lg: { span: 5 }
      },
      wrapperCol: {
        sm: { span: 24 },
        md: { span: 17 },
        lg: { span: 19 }
      }
    }

    // Only show error after a field is touched.
    const fieldError = fieldKey => isFieldTouched(fieldKey) && getFieldError(fieldKey)
    const allocationMissing = Boolean(data && !data.allocation)

    return (
      <div>

        <ResourceVerificationStatus resourceStatus={data && data.status} style={{ marginBottom: '15px' }}/>

        {allocationMissing ? (
          <Alert
            message="This project is not allocated to any teams"
            description="Give the project a Name and Description below and enter Allocated team(s) as appropriate. Once complete, click Save to allocate this project."
            type="warning"
            showIcon
            style={{ marginBottom: '20px', marginTop: '5px' }}
          />
        ) : null}

        <Form {...formConfig} onSubmit={this.handleSubmit}>
          <FormErrorMessage message={formErrorMessage} />
          <Card style={{ marginBottom: '20px' }}>
            <Alert
              message="Project and service account"
              description="Retrieve these values from your GCP project. Providing these gives Kore the ability to create clusters within the project."
              type="info"
              style={{ marginBottom: '20px' }}
            />
            <Form.Item label="Project name" validateStatus={fieldError('project') ? 'error' : ''} help={fieldError('project') || 'The GCP project that Kore will be able to build clusters within.'}>
              {getFieldDecorator('project', {
                rules: [{ required: true, message: 'Please enter your project name!' }],
                initialValue: data && data.spec.project
              })(
                <Input placeholder="Project" />,
              )}
            </Form.Item>
            <Form.Item label="Service Account JSON" labelCol={{ span: 24 }} wrapperCol={{ span: 24 }} validateStatus={fieldError('account') ? 'error' : ''} help={fieldError('account') || 'The Service Account key in JSON format, with GKE admin permissions on the GCP project'}>
              {getFieldDecorator('account', {
                rules: [{ required: true, message: 'Please enter your Service Account!' }],
                initialValue: data && data.spec.account
              })(
                <Input.TextArea autoSize={{ minRows: 4, maxRows: 10  }} placeholder="Service Account JSON" />,
              )}
            </Form.Item>
          </Card>

          <Card style={{ marginBottom: '20px' }}>
            <Alert
              message="Help Kore teams understand this project"
              description="Give this project a name and description to help teams choose the correct one."
              type="info"
              style={{ marginBottom: '20px' }}
            />
            <Form.Item label="Name" validateStatus={fieldError('name') ? 'error' : ''} help={fieldError('name') || 'The name for the project eg. MyOrg project-one'}>
              {getFieldDecorator('name', {
                rules: [{ required: true, message: 'Please enter the name!' }],
                initialValue: data && data.allocation && data.allocation.spec.name
              })(
                <Input placeholder="Name" />,
              )}
            </Form.Item>
            <Form.Item label="Description" validateStatus={fieldError('summary') ? 'error' : ''} help={fieldError('summary') || 'A description of the project to help when choosing between them'}>
              {getFieldDecorator('summary', {
                rules: [{ required: true, message: 'Please enter the description!' }],
                initialValue: data && data.allocation && data.allocation.spec.summary
              })(
                <Input placeholder="Description" />,
              )}
            </Form.Item>
          </Card>

          <Card style={{ marginBottom: '20px' }}>
            <Alert
              message="Make this project available to teams in Kore"
              description="This will give teams the ability to create clusters within the project."
              type="info"
              style={{ marginBottom: '20px' }}
            />

            {allTeams.items.length === 0 ? (
              <Alert
                message="This project will be available to all teams"
                description="No teams exist in Kore yet, therefore currently this project will be available to all teams created in the future. If you wish to restrict this please return here and allocate to teams once they have been created."
                type="warning"
                showIcon
              />
            ) : (
              <Form.Item label="Allocate team(s)" extra="If nothing selected then this project will be available to ALL teams">
                {getFieldDecorator('allocations', { initialValue: allocations })(
                  <Select
                    mode="multiple"
                    style={{ width: '100%' }}
                    placeholder={allocationMissing ? 'No teams' : 'All teams'}
                    onChange={this.onAllocationsChange}
                  >
                    {allTeams.items.map(t => (
                      <Select.Option key={t.metadata.name} value={t.metadata.name}>{t.spec.summary}</Select.Option>
                    ))}
                  </Select>
                )}
              </Form.Item>
            )}

          </Card>
          <Form.Item style={{ marginBottom: '0' }}>
            <Button type="primary" htmlType="submit" loading={submitting} disabled={this.disableButton(getFieldsError())}>{saveButtonText || 'Save'}</Button>
            {inlineVerificationFailed ? (
              <Button onClick={this.continueWithoutVerification} disabled={this.disableButton(getFieldsError())} style={{ marginLeft: '10px' }}>Continue without verification</Button>
            ) : null}
          </Form.Item>
        </Form>
      </div>
    )
  }
}

const WrappedGKECredentialsForm = Form.create({ name: 'gke_credentials' })(GKECredentialsForm)

export default WrappedGKECredentialsForm
