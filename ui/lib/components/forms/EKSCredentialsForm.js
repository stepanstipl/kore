import * as React from 'react'
import PropTypes from 'prop-types'
import copy from '../../utils/object-copy'
import canonical from '../../utils/canonical'
import V1alpha1EKSCredentials from '../../kore-api/model/V1alpha1EKSCredentials'
import V1alpha1EKSCredentialsSpec from '../../kore-api/model/V1alpha1EKSCredentialsSpec'
import V1Allocation from '../../kore-api/model/V1Allocation'
import V1AllocationSpec from '../../kore-api/model/V1AllocationSpec'
import V1ObjectMeta from '../../kore-api/model/V1ObjectMeta'
import V1Ownership from '../../kore-api/model/V1Ownership'
import ResourceVerificationStatus from '../../components/ResourceVerificationStatus'
import KoreApi from '../../kore-api'
import { Button, Form, Input, Alert, Select, message, Typography, Card } from 'antd'
const { Paragraph, Text } = Typography

class EKSCredentialsForm extends React.Component {
  static propTypes = {
    form: PropTypes.any.isRequired,
    team: PropTypes.string.isRequired,
    allTeams: PropTypes.object,
    data: PropTypes.object,
    handleSubmit: PropTypes.func.isRequired,
    saveButtonText: PropTypes.string,
    inlineVerification: PropTypes.bool
  }

  constructor(props) {
    super(props)
    let allocations = []
    if (props.data && props.data.allocation) {
      allocations = props.data.allocation.spec.teams.filter(a => a !== '*')
    }
    this.state = {
      submitting: false,
      formErrorMessage: false,
      allocations,
      inlineVerificationFailed: false
    }
  }

  componentDidMount() {
    // To disabled submit button at the beginning.
    this.props.form.validateFields()
  }

  disableButton = fieldsError => {
    if (this.state.submitting) {
      return true
    }
    return Object.keys(fieldsError).some(field => fieldsError[field])
  }

  onAllocationsChange = value => {
    const state = copy(this.state)
    state.allocations = value
    this.setState(state)
  }

  async verify(eksCredentials, tryCount) {
    const messageKey = 'verify'
    tryCount = tryCount || 0
    if (tryCount === 0) {
      message.loading({ content: 'Verifying AWS account credentials', key: messageKey, duration: 0 })
    }
    if (tryCount === 3) {
      message.error({ content: 'AWS account credentials verification failed', key: messageKey })
      const state = copy(this.state)
      state.inlineVerificationFailed = true
      state.submitting = false
      state.formErrorMessage = (
        <div>
          <Paragraph>The credentials have been saved but could not be verified, see the error below. Please try again or click &quot;Continue without verification&quot;.</Paragraph>
          {(eksCredentials.status.conditions || []).map((c, idx) =>
            <Paragraph key={idx} style={{ marginBottom: '0' }}>
              <Text strong>{c.message}</Text>
              <br/>
              <Text>{c.detail}</Text>
            </Paragraph>
          )}
        </div>
      )
      this.setState(state)
    } else {
      setTimeout(async () => {
        const api = await KoreApi.client()
        const eksCredentialsResult = await api.GetEKSCredentials(this.props.team, eksCredentials.metadata.name)
        eksCredentialsResult.allocation = eksCredentials.allocation
        if (eksCredentialsResult.status.status === 'Success') {
          message.success({ content: 'AWS account credentials verification successful', key: messageKey })
          return await this.props.handleSubmit(eksCredentialsResult)
        } else {
          return await this.verify(eksCredentialsResult, tryCount + 1)
        }
      }, 2000)
    }
  }

  setFormSubmitting = (submitting, formErrorMessage = false) => {
    this.setState({
      ...this.state,
      submitting,
      formErrorMessage
    })
  }

  getMetadataName = values => {
    const data = this.props.data
    return (data && data.metadata && data.metadata.name) || canonical(values.name)
  }

  getEKSCredentialsResource = values => {
    const resource = new V1alpha1EKSCredentials()
    resource.setApiVersion('aws.compute.kore.appvia.io/v1alpha1')
    resource.setKind('EKSCredentials')

    const meta = new V1ObjectMeta()
    meta.setName(this.getMetadataName(values))
    meta.setNamespace(this.props.team)
    resource.setMetadata(meta)

    const spec = new V1alpha1EKSCredentialsSpec()
    spec.setAccountID(values.accountID)
    spec.setAccessKeyID(values.accessKeyID)
    spec.setSecretAccessKey(values.secretAccessKey)

    resource.setSpec(spec)

    return resource
  }

  getAllocationResource = values => {
    const metadataName = this.getMetadataName(values)

    const resource = new V1Allocation()
    resource.setApiVersion('config.kore.appvia.io/v1')
    resource.setKind('Allocation')

    const meta = new V1ObjectMeta()
    meta.setName(metadataName)
    meta.setNamespace(this.props.team)
    resource.setMetadata(meta)

    const spec = new V1AllocationSpec()
    spec.setName(values.name)
    spec.setSummary(values.summary)
    spec.setTeams(this.state.allocations.length > 0 ? this.state.allocations : ['*'])
    const owner = new V1Ownership()
    owner.setGroup('aws.compute.kore.appvia.io')
    owner.setVersion('v1alpha1')
    owner.setKind('EKSCredentials')
    owner.setName(metadataName)
    owner.setNamespace(this.props.team)
    spec.setResource(owner)

    resource.setSpec(spec)

    return resource
  }

