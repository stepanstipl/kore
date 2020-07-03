import * as React from 'react'
import { Table } from 'antd'
import PropTypes from 'prop-types'

export default class CostTable extends React.Component {
  static propTypes = {
    costs: PropTypes.object.isRequired,
    hideHourly: PropTypes.bool,
    hideMonthly: PropTypes.bool,
    alwaysShowHeader: PropTypes.bool,
    style: PropTypes.object
  }

  render() {
    const { costs, hideHourly, hideMonthly, alwaysShowHeader, style } = this.props
    if (!costs) {
      return null
    }
    const formatter = new Intl.NumberFormat(undefined, { maximumSignificantDigits: 3 }).format
    const columns = []
    let showHeader = true
    if (costs.minCost === costs.maxCost && costs.minCost === costs.typicalCost) {
      columns.push({ title: 'Approximate Cost', dataIndex: 'typicalCost', key: 'typicalCost', width: '100%' })
      showHeader = alwaysShowHeader
    } else {
      columns.push({ title: 'Minimum', dataIndex: 'minCost', key: 'minCost', width: '33%' })
      columns.push({ title: 'Initial', dataIndex: 'typicalCost', key: 'typicalCost', width: '33%' })
      columns.push({ title: 'Maximum', dataIndex: 'maxCost', key: 'maxCost', width: '34%' })
    }
    const rows = []
    if (!hideHourly) {
      rows.push({ key: 'hourly', minCost: `$${formatter(costs.minCost/1000000)}/hr`, typicalCost: `$${formatter(costs.typicalCost/1000000)}/hr`, maxCost: `$${formatter(costs.maxCost/1000000)}/hr` })
    }
    if (!hideMonthly) {
      rows.push({ key: 'monthly', minCost: `$${formatter((costs.minCost*730.5)/1000000)}/mo`, typicalCost: `$${formatter((costs.typicalCost*730.5)/1000000)}/mo`, maxCost: `$${formatter((costs.maxCost*730.5)/1000000)}/mo` })
    }
    return <Table style={style} size="small" showHeader={showHeader} pagination={false} dataSource={rows} columns={columns} />
  }
}
