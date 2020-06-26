import * as React from 'react'
import PropTypes from 'prop-types'
import { Alert, Button, Card, Checkbox, Form, Icon, Input, Select } from 'antd'

import ServiceOptionsForm from '../../services/ServiceOptionsForm'
import FormErrorMessage from '../../forms/FormErrorMessage'
import KoreApi from '../../../kore-api'
import { getKoreLabel } from '../../../utils/crd-helpers'
import { errorMessage, loadingMessage } from '../../../utils/message'

class ApplicationServiceForm extends React.Component {
  static propTypes = {
    form: PropTypes.any.isRequired,
    team: PropTypes.object.isRequired,
    cluster: PropTypes.object.isRequired,
    teamServices: PropTypes.array.isRequired,
    handleSubmit: PropTypes.func.isRequired,
    handleCancel: PropTypes.func.isRequired
  }

  static initialState = {
    submitButtonText: 'Save',
    submitting: false,
    formErrorMessage: false,
    selectedServiceKind: false,
    selectedServicePlan: false,
    dataLoading: true,
    planValues: null,
    validationErrors: null,
    createNewNamespace: false
  }

  state = { ...ApplicationServiceForm.initialState }

  async fetchComponentData() {
    const api = await KoreApi.client()
    let [ serviceKinds, servicePlans, namespaceClaims ] = await Promise.all([
      api.ListServiceKinds(),
      api.ListServicePlans(),
      api.ListNamespaces(this.props.team.metadata.name)
    ])
    serviceKinds = serviceKinds.items.filter(sk => getKoreLabel(sk, 'platform') === 'Kubernetes' && sk.spec.enabled)
    servicePlans = servicePlans.items
    namespaceClaims = namespaceClaims.items.filter(nc => nc.spec.cluster.name === this.props.cluster.metadata.name)
    return { serviceKinds, servicePlans, namespaceClaims }
  }

  componentDidMountComplete = null
  componentDidMount() {
    // Assign the promise chain to a variable so tests can wait for it to complete.
    this.componentDidMountComplete = Promise.resolve().then(async () => {
      const data = await this.fetchComponentData()
      this.setState({ ...data, dataLoading: false })
      if (data.serviceKinds && data.serviceKinds.length === 1) {
        this.handleSelectKind(data.serviceKinds[0].metadata.name)
      }
    })
  }

  getServiceResource = (values) => {
    const team = this.props.team.metadata.name
    const cluster = this.props.cluster
    const selectedServicePlan = this.state.servicePlans.find(p => p.metadata.name === values.servicePlan)

    return KoreApi.resources()
      .team(team)
      .generateServiceResource(cluster, values, selectedServicePlan, this.state.planValues)
  }

  validatedFormsFields = (callback) => {
    this.props.form.validateFields((serviceErr, serviceValues) => {
      this.serviceOptionsForm.props.form.validateFields((optionsErr, optionsValues) => {
        const err = serviceErr || optionsErr ? { ...serviceErr, ...optionsErr } : null
        callback(err, { ...serviceValues, ...optionsValues })
      })
    })
  }

  handleSubmit = async (e) => {
    e.preventDefault()

    this.setState({ submitting: true, formErrorMessage: false })

    this.validatedFormsFields(async (err, values) => {
      if (err) {
        errorMessage('Validation failed')
        this.setState({ submitting: false, formErrorMessage: 'Validation failed' })
        return
      }
      try {
        const existing = await (await KoreApi.client()).GetService(this.props.team.metadata.name, values.serviceName)
        if (existing) {
          errorMessage('Validation failed')
          return this.setState({
            submitting: false,
            formErrorMessage: `A service with the name "${values.serviceName}" already exists in this team.`
          })
        }

        const service = await (await KoreApi.client()).UpdateService(this.props.team.metadata.name, values.serviceName, this.getServiceResource(values))
        loadingMessage('Application service requested...')

        return this.props.handleSubmit(service)
      } catch (err) {
        console.error('Error saving application service', err)
        errorMessage('Error requesting application service, please try again.')
        this.setState({
          submitting: false,
          formErrorMessage: (err.fieldErrors && err.message) ? err.message : 'An error occurred requesting the application service, please try again',
          validationErrors: err.fieldErrors // This will be undefined on non-validation errors, which is fine.
        })
      }
    })
  }

  handleSelectKind = (kind) => {
    this.setState({
      selectedServiceKind: kind,
      planValues: null,
      validationErrors: null
    })
  }

