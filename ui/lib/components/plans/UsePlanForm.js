import * as React from 'react'
import PropTypes from 'prop-types'
import { set } from 'lodash'

import copy from '../../utils/object-copy'
import KoreApi from '../../kore-api'
import PlanViewEdit from './PlanViewEdit'
import { Icon } from 'antd'
import { errorMessage } from '../../utils/message'

/**
 * UsePlanForm is for *using* a plan to configure a cluster, service or service credential.
 *
 * To *manage* a plan (create, view, edit the plan itself), use Manage(Service/Cluster)PlanForm.
 */
class UsePlanForm extends React.Component {
  static propTypes = {
    team: PropTypes.object.isRequired,
    cluster: PropTypes.object,
    resourceType: PropTypes.oneOf(['cluster', 'service', 'servicecredential']).isRequired,
    kind: PropTypes.string.isRequired,
    plan: PropTypes.string.isRequired,
    planValues: PropTypes.object,
    onPlanValuesChange: PropTypes.func,
    validationErrors: PropTypes.array,
    mode: PropTypes.oneOf(['create', 'edit', 'view']).isRequired,
  }

  static initialState = {
    dataLoading: true,
    schema: null,
    editableParams: [],
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
    this.setState({ dataLoading: true })

    let planDetails, schema, editableParams, planValues

    try {
      switch (this.props.resourceType) {
      case 'cluster':
        planDetails = await (await KoreApi.client()).GetTeamPlanDetails(this.props.team.metadata.name, this.props.plan);
        [schema, editableParams, planValues] = [planDetails.schema, planDetails.editableParams, planDetails.plan.configuration]
        break
      case 'service':
        planDetails = await (await KoreApi.client()).GetServicePlanDetails(this.props.plan, this.props.team.metadata.name, this.props.cluster.metadata.name);
        [schema, editableParams, planValues] = [planDetails.schema, planDetails.editableParams, planDetails.configuration]
        break
      case 'servicecredential':
        planDetails = await (await KoreApi.client()).GetServicePlanDetails(this.props.plan, this.props.team.metadata.name, this.props.cluster.metadata.name)
        schema = planDetails.credentialSchema
        editableParams = ['*']
        planValues = {}
        break
      }
    } catch (err) {
      errorMessage(`Failed to load plan: ${err}`)
      return
    }

    if (schema && typeof schema === 'string') {
      schema = JSON.parse(schema)
    }

    this.setState({
      ...this.state,
      schema: schema || { properties:[] },
      editableParams: editableParams || [],
      // Overwrite plan values only if it's still set to the default value
      planValues: this.state.planValues === UsePlanForm.initialState.planValues ? copy(planValues || {}) : this.state.planValues,
      dataLoading: false
    })

    this.props.onPlanValuesChange && this.props.onPlanValuesChange({ ...this.state.planValues })
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
      this.props.onPlanValuesChange && this.props.onPlanValuesChange({ ...planValues })
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
          editableParams={this.state.editableParams}
          onPlanValueChange={(n, v) => this.onValueChange(n, v)}
          validationErrors={this.props.validationErrors}
        />
      </>
    )
  }
}

export default UsePlanForm

