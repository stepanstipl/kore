import * as React from 'react'
import PropTypes from 'prop-types'
import { Alert } from 'antd'
import { debounce } from 'lodash'

import KoreApi from '../../kore-api'
import CostBreakdown from './CostBreakdown'

export default class CostEstimate extends React.Component {
  static propTypes = {
    resourceType: PropTypes.oneOf(['service','cluster']).isRequired,
    kind: PropTypes.string.isRequired,
    planValues: PropTypes.object,
    noPriceDataError: PropTypes.string
  }

  state = {
    estimating: false,
    estimate: null,
    planValuesChangedSinceEstimate: false,
    estimateError: null
  }

  componentDidMount() {
    this.componentDidUpdate(null)
  }

  componentDidUpdate(prevProps) {
    // If we don't have planValues yet, or we know we don't have prices available, 
    // we can't do anything.
    if (!this.props.planValues || this.pricingUnavailable()) {
      return
    }

    // Prepare an initial estimate as soon as we've got plan values, unless we've
    // got an error already (so we don't endlessly loop this on an error case)
    if (!(this.state.estimate || this.state.estimateError)) {
      this.estimate()
      return
    }

    // Update the estimate (with a debounce) when the plan values change
    if (!prevProps || prevProps.planValues !== this.props.planValues) {
      if (!this.state.estimating) {
        this.setState({ estimating: true })
      }
      this.debouncedEstimate()
    }
  }

  debouncedEstimate = debounce(() => this.estimate(), 1000)

  estimate = async () => {
    if (!this.props.planValues) {
      return
    }
    if (this.props.resourceType !== 'cluster') {
      this.setState({ estimating: false, estimate: null, estimateError: { message: 'Only cluster estimation supported at this time' } })
      return
    }
    if (!this.state.estimating) {
      this.setState({ estimating: true })
    }
    const currPlan = KoreApi.resources().generatePlanResource(this.props.kind, { configuration: { ...this.props.planValues } })
    const api = await KoreApi.client()
    try {
      const estimate = await api.costestimates.EstimateClusterPlanCost(currPlan)
      this.setState({ estimate, estimateError: null, planValuesChangedSinceEstimate: false, estimating: false })
    } catch (err) {
      this.setState({ estimate: null, estimateError: err, planValuesChangedSinceEstimate: false, estimating: false })
    }
  }

  pricingUnavailable = () => {
    const { estimateError } = this.state

    // If the pricing metadata is not available, we get a specific field error back:
    return estimateError && 
      estimateError.fieldErrors && 
      estimateError.fieldErrors[0].field === 'prices' && 
      estimateError.fieldErrors[0].errCode === 'mustExist'    
  }

  render() {
    const { noPriceDataError } = this.props
    const { estimate, estimateError, estimating } = this.state
    const pricingUnavailable = this.pricingUnavailable()

    return (
      <>
        <CostBreakdown costs={estimate} style={{ marginBottom: '20px' }} loading={estimating} />
        {!estimateError ? null : pricingUnavailable ? (
          <Alert type="error" message={noPriceDataError ? noPriceDataError : 'Pricing data is not currently available. Ask your Kore Administrator to configure the Kore Costs feature to enable cost estimates.'} /> 
        ) : (
          <>
            <Alert 
              type="error" 
              message={<>
                Cannot prepare estimate yet - {estimateError.message}{estimateError.fieldErrors ? ':' : ''}
                <ul>
                  {estimateError.fieldErrors ? estimateError.fieldErrors.map((fe, ind) => <li key={ind}>{fe.message}</li>) : null}
                </ul>
              </>}
            />
          </>
        )}
      </>
    )
  }
}