import * as React from 'react'
import PropTypes from 'prop-types'
import Router from 'next/router'
import { Button, Form, message } from 'antd'

import redirect from '../../../utils/redirect'
import ServiceKindSelector from '../../services/ServiceKindSelector'
import ServiceOptionsForm from '../../services/ServiceOptionsForm'
import FormErrorMessage from '../../forms/FormErrorMessage'
import KoreApi from '../../../kore-api'
import V1ServiceSpec from '../../../kore-api/model/V1ServiceSpec'
import V1Service from '../../../kore-api/model/V1Service'
import V1ObjectMeta from '../../../kore-api/model/V1ObjectMeta'

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
      selectedServiceKind: '',
      dataLoading: true,
      servicePlanOverride: null,
      validationErrors: null
    }
  }

  async fetchComponentData() {
    const servicePlans = await (await KoreApi.client()).ListServicePlans()
    return { servicePlans }
  }

  componentDidMountComplete = null
  componentDidMount() {
    // Assign the promise chain to a variable so tests can wait for it to complete.
    this.componentDidMountComplete = Promise.resolve().then(async () => {
      const { servicePlans } = await this.fetchComponentData()
      this.setState({
        servicePlans: servicePlans,
        dataLoading: false
      })
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

  handleSubmit = e => {
    e.preventDefault()

    this.serviceOptionsForm.props.form.validateFields(async (err, values) => {
      if (err) {
        this.setState({
          ...this.state,
          formErrorMessage: 'Validation failed'
        })
        return
      }
      this.setState({
        ...this.state,
        submitting: true,
        formErrorMessage: false
      })
      try {
        await (await KoreApi.client()).UpdateService(
          this.props.team.metadata.name,
          values.serviceName,
          this.getServiceResource(values))
        message.loading('Service build requested...')
        return redirect({
          router: Router,
          path: `/teams/${this.props.team.metadata.name}`
        })
      } catch (err) {
        this.setState({
          ...this.state,
          submitting: false,
          formErrorMessage: (err.fieldErrors && err.message) ? err.message : 'An error occurred requesting the service, please try again',
          validationErrors: err.fieldErrors // This will be undefined on non-validation errors, which is fine.
        })
      }
    })
  }

  handleSelectKind = kind => {
    this.setState({
      selectedServiceKind: kind,
      servicePlanOverride: null,
      validationErrors: null
    })
  }

  handleServicePlanOverride = servicePlanOverrides => {
    this.setState({
      servicePlanOverride: servicePlanOverrides
    })
  }

  serviceBuildForm = () => {
    const { submitting, selectedServiceKind, formErrorMessage } = this.state
    const filteredServicePlans = this.state.servicePlans.items.filter(p => p.spec.kind === selectedServiceKind)
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

    return (
      <Form {...formConfig} onSubmit={this.handleSubmit}>
        <FormErrorMessage message={formErrorMessage} />
        <ServiceOptionsForm
          team={this.props.team}
          selectedServiceKind={selectedServiceKind}
          servicePlans={filteredServicePlans}
          teamServices={this.props.teamServices}
          onServicePlanOverridden={this.handleServicePlanOverride}
          validationErrors={this.state.validationErrors}
          wrappedComponentRef={inst => this.serviceOptionsForm = inst}
        />
        <Form.Item style={{ marginTop: '20px' }}>
          <Button type="primary" htmlType="submit" loading={submitting}>
            {this.state.submitButtonText}
          </Button>
        </Form.Item>
      </Form>
    )
  }

  render() {
    if (this.state.dataLoading || !this.props.team) {
      return null
    }

    const { selectedServiceKind } = this.state

    return (
      <div>
        <ServiceKindSelector
          selectedServiceKind={selectedServiceKind}
          handleSelectKind={this.handleSelectKind} />
        {selectedServiceKind ? <this.serviceBuildForm /> : null}
      </div>
    )
  }
}

const WrappedServiceBuildForm = Form.create({ name: 'new_team_service_build' })(ServiceBuildForm)

export default WrappedServiceBuildForm
