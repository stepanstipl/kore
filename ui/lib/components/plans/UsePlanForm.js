import * as React from 'react'
import PropTypes from 'prop-types'
import { set } from 'lodash'

import copy from '../../utils/object-copy'
import KoreApi from '../../kore-api'
import PlanViewEdit from './PlanViewEdit'
import { Icon } from 'antd'

/**
 * UsePlanForm is for *using* a plan to configure a cluster, service or service credential.
 *
 * To *manage* a plan (create, view, edit the plan itself), use Manage(Service/Cluster)PlanForm.
 */
class UsePlanForm extends React.Component {
  static propTypes = {
    team: PropTypes.object.isRequired,
    resourceType: PropTypes.oneOf(['cluster', 'service', 'servicecredential']).isRequired,
    kind: PropTypes.string.isRequired,
    plan: PropTypes.string.isRequired,
    planValues: PropTypes.object,
    onPlanChange: PropTypes.func,
    validationErrors: PropTypes.array,
    mode: PropTypes.oneOf(['create', 'edit', 'view']).isRequired,
  }

  static initialState = {
    dataLoading: true,
    schema: null,
    parameterEditable: {},
    planValues: {},
  }

  constructor(props) {
    super(props)
    // Use passed-in plan values if we have them.
    const planValues = this.props.planValues ? this.props.planValues : UsePlanForm.initialState.planValues
    this.state = { 
      ...UsePlanForm.initialState,
      planValues
    }
  }

  componentDidMountComplete = null
  componentDidMount() {
    this.componentDidMountComplete = this.fetchComponentData()
  }

  componentDidUpdateComplete = null
  componentDidUpdate(prevProps) {
    if (this.props.plan !== prevProps.plan || this.props.team !== prevProps.team) {
      this.setState({ ...UsePlanForm.initialState })
      this.componentDidUpdateComplete = this.fetchComponentData()
    }
    if (this.props.planValues !== prevProps.planValues) {
      this.setState({
        planValues: this.props.planValues
      })
    }
  }

  async fetchComponentData() {
    let planDetails, schema, parameterEditable, planValues

    switch (this.props.resourceType) {
    case 'cluster':
      planDetails = await (await KoreApi.client()).GetTeamPlanDetails(this.props.team.metadata.name, this.props.plan);
      [schema, parameterEditable, planValues] = [planDetails.schema, planDetails.parameterEditable, planDetails.plan.configuration]
      break
    case 'service':
      planDetails = await (await KoreApi.client()).GetTeamServicePlanDetails(this.props.team.metadata.name, this.props.plan);
      [schema, parameterEditable, planValues] = [planDetails.schema, planDetails.parameterEditable, planDetails.servicePlan.configuration]
      break
    case 'servicecredential':
      schema = await (await KoreApi.client()).GetServiceCredentialSchema(this.props.team.metadata.name, this.props.plan)
      parameterEditable = { '*': true }
      planValues = {}
      break
    }

    if (schema && typeof schema === 'string') {
      schema = JSON.parse(schema)
    }

    this.setState({
      ...this.state,
      schema: schema || { properties:[] },
      parameterEditable: parameterEditable || {},
      // Overwrite plan values only if it's still set to the default value
      planValues: this.state.planValues === UsePlanForm.initialState.planValues ? copy(planValues || {}) : this.state.planValues,
      dataLoading: false
    })
  }

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
      // Fire a copy of the plan values out if anyone is listening.
      this.props.onPlanChange && this.props.onPlanChange({ ...planValues })
      return { planValues }
    })
  }

  render() {
    if (this.state.dataLoading) {
      return (
        <Icon type="loading" />
      )
    }

    return (
      <>
        <PlanViewEdit
          resourceType={this.props.resourceType}
          mode={this.props.mode}
          manage={false}
          team={this.props.team}
          kind={this.props.kind}
          plan={this.state.planValues}
          schema={this.state.schema}
          parameterEditable={this.state.parameterEditable}
          onPlanValueChange={(n, v) => this.onValueChange(n, v)}
          validationErrors={this.props.validationErrors}
        />
      </>
    )
  }
}

export default UsePlanForm

