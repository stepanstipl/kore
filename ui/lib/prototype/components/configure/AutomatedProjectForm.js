import React from 'react'
import PropTypes from 'prop-types'
import { Button, Form, Input } from 'antd'

import FormErrorMessage from '../../../../lib/components/forms/FormErrorMessage'

class AutomatedProjectForm extends React.Component {
  static propTypes = {
    form: PropTypes.object.isRequired,
    handleSubmit: PropTypes.func.isRequired,
    handleCancel: PropTypes.func.isRequired,
  }

  state = {
    submitting: false,
    formErrorMessage: false
  }

  handleSubmit = (e) => {
    e.preventDefault()
    this.setState({ submitting: true })
    this.props.form.validateFields((err, values) => {
      if (err) {
        this.setState({ submitting: false, formErrorMessage: 'Validation failed' })
        return
      }
      this.props.handleSubmit(values)
    })
  }

  fieldError = fieldKey => this.props.form.isFieldTouched(fieldKey) && this.props.form.getFieldError(fieldKey)

  render() {
    const { getFieldDecorator } = this.props.form
    const formConfig = {
      layout: 'horizontal',
      labelAlign: 'left',
      hideRequiredMark: true,
      labelCol: { span: 8 },
      wrapperCol: { span: 16 }
    }
    return (
      <>
        <FormErrorMessage message={this.state.formErrorMessage} />
        <Form { ...formConfig } onSubmit={this.handleSubmit}>
          <Form.Item label="Title" help={this.fieldError('title') || ''}>
            {getFieldDecorator('title', { rules: [{ required: true, message: 'Please enter the title!' }] })(
              <Input placeholder="Title" />
            )}
          </Form.Item>
          <Form.Item label="Description" help={this.fieldError('description') || ''}>
            {getFieldDecorator('description', { rules: [{ required: true, message: 'Please enter the description!' }] })(
              <Input placeholder="Description" />
            )}
          </Form.Item>
          <Form.Item label="Prefix" help={this.fieldError('prefix') || 'The project ID prefix'}>
            {getFieldDecorator('prefix', { initialValue: 'kore' })(
              <Input placeholder="Prefix" />
            )}
          </Form.Item>
          <Form.Item label="Suffix" help={this.fieldError('suffix') || 'The project ID suffix'}>
            {getFieldDecorator('suffix', { rules: [{ required: true, message: 'Please enter the suffix!' }] })(
              <Input placeholder="Suffix" />
            )}
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={this.state.submitting}>Save</Button>
            <Button type="link" onClick={this.props.handleCancel}>Cancel</Button>
          </Form.Item>
        </Form>
      </>
    )
  }
}

const WrappedAutomatedProjectForm = Form.create({ name: 'project_form' })(AutomatedProjectForm)

export default WrappedAutomatedProjectForm
