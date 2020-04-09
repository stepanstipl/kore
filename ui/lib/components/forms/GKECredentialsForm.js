import V1alpha1GKECredentials from '../../kore-api/model/V1alpha1GKECredentials'
import V1alpha1GKECredentialsSpec from '../../kore-api/model/V1alpha1GKECredentialsSpec'
import V1ObjectMeta from '../../kore-api/model/V1ObjectMeta'
import VerifiedAllocatedResourceForm from '../../components/forms/VerifiedAllocatedResourceForm'
import KoreApi from '../../kore-api'
import { Form, Input, Alert, Card } from 'antd'

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

  allocationFormFieldsInfo = {
    allocationMissing: {
      infoMessage: 'This project credential is not allocated to any teams',
      infoDescription: "Give the project credential a Name and Description below and enter Allocated team(s) as appropriate. Once complete, click Save to allocate it."
    },
    nameSection: {
      infoMessage: 'Help Kore teams understand this project credential',
      infoDescription: 'Give this project credential a name and description to help teams choose the correct one.',
      nameHelp: 'The name for the project credential eg. MyOrg project-one',
      descriptionHelp: 'A description of the project credential to help when choosing between them'
    },
    allocationSection: {
      infoMessage: 'Make this project credential available to teams in Kore',
      infoDescription: 'This will give teams the ability to create clusters within the project.',
      allTeamsWarningMessage: 'This project credential will be available to all teams',
      allTeamsWarningDescription: 'No teams exist in Kore yet, therefore currently this project credential will be available to all teams created in the future. If you wish to restrict this please return here and allocate to teams once they have been created.',
      allocateExtra: 'If nothing selected then this project will credential be available to ALL teams'
    }
  }

  resourceFormFields = () => {
    const { form, data } = this.props
    return (
      <Card style={{ marginBottom: '20px' }}>
        <Alert
          message="Project name and service account"
          description="Retrieve these values from your GCP project. Providing these gives Kore the ability to create clusters within the project."
          type="info"
          style={{ marginBottom: '20px' }}
        />
        <Form.Item label="Project name" validateStatus={this.fieldError('project') ? 'error' : ''} help={this.fieldError('project') || 'The GCP project that Kore will be able to build clusters within.'}>
          {form.getFieldDecorator('project', {
            rules: [{ required: true, message: 'Please enter your project name!' }],
            initialValue: data && data.spec.project
          })(
            <Input placeholder="Project" />,
          )}
        </Form.Item>
        <Form.Item label="Service Account JSON" labelCol={{ span: 24 }} wrapperCol={{ span: 24 }} validateStatus={this.fieldError('account') ? 'error' : ''} help={this.fieldError('account') || 'The Service Account key in JSON format, with GKE admin permissions on the GCP project'}>
          {form.getFieldDecorator('account', {
            rules: [{ required: true, message: 'Please enter your Service Account!' }],
            initialValue: data && data.spec.account
          })(
            <Input.TextArea autoSize={{ minRows: 4, maxRows: 10  }} placeholder="Service Account JSON" />,
          )}
        </Form.Item>
      </Card>
    )
  }
}

const WrappedGKECredentialsForm = Form.create({ name: 'gke_credentials' })(GKECredentialsForm)

export default WrappedGKECredentialsForm
