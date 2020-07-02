import * as React from 'react'
import PropTypes from 'prop-types'
import { Button, Form, Input } from 'antd'

import FormErrorMessage from '../../forms/FormErrorMessage'
import KoreApi from '../../../kore-api'
import { patterns } from '../../../utils/validation'
import { loadingMessage } from '../../../utils/message'

class NamespaceClaimForm extends React.Component {
  static propTypes = {
    form: PropTypes.any.isRequired,
    team: PropTypes.string.isRequired,
    cluster: PropTypes.object.isRequired,
    handleSubmit: PropTypes.func.isRequired,
    handleCancel: PropTypes.func.isRequired
  }

  state = {
    submitting: false,
    formErrorMessage: false
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

  cancel = () => {
    this.props.form.resetFields()
    this.props.handleCancel()
  }

  getMetadataName = (values) => `${this.props.cluster.metadata.name}-${values.name}`

  handleSubmit = e => {
    e.preventDefault()

    this.setState({ submitting: true, formErrorMessage: false })

    return this.props.form.validateFields(async (err, values) => {
      if (err) {
        return this.setState({ submitting: false, formErrorMessage: 'Validation failed' })
      }

      try {
        const resourceName = this.getMetadataName(values)
        const existingNc = await (await KoreApi.client()).GetNamespace(this.props.team, resourceName)
        if (existingNc) {
          return this.setState({
            submitting: false,
            formErrorMessage:`A namespace with the name "${values.name}" already exists on cluster "${this.props.cluster.metadata.name}"`
          })
        }
        const nsClaimResource = KoreApi.resources().team(this.props.team).generateNamespaceClaimResource(this.props.cluster, resourceName, values)
        const nsClaimResult = await (await KoreApi.client()).UpdateNamespace(this.props.team, resourceName, nsClaimResource)
        this.props.form.resetFields()
        this.setState({ submitting: false })
        loadingMessage(`Namespace "${values.name}" requested`)
        await this.props.handleSubmit(nsClaimResult)
      } catch (err) {
        console.error('Error submitting form', err)
        this.setState({ submitting: false, formErrorMessage: 'An error occurred creating the namespace, please try again' })
      }
    })
  }

  render() {
    const { getFieldDecorator, getFieldsError, getFieldError, isFieldTouched } = this.props.form
    const { submitting, formErrorMessage } = this.state
    const formConfig = {
      layout: 'horizontal',
      labelAlign: 'left',
      hideRequiredMark: true,
      labelCol: {
        sm: { span: 24 },
        md: { span: 6 },
        lg: { span: 4 }
      },
      wrapperCol: {
        span: 12
      }
    }

    // Only show error after a field is touched.
    const nameError = isFieldTouched('name') && getFieldError('name')

    return (
      <Form {...formConfig} onSubmit={this.handleSubmit} style={{ marginBottom: '30px' }}>
        <FormErrorMessage message={formErrorMessage} />
        <Form.Item label="Name" validateStatus={nameError ? 'error' : ''} help={nameError || ''}>
          {getFieldDecorator('name', {
            rules: [
              { required: true, message: 'Please enter the namespace name!' },
              { ...patterns.uriCompatible63CharMax }
            ]
          })(
            <Input placeholder="Name" />,
          )}
        </Form.Item>
        <Form.Item>
          <Button id="save" type="primary" htmlType="submit" loading={submitting} disabled={this.disableButton(getFieldsError())}>Save</Button>
          <Button type="link" onClick={this.cancel}>Cancel</Button>
        </Form.Item>
      </Form>
    )
  }
}

const WrappedNamespaceClaimForm = Form.create({ name: 'namespace_claim' })(NamespaceClaimForm)

export default WrappedNamespaceClaimForm
