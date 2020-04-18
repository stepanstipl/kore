import * as React from 'react'
import PropTypes from 'prop-types'
import { Form, Checkbox } from 'antd'
import { set } from 'lodash'

import KoreApi from '../../kore-api'
import copy from '../../utils/object-copy'
import PlanOption from './PlanOption'

class PlanOptionsForm extends React.Component {
  static propTypes = {
    team: PropTypes.object.isRequired,
    plan: PropTypes.string.isRequired,
    planValues: PropTypes.object,
    onPlanChange: PropTypes.func,
    validationErrors: PropTypes.array,
    mode: PropTypes.oneOf(['create','edit','view']).isRequired
  }
  static initialState = {
    dataLoading: true,
    schema: null,
    parameterEditable: {},
    planSpec: null,
    planValues: {},
  }

  constructor(props) {
    super(props)
    // Use passed-in plan values if we have them.
    const planValues = this.props.planValues ? this.props.planValues : PlanOptionsForm.initialState.planValues
    this.state = { 
      ...PlanOptionsForm.initialState,
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
      this.setState({ ...PlanOptionsForm.initialState })
      this.componentDidUpdateComplete = this.fetchComponentData()
    }
    if (this.props.mode !== prevProps.mode) {
      this.setState({ showReadOnly: this.props.mode === 'view' })
    }
    if (this.props.planValues !== prevProps.planValues) {
      this.setState({
        planValues: this.props.planValues
      })
    }
  }

  async fetchComponentData() {
    const planDetails = await (await KoreApi.client()).GetTeamPlanDetails(this.props.team.metadata.name, this.props.plan)
    this.setState({
      ...this.state,
      schema: JSON.parse(planDetails.schema),
      parameterEditable: planDetails.parameterEditable,
      planSpec: planDetails.plan,
      // Overwrite plan values only if it's still set to the default value
      planValues: this.state.planValues === PlanOptionsForm.initialState.planValues ? copy(planDetails.plan.configuration) : this.state.planValues,
      showReadOnly: this.props.mode === 'view',
      dataLoading: false
    })
  }

  onValueChange(name, value) {
    // Texture this back into a state update using the nifty lodash set function:
    const newPlanValues = set({ ...this.state.planValues }, name, value)
    this.setState({
      planValues: newPlanValues
    })
    this.props.onPlanChange && this.props.onPlanChange(newPlanValues)
  }

  handleShowReadOnlyChange = (checked) => {
    this.setState({
      showReadOnly: checked
    })
  }

  render() {
    if (this.state.dataLoading) {
      return (
        <div>Loading plan details...</div>
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
            <PlanOption 
              key={name} 
              name={name} 
              property={this.state.schema.properties[name]} 
              value={this.state.planValues[name]} 
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

export default PlanOptionsForm

