import { Form } from 'antd'
import PropTypes from 'prop-types'

import KoreApi from '../../kore-api'
import V1ServicePlan from '../../kore-api/model/V1ServicePlan'
import V1ServicePlanSpec from '../../kore-api/model/V1ServicePlanSpec'
import V1ObjectMeta from '../../kore-api/model/V1ObjectMeta'
import canonical from '../../utils/canonical'
import ManagePlanForm from './ManagePlanForm'

class ManageServicePlanForm extends ManagePlanForm {
  static propTypes = {
    kind: PropTypes.string.isRequired
  }

  resourceType = () => 'service'
  
  async fetchComponentData() {
    // Some services have plan-specific schemas, in which case we have to use that schema
    // instead of any default schema registered on the service kind.
    const { data } = this.props
    if (data && data.spec && data.spec.schema) {
      this.setState({
        schema: JSON.parse(data.spec.schema),
        dataLoading: false
      })
      return
    }
    const serviceKind = await (await KoreApi.client()).GetServiceKind(this.props.kind)
    // Use a default empty schema if no schema provided.
    const schema = serviceKind.spec.schema ? JSON.parse(serviceKind.spec.schema) : { properties: [] }
    this.setState({
      schema,
      dataLoading: false
    })
  }

  getMetadataName = (values) => {
    const { data, kind } = this.props
    return (data && data.metadata && data.metadata.name) || `${kind}-${canonical(values.summary)}`
  }

  generateServicePlanResource = values => {
    const metadataName = this.getMetadataName(values)

    const servicePlanResource = new V1ServicePlan()
    servicePlanResource.setApiVersion('services.kore.appvia.io/v1')
    servicePlanResource.setKind('ServicePlan')

    const meta = new V1ObjectMeta()
    meta.setName(metadataName)
    meta.setNamespace('kore')
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
    return { ...defaultConfiguration, ...this.state.planValues }
  }

  process = async (err, values) => {
    if (err) {
      this.setFormSubmitting(false, 'Validation failed')
      return
    }
    try {
      const api = await KoreApi.client()
      const metadataName = this.getMetadataName(values)
      const servicePlanResult = await api.UpdateServicePlan(metadataName, this.generateServicePlanResource({ ...values, configuration: this.generateServicePlanConfiguration() }))
      this.setFormSubmitting(false, null, [])
      return await this.props.handleSubmit(servicePlanResult)
    } catch (err) {
      console.error('Error submitting form', err)
      const message = (err.fieldErrors && err.message) ? err.message : 'An error occurred saving the plan, please try again'
      this.setFormSubmitting(false, message, err.fieldErrors)
    }
  }
}

const WrappedManageServicePlanForm = Form.create({ name: 'servicePlan' })(ManageServicePlanForm)

export default WrappedManageServicePlanForm

