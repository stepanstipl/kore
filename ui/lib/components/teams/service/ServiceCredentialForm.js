import * as React from 'react'
import PropTypes from 'prop-types'
import { patterns } from '../../../utils/validation'
import { Button, Form, Input, Alert, Select, Collapse } from 'antd'
import KoreApi from '../../../kore-api'
import UsePlanForm from '../../plans/UsePlanForm'
import V1ServiceCredentials from '../../../kore-api/model/V1ServiceCredentials'
import V1ServiceCredentialsSpec from '../../../kore-api/model/V1ServiceCredentialsSpec'
import { NewV1ObjectMeta, NewV1Ownership } from '../../../utils/model'

class ServiceCredentialForm extends React.Component {
  static propTypes = {
    form: PropTypes.any.isRequired,
    team: PropTypes.object.isRequired,
    clusters: PropTypes.object,
    services: PropTypes.object,
    handleSubmit: PropTypes.func.isRequired,
    handleCancel: PropTypes.func.isRequired
  }

  constructor(props) {
    super(props)
    this.state =  {
      clusters: props.clusters,
      services: props.services,
      namespaceClaims: { items: [] },
      servicePlan: null,
      submitting: false,
      formErrorMessage: false,
      dataLoading: true,
      validationErrors: null,
      config: null,
      planSchemaFound: false
    }
  }

