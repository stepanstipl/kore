import * as React from 'react'
import PropTypes from 'prop-types'
import { Form, Checkbox } from 'antd'
import { set } from 'lodash'

import KoreApi from '../../kore-api'
import copy from '../../utils/object-copy'
import ServicePlanOption from './ServicePlanOption'

class ServicePlanOptionsForm extends React.Component {
  static propTypes = {
    team: PropTypes.object.isRequired,
    servicePlan: PropTypes.string.isRequired,
    servicePlanValues: PropTypes.object,
    onServicePlanChange: PropTypes.func,
    validationErrors: PropTypes.array,
    mode: PropTypes.oneOf(['create','edit','view']).isRequired
  }
  static initialState = {
    dataLoading: true,
    schema: null,
    parameterEditable: {},
    servicePlanSpec: null,
    servicePlanValues: {},
  }

  constructor(props) {
    super(props)
    // Use passed-in service plan values if we have them.
    const servicePlanValues = this.props.servicePlanValues ? this.props.servicePlanValues : ServicePlanOptionsForm.initialState.servicePlanValues
    this.state = { 
      ...ServicePlanOptionsForm.initialState,
      servicePlanValues
    }
  }

  componentDidMountComplete = null
  componentDidMount() {
    this.componentDidMountComplete = this.fetchComponentData()
  }

  componentDidUpdateComplete = null
  componentDidUpdate(prevProps) {
    if (this.props.servicePlan !== prevProps.servicePlan || this.props.team !== prevProps.team) {
      this.setState({ ...ServicePlanOptionsForm.initialState })
      this.componentDidUpdateComplete = this.fetchComponentData()
    }
    if (this.props.mode !== prevProps.mode) {
      this.setState({ showReadOnly: this.props.mode === 'view' })
    }
    if (this.props.servicePlanValues !== prevProps.servicePlanValues) {
      this.setState({
        servicePlanValues: this.props.servicePlanValues
      })
    }
  }

  async fetchComponentData() {
    const servicePlanDetails = await (await KoreApi.client()).GetTeamServicePlanDetails(this.props.team.metadata.name, this.props.servicePlan)
    this.setState({
      ...this.state,
      schema: JSON.parse(servicePlanDetails.schema),
      parameterEditable: servicePlanDetails.parameterEditable || {},
      servicePlanSpec: servicePlanDetails.servicePlan,
      // Overwrite service plan values only if it's still set to the default value
      servicePlanValues: this.state.servicePlanValues === ServicePlanOptionsForm.initialState.servicePlanValues ? copy(servicePlanDetails.servicePlan.configuration) : this.state.servicePlanValues,
      showReadOnly: this.props.mode === 'view',
      dataLoading: false
    })
  }

  onValueChange(name, value) {
    // Texture this back into a state update using the nifty lodash set function:
    const newServicePlanValues = set({ ...this.state.servicePlanValues }, name, value)
    this.setState({
      servicePlanValues: newServicePlanValues
    })
    this.props.onServicePlanChange && this.props.onServicePlanChange(newServicePlanValues)
  }

  handleShowReadOnlyChange = (checked) => {
    this.setState({
      showReadOnly: checked
    })
  }

  render() {
    if (this.state.dataLoading) {
      return (
        <div>Loading service plan details...</div>
      )
    }

    return (
      <>
        {this.props.mode !== 'view' ? (
          <Form.Item label="Show read-only parameters">
            <Checkbox onChange={(v) => this.handleShowReadOnlyChange(v.target.checked)} checked={this.state.showReadOnly} />
          </Form.Item>
        ): null}
        {Object.keys(this.state.schema.properties).map((name) => {
          const editable = this.props.mode !== 'view' &&
            this.state.parameterEditable[name] === true &&
            (this.props.mode === 'create' || !this.state.schema.properties[name].immutable) // Disallow editing of params which can only be set at create time.

          return (
            <ServicePlanOption 
              key={name} 
              name={name} 
              property={this.state.schema.properties[name]} 
              value={this.state.servicePlanValues[name]} 
              hideNonEditable={!this.state.showReadOnly} 
              editable={editable} 
              onChange={(n, v) => this.onValueChange(n, v)}
              validationErrors={this.props.validationErrors} />
          )
        })}
      </>
    )
  }
}

export default ServicePlanOptionsForm

