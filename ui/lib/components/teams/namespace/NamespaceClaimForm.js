import * as React from 'react'
import PropTypes from 'prop-types'
import { Button, Form, Input } from 'antd'

import FormErrorMessage from '../../forms/FormErrorMessage'
import copy from '../../../utils/object-copy'
import Generic from '../../../crd/Generic'
import apiRequest from '../../../utils/api-request'
import apiPaths from '../../../utils/api-paths'
import { patterns } from '../../../utils/validation'

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

  handleSubmit = e => {
    e.preventDefault()

    this.setState({ submitting: true, formErrorMessage: false })

    return this.props.form.validateFields(async (err, values) => {
      if (err) {
        this.setState({ submitting: false, formErrorMessage: 'Validation failed' })
      }

      try {
        const cluster = this.props.cluster
        const name = values.name
        const resourceName = `${cluster.metadata.name}-${name}`

        const existingNc = await apiRequest(null, 'get', `${apiPaths.team(this.props.team).namespaceClaims}/${resourceName}`)
        if (Object.keys(existingNc).length) {
          const state = copy(this.state)
          state.submitting = false
          state.formErrorMessage = `A namespace with the name "${name}" already exists on cluster "${cluster.metadata.name}"`
          return this.setState(state)
        }

        const [ group, version ] = cluster.apiVersion.split('/')
        const spec = {
          name,
          cluster: {
            group,
            version,
            kind: cluster.kind,
            name: cluster.metadata.name,
            namespace: this.props.team
          }
        }
        const nsClaimResource = Generic({
          apiVersion: 'namespaceclaims.clusters.compute.kore.appvia.io/v1alpha1',
          kind: 'NamespaceClaim',
          name: resourceName,
          spec
        })
        const nsClaimResult = await apiRequest(null, 'put', `${apiPaths.team(this.props.team).namespaceClaims}/${resourceName}`, nsClaimResource)
        this.props.form.resetFields()
        const state = copy(this.state)
        state.submitting = false
        this.setState(state)
        await this.props.handleSubmit(nsClaimResult)
      } catch (err) {
        console.error('Error submitting form', err)
        const state = copy(this.state)
        state.submitting = false
        state.formErrorMessage = 'An error occurred creating the namespace, please try again'
        this.setState(state)
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
          <Button type="primary" htmlType="submit" loading={submitting} disabled={this.disableButton(getFieldsError())}>Save</Button>
          <Button type="link" onClick={this.cancel}>Cancel</Button>
        </Form.Item>
      </Form>
    )
  }
}

const WrappedNamespaceClaimForm = Form.create({ name: 'namespace_claim' })(NamespaceClaimForm)

export default WrappedNamespaceClaimForm
