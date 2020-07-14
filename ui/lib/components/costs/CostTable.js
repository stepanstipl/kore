import * as React from 'react'
import { Table } from 'antd'
import PropTypes from 'prop-types'
import { formatHourlyCost, formatDailyCost, formatMonthlyCost } from '../../utils/cost-formatters'

export default class CostTable extends React.Component {
  static propTypes = {
    costs: PropTypes.object.isRequired,
    hideHourly: PropTypes.bool,
    hideDaily: PropTypes.bool,
    hideMonthly: PropTypes.bool,
    alwaysShowHeader: PropTypes.bool,
    style: PropTypes.object
  }

  render() {
    const { costs, hideHourly, hideDaily, hideMonthly, alwaysShowHeader, style } = this.props
    if (!costs) {
      return null
    }
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
      rows.push({ key: 'hourly', minCost: formatHourlyCost(costs.minCost), typicalCost: formatHourlyCost(costs.typicalCost), maxCost: formatHourlyCost(costs.maxCost) })
    }
    if (!hideDaily) {
      rows.push({ key: 'daily', minCost: formatDailyCost(costs.minCost), typicalCost: formatDailyCost(costs.typicalCost), maxCost: formatDailyCost(costs.maxCost) })
    }
    if (!hideMonthly) {
      rows.push({ key: 'monthly', minCost: formatMonthlyCost(costs.minCost), typicalCost: formatMonthlyCost(costs.typicalCost), maxCost: formatMonthlyCost(costs.maxCost) })
    }
    return <Table style={style} size="small" showHeader={showHeader} pagination={false} dataSource={rows} columns={columns} />
  }
}
