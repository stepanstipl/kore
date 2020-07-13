import * as React from 'react'
import { Card, Alert, Form, Input, Button, Icon } from 'antd'
import PropTypes from 'prop-types'
import { set } from 'lodash'

import FormErrorMessage from '../forms/FormErrorMessage'
import canonical from '../../utils/canonical'
import PlanViewEdit from './PlanViewEdit'
import CostEstimate from '../costs/CostEstimate'

/**
 * ManagePlanForm is an abstract base for *managing* a plan itself (i.e. viewing, editing or creating PLANS).
 *
 * To *use* a plan to create/manage a resource, use UsePlanForm.
 */
export default class ManagePlanForm extends React.Component {
  static propTypes = {
    form: PropTypes.any.isRequired,
    data: PropTypes.object,
    kind: PropTypes.string.isRequired,
    handleSubmit: PropTypes.func.isRequired,
    mode: PropTypes.oneOf(['create', 'edit', 'view']).isRequired,
    resourceType: PropTypes.oneOf(['service','cluster']).isRequired
  }

  state = {
    dataLoading: true,
    submitting: false,
    schema: null,
    planValues: this.props.data && this.props.data.spec.configuration || {},
    formErrorMessage: false,
    validationErrors: [],
  }

  componentDidMountComplete = null
  componentDidMount() {
    this.componentDidMountComplete = this.fetchComponentData().then(() => {
      if (this.props.mode === 'create' && !this.props.data) {
        // Copy default values into the plan values as a starting point.
        const planValues = {}
        const schemaProps = this.state.schema.properties
        Object.keys(schemaProps).forEach((prop) => {
          if (schemaProps[prop].default !== undefined) {
            planValues[prop] = schemaProps[prop].default
          }
        })
        this.setState({ planValues })
      }
      // To disabled submit button at the beginning.
      this.props.form.validateFields()
    })
  }

  async fetchComponentData() {
    throw 'Must be overridden to populate planValues and schema in state'
  }

  /** 
   * estimateSupported should return true in an overriding class where
   * cost estimates are supported for that provider.
   */
  estimateSupported = () => false

  onValueChange(name, value) {
    this.setState((state) => {
      let planValues = {
        ...state.planValues
      }
      if (value !== undefined) {
        // Texture this back into a state update using the nifty lodash set function:
        planValues = set(planValues, name, value)
      } else {
        // Property set to undefined, so remove it completely from the plan values.
        delete planValues[name]
      }
      return { planValues, planValuesChangedSinceEstimate: true }
    })
  }

  disableButton = (fieldsError) => {
    if (this.state.submitting) {
      return true
    }
    return Object.keys(fieldsError).some(field => fieldsError[field])
  }

  setFormSubmitting = (submitting = true, formErrorMessage = false, validationErrors = []) => {
    this.setState({
      submitting,
      formErrorMessage,
      validationErrors
    })
  }

  getMetadataName = (values) => {
    const data = this.props.data
    return (data && data.metadata && data.metadata.name) || canonical(values.description)
  }

  process = async (err, values) => {
    throw `Must be overridden to save the plan or update the validationErrors in state ${err} ${values}`
  }

  resourceType = () => {
    throw 'Must be overridden to return cluster or service'
  }

  handleSubmit = e => {
    e.preventDefault()
    this.setFormSubmitting()
    this.props.form.validateFields(this.process)
  }

  // Only show error after a field is touched.
  fieldError = (fieldKey) => this.props.form.isFieldTouched(fieldKey) && this.props.form.getFieldError(fieldKey)

  validationErrors = (name) => {
    const { validationErrors } = this.state
    if (!validationErrors) {
      return null
    }
    const valErrors = validationErrors.filter(v => v.field === name)
    if (valErrors.length === 0) {
      return null
    }
    return valErrors.map((ve, i) => <Alert key={`${name}.${i}`} type="error" message={ve.message} style={{ marginTop: '10px' }} />)
  }

  /**
   * Override formHeader to provide a custom header to the form.
   */
  formHeader = (formErrorMessage, mode, data) => {
    return this.defaultFormHeader(formErrorMessage, mode, data)
  }

  defaultFormHeader = (formErrorMessage, mode, data) => {
    const { getFieldDecorator } = this.props.form
    return (
      <>
        <FormErrorMessage message={formErrorMessage} />

        <Card style={{ marginBottom: '20px' }}>
          {mode ==='view' ? null : <Alert
            message="Help Kore teams understand this plan"
            description="Give this plan a summary and description to help teams choose the correct one."
            type="info"
            style={{ marginBottom: '20px' }}
          />}
          <Form.Item label="Summary" validateStatus={this.fieldError('summary') ? 'error' : ''} help={this.fieldError('summary') || 'Summary of the plan'}>
            {getFieldDecorator('summary', {
              rules: [{ required: true, message: 'Please enter the name!' }],
              initialValue: data && data.spec.summary
            })(
              <Input placeholder="Summary" readOnly={mode==='view'} />,
            )}
            {this.validationErrors('summary')}
            {this.validationErrors('name')}
          </Form.Item>
          <Form.Item label="Description" validateStatus={this.fieldError('description') ? 'error' : ''} help={this.fieldError('description') || 'Description of the plan'}>
            {getFieldDecorator('description', {
              rules: [{ required: true, message: 'Please enter the description!' }],
              initialValue: data && data.spec.description
            })(
              <Input placeholder="Description" readOnly={mode==='view'} />,
            )}
            {this.validationErrors('description')}
          </Form.Item>
        </Card>
      </>
    )
  }

  render() {
    const { form, kind, data, mode } = this.props
    const { dataLoading, submitting, schema, planValues, formErrorMessage, validationErrors } = this.state
    const { getFieldsError } = form
    const resourceType = this.resourceType()

    const formConfig = {
      layout: 'horizontal',
      labelAlign: 'left',
      hideRequiredMark: mode === 'view',
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
        <Icon type="loading" />
      )
    }

    return (
      <Form {...formConfig} onSubmit={this.handleSubmit}>

        {this.formHeader(formErrorMessage, mode, data)}

        {!this.estimateSupported() ? null : (
          <Card style={{ marginBottom: '20px' }}>
            <Alert
              message="Cost Estimate"
              description={<>See an <b>approximate</b> cost estimate for usage of this plan</>}
              type="success"
              style={{ marginBottom: '20px' }}
            />
            <CostEstimate 
              planValues={planValues} 
              resourceType={resourceType} 
              kind={kind} 
              noPriceDataError="Pricing information is not available. Please check that the Kore Costs feature has been configured and enabled."
            />
          </Card>
        )}

        <Card style={{ marginBottom: '20px' }}>
          {mode === 'view' ? null : <Alert
            message="Plan configuration"
            description="Set the plan configuration"
            type="info"
            style={{ marginBottom: '20px' }}
          />}
          <PlanViewEdit
            resourceType={resourceType}
            mode={mode}
            manage={true}
            kind={kind}
            plan={planValues}
            schema={schema}
            editableParams={['*']} // everything editable when managing plans
            onPlanValueChange={(n, v) => this.onValueChange(n, v)}
            validationErrors={validationErrors}
          />
        </Card>

        {mode === 'view' ? null :
          <>
            <FormErrorMessage message={formErrorMessage} />
            <Form.Item>
              <Button id="plan_save" type="primary" htmlType="submit" loading={submitting} disabled={this.disableButton(getFieldsError())}>Save</Button>
            </Form.Item>
          </>
        }
      </Form>
    )
  }
}