  setPlanValues = planValues => {
    this.setState({ planValues })
  }

  disableButton = () => {
    if (!this.state.selectedServicePlan) {
      return true
    }
    return false
  }

  cancel = () => {
    this.props.form.resetFields()
    this.setState({ ...ApplicationServiceForm.initialState })
    this.props.handleCancel()
  }

  createNewNamespace = (checked) => {
    this.setState({ createNewNamespace: checked })
    this.props.form.resetFields('namespace')
  }

  render() {
    if (this.state.dataLoading || !this.props.team) {
      return <Icon type="loading" />
    }
    const formConfig = {
      layout: 'horizontal',
      labelAlign: 'left',
      hideRequiredMark: true,
      labelCol: {
        sm: { span: 24 },
        md: { span: 24 },
        lg: { span: 6 }
      },
      wrapperCol: {
        sm: { span: 24 },
        md: { span: 24 },
        lg: { span: 18 }
      }
    }

    const { getFieldDecorator } = this.props.form
    const { serviceKinds, selectedServiceKind, namespaceClaims, formErrorMessage, submitting, createNewNamespace } = this.state

    let filteredServicePlans = []
    if (selectedServiceKind) {
      filteredServicePlans = this.state.servicePlans.filter(p => p.spec.kind === selectedServiceKind)
    }

    return (
      <Form {...formConfig} onSubmit={this.handleSubmit}>
        <FormErrorMessage message={formErrorMessage} />
        <Card style={{ marginBottom: '20px' }}>
          <Alert
            message="Application service"
            description="Select the service you would like to use."
            type="info"
            showIcon
            style={{ marginBottom: '20px' }}
          />

          <Form.Item label="Service type">
            {getFieldDecorator('serviceKind', {
              rules: [{ required: true, message: 'Please select your service type!' }],
              initialValue: selectedServiceKind
            })(
              <Select
                onChange={this.handleSelectKind}
                placeholder="Choose service type"
                optionFilterProp="children"
                filterOption={(input, option) =>
                  option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
                }
              >
                {serviceKinds.map(k => <Select.Option key={k.metadata.name} value={k.metadata.name}>{k.spec.displayName || k.metadata.name}</Select.Option>)}
              </Select>
            )}
          </Form.Item>

          {selectedServiceKind && (
            <ServiceOptionsForm
              team={this.props.team}
              cluster={this.props.cluster}
              selectedServiceKind={selectedServiceKind}
              servicePlans={filteredServicePlans}
              teamServices={this.props.teamServices}
              onServicePlanSelected={(selectedServicePlan) => this.setState({ selectedServicePlan })}
              onPlanValuesChange={this.setPlanValues}
              validationErrors={this.state.validationErrors}
              wrappedComponentRef={inst => this.serviceOptionsForm = inst}
            />
          )}
        </Card>

        <Card>
          <Alert
            message="Target namespace"
            description="The namespace you would like service to be deploying into. Select from the existing cluster namespaces, or create a new one."
            type="info"
            showIcon
            style={{ marginBottom: '20px' }}
          />
          <Form.Item label="Namespace">
            {getFieldDecorator('namespace', {
              rules: [{ required: !createNewNamespace, message: 'Please select the target namespace!' }]
            })(
              <Select placeholder="Choose existing namespace" disabled={createNewNamespace}>
                {namespaceClaims.map(nc => <Select.Option key={nc.spec.name} value={nc.spec.name}>{nc.spec.name}</Select.Option>)}
              </Select>
            )}
            <Checkbox onChange={(e) => this.createNewNamespace(e.target.checked)}>Create a new namespace</Checkbox>
            {createNewNamespace && (
              <Form.Item style={{ marginBottom: 0 }}>
                {getFieldDecorator('createNamespace', {
                  rules: [{ required: true, message: 'Please enter the new namespace!' }],
                })(
                  <Input placeholder="New namespace name" />
                )}
              </Form.Item>
            )}
          </Form.Item>
        </Card>

        <Form.Item style={{ marginTop: '20px', marginBottom: 0 }}>
          <Button type="primary" htmlType="submit" loading={submitting} disabled={this.disableButton()}>{this.state.submitButtonText}</Button>
          <Button type="link" onClick={this.cancel}>Cancel</Button>
        </Form.Item>
      </Form>
    )
  }
}

const WrappedApplicationServiceForm = Form.create({ name: 'cluster_application' })(ApplicationServiceForm)

export default WrappedApplicationServiceForm
