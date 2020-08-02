import * as React from 'react'
import PropTypes from 'prop-types'
import { Icon, Table, Typography, Alert, Row, Col, Checkbox } from 'antd'
import moment from 'moment'

import { formatCost } from '../../utils/cost-formatters'
import IconTooltip from '../utils/IconTooltip'
import { getCloudInfo } from '../../utils/cloud'

export default class TeamCostSummary extends React.Component {
  static propTypes = {
    summary: PropTypes.object,
    style: PropTypes.object,
    loading: PropTypes.bool
  }

  state = {
    selectedAsset: null,
    showZeroCostLineItems: false
  }

  projectMonthlyCost = (from, to, cost) => {
    const costPerDay = cost / (to.diff(from, 'days') + 1)
    const daysInMonth = moment(to).endOf('month').diff(moment(from).startOf('month'), 'days') + 1
    return costPerDay * daysInMonth
  } 

  componentDidUpdate = (prevProps) => {
    if (this.props.summary !== prevProps.summary) {
      this.setState({ selectedAsset: null })
    }
  }

  render() {
    const { summary, style, loading } = this.props
    const { selectedAsset, showZeroCostLineItems } = this.state

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

    if ((!summary.cost || summary.cost === 0) && (!summary.assetCosts || summary.assetCosts.length === 0)) {
      return (
        <div style={style}>
          <Alert type="warning" showIcon={true} message="No costs information exists in Kore for this team and time period. Ensure you have a costs provider set up to enable costs visibility." />
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
            <Typography.Paragraph style={{ fontSize: '14px', marginBottom: 0 }} type="secondary">Team Costs for {month}</Typography.Paragraph>
            <Typography.Text strong style={{ fontSize: '50px' }}>{summary.cost ? formatCost(summary.cost) : '$0.00' }</Typography.Text>
          </Col>
          {!isCurrentMonth || !isSingleMonth ? null : (
            <Col span={12}>
              <Typography.Paragraph style={{ fontSize: '14px', marginBottom: 0 }} type="secondary">Projected Team Costs for {month} <IconTooltip icon="info-circle" text="This figure is projected from the usage so far this month, it could increase or decrease with usage changes." /></Typography.Paragraph>
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

        <Typography.Title level={3}>Asset breakdown</Typography.Title>
        <Table 
          style={{ marginBottom: '20px' }} 
          size="small" 
          pagination={false} 
          rowSelection={{ type: 'radio', onSelect: (r) => this.setState({ selectedAsset: r }) }}
          columns={[
            { title: 'Asset', dataIndex: 'assetName', width: '30%' },
            { title: 'Type', dataIndex: 'assetType', width: '15%' },
            { title: 'Provider', dataIndex: 'provider', width: '15%', render: (v) => v && v.length > 0 ? getCloudInfo(v).cloud : 'n/a' },
            { title: 'Cost', dataIndex: 'cost', render: (v) => v ? formatCost(v) : '$0.00', width: '20%' },
            { 
              title: 'Projected', 
              key: 'costproj', 
              width: '20%',
              render: (_, row) => row.cost ? formatCost(this.projectMonthlyCost(fromMoment, toActual, row.cost)) : '$0.00' 
            }
          ]}
          rowKey="assetIdentifier"
          dataSource={summary.assetCosts}
        />

        {!selectedAsset ? (
          <Alert type="info" showIcon={true} message="Select an asset to inspect individual cost line items" />
        ) : (
          <>
            <Typography.Title level={3}>Asset detail - {selectedAsset.assetName}</Typography.Title>
            
            <Table 
              style={{ marginBottom: '20px' }} 
              size="small" 
              pagination={true} 
              columns={[
                { title: 'Cost Type', dataIndex: 'description' },
                { title: 'From', dataIndex: 'usageStartTime', render: (v) => moment(v).utc(false).format('DD/MM/YY HH:mm') },
                { title: 'To', dataIndex: 'usageEndTime', render: (v) => moment(v).utc(false).format('DD/MM/YY HH:mm') },
                { title: 'Amount', dataIndex: 'usageAmount', render: (v, r) => `${v} ${r.usageUnit}`, 
                  sortDirections: ['descend', 'ascend'], 
                  sorter: (a, b) => a.usageAmount && b.usageAmount ? a.usageAmount - b.usageAmount : a.usageAmount ? 1 : -1
                },
                { title: 'Cost', dataIndex: 'cost', render: (v) => v ? formatCost(v) : '$0.00', width: '20%', 
                  sortDirections: ['descend', 'ascend'], 
                  defaultSortOrder: 'descend',
                  sorter: (a, b) => a.cost && b.cost ? a.cost - b.cost : a.cost ? 1 : -1
                },
              ]}
              dataSource={selectedAsset.details.filter((detail) => showZeroCostLineItems || detail.cost > 0)}
              footer={() => <><Checkbox onChange={(e) => this.setState({ showZeroCostLineItems: e.target.checked })} checked={showZeroCostLineItems} />&nbsp;&nbsp;Show zero-cost line items</>}
            />
          </>          
        )}
      </div>
    )
  }
}
