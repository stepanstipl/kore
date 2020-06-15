import * as React from 'react'
import PropTypes from 'prop-types'
import { Button, Form, Input } from 'antd'

import FormErrorMessage from '../../forms/FormErrorMessage'
import KoreApi from '../../../kore-api'
import V1NamespaceClaim from '../../../kore-api/model/V1NamespaceClaim'
import V1NamespaceClaimSpec from '../../../kore-api/model/V1NamespaceClaimSpec'
import { NewV1ObjectMeta, NewV1Ownership } from '../../../utils/model'
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

  generateNamespaceClaimResource = (values) => {
    const cluster = this.props.cluster
    const resource = new V1NamespaceClaim()
    resource.setApiVersion('namespaceclaims.clusters.compute.kore.appvia.io/v1alpha1')
    resource.setKind('NamespaceClaim')

    resource.setMetadata(NewV1ObjectMeta(this.getMetadataName(values), this.props.team))

    const spec = new V1NamespaceClaimSpec()
    spec.setName(values.name)

    const [ group, version ] = cluster.apiVersion.split('/')
    spec.setCluster(NewV1Ownership({
      group,
      version,
      kind: cluster.kind,
      name: cluster.metadata.name,
      namespace: this.props.team
    }))
    resource.setSpec(spec)

    return resource
  }

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
        const nsClaimResult = await (await KoreApi.client()).UpdateNamespace(this.props.team, resourceName, this.generateNamespaceClaimResource(values))
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
          <Button type="primary" htmlType="submit" loading={submitting} disabled={this.disableButton(getFieldsError())}>Save</Button>
          <Button type="link" onClick={this.cancel}>Cancel</Button>
        </Form.Item>
      </Form>
    )
  }
}

const WrappedNamespaceClaimForm = Form.create({ name: 'namespace_claim' })(NamespaceClaimForm)

export default WrappedNamespaceClaimForm
