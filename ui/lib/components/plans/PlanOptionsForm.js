import * as React from 'react'
import PropTypes from 'prop-types'
import { Collapse, Checkbox } from 'antd'
import { set } from 'lodash'

import KoreApi from '../../kore-api'
import copy from '../../utils/object-copy'
import PlanOption from './PlanOption'

class PlanOptionsForm extends React.Component {
  static propTypes = {
    team: PropTypes.object.isRequired,
    plan: PropTypes.string.isRequired,
    onPlanChange: PropTypes.func,
    validationErrors: PropTypes.array
  }
  static initialState = {
    dataLoading: true,
    schema: null,
    parameterEditable: {},
    planSpec: null,
    planOverrides: {},
  }

  constructor(props) {
    super(props)
    this.state = { ...PlanOptionsForm.initialState }
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
  }

  async fetchComponentData() {
    const planDetails = await (await KoreApi.client()).GetTeamPlanDetails(this.props.team.metadata.name, this.props.plan)
    this.setState({
      ...this.state,
      schema: JSON.parse(planDetails.schema),
      parameterEditable: planDetails.parameterEditable,
      planSpec: planDetails.plan,
      planValues: copy(planDetails.plan.configuration),
      showReadOnly: false,
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

  handleShowReadOnlyChange = (e) => {
    this.setState({
      showReadOnly: e.target.checked
    })
  }

  render() {
    if (this.state.dataLoading) {
      return (
        <div>Loading plan details...</div>
      )
    }

    return (
      <Collapse>
        <Collapse.Panel header="Customize cluster parameters">
          <Checkbox onChange={this.handleShowReadOnlyChange} checked={this.state.showReadOnly}>Show read-only</Checkbox>
          {Object.keys(this.state.schema.properties).map((name) => 
            <PlanOption 
              key={name} 
              name={name} 
              property={this.state.schema.properties[name]} 
              value={this.state.planValues[name]} 
              hideNonEditable={!this.state.showReadOnly} 
              editable={this.state.parameterEditable[name] === true} 
              onChange={(n, v) => this.onValueChange(n, v)}
              validationErrors={this.props.validationErrors} />
          )}
        </Collapse.Panel>
      </Collapse>
    )
  }
}

export default PlanOptionsForm

