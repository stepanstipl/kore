import * as React from 'react'
import PropTypes from 'prop-types'
import Router from 'next/router'
import { Alert, Button, Card, Checkbox, Col, Collapse, Form, Input, message, Row, Typography, Select } from 'antd'
const { Paragraph, Text } = Typography
const { Panel } = Collapse

import redirect from '../../../utils/redirect'
import CloudSelector from '../../common/CloudSelector'
import ServiceOptionsForm from '../../services/ServiceOptionsForm'
import FormErrorMessage from '../../forms/FormErrorMessage'
import KoreApi from '../../../kore-api'
import asyncForEach from '../../../utils/async-foreach'
import V1ServiceSpec from '../../../kore-api/model/V1ServiceSpec'
import V1Service from '../../../kore-api/model/V1Service'
import V1ObjectMeta from '../../../kore-api/model/V1ObjectMeta'
import V1ServiceCredentials from '../../../kore-api/model/V1ServiceCredentials'
import V1ServiceCredentialsSpec from '../../../kore-api/model/V1ServiceCredentialsSpec'
import { NewV1ObjectMeta, NewV1Ownership } from '../../../utils/model'
import { getKoreLabel } from '../../../utils/crd-helpers'

class ServiceBuildForm extends React.Component {
  static propTypes = {
    form: PropTypes.any.isRequired,
    skipButtonText: PropTypes.string,
    team: PropTypes.object.isRequired,
    teamServices: PropTypes.array.isRequired,
    user: PropTypes.object.isRequired
  }

  constructor(props) {
    super(props)
    this.state = {
      submitButtonText: 'Save',
      skipButtonText: this.props.skipButtonText || 'Skip',
      submitting: false,
      formErrorMessage: false,
      selectedCloud: false,
      selectedServiceKind: '',
      selectedServicePlan: false,
      dataLoading: true,
      servicePlanOverride: null,
      validationErrors: null,
      bindingsToCreate: [],
      planSchemaFound: false
    }
  }

