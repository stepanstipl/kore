import * as React from 'react'
import PropTypes from 'prop-types'
import { Alert, Button, Card, Checkbox, Col, Collapse, Form, Icon, Input, Row, Typography, Select } from 'antd'
const { Paragraph, Text } = Typography
const { Panel } = Collapse

import ServiceOptionsForm from '../../services/ServiceOptionsForm'
import FormErrorMessage from '../../forms/FormErrorMessage'
import KoreApi from '../../../kore-api'
import asyncForEach from '../../../utils/async-foreach'
import { getKoreLabel } from '../../../utils/crd-helpers'
import { errorMessage, loadingMessage } from '../../../utils/message'
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

class ServiceBuildForm extends React.Component {
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
    selectedCloud: false,
    selectedServiceKind: false,
    selectedServicePlan: false,
    dataLoading: true,
    planValues: null,
    validationErrors: null,
    bindingsToCreate: [],
    planSchemaFound: false
  }

  constructor(props) {
    super(props)

    const [selectedCloud] = Object.entries(publicRuntimeConfig.clusterProviderMap).find(([ , provider]) => provider === this.props.cluster.spec.kind)

    this.state = { ...ServiceBuildForm.initialState, selectedCloud }
  }

  async fetchComponentData() {
    const api = await KoreApi.client()
    const [ serviceKinds, servicePlans, clusters, namespaceClaims ] = await Promise.all([
      api.ListServiceKinds(),
      api.ListServicePlans(),
      api.ListClusters(this.props.team.metadata.name),
      api.ListNamespaces(this.props.team.metadata.name)
    ])
    const { cluster } = this.props
    const bindingsData = {}
    namespaceClaims.items.filter(nc => nc.spec.cluster.name === cluster.metadata.name).forEach((ns) => {
      if (bindingsData[cluster.metadata.name]) {
        bindingsData[cluster.metadata.name].children.push({ title: ns.spec.name, value: ns.metadata.name })
      } else {
        bindingsData[cluster.metadata.name] = {
          title: cluster.metadata.name,
          value: cluster.metadata.name,
          selectable: false,
          children: [{ title: ns.spec.name, value: ns.metadata.name }]
        }
      }
    })
    const bindingSelectData = Object.keys(bindingsData).map(bd => bindingsData[bd])
    return { serviceKinds, servicePlans, clusters, namespaceClaims, bindingSelectData }
  }

  componentDidMountComplete = null
  componentDidMount() {
    // Assign the promise chain to a variable so tests can wait for it to complete.
    this.componentDidMountComplete = Promise.resolve().then(async () => {
      const data = await this.fetchComponentData()
      this.setState({ ...data, dataLoading: false })
    })
  }

  getServiceResource = (values) => {
    const team = this.props.team.metadata.name
    const cluster = this.props.cluster
    const selectedServicePlan = this.state.servicePlans.items.find(p => p.metadata.name === values.servicePlan)

    return KoreApi.resources()
      .team(team)
      .generateServiceResource(cluster, values, selectedServicePlan, this.state.planValues)
  }

  hasNamespaceBindingErrors = async () => {
    let namespaceBindingErrors = false

    const fieldErrors = this.props.form.getFieldsError()
    namespaceBindingErrors = Object.keys(fieldErrors).some(field => fieldErrors[field])

    await asyncForEach(this.state.bindingsToCreate, async (b) => {
      const secretName = this.props.form.getFieldValue(`${b}-secretName`)
      if (!secretName) {
        namespaceBindingErrors = true
        this.props.form.setFields({
          [`${b}-secretName`]: { errors: [new Error('Please enter the secret name or un-check namespace!')] }
        })
      } else {
        try {
          const existing = await (await KoreApi.client()).GetServiceCredentials(this.props.team.metadata.name, secretName)
          if (existing) {
            namespaceBindingErrors = true
            this.props.form.setFields({
              [`${b}-secretName`]: { value: secretName, errors: [new Error('A secret with this name already exists')] }
            })
          }
        } catch (error) {
          console.error('Error checking for existing service binding', error)
        }
      }
    })

    return namespaceBindingErrors
  }

  handleSubmit = async (e) => {
    e.preventDefault()

    this.setState({ submitting: true, formErrorMessage: false })

    const hasBindingErrors = await this.hasNamespaceBindingErrors()
    if (hasBindingErrors) {
      this.setState({ submitting: false, formErrorMessage: 'Validation failed' })
      return
    }

    this.serviceOptionsForm.props.form.validateFields(async (err, values) => {
      if (err) {
        this.setState({ submitting: false, formErrorMessage: 'Validation failed' })
        return
      }
      try {
        const team = this.props.team.metadata.name
        const api = await KoreApi.client()
        const service = await api.UpdateService(team, values.serviceName, this.getServiceResource(values))
        loadingMessage('Service build requested...')

        if (this.state.bindingsToCreate.length > 0) {
          await asyncForEach(this.state.bindingsToCreate, async (bindingNamespace) => {
            const namespaceClaim = this.state.namespaceClaims.items.find(ns => ns.metadata.name === bindingNamespace)
            try {
              const secretName = this.props.form.getFieldValue(`${bindingNamespace}-secretName`)
              const cluster = this.state.clusters.items.find(c => c.metadata.name === namespaceClaim.spec.cluster.name)
              const credentialName = `${cluster.metadata.name}-${namespaceClaim.spec.name}-${secretName}`
              const serviceCredsResource = KoreApi.resources()
                .team(team)
                .generateServiceCredentialsResource(credentialName, secretName, {}, service, cluster, namespaceClaim.spec.name)
              await api.UpdateServiceCredentials(team, credentialName, serviceCredsResource)
              loadingMessage(`Service access for namespace "${namespaceClaim.spec.name}" requested...`)
            } catch (error) {
              console.error('Error creating service binding', error)
              errorMessage(`Failed to create service access for namespace "${namespaceClaim.spec.name}"`)
            }
          })
        }

        return this.props.handleSubmit(service)
      } catch (err) {
        console.error('Error saving service', err)
        this.setState({
          submitting: false,
          formErrorMessage: (err.fieldErrors && err.message) ? err.message : 'An error occurred requesting the service, please try again',
          validationErrors: err.fieldErrors // This will be undefined on non-validation errors, which is fine.
        })
      }
    })
  }

  handleSelectKind = (kind) => {
    this.setState({
      selectedServiceKind: kind,
      selectedServicePlan: false,
      planValues: null,
      validationErrors: null
    })
  }

  setPlanValues = planValues => {
    this.setState({ planValues })
  }

  handleServicePlanSelected = async (plan) => {
    this.setState({ selectedServicePlan: plan })
    try {
      // check if there is a schema for the binding for the selected service kind/plan
      const schema = (await (await KoreApi.client()).GetServicePlanDetails(plan, this.props.team.metadata.name, this.props.cluster.metadata.name)).schema
      this.setState({ planSchemaFound: Boolean(schema) })
    } catch (err) {
      console.error('Error getting service credentials schema for plan', err)
    }
  }

  onChange = (checked, value) => () => {
    this.props.form.resetFields(`${value}-secretName`)
    if (checked) {
      this.setState({ bindingsToCreate: this.state.bindingsToCreate.concat([value]) })
    } else {
      this.setState({ bindingsToCreate: this.state.bindingsToCreate.filter(v => v !== value) })
    }
  }

  disableButton = () => {
    if (!this.state.selectedCloud || !this.state.selectedServiceKind) {
      return true
    }
    return false
  }

  cancel = () => {
    this.props.form.resetFields()
    this.setState({ ...ServiceBuildForm.initialState })
    this.props.handleCancel()
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
    const { selectedCloud, serviceKinds, selectedServiceKind, selectedServicePlan, formErrorMessage, submitting, bindingSelectData, planSchemaFound } = this.state
    let filteredServiceKinds = []
    if (selectedCloud) {
      filteredServiceKinds = serviceKinds.items.filter(sk => getKoreLabel(sk, 'platform') === selectedCloud && sk.spec.enabled)
    }

    let filteredServicePlans = []
    let selectedServiceKindObject = false
    if (selectedServiceKind) {
      filteredServicePlans = this.state.servicePlans.items.filter(p => p.spec.kind === selectedServiceKind)
      selectedServiceKindObject = serviceKinds.items.find(sk => sk.metadata.name === selectedServiceKind)
    }

    let selectedServicePlanObject = false
    if (selectedServiceKind && selectedServicePlan) {
      selectedServicePlanObject = filteredServicePlans.find(sp => sp.metadata.name === selectedServicePlan)
    }

    const showCredentialsForm = selectedServiceKindObject
      && selectedServiceKindObject.spec.serviceAccessEnabled
      && selectedServicePlanObject
      && !selectedServicePlanObject.spec.serviceAccessDisabled
      && !planSchemaFound
      && bindingSelectData.length > 0

    return (
      <div>
        <Form {...formConfig} onSubmit={this.handleSubmit}>
          {selectedCloud && (
            <>
              <Card style={{ marginBottom: '20px' }}>
                <Alert
                  message="Cloud service"
                  description="Select the cloud service you would like to use."
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
                      {filteredServiceKinds.map(k => <Select.Option key={k.metadata.name} value={k.metadata.name}>{k.spec.displayName || k.metadata.name}</Select.Option>)}
                    </Select>
                  )}
                  {selectedServiceKind && (
                    <Alert
                      style={{ margin: '10px 0' }}
                      type="info"
                      message={selectedServiceKindObject.spec.displayName}
                      description={<>
                        <Paragraph>{selectedServiceKindObject.spec.description}</Paragraph>
                        {Boolean(selectedServiceKindObject.spec.documentationURL) && (
                          <Paragraph style={{ marginBottom: 0 }}>Documentation: <a target="_blank" rel="noopener noreferrer" href={selectedServiceKindObject.spec.documentationURL}>{selectedServiceKindObject.spec.documentationURL}</a></Paragraph>
                        )}
                      </>}
                    />
                  )}
                </Form.Item>

                <FormErrorMessage message={formErrorMessage} />
                {selectedServiceKind && (
                  <ServiceOptionsForm
                    team={this.props.team}
                    cluster={this.props.cluster}
                    selectedServiceKind={selectedServiceKind}
                    servicePlans={filteredServicePlans}
                    teamServices={this.props.teamServices}
                    onServicePlanSelected={this.handleServicePlanSelected}
                    onPlanValuesChange={this.setPlanValues}
                    validationErrors={this.state.validationErrors}
                    wrappedComponentRef={inst => this.serviceOptionsForm = inst}
                  />
                )}
              </Card>
              {showCredentialsForm && (
                <Collapse defaultActiveKey={['bindings']}>
                  <Panel header="Create service access" key="bindings">
                    <Alert
                      message="Add service access for your already existing cluster namespaces, check the required namespaces below. Alternatively, this can also be done after your service is created"
                      type="info"
                      showIcon
                      style={{ marginBottom: '20px' }}
                    />

                    {bindingSelectData.map(c => (
                      <Row key={c.value} style={{ marginBottom: '10px', padding: '10px' }}>
                        <Col>
                          <Paragraph><Text strong>Cluster</Text><Text style={{ fontFamily: 'monospace', marginLeft: '10px' }}>{c.title}</Text></Paragraph>
                          {c.children.map(ns => {
                            const checked = this.state.bindingsToCreate.includes(ns.value)
                            return (
                              <Form.Item
                                style={{ marginBottom: 0 }}
                                key={ns.value}
                                colon={false}
                                label={<Checkbox key={ns.value} onChange={(e) => this.onChange(e.target.checked, ns.value)()}>{ns.title}</Checkbox>}
                                labelCol={{ span: 24 }}
                                wrapperCol={{ span: 12 }}
                              >
                                {getFieldDecorator(`${ns.value}-secretName`, {
                                  rules: [
                                    { required: checked, message: 'Please enter the secret name or un-check namespace!' },
                                    { pattern: '^[a-z][a-z0-9-]{0,38}[a-z0-9]$', message: 'Name must consist of lower case alphanumeric characters or "-", it must start with a letter and end with an alphanumeric and must be no longer than 40 characters' },
                                  ]
                                })(
                                  <Input disabled={!checked} placeholder="Secret name" />
                                )}
                              </Form.Item>
                            )
                          })}
                        </Col>
                      </Row>
                    ))}

                  </Panel>
                </Collapse>
              )}
            </>
          )}
          <Form.Item style={{ marginTop: '20px', marginBottom: 0 }}>
            <Button type="primary" htmlType="submit" loading={submitting} disabled={this.disableButton()}>
              {this.state.submitButtonText}
            </Button>
            <Button type="link" onClick={this.cancel}>Cancel</Button>
          </Form.Item>
        </Form>
      </div>
    )
  }
}

const WrappedServiceBuildForm = Form.create({ name: 'new_team_service_build' })(ServiceBuildForm)

export default WrappedServiceBuildForm