  componentDidMountComplete = null
  componentDidMount() {

    // Assign the promise chain to a variable so tests can wait for it to complete.
    this.componentDidMountComplete = Promise.resolve().then(async () => {
      let { services, clusters, namespaceClaims, servicePlan } = this.state

      const api = await KoreApi.client()
      if (!services) {
        services = await api.ListServices(this.props.team.metadata.name)
      }

      if (services && services.items && services.items.length === 1) {
        servicePlan = await api.GetServicePlan(services.items[0].spec.plan)
      }

      if (!clusters) {
        clusters = await api.ListClusters(this.props.team.metadata.name) // eslint-disable-line require-atomic-updates
      }

      if (clusters && clusters.items && clusters.items.length === 1) {
        const allNamespaceClaims = await api.ListNamespaces(this.props.team.metadata.name)
        namespaceClaims = { items: allNamespaceClaims.items.filter(nc => nc.spec.cluster.name === clusters.items[0].metadata.name) }
      }

      this.setState({
        services: services || { items: [] },
        clusters: clusters || { items: [] },
        namespaceClaims: namespaceClaims || [],
        servicePlan,
        dataLoading: false
      })

      this.props.form.validateFields()
    })
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

  onClusterChange = () => {
    Promise.resolve().then(async () => {
      const clusterName = this.props.form.getFieldValue('cluster')
      const api = await KoreApi.client()
      const namespaceClaims = await api.ListNamespaces(this.props.team.metadata.name)
      this.setState({
        namespaceClaims: { items: namespaceClaims.items.filter(nc => nc.spec.cluster.name === clusterName) }
      })

      this.props.form.resetFields([ 'namespace' ])
    })
  }

  onServiceChange = () => {
    Promise.resolve().then(async () => {
      const { services, servicePlan: prevServicePlan  } = this.state
      const serviceName = this.props.form.getFieldValue('service')
      const service = services.items.find(s => s.metadata.name === serviceName)

      if (prevServicePlan && prevServicePlan.metadata.name === service.spec.plan) {
        return
      }
      
      const api = await KoreApi.client()
      const servicePlan = await api.GetServicePlan(service.spec.plan)
      this.setState({
        servicePlan: servicePlan,
        config: null,
        validationErrors: null
      })
    })
  }

  handleConfigurationUpdate = c => {
    this.setState({
      config: c
    })
  }

  handleSubmit = e => {
    e.preventDefault()

    const { clusters, services, namespaceClaims, config } = this.state

    this.setState({
      submitting: true,
      formErrorMessage: false
    })

    return this.props.form.validateFields(async (err, values) => {
      if (err) {
        return this.setState({
          submitting: false
        })
      }

      if (!values.namespace) {
        return this.setState({
          submitting: false,
          formErrorMessage: 'Please select the namespace!',
          validationErrors: null
        })
      }

      const cluster = clusters.items.find(c => c.metadata.name === values.cluster)
      const service = services.items.find(s => s.metadata.name === values.service)
      const namespaceClaim = namespaceClaims.items.find(n => n.spec.name === values.namespace)
      const name = values.name

      try {
        const existing = await (await KoreApi.client()).GetServiceCredentials(this.props.team.metadata.name, name)
        if (existing) {
          return this.setState({
            submitting: false,
            formErrorMessage: `A service credential with the name "${name}" already exists`,
            validationErrors: null
          })
        }
      } catch(err) {
        // TODO: we should differentiate between 404 and other errors
      }

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
      serviceCredentialSpec.setConfiguration(config)

      serviceCredential.setSpec(serviceCredentialSpec)

      try {
        const result = await (await KoreApi.client()).UpdateServiceCredentials(
          this.props.team.metadata.name,
          name,
          serviceCredential)

        this.props.form.resetFields()
        this.setState({
          submitting: false
        })
        await this.props.handleSubmit(result)
      } catch (err) {
        this.setState({
          ...this.state,
          submitting: false,
          formErrorMessage: (err.fieldErrors && err.message) ? err.message : 'An error occurred creating the service credential, please try again',
          validationErrors: err.fieldErrors // This will be undefined on non-validation errors, which is fine.
        })
      }
    })
  }

  render() {
    if (this.state.dataLoading) {
      return null
    }

    const { getFieldDecorator, getFieldsError, getFieldError, isFieldTouched } = this.props.form
    const { clusters, services, namespaceClaims, submitting, formErrorMessage, servicePlan, validationErrors } = this.state
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
    const clusterError = isFieldTouched('cluster') && getFieldError('cluster')
    const serviceError = isFieldTouched('service') && getFieldError('service')
    const namespaceError = (isFieldTouched('cluster') || isFieldTouched('namespace')) && getFieldError('namespace')

    const FormErrorMessage = () => {
      if (formErrorMessage) {
        return (
          <Alert
            message={formErrorMessage}
            type="error"
            showIcon
            closable
            style={{ marginBottom: '20px' }}
          />
        )
      }
      return null
    }

    return (
      <Form {...formConfig} onSubmit={this.handleSubmit} style={{ marginBottom: '30px' }}>
        <FormErrorMessage />
        <Form.Item label="Name" validateStatus={nameError ? 'error' : ''} help={nameError || ''}>
          {getFieldDecorator('name', {
            rules: [
              { required: true, message: 'Please enter the service credential name!' },
              { ...patterns.uriCompatible63CharMax }
            ]
          })(
            <Input placeholder="Name" />,
          )}
        </Form.Item>
        <Form.Item label="Service" validateStatus={serviceError ? 'error' : ''} help={serviceError || ''}>
          {getFieldDecorator('service', {
            rules: [{ required: true, message: 'Please select the service!' }],
            initialValue: services.items.length === 1 ? services.items[0].metadata.name : undefined
          })(
            <Select placeholder="Service" onChange={this.onServiceChange}>
              {services.items.map(s => (
                <Select.Option key={s.metadata.name} value={s.metadata.name}>{s.metadata.name}</Select.Option>
              ))}
            </Select>
          )}
        </Form.Item>
        <Form.Item label="Cluster" validateStatus={clusterError ? 'error' : ''} help={clusterError || ''}>
          {getFieldDecorator('cluster', {
            rules: [{ required: true, message: 'Please select the cluster!' }],
            initialValue: clusters.items.length === 1 ? clusters.items[0].metadata.name : undefined
          })(
            <Select placeholder="Cluster" onChange={this.onClusterChange}>
              {clusters.items.map(c => (
                <Select.Option key={c.metadata.name} value={c.metadata.name}>{c.metadata.name}</Select.Option>
              ))}
            </Select>
          )}
        </Form.Item>
        <Form.Item label="Namespace" validateStatus={namespaceError ? 'error' : ''} help={namespaceError || ''}>
          {getFieldDecorator('namespace', {
            rules: [{ required: true, message: 'Please select the namespace!' }],
            initialValue: namespaceClaims.items.length === 1 ? namespaceClaims.items[0].spec.name : undefined
          })(
            <Select placeholder="Namespace">
              {namespaceClaims.items.map(n => (
                <Select.Option key={n.spec.name} value={n.spec.name}>{n.spec.name}</Select.Option>
              ))}
            </Select>
          )}
        </Form.Item>
        {servicePlan && this.state.planSchemaFound ? (
          <Collapse style={{ marginBottom: '24px' }} >
            <Collapse.Panel header="Customize service binding parameters" forceRender={true}>
              <UsePlanForm
                team={this.props.team}
                resourceType="servicecredential"
                kind={servicePlan.spec.kind}
                plan={servicePlan.metadata.name}
                mode="create"
                validationErrors={validationErrors}
                onPlanChange={this.handleConfigurationUpdate}
                schemaFound={(found) => this.setState({ planSchemaFound: found })}
              />
            </Collapse.Panel>
          </Collapse>
        ) : null}
        <Form.Item>
          <Button type="primary" htmlType="submit" loading={submitting} disabled={this.disableButton(getFieldsError())}>Save</Button>
          <Button type="link" onClick={this.cancel}>Cancel</Button>
        </Form.Item>
      </Form>
    )
  }
}

const WrappedServiceCredentialForm = Form.create({ name: 'service_credential' })(ServiceCredentialForm)

export default WrappedServiceCredentialForm
