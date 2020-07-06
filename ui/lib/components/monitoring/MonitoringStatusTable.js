import React from 'react'
import Link from 'next/link'
import moment from 'moment'
import PropTypes from 'prop-types'
import { Form, Icon, Input, Modal, Table, Tag, Tooltip } from 'antd'

import MonitoringSummary from './MonitoringSummary'
import MonitoringTable from './MonitoringTable'
import KoreApi from '../../kore-api'
import FormErrorMessage from '../forms/FormErrorMessage'
import { errorMessage, successMessage } from '../../utils/message'

export default class MonitoringStatusTable extends React.Component {
  static propTypes = {
    alerts: PropTypes.object.isRequired,
    refreshData: PropTypes.func,
    severity: PropTypes.array,
    status: PropTypes.array,
  }

  static summaries = {
    'Application Availability': 'Relates to application deployments and availablity',
    'Authentication Proxy': 'Related to Single Sign on to clusters',
    'Infrastructure': 'Relates to alert on the cluster infrastructure',
    'Monitoring': 'Relates to components used to monitor the infrastructure',
  }

  static severity = {
    'Warning': 'orange',
    'Critical': 'red',
    'Silenced': 'blue',
    'OK': 'green',
  }

  columns = [
    {
      title: 'Category',
      dataIndex: 'category',
      key: 'category',
    },
    {
      title: 'Summary',
      dataIndex: 'category',
      key: 'category',
      render: (text, record) => (
        <>
          <span>{MonitoringStatusTable.summaries[text]}</span>
        </>
      )
    },
    {
      title: 'Status',
      width: 100,
      render: (text, record) => (
        <>
          <Tooltip placement="top" title="The current state of the alert">
            <Tag color={MonitoringStatusTable.severity[record.status]}><Icon type="info-circle" style={{ marginRight: '5px' }}/>{record.status}</Tag>
          </Tooltip>
        </>
      )
    },
  ]

  filterOnRules = () => {
    if (!this.props.alerts) {
      return []
    }
    var filtered = []
    var matches = new Map()

    this.props.alerts.items.map((record, index) => {
      if ((this.props.severity || []).includes(record.status.rule.spec.severity)) {

        if ((this.props.status || []).includes(record.status.status)) {

          var category = record.status.rule.metadata.labels['category']
          if (category !== null) {
            if (matches[category] === undefined) {
              matches.set(category, {
                'alerts': [record],
                'count': 1,
                'status': record.status.rule.spec.severity,
              })
            } else {
              matches[category]['count'] += 1
              matches[category]['alerts'].push(record)
              if (record.status.rule.spec.severity === 'Critical') {
                matches[category]['status'] = record.status.rule.spec.severity
              }
            }
          }
        }
      }
    })
    if (matches.length <= 0) {
      return filtered
    }

    matches.forEach((value, key) => {
      value['category'] = key
      filtered.push(value)
    })

    return filtered
  }

  render() {
    return (
      <>
        <Table
          rowKey={(a) => a.category}
          dataSource={this.filterOnRules()}
          columns={this.columns}
          expandedRowRender={a => <MonitoringTable alerts={{'items': a.alerts}} severity={['Critical', 'Warning']} status={['Active', 'Silenced']}/> }
        />
      </>
    )
  }
}

