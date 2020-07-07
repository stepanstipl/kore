import React from 'react'
import PropTypes from 'prop-types'
import { Alert, Card, Button, Form, Input, Typography } from 'antd'
const { Text } = Typography

import FormErrorMessage from '../../forms/FormErrorMessage'
import IconTooltip from '../../utils/IconTooltip'
import { patterns } from '../../../utils/validation'

class AutomatedCloudAccountForm extends React.Component {
  static propTypes = {
    form: PropTypes.object.isRequired,
    data: PropTypes.object,
    alertTitle: PropTypes.string.isRequired,
    alertDescription: PropTypes.string.isRequired,
    handleSubmit: PropTypes.func.isRequired,
    handleCancel: PropTypes.func.isRequired,
  }

  state = {
    submitting: false,
    formErrorMessage: false
  }

  componentDidMount() {
    // To disabled submit button at the beginning.
    this.props.form.validateFields()
  }

  disableButton = fieldsError => this.state.submitting ? true : Object.keys(fieldsError).some(field => fieldsError[field])

  handleCancel = () => {
    this.props.form.resetFields()
    this.props.handleCancel()
  }

  handleSubmit = (e) => {
    e.preventDefault()
    this.setState({ submitting: true })
    this.props.form.validateFields((err, values) => {
      if (err) {
        this.setState({ submitting: false, formErrorMessage: 'Validation failed' })
        return
      }
      this.setState({ submitting: false, formErrorMessage: false })
      this.props.form.resetFields()
      this.props.handleSubmit(values)
    })
  }

  fieldError = fieldKey => this.props.form.isFieldTouched(fieldKey) && this.props.form.getFieldError(fieldKey)

  render() {
    const data = this.props.data
    const { getFieldDecorator, getFieldsError } = this.props.form
    const formConfig = {
      layout: 'inline',
      labelAlign: 'left',
      hideRequiredMark: true,
      labelCol: { span: 6 },
      wrapperCol: { span: 16 }
    }

    return (
      <>
        <FormErrorMessage message={this.state.formErrorMessage} />
        <Form { ...formConfig } onSubmit={this.handleSubmit}>
          <Form.Item className="inline-form-item-block" label="Name" validateStatus={this.fieldError('name') ? 'error' : ''}  help={this.fieldError('name') || ''}>
            {getFieldDecorator('name', { rules: [{ required: true, message: 'Please enter the name!' }], initialValue: data && data.name })(
              <Input placeholder="Name" />
            )}
          </Form.Item>
          <Form.Item className="inline-form-item-block" label="Description" validateStatus={this.fieldError('description') ? 'error' : ''} help={this.fieldError('description') || ''}>
            {getFieldDecorator('description', { rules: [{ required: true, message: 'Please enter the description!' }], initialValue: data && data.description })(
              <Input placeholder="Description" />
            )}
          </Form.Item>
          <Card style={{ marginTop: '10px', marginBottom: '20px' }}>
            <Alert
              style={{ marginBottom: '10px' }}
              showIcon={true}
              message={this.props.alertTitle}
              description={this.props.alertDescription}
            />
            <Form.Item style={{ marginRight: '-40px' }} labelCol={{ span: 16 }} label="Prefix" validateStatus={this.fieldError('prefix') ? 'error' : ''} help={this.fieldError('prefix') ? <IconTooltip icon="close-circle" color="red" text={this.fieldError('prefix')} /> : ''}>
              {getFieldDecorator('prefix', { rules: [{ ...patterns.uriCompatible10CharMax }], initialValue: (data && data.prefix) || 'kore' })(
                <Input placeholder="Prefix" maxLength={10} />
              )}
            </Form.Item>
            <Form.Item style={{ marginRight: '30px', paddingTop: '10px' }} label="&nbsp;" labelCol={{ span: 24 }} colon={false}>
              <div style={{ width: '90px' }}>
                <Text style={{ fontStyle: 'italic' }}>-team-name-</Text>
              </div>
            </Form.Item>
            <Form.Item labelCol={{ span: 16 }}  label="Suffix" validateStatus={this.fieldError('suffix') ? 'error' : ''} help={this.fieldError('suffix') ? <IconTooltip icon="close-circle" color="red" text={this.fieldError('suffix')} /> : ''}>
              {getFieldDecorator('suffix', { rules: [{ ...patterns.uriCompatible10CharMax }], initialValue: data && data.suffix })(
                <Input placeholder="Suffix" maxLength={10}/>
              )}
            </Form.Item>
          </Card>

          <Form.Item className="inline-form-item-block">
            <Button type="primary" htmlType="submit" loading={this.state.submitting} disabled={this.disableButton(getFieldsError())}>Save</Button>
            <Button type="link" onClick={this.handleCancel}>Cancel</Button>
          </Form.Item>
        </Form>
      </>
    )
  }
}

const WrappedAutomatedCloudAccountForm = Form.create({ name: 'automated_cloud_account_form' })(AutomatedCloudAccountForm)

export default WrappedAutomatedCloudAccountForm
