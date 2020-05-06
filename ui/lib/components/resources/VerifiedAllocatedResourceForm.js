import * as React from 'react'
import PropTypes from 'prop-types'
import { message, Typography, Form, Card, Alert, Button, Input, Select } from 'antd'
const { Paragraph, Text } = Typography

import canonical from '../../utils/canonical'
import ResourceVerificationStatus from './ResourceVerificationStatus'
import FormErrorMessage from '../forms/FormErrorMessage'
import AllocationHelpers from '../../utils/allocation-helpers'

class VerifiedAllocatedResourceForm extends React.Component {
  static propTypes = {
    form: PropTypes.any.isRequired,
    team: PropTypes.string.isRequired,
    allTeams: PropTypes.object,
    data: PropTypes.object,
    handleSubmit: PropTypes.func.isRequired,
    saveButtonText: PropTypes.string,
    inlineVerification: PropTypes.bool,
    autoAllocateToAllTeams: PropTypes.bool
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
    this.setState({
      ...this.state,
      allocations: value
    })
  }

  getResource = () => {
    throw new Error('getResource must be implemented')
  }

  putResource = () => {
    throw new Error('putResource must be implemented')
  }

  async verify(resource, tryCount) {
    const messageKey = 'verify'
    tryCount = tryCount || 0
    if (tryCount === 0) {
      message.loading({ content: 'Verifying credentials', key: messageKey, duration: 0 })
    }
    if (tryCount === 3) {
      message.error({ content: 'Credentials verification failed', key: messageKey })
      this.setState({
        ...this.state,
        inlineVerificationFailed: true,
        submitting: false,
        formErrorMessage: (
          <>
            <Paragraph>The credentials have been saved but could not be verified, see the error below. Please try again or click &quot;Continue without verification&quot;.</Paragraph>
            {(resource.status.conditions || []).map((c, idx) =>
              <Paragraph key={idx} style={{ marginBottom: '0' }}>
                <Text strong>{c.message}</Text>
                <br/>
                <Text>{c.detail}</Text>
              </Paragraph>
            )}
          </>
        )
      })
    } else {
      setTimeout(async () => {
        const resourceResult = await this.getResource(resource.metadata.name)
        if (resourceResult.status.status === 'Success') {
          message.success({ content: 'Credentials verification successful', key: messageKey })
          return await this.props.handleSubmit(resourceResult)
        }
        return await this.verify(resourceResult, tryCount + 1)
      }, 2000)
    }
  }

  setFormSubmitting = (submitting = true, formErrorMessage = false) => {
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

  storeAllocation = async (resource, values) => {
    return await AllocationHelpers.storeAllocation({
      resourceToAllocate: resource,
      teams: this.state.autoAllocateToAllTeams ? '*' : this.state.allocations,
      name: values.name,
      summary: values.summary
    })
  }

  _process = async (err, values) => {
    if (err) {
      this.setFormSubmitting(false, 'Validation failed')
      return
    }
    try {
      const resourceResult = await this.putResource(values)
      if (this.props.inlineVerification) {
        return await this.verify(resourceResult)
      }
      return await this.props.handleSubmit(resourceResult)
    } catch (err) {
      console.error('Error submitting form', err)
      this.setFormSubmitting(false, 'An error occurred saving the form, please try again')
    }
  }

  handleSubmit = e => {
    e.preventDefault()
    this.setFormSubmitting()
    this.props.form.validateFields(this._process)
  }

  continueWithoutVerification = async () => {
    try {
      const metadataName = this.getMetadataName(this.props.form.getFieldsValue())
      const resourceResult = await this.getResource(metadataName)
      await this.props.handleSubmit(resourceResult)
    } catch (err) {
      console.error('Error getting data', err)
      this.setFormSubmitting(false, 'An error occurred saving the form, please try again')
    }
  }

  // Only show error after a field is touched.
  fieldError = fieldKey => this.props.form.isFieldTouched(fieldKey) && this.props.form.getFieldError(fieldKey)

  resourceFormFields = () => {
    throw new Error('resourceFormFields must be implemented')
  }

  render() {
    const { form, data, allTeams, saveButtonText, autoAllocateToAllTeams } = this.props
    const { getFieldDecorator, getFieldsError } = form
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

    const noAllocation = Boolean(data && !data.allocation)
    const { allocationMissing, nameSection, allocationSection } = this.allocationFormFieldsInfo

    return (
      <>
        <ResourceVerificationStatus resourceStatus={data && data.status} style={{ marginBottom: '15px' }}/>

        {noAllocation ? (
          <Alert
            message={allocationMissing.infoMessage}
            description={allocationMissing.infoDescription}
            type="warning"
            showIcon
            style={{ marginBottom: '20px', marginTop: '5px' }}
          />
        ) : null}

        <Form {...formConfig} onSubmit={this.handleSubmit}>
          <FormErrorMessage message={formErrorMessage} />

          <this.resourceFormFields />

          <Card style={{ marginBottom: '20px' }}>
            <Alert
              message={nameSection.infoMessage}
              description={nameSection.infoDescription}
              type="info"
              style={{ marginBottom: '20px' }}
            />
            <Form.Item label="Name" validateStatus={this.fieldError('name') ? 'error' : ''} help={this.fieldError('name') || nameSection.nameHelp}>
              {getFieldDecorator('name', {
                rules: [{ required: true, message: 'Please enter the name!' }],
                initialValue: data && data.allocation && data.allocation.spec.name
              })(
                <Input placeholder="Name" />,
              )}
            </Form.Item>
            <Form.Item label="Description" validateStatus={this.fieldError('summary') ? 'error' : ''} help={this.fieldError('summary') || nameSection.descriptionHelp}>
              {getFieldDecorator('summary', {
                rules: [{ required: true, message: 'Please enter the description!' }],
                initialValue: data && data.allocation && data.allocation.spec.summary
              })(
                <Input placeholder="Description" />,
              )}
            </Form.Item>
          </Card>

          {!autoAllocateToAllTeams ? (
            <Card style={{ marginBottom: '20px' }}>
              <Alert
                message={allocationSection.infoMessage}
                description={allocationSection.infoDescription}
                type="info"
                style={{ marginBottom: '20px' }}
              />

              {allTeams.items.length === 0 ? (
                <Alert
                  message={allocationSection.allTeamsWarningMessage}
                  description={allocationSection.allTeamsWarningDescription}
                  type="warning"
                  showIcon
                />
              ) : (
                <Form.Item label="Allocate team(s)" extra={allocationSection.allocateExtra}>
                  {getFieldDecorator('allocations', { initialValue: allocations })(
                    <Select
                      mode="multiple"
                      style={{ width: '100%' }}
                      placeholder={noAllocation ? 'No teams' : 'All teams'}
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
          ) : null}

          <Form.Item style={{ marginBottom: '0' }}>
            <Button id="save" type="primary" htmlType="submit" loading={submitting} disabled={this.disableButton(getFieldsError())}>{saveButtonText || 'Save'}</Button>
            {inlineVerificationFailed ? (
              <Button id="continue-without-verification" onClick={this.continueWithoutVerification} disabled={this.disableButton(getFieldsError())} style={{ marginLeft: '10px' }}>Continue without verification</Button>
            ) : null}
          </Form.Item>
        </Form>
      </>
    )
  }
}

export default VerifiedAllocatedResourceForm
