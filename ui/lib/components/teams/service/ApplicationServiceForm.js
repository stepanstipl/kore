import * as React from 'react'
import PropTypes from 'prop-types'
import { Alert, Button, Card, Checkbox, Form, Icon, Input, message, Select } from 'antd'

import ServiceOptionsForm from '../../services/ServiceOptionsForm'
import FormErrorMessage from '../../forms/FormErrorMessage'
import KoreApi from '../../../kore-api'
import V1ServiceSpec from '../../../kore-api/model/V1ServiceSpec'
import V1Service from '../../../kore-api/model/V1Service'
import { NewV1ObjectMeta, NewV1Ownership } from '../../../utils/model'
import { getKoreLabel } from '../../../utils/crd-helpers'

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
    servicePlanOverride: null,
    validationErrors: null,
    createNewNamespace: false
  }

  state = { ...ApplicationServiceForm.initialState }

  async fetchComponentData() {
    const api = await KoreApi.client()
    const [ serviceKinds, servicePlans, namespaceClaims ] = await Promise.all([
      api.ListServiceKinds(),
      api.ListServicePlans(),
      api.ListNamespaces(this.props.team.metadata.name)
    ])
    serviceKinds.items = serviceKinds.items.filter(sk => getKoreLabel(sk, 'platform') === 'Kubernetes' && sk.spec.enabled)
    namespaceClaims.items = namespaceClaims.items.filter(nc => nc.spec.cluster.name === this.props.cluster.metadata.name)
    return { serviceKinds, servicePlans, namespaceClaims }
  }

  componentDidMountComplete = null
  componentDidMount() {
    // Assign the promise chain to a variable so tests can wait for it to complete.
    this.componentDidMountComplete = Promise.resolve().then(async () => {
      const data = await this.fetchComponentData()
      this.setState({ ...data, dataLoading: false })
    })
  }

  generateServiceResource = (values) => {
    const cluster = this.props.cluster
    const selectedServicePlan = this.state.servicePlans.items.find(p => p.metadata.name === values.servicePlan)

    const serviceResource = new V1Service()
    serviceResource.setApiVersion('services.compute.kore.appvia.io/v1')
    serviceResource.setKind('Service')

    serviceResource.setMetadata(NewV1ObjectMeta(values.serviceName, this.props.team.metadata.name))

    const serviceSpec = new V1ServiceSpec()
    serviceSpec.setKind(selectedServicePlan.spec.kind)
    serviceSpec.setPlan(selectedServicePlan.metadata.name)
    serviceSpec.setClusterNamespace(values.createNamespace || values.namespace)

    serviceSpec.setCluster(NewV1Ownership({
      group: cluster.apiVersion.split('/')[0],
      version: cluster.apiVersion.split('/')[1],
      kind: cluster.kind,
      name: cluster.metadata.name,
      namespace: this.props.team.metadata.name
    }))
    if (this.state.servicePlanOverride) {
      serviceSpec.setConfiguration(this.state.servicePlanOverride)
    } else {
      serviceSpec.setConfiguration({ ...selectedServicePlan.spec.configuration })
    }

    serviceResource.setSpec(serviceSpec)
    return serviceResource
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
        message.error('Validation failed')
        this.setState({ submitting: false, formErrorMessage: 'Validation failed' })
        return
      }
      try {
        const existing = await (await KoreApi.client()).GetService(this.props.team.metadata.name, values.serviceName)
        if (existing) {
          message.error('Validation failed')
          return this.setState({
            submitting: false,
            formErrorMessage: `A service with the name "${values.serviceName}" already exists in this team.`
          })
        }

        const service = await (await KoreApi.client()).UpdateService(this.props.team.metadata.name, values.serviceName, this.generateServiceResource(values))
        message.loading('Application service requested...')

        return this.props.handleSubmit(service)
      } catch (err) {
        console.error('Error saving application service', err)
        message.error('Error requesting application service, please try again.')
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
      servicePlanOverride: null,
      validationErrors: null
    })
  }

  handleServicePlanOverride = servicePlanOverrides => {
    this.setState({ servicePlanOverride: servicePlanOverrides })
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
      filteredServicePlans = this.state.servicePlans.items.filter(p => p.spec.kind === selectedServiceKind)
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
            })(
              <Select
                onChange={this.handleSelectKind}
                placeholder="Choose service type"
                showSearch
                optionFilterProp="children"
                filterOption={(input, option) =>
                  option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
                }
              >
                {serviceKinds.items.map(k => <Select.Option key={k.metadata.name} value={k.metadata.name}>{k.spec.displayName || k.metadata.name}</Select.Option>)}
              </Select>
            )}
          </Form.Item>

          {selectedServiceKind && (
            <ServiceOptionsForm
              team={this.props.team}
              selectedServiceKind={selectedServiceKind}
              servicePlans={filteredServicePlans}
              teamServices={this.props.teamServices}
              onServicePlanSelected={(selectedServicePlan) => this.setState({ selectedServicePlan })}
              onServicePlanOverridden={this.handleServicePlanOverride}
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
                {namespaceClaims.items.map(nc => <Select.Option key={nc.spec.name} value={nc.spec.name}>{nc.spec.name}</Select.Option>)}
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
