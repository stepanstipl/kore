import * as React from 'react'
import PropTypes from 'prop-types'
import KoreApi from '../../kore-api'
import CostBreakdown from './CostBreakdown'
import { Button, Alert } from 'antd'

export default class CostEstimate extends React.Component {
  static propTypes = {
    resourceType: PropTypes.oneOf(['service','cluster']).isRequired,
    kind: PropTypes.string.isRequired,
    planValues: PropTypes.object,
    estimateInit: PropTypes.bool
  }

  state = {
    estimate: null,
    planValuesChangedSinceEstimate: false,
    estimateError: null
  }

  componentDidUpdate(prevProps) {
    if (prevProps.planValues !== this.props.planValues) {
      this.setState({ planValuesChangedSinceEstimate: true })
    }
  }

  estimate = async () => {
    if (!this.props.planValues) {
      return
    }
    if (this.props.resourceType !== 'cluster') {
      this.setState({ estimate: null, estimateError: { message: 'Only cluster estimation supported at this time' } })
      return
    }
    const currPlan = KoreApi.resources().generatePlanResource(this.props.kind, { configuration: { ...this.props.planValues } })
    const api = await KoreApi.client()
    try {
      const estimate = await api.costestimates.EstimateClusterPlanCost(currPlan)
      this.setState({ estimate, estimateError: null, planValuesChangedSinceEstimate: false })
    } catch (err) {
      this.setState({ estimate: null, estimateError: err, planValuesChangedSinceEstimate: false })
    }
  }

  render() {
    const { estimate, planValuesChangedSinceEstimate, estimateError } = this.state
    if (!estimate && !estimateError && this.props.estimateInit && this.props.planValues) {
      this.estimate()
    }
    return (
      <>
        <CostBreakdown costs={estimate} style={{ marginBottom: '20px' }} />
        {!estimateError ? null : (
          <>
            <Alert type="warn" message={`Unable to prepare estimate - ${estimateError.message}`} />
            <ul>
              {estimateError.fieldErrors ? estimateError.fieldErrors.map((fe, ind) => <li key={ind}>{fe.message}</li>) : null}
            </ul>
          </>
        )}
        {!estimate ? <Button onClick={() => this.estimate()}>Prepare estimate</Button> : null}
        {estimate && planValuesChangedSinceEstimate ? <Button onClick={() => this.estimate()}><b>Plan values have changed</b> - click to refresh estimate</Button> : null}
      </>
    )
  }
}