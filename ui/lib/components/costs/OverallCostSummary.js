import * as React from 'react'
import PropTypes from 'prop-types'
import { Icon, Table, Typography, Alert, Row, Col } from 'antd'
import moment from 'moment'

import { formatCost } from '../../utils/cost-formatters'
import IconTooltip from '../utils/IconTooltip'

export default class OverallCostSummary extends React.Component {
  static propTypes = {
    summary: PropTypes.object,
    style: PropTypes.object,
    loading: PropTypes.bool,
    onTeamDetail: PropTypes.func
  }

  projectMonthlyCost = (from, to, cost) => {
    const costPerDay = cost / (to.diff(from, 'days') + 1)
    const daysInMonth = moment(to).endOf('month').diff(moment(from).startOf('month'), 'days') + 1
    return costPerDay * daysInMonth
  } 

  render() {
    const { summary, style, loading, onTeamDetail } = this.props
    if (loading) {
      return (
        <div style={style}>
          <Icon type="loading" />
        </div>
      )
    }

    if (!summary) {
      return null
    }

    if ((!summary.cost || summary.cost === 0) && (!summary.teamCosts || summary.teamCosts.length === 0)) {
      return (
        <div style={style}>
          <Alert type="warning" showIcon={true} message="No costs information exists in Kore for this time period. Ensure you have a costs provider set up to enable costs visibility." />
        </div>
      )
    }

    const fromMoment = moment(summary.usageStartTime).utc(false)
    const toMoment = moment(summary.usageEndTime).utc(false)
    const isSingleMonth = fromMoment.format('MMM') === toMoment.format('MMM')
    const isCurrentMonth = toMoment.format('MMM') === moment().utc(false).format('MMM')
    const toActual = isCurrentMonth ? moment().utc(false) : toMoment.utc(false)
    const month = isSingleMonth ? fromMoment.format('MMMM YYYY') : `${fromMoment.format('MMMM YYYY')} to ${toMoment.format('MMMM YYYY')}`

    return (
      <div style={style}>
        <Row gutter={16}>
          <Col span={12}>
            <Typography.Paragraph style={{ fontSize: '14px', marginBottom: 0 }} type="secondary">Accrued Costs for {month}</Typography.Paragraph>
            <Typography.Text strong style={{ fontSize: '50px' }}>{summary.cost ? formatCost(summary.cost) : '$0.00' }</Typography.Text>
          </Col>
          {!isCurrentMonth || !isSingleMonth ? null : (
            <Col span={12}>
              <Typography.Paragraph style={{ fontSize: '14px', marginBottom: 0 }} type="secondary">Projected Costs for {month} <IconTooltip icon="info-circle" text="This figure is projected from the usage so far this month, it could increase or decrease with usage changes." /></Typography.Paragraph>
              <Typography.Text strong style={{ fontSize: '50px' }}>{summary.cost ? formatCost(this.projectMonthlyCost(fromMoment, toActual, summary.cost)) : '$0.00' }</Typography.Text>

              {/* TODO: 
              <Statistic
                title="compared to May 2020"
                value={predictedCostChangePercent}
                precision={1}
                prefix={predictedCostChangePercent >= 0 ? <Icon type="arrow-up" /> : <Icon type="arrow-down" />}
                suffix="%"
                style={{ display: 'inline-block', marginLeft: '20px' }}
              /> 
              */}
            </Col>
          )}
        </Row>

        <Typography.Title level={3} style={{ marginTop: '20px' }}>Team breakdown</Typography.Title>
        <Table 
          style={{ marginTop: '20px' }} 
          size="small" 
          pagination={true} 
          rowSelection={{ type: 'radio', onSelect: (r) => onTeamDetail && onTeamDetail(r.teamIdentifier) }}
          columns={[
            { title: 'Team', dataIndex: 'teamName', width: '40%' },
            { title: 'Accrued Costs', dataIndex: 'cost', render: (v) => v ? formatCost(v) : '$0.00', width: '30%',
              sortDirections: ['descend', 'ascend'], 
              defaultSortOrder: 'descend',
              sorter: (a, b) => a.cost && b.cost ? a.cost - b.cost : a.cost ? 1 : -1
            },
            { 
              title: 'Projected Month Total', 
              key: 'costproj', 
              width: '30%',
              render: (_, row) => row.cost ? formatCost(this.projectMonthlyCost(fromMoment, toActual, row.cost)) : '$0.00' 
            }
          ]}
          rowKey="teamIdentifier"
          dataSource={summary.teamCosts}
        />
      </div>
    )
  }
}
