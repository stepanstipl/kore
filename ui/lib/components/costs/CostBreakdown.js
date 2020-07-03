import * as React from 'react'
import { Table, Switch } from 'antd'
import PropTypes from 'prop-types'

export default class CostBreakdown extends React.Component {
  static propTypes = {
    costs: PropTypes.object,
    style: PropTypes.object
  }
  state = {
    monthly: true
  }

  render() {
    const { costs, style } = this.props
    const { monthly } = this.state
    if (!costs || !costs.costElements || costs.costElements.length === 0) {
      return null
    }
    const currFormatter = new Intl.NumberFormat(undefined, { maximumSignificantDigits: monthly ? 4 : 3 }).format
    const formatter = !monthly ? (c) => `$${currFormatter(c/1000000)}/hr` : (c) => `$${currFormatter((c*730.5)/1000000)}/mo`
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
        footer={() => <>Cost basis: <Switch className="double-on-switch" checkedChildren="Monthly" unCheckedChildren="Hourly" onChange={(c) => this.setState({ monthly: c })} defaultChecked={true} /></>} 
      />
    )
  }
}
