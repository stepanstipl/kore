import * as React from 'react'
import { Card, Alert, Form, Input, Button } from 'antd'
import PropTypes from 'prop-types'
import { set } from 'lodash'

import KoreApi from '../../kore-api'
import V1ServicePlan from '../../kore-api/model/V1ServicePlan'
import V1ServicePlanSpec from '../../kore-api/model/V1ServicePlanSpec'
import V1ObjectMeta from '../../kore-api/model/V1ObjectMeta'
import ServicePlanOption from './ServicePlanOption'
import FormErrorMessage from '../forms/FormErrorMessage'
import canonical from '../../utils/canonical'

class ServicePlanForm extends React.Component {
  static propTypes = {
    form: PropTypes.any.isRequired,
    data: PropTypes.object,
    kind: PropTypes.string.isRequired,
    validationErrors: PropTypes.array,
    handleSubmit: PropTypes.func.isRequired,
    handleValidationErrors: PropTypes.func.isRequired
  }

  state = {
    dataLoading: true,
    submitting: false,
    schema: null,
    servicePlanValues: this.props.data && this.props.data.spec.configuration || {},
    formErrorMessage: false
  }

  componentDidMountComplete = null
  componentDidMount() {
    this.componentDidMountComplete = this.fetchComponentData().then(() => {
      // To disabled submit button at the beginning.
      this.props.form.validateFields()
    })
  }

  async fetchComponentData() {
    const schema = await (await KoreApi.client()).GetServicePlanSchema(this.props.kind)
    this.setState({
      schema,
      dataLoading: false
    })
  }

  onValueChange(name, value) {
    // Texture this back into a state update using the nifty lodash set function:
    const newServicePlanValues = set({ ...this.state.servicePlanValues }, name, value)
    this.setState({
      servicePlanValues: newServicePlanValues
    })
  }

  disableButton = fieldsError => {
    if (this.state.submitting) {
      return true
    }
    return Object.keys(fieldsError).some(field => fieldsError[field])
  }

  setFormSubmitting = (submitting = true, formErrorMessage = false) => {
    this.setState({
      submitting,
      formErrorMessage
    })
  }

  getMetadataName = values => {
    const data = this.props.data
    return (data && data.metadata && data.metadata.name) || canonical(values.description)
  }

  generateServicePlanResource = values => {
    const metadataName = this.getMetadataName(values)

    const servicePlanResource = new V1ServicePlan()
    servicePlanResource.setApiVersion('services.kore.appvia.io/v1')
    servicePlanResource.setKind('ServicePlan')

    const meta = new V1ObjectMeta()
    meta.setName(metadataName)
    servicePlanResource.setMetadata(meta)

    const spec = new V1ServicePlanSpec()
    spec.setKind(this.props.kind)
    spec.setDescription(values.description)
    spec.setSummary(values.summary)
    spec.setConfiguration(values.configuration)
    servicePlanResource.setSpec(spec)

    return servicePlanResource
  }

  generateServicePlanConfiguration = () => {
    const properties = this.state.schema.properties
    const defaultConfiguration = {}
    Object.keys(properties).forEach(p => properties[p].type === 'boolean' ? defaultConfiguration[p] = false : null)
    return { ...defaultConfiguration, ...this.state.servicePlanValues }
  }

  _process = async (err, values) => {
    if (err) {
      this.setFormSubmitting(false, 'Validation failed')
      return
    }
    try {
      const api = await KoreApi.client()
      const metadataName = this.getMetadataName(values)
      const servicePlanResult = await api.UpdateServicePlan(metadataName, this.generateServicePlanResource({ ...values, configuration: this.generateServicePlanConfiguration() }))
      return await this.props.handleSubmit(servicePlanResult)
    } catch (err) {
      console.error('Error submitting form', err)

      const message = (err.fieldErrors && err.message) ? err.message : 'An error occurred saving the service plan, please try again'
      await this.props.handleValidationErrors(err.fieldErrors)
      this.setFormSubmitting(false, message)
    }
  }

  handleSubmit = e => {
    e.preventDefault()
    this.setFormSubmitting()
    this.props.form.validateFields(this._process)
  }

  // Only show error after a field is touched.
  fieldError = fieldKey => this.props.form.isFieldTouched(fieldKey) && this.props.form.getFieldError(fieldKey)

  render() {
    const { form, data, validationErrors } = this.props
    const { dataLoading, submitting, schema, servicePlanValues, formErrorMessage } = this.state
    const { getFieldDecorator, getFieldsError } = form

    const formConfig = {
      layout: 'horizontal',
      labelAlign: 'left',
      labelCol: {
        sm: { span: 24 },
        md: { span: 8 }
      },
      wrapperCol: {
        sm: { span: 24 },
        md: { span: 16 }
      }
    }

    if (dataLoading) {
      return (
        <div>Loading...</div>
      )
    }

    return (
      <Form {...formConfig} onSubmit={this.handleSubmit}>

        <FormErrorMessage message={formErrorMessage} />

        <Card style={{ marginBottom: '20px' }}>
          <Alert
            message="Help Kore teams understand this service plan"
            description="Give this service plan a summary and description to help teams choose the correct one."
            type="info"
            style={{ marginBottom: '20px' }}
          />
          <Form.Item label="Summary" validateStatus={this.fieldError('summary') ? 'error' : ''} help={this.fieldError('summary') || 'Summary of the service plan'}>
            {getFieldDecorator('summary', {
              rules: [{ required: true, message: 'Please enter the name!' }],
              initialValue: data && data.spec.summary
            })(
              <Input placeholder="Summary" />,
            )}
          </Form.Item>
          <Form.Item label="Description" validateStatus={this.fieldError('description') ? 'error' : ''} help={this.fieldError('description') || 'Description of the service plan'}>
            {getFieldDecorator('description', {
              rules: [{ required: true, message: 'Please enter the description!' }],
              initialValue: data && data.spec.description
            })(
              <Input placeholder="Description" />,
            )}
          </Form.Item>
        </Card>

        <Card style={{ marginBottom: '20px' }}>
          <Alert
            message="Service plan configuration"
            description="Set the service plan configuration"
            type="info"
            style={{ marginBottom: '20px' }}
          />

          {Object.keys(schema.properties).map(property =>
            <ServicePlanOption
              key={property}
              name={property}
              property={schema.properties[property]}
              value={servicePlanValues[property]}
              editable={true}
              hideNonEditable={false}
              onChange={(n, v) => this.onValueChange(n, v)}
              validationErrors={validationErrors}
            />
          )}

        </Card>

        <Form.Item>
          <Button type="primary" htmlType="submit" loading={submitting} disabled={this.disableButton(getFieldsError())}>Save</Button>
        </Form.Item>
      </Form>
    )
  }
}

const WrappedServicePlanForm = Form.create({ name: 'servicePlan' })(ServicePlanForm)

export default WrappedServicePlanForm

