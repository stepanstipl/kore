import * as React from 'react'
import PropTypes from 'prop-types'
import CostTable from './CostTable'
import { Form } from 'antd'

export default class NodePoolCost extends React.Component {
  static propTypes = {
    prices: PropTypes.object,
    nodePool: PropTypes.object,
    help: PropTypes.string.isRequired,
    zoneMultiplier: PropTypes.number,
    priceType: PropTypes.string,
  }

  calculatePoolCost = (nodePool, prices) => {
    if (!prices || !nodePool) {
      return null
    }
    const priceType = this.props.priceType || 'OnDemand'
    const nodePrice = prices[priceType]
    if (!nodePrice) {
      return null
    }

    const nodePoolCosts = { minCost: null, maxCost: null, typicalCost: null }
    const zoneMultiplier = this.props.zoneMultiplier ? this.props.zoneMultiplier : 1
    const size = nodePool.size || nodePool.desiredSize // annoying diff between GKE and EKS naming
    if (nodePool.enableAutoscaler) {
      nodePoolCosts.minCost = nodePool.minSize * nodePrice * zoneMultiplier
      nodePoolCosts.maxCost = nodePool.maxSize * nodePrice * zoneMultiplier
      nodePoolCosts.typicalCost = size * nodePrice * zoneMultiplier
    } else {
      nodePoolCosts.typicalCost = size * nodePrice * zoneMultiplier
      nodePoolCosts.minCost = nodePoolCosts.typicalCost
      nodePoolCosts.maxCost = nodePoolCosts.typicalCost
    }
    return nodePoolCosts
  }

  render() {
    const { prices, nodePool, help } = this.props
    const nodePoolCosts = this.calculatePoolCost(nodePool, prices)
    if (!nodePoolCosts) {
      return null
    }
    return (
      <Form.Item label="Approximate pool cost" help={help}>
        <CostTable costs={nodePoolCosts} alwaysShowHeader={false} />
      </Form.Item>
    )    
  }
}