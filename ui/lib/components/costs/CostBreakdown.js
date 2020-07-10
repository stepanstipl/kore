import * as React from 'react'
import { Table, Radio, Typography } from 'antd'
import PropTypes from 'prop-types'
import { formatMonthlyCost, formatDailyCost, formatHourlyCost } from '../../utils/cost-formatters'

export default class CostBreakdown extends React.Component {
  static propTypes = {
    costs: PropTypes.object,
    style: PropTypes.object
  }
  state = {
    basis: 'month'
  }

  render() {
    const { costs, style } = this.props
    const { basis } = this.state
    if (!costs || !costs.costElements || costs.costElements.length === 0) {
      return null
    }
    let formatter = null
    let basisWarn = null
    switch (basis) {
      case 'month': {
        formatter = formatMonthlyCost
        basisWarn = <>Based on an average 730 hour month. Sustained usage discounts (if available) are <b>not included</b> in this monthly cost estimation</>
        break
      }
      case 'day': {
        formatter = formatDailyCost
        basisWarn = <>Sustained usage discounts (if available) are <b>not included</b> in this daily cost estimation</>
        break
      }
      case 'hour': {
        formatter = formatHourlyCost
        break
      }
    }
    const columns = []
    columns.push({ title: 'Component', dataIndex: 'component', key: 'component', width: '40%' })
    columns.push({ title: 'Minimum', dataIndex: 'minCost', key: 'minCost', width: '20%' })
    columns.push({ title: 'Typical', dataIndex: 'typicalCost', key: 'typicalCost', width: '20%' })
    columns.push({ title: 'Maximum', dataIndex: 'maxCost', key: 'maxCost', width: '20%' })
    const rows = []
    costs.costElements.forEach((el) => {
      rows.push({ key: el.name, component: el.name, minCost: formatter(el.minCost), typicalCost: formatter(el.typicalCost), maxCost: formatter(el.maxCost) })
    })
    rows.push({ key: 'Total', component: <><b>Total</b></>, minCost: formatter(costs.minCost), typicalCost: formatter(costs.typicalCost), maxCost: formatter(costs.maxCost) })
    return (
      <Table 
        style={style} 
        size="small" 
        pagination={false} 
        dataSource={rows} 
        columns={columns} 
        footer={() => (
          <>
            <Radio.Group size="small" buttonStyle="solid" onChange={(e) => this.setState({ basis: e.target.value })} defaultValue="month">
              <Radio.Button value="hour">Hourly</Radio.Button>
              <Radio.Button value="day">Daily</Radio.Button>
              <Radio.Button value="month">Monthly</Radio.Button>
            </Radio.Group>
            {basisWarn ? <Typography.Paragraph type="secondary" style={{ paddingTop: '10px' }}>{basisWarn}</Typography.Paragraph> : null}
          </>
        )} 
      />
    )
  }
}
