import * as React from 'react'
import { Card, Alert, Form, Input, Button } from 'antd'
import PropTypes from 'prop-types'
import { set } from 'lodash'

import KoreApi from '../../kore-api'
import V1Plan from '../../kore-api/model/V1Plan'
import V1PlanSpec from '../../kore-api/model/V1PlanSpec'
import V1ObjectMeta from '../../kore-api/model/V1ObjectMeta'
import PlanOption from './PlanOption'
import FormErrorMessage from '../forms/FormErrorMessage'
import canonical from '../../utils/canonical'

class PlanForm extends React.Component {
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
    planValues: this.props.data && this.props.data.spec.configuration || {},
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
    const schema = await (await KoreApi.client()).GetPlanSchema(this.props.kind)
    this.setState({
      schema,
      dataLoading: false
    })
  }

  onValueChange(name, value) {
    // Texture this back into a state update using the nifty lodash set function:
    const newPlanValues = set({ ...this.state.planValues }, name, value)
    this.setState({
      planValues: newPlanValues
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

  generatePlanResource = values => {
    const metadataName = this.getMetadataName(values)

    const planResource = new V1Plan()
    planResource.setApiVersion('config.kore.appvia.io/v1')
    planResource.setKind('Plan')

    const meta = new V1ObjectMeta()
    meta.setName(metadataName)
    planResource.setMetadata(meta)

    const spec = new V1PlanSpec()
    spec.setKind(this.props.kind)
    spec.setDescription(values.description)
    spec.setSummary(values.summary)
    spec.setConfiguration(values.configuration)
    planResource.setSpec(spec)

    return planResource
  }

  generatePlanConfiguration = () => {
    const properties = this.state.schema.properties
    const defaultConfiguration = {}
    Object.keys(properties).forEach(p => properties[p].type === 'boolean' ? defaultConfiguration[p] = false : null)
    return { ...defaultConfiguration, ...this.state.planValues }
  }

  _process = async (err, values) => {
    if (err) {
      this.setFormSubmitting(false, 'Validation failed')
      return
    }
    try {
      const api = await KoreApi.client()
      const metadataName = this.getMetadataName(values)
      const planResult = await api.UpdatePlan(metadataName, this.generatePlanResource({ ...values, configuration: this.generatePlanConfiguration() }))
      return await this.props.handleSubmit(planResult)
    } catch (err) {
      console.error('Error submitting form', err)

      const message = (err.fieldErrors && err.message) ? err.message : 'An error occurred saving the plan, please try again'
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
    const { form, kind, data, validationErrors } = this.props
    const { dataLoading, submitting, schema, planValues, formErrorMessage } = this.state
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
            message="Help Kore teams understand this plan"
            description="Give this plan a summary and description to help teams choose the correct one."
            type="info"
            style={{ marginBottom: '20px' }}
          />
          <Form.Item label="Summary" validateStatus={this.fieldError('summary') ? 'error' : ''} help={this.fieldError('summary') || 'Summary of the plan'}>
            {getFieldDecorator('summary', {
              rules: [{ required: true, message: 'Please enter the name!' }],
              initialValue: data && data.spec.summary
            })(
              <Input placeholder="Summary" />,
            )}
          </Form.Item>
          <Form.Item label="Description" validateStatus={this.fieldError('description') ? 'error' : ''} help={this.fieldError('description') || 'Description of the plan'}>
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
            message="Plan configuration"
            description="Set the plan configuration"
            type="info"
            style={{ marginBottom: '20px' }}
          />

          {Object.keys(schema.properties).map(property =>
            <PlanOption
              mode="manage"
              resourceType="cluster"
              kind={kind}
              plan={planValues}
              key={property}
              name={property}
              property={schema.properties[property]}
              value={planValues[property]}
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

const WrappedPlanForm = Form.create({ name: 'plan' })(PlanForm)

export default WrappedPlanForm

