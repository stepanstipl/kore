import * as React from 'react'
import PropTypes from 'prop-types'
import copy from '../../utils/object-copy'
import canonical from '../../utils/canonical'
import GKECredentials from '../../crd/GKECredentials'
import ResourceVerificationStatus from '../../components/ResourceVerificationStatus'
import Allocation from '../../crd/Allocation'
import apiRequest from '../../utils/api-request'
import apiPaths from '../../utils/api-paths'
import { Button, Form, Input, Alert, Select, message, Typography, Card } from 'antd'

const { Paragraph, Text } = Typography

class GKECredentialsForm extends React.Component {
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

  async verify(gke, tryCount) {
    const messageKey = 'verify'
    tryCount = tryCount || 0
    if (tryCount === 0) {
      message.loading({ content: 'Verifying GCP project Service Account', key: messageKey, duration: 0 })
    }
    if (tryCount === 3) {
      message.error({ content: 'GCP project Service Account verification failed', key: messageKey })
      const state = copy(this.state)
      state.inlineVerificationFailed = true
      state.submitting = false
      state.formErrorMessage = (
        <div>
          <Paragraph>The credentials have been saved but could not be verified, see the error below. Please try again or click &quot;Continue without verification&quot;.</Paragraph>
          {(gke.status.conditions || []).map((c, idx) =>
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
        const gkeResult = await apiRequest(null, 'get', `${apiPaths.team(this.props.team).gkeCredentials}/${gke.metadata.name}`)
        if (gkeResult.status.status === 'Success') {
          message.success({ content: 'GCP project Service Account verification successful', key: messageKey })
          return await this.props.handleSubmit(gkeResult)
        } else {
          return await this.verify(gkeResult, tryCount + 1)
        }
      }, 2000)
    }
  }

  continueWithoutVerification = async () => {
    await this.props.handleSubmit()
  }

  handleSubmit = e => {
    e.preventDefault()

    const state = copy(this.state)
    state.submitting = true
    state.formErrorMessage = false
    this.setState(state)

    return this.props.form.validateFields(async (err, values) => {
      if (!err) {
        try {
          const data = this.props.data
          const canonicalName = (data && data.metadata && data.metadata.name) || canonical(values.name)
          const gkeCredsResource = GKECredentials(canonicalName, values)
          const gkeResult = await apiRequest(null, 'put', `${apiPaths.team(this.props.team).gkeCredentials}/${canonicalName}`, gkeCredsResource)
          const [ group, version ] = gkeCredsResource.apiVersion.split('/')
          const allocationSpec = {
            name: values.name,
            summary: values.summary,
            resource: {
              group,
              version,
              kind: gkeCredsResource.kind,
              namespace: this.props.team,
              name: canonicalName
            },
            teams: this.state.allocations.length > 0 ? this.state.allocations : ['*']
          }
          const allocationResult = await apiRequest(null, 'put', `${apiPaths.team(this.props.team).allocations}/${canonicalName}`, Allocation(canonicalName, allocationSpec))
          gkeResult.allocation = allocationResult

          if (this.props.inlineVerification) {
            await this.verify(gkeResult)
          } else {
            await this.props.handleSubmit(gkeResult)
          }

        } catch (err) {
          console.error('Error submitting form', err)
          const state = copy(this.state)
          state.submitting = false
          state.formErrorMessage = 'An error occurred saving the project, please try again'
          this.setState(state)
        }
      }
    })
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
            message="This project is not allocated to any teams"
            description="Give the project a Name and Description below and enter Allocated team(s) as appropriate. Once complete, click Save to allocate this project."
            type="warning"
            showIcon
            style={{ marginBottom: '20px', marginTop: '5px' }}
          />
        ) : null}

        <Form {...formConfig} onSubmit={this.handleSubmit}>
          <FormErrorMessage />
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
              <Form.Item label="Allocate team(s)" extra="If nothing selected then this integration will be available to ALL teams">
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

const WrappedGKECredentialsForm = Form.create({ name: 'gke_credentials' })(GKECredentialsForm)

export default WrappedGKECredentialsForm