  handleSubmit = e => {
    e.preventDefault()

    this.setFormSubmitting(true)

    return this.props.form.validateFields(async (err, values) => {
      if (err) {
        this.setFormSubmitting(false, 'Validation failed')
        return
      }

      const metadataName = this.getMetadataName(values)
      const eksCredentialsResource = this.getEKSCredentialsResource(values)
      const allocationResource = this.getAllocationResource(values)

      try {
        const api = await KoreApi.client()
        const eksResult = await api.UpdateEKSCredentials(this.props.team, metadataName, eksCredentialsResource)
        eksResult.allocation = await api.UpdateAllocation(this.props.team, metadataName, allocationResource)

        if (this.props.inlineVerification) {
          await this.verify(eksResult)
        } else {
          await this.props.handleSubmit(eksResult)
        }
      } catch (err) {
        console.error('Error submitting form', err)
        this.setFormSubmitting(false, 'An error occurred saving the account, please try again')
      }
    })
  }

  continueWithoutVerification = async () => {
    try {
      const metadataName = this.getMetadataName(this.props.form.getFieldsValue())
      const api = await KoreApi.client()
      const eksCredentialsResult = await api.GetEKSCredentials(this.props.team, metadataName)
      eksCredentialsResult.allocation = await api.GetAllocation(this.props.team, metadataName)
      await this.props.handleSubmit(eksCredentialsResult)
    } catch (err) {
      console.error('Error getting data', err)
      this.setFormSubmitting(false, 'An error occurred saving the account, please try again')
    }
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
        md: { span: 8 },
        lg: { span: 6 }
      },
      wrapperCol: {
        sm: { span: 24 },
        md: { span: 16 },
        lg: { span: 18 }
      }
    }

    // Only show error after a field is touched.
    const fieldError = fieldKey => isFieldTouched(fieldKey) && getFieldError(fieldKey)
    const allocationMissing = Boolean(data && !data.allocation)

    const FormErrorMessage = () => {
      if (formErrorMessage) {
        return (
          <Alert
            message={formErrorMessage}
            type="error"
            showIcon
            closable
            style={{ marginBottom: '20px'}}
          />
        )
      }
      return null
    }

    return (
      <div>

        <ResourceVerificationStatus resourceStatus={data && data.status} style={{ marginBottom: '15px' }}/>

        {allocationMissing ? (
          <Alert
            message="This account is not allocated to any teams"
            description="Give the account a Name and Description below and enter Allocated team(s) as appropriate. Once complete, click Save to allocate this account."
            type="warning"
            showIcon
            style={{ marginBottom: '20px', marginTop: '5px' }}
          />
        ) : null}

        <Form {...formConfig} onSubmit={this.handleSubmit}>
          <FormErrorMessage />
          <Card style={{ marginBottom: '20px' }}>
            <Alert
              message="Account and access key"
              description="Retrieve these values from your AWS account. Providing these gives Kore the ability to create clusters within the account."
              type="info"
              style={{ marginBottom: '20px' }}
            />
            <Form.Item label="Account ID" validateStatus={fieldError('accountID') ? 'error' : ''} help={fieldError('accountID') || 'The AWS account that Kore will be able to build clusters within.'}>
              {getFieldDecorator('accountID', {
                rules: [{ required: true, message: 'Please enter your account ID!' }],
                initialValue: data && data.spec.accountID
              })(
                <Input placeholder="Account ID" />,
              )}
            </Form.Item>
            <Form.Item label="Access key ID" validateStatus={fieldError('accessKeyID') ? 'error' : ''} help={fieldError('accessKeyID') || 'The Access key ID part of the access key'}>
              {getFieldDecorator('accessKeyID', {
                rules: [{ required: true, message: 'Please enter your Access key ID!' }],
                initialValue: data && data.spec.accessKeyID
              })(
                <Input placeholder="Access key ID" />,
              )}
            </Form.Item>
            <Form.Item label="Secret access key" validateStatus={fieldError('secretAccessKey') ? 'error' : ''} help={fieldError('secretAccessKey') || 'The Secret access key part of the access key'}>
              {getFieldDecorator('secretAccessKey', {
                rules: [{ required: true, message: 'Please enter your Secret access key!' }],
                initialValue: data && data.spec.secretAccessKey
              })(
                <Input placeholder="Secret access key" type="password" />,
              )}
            </Form.Item>
          </Card>

          <Card style={{ marginBottom: '20px' }}>
            <Alert
              message="Help Kore teams understand this account"
              description="Give this account a name and description to help teams choose the correct one."
              type="info"
              style={{ marginBottom: '20px' }}
            />
            <Form.Item label="Name" validateStatus={fieldError('name') ? 'error' : ''} help={fieldError('name') || 'The name for the account eg. MyOrg project-one'}>
              {getFieldDecorator('name', {
                rules: [{ required: true, message: 'Please enter the name!' }],
                initialValue: data && data.allocation && data.allocation.spec.name
              })(
                <Input placeholder="Name" />,
              )}
            </Form.Item>
            <Form.Item label="Description" validateStatus={fieldError('summary') ? 'error' : ''} help={fieldError('summary') || 'A description of the account to help when choosing between them'}>
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
              message="Make this account available to teams in Kore"
              description="This will give teams the ability to create clusters within the account."
              type="info"
              style={{ marginBottom: '20px' }}
            />

            {allTeams.items.length === 0 ? (
              <Alert
                message="This account will be available to all teams"
                description="No teams exist in Kore yet, therefore currently this account will be available to all teams created in the future. If you wish to restrict this please return here and allocate to teams once they have been created."
                type="warning"
                showIcon
              />
            ) : (
              <Form.Item label="Allocate team(s)" extra="If nothing selected then this account will be available to ALL teams">
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
          <Form.Item style={{ marginBottom: '0'}}>
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

const WrappedEKSCredentialsForm = Form.create({ name: 'eks_credentials' })(EKSCredentialsForm)

export default WrappedEKSCredentialsForm