  async fetchComponentData() {
    const api = await KoreApi.client()
    const [ serviceKinds, servicePlans, clusters, namespaceClaims ] = await Promise.all([
      api.ListServiceKinds(),
      api.ListServicePlans(),
      api.ListClusters(this.props.team.metadata.name),
      api.ListNamespaces(this.props.team.metadata.name)
    ])
    const bindingsData = {}
    namespaceClaims.items.forEach(ns => {
      const cluster = clusters.items.find(c => c.metadata.name === ns.spec.cluster.name)
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
    const selectedServicePlan = this.state.servicePlans.items.find(p => p.metadata.name === values.servicePlan)

    const serviceResource = new V1Service()
    serviceResource.setApiVersion('services.compute.kore.appvia.io/v1')
    serviceResource.setKind('Service')

    const meta = new V1ObjectMeta()
    meta.setName(values.serviceName)
    meta.setNamespace(this.props.team.metadata.name)
    serviceResource.setMetadata(meta)

    const serviceSpec = new V1ServiceSpec()
    serviceSpec.setKind(selectedServicePlan.spec.kind)
    serviceSpec.setPlan(selectedServicePlan.metadata.name)
    if (this.state.servicePlanOverride) {
      serviceSpec.setConfiguration(this.state.servicePlanOverride)
    } else {
      serviceSpec.setConfiguration({ ...selectedServicePlan.spec.configuration })
    }

    serviceResource.setSpec(serviceSpec)
    return serviceResource
  }

  getServiceCredentialsResource = (name, secretName, service, cluster, namespaceClaim) => {
    const serviceCredential = new V1ServiceCredentials()
    serviceCredential.setApiVersion('servicecredentials.services.kore.appvia.io/v1')
    serviceCredential.setKind('ServiceCredentials')
    serviceCredential.setMetadata(NewV1ObjectMeta(name, this.props.team.metadata.name))

    const serviceCredentialSpec = new V1ServiceCredentialsSpec()
    serviceCredentialSpec.setKind(service.spec.kind)
    serviceCredentialSpec.setService(NewV1Ownership({
      group: service.apiVersion.split('/')[0],
      version: service.apiVersion.split('/')[1],
      kind: service.kind,
      name: service.metadata.name,
      namespace: this.props.team.metadata.name
    }))
    serviceCredentialSpec.setCluster(NewV1Ownership({
      group: cluster.apiVersion.split('/')[0],
      version: cluster.apiVersion.split('/')[1],
      kind: cluster.kind,
      name: cluster.metadata.name,
      namespace: this.props.team.metadata.name
    }))
    serviceCredentialSpec.setClusterNamespace(namespaceClaim.spec.name)
    serviceCredentialSpec.setSecretName(secretName)
    serviceCredentialSpec.setConfiguration({})

    serviceCredential.setSpec(serviceCredentialSpec)

    return serviceCredential
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
          console.log('checking existing credentials', existing)
          if (existing) {
            console.log('setting exsiting error on', `${b}-secretName`)
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
        const api = await KoreApi.client()
        const service = await api.UpdateService(this.props.team.metadata.name, values.serviceName, this.getServiceResource(values))
        message.loading('Service build requested...')

        if (this.state.bindingsToCreate.length > 0) {
          await asyncForEach(this.state.bindingsToCreate, async (bindingNamespace) => {
            const namespaceClaim = this.state.namespaceClaims.items.find(ns => ns.metadata.name === bindingNamespace)
            try {
              const secretName = this.props.form.getFieldValue(`${bindingNamespace}-secretName`)
              const cluster = this.state.clusters.items.find(c => c.metadata.name === namespaceClaim.spec.cluster.name)
              const credentialName = `${cluster.metadata.name}-${namespaceClaim.spec.name}-${secretName}`
              const resource = this.getServiceCredentialsResource(credentialName, secretName, service, cluster, namespaceClaim)
              await api.UpdateServiceCredentials(this.props.team.metadata.name, credentialName, resource)
              message.loading(`Service binding for namespace "${namespaceClaim.spec.name}" requested...`)
            } catch (error) {
              console.error('Error creating service binding', error)
              message.error(`Failed to create service binding for namespace "${namespaceClaim.spec.name}"`)
            }
          })
        }

        return redirect({
          router: Router,
          path: `/teams/${this.props.team.metadata.name}/services/${values.serviceName}`
        })
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

  handleSelectCloud = cloud => {
    this.setState({
      selectedCloud: cloud,
      planOverride: null,
      validationErrors: null
    })
  }

  handleSelectKind = (kind) => {
    this.setState({
      selectedServiceKind: kind,
      selectedServicePlan: false,
      servicePlanOverride: null,
      validationErrors: null
    })
  }

  handleServicePlanOverride = servicePlanOverrides => {
    this.setState({ servicePlanOverride: servicePlanOverrides })
  }

  handleServicePlanSelected = async (plan) => {
    this.setState({ selectedServicePlan: plan })
    try {
      // check if there is a schema for the binding for the selected service kind/plan
      const schema = await (await KoreApi.client()).GetServiceCredentialSchema(plan)
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

  render() {
    if (this.state.dataLoading || !this.props.team) {
      return null
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

    return (
      <div>
        <CloudSelector showCustom={false} selectedCloud={selectedCloud} handleSelectCloud={this.handleSelectCloud} enabledCloudList={['AWS']}/>
        {selectedCloud && (
          <Form {...formConfig} onSubmit={this.handleSubmit}>
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
                  selectedServiceKind={selectedServiceKind}
                  servicePlans={filteredServicePlans}
                  teamServices={this.props.teamServices}
                  onServicePlanSelected={this.handleServicePlanSelected}
                  onServicePlanOverridden={this.handleServicePlanOverride}
                  validationErrors={this.state.validationErrors}
                  wrappedComponentRef={inst => this.serviceOptionsForm = inst}
                />
              )}
            </Card>
            {selectedServicePlan && !planSchemaFound && bindingSelectData.length > 0 && (
              <Collapse>
                <Panel header="Optional: Create service bindings" key="bindings">
                  <Alert
                    message="Add service bindings for your already existing cluster namespaces, check the required namespaces below. Alternatively, this can also be done after your service is created"
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
            <Form.Item style={{ marginTop: '20px', marginBottom: 0 }}>
              <Button type="primary" htmlType="submit" loading={submitting} disabled={this.disableButton()}>
                {this.state.submitButtonText}
              </Button>
            </Form.Item>
          </Form>
        )}
      </div>
    )
  }
}

const WrappedServiceBuildForm = Form.create({ name: 'new_team_service_build' })(ServiceBuildForm)

export default WrappedServiceBuildForm
