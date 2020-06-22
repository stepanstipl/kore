import React from 'react'
import Link from 'next/link'
import PropTypes from 'prop-types'
import { List, Icon, Tooltip, Table, Tag, Button } from 'antd'

import MonitoringAlert from './MonitoringAlert'
import MonitoringSummary from './MonitoringSummary'

export default class MonitoringTable extends React.Component {
  static propTypes = {
    alerts: PropTypes.object.isRequired,
    severity: PropTypes.array,
    status: PropTypes.array,
  }

  static columns = [
    {
        title: 'Severity',
        dataIndex: 'status.rule.spec.severity',
        key: 'severity',
        onFilter: (value, record) => record.status.rule.spec.severity(value) === 0,
        sorter: (a, b) => a.status.rule.spec.severity.length - b.status.rule.spec.severity.length,
        sortDirections: ['descend'],
        render: (text) => (
          <>
            <Tag key="middle" color={text == "Critical" ? "red" : "orange"}>
              {text}
            </Tag>
          </>
        ),
    },
    {
        title: 'Triggered',
        dataIndex: 'metadata.creationTimestamp',
        key: 'triggered',
    },
    {
        title: 'Rule Name',
        dataIndex: 'status.rule.metadata.name',
        key: 'name',
        render: (text) => (
          <Link
            key="view"
            passHref
            href="/docs"
          >
            <a>
              <Tooltip placement="top" title="View the definition of this rule">
                <Icon type="info-circle" />  {text}
              </Tooltip>
            </a>
          </Link>
        )
    },
    {
        title: 'Summary',
        dataIndex: 'spec.summary',
        key: 'summary',
        render: (text, record) => <MonitoringSummary record={record}/>
    },
    {
        title: 'State',
        render: (text, record) => (
          <>
            <Tooltip placement="top" title="The current state of the alert">
              <Tag color="green"><Icon type="info-circle" />  {record.status.status}</Tag>
            </Tooltip>
          </>
        )
    },
    /*
    {
        title: 'Actions',
        render: (text, record) => (
          <>
          {record.status.status == "Silenced"
            ? <Button type="primary" shape="round"> Unsilence</Button>
            : <Button type="primary" shape="round"> Silence</Button>
          }
          </>
        )
    },
    */
  ]

  filterOnRules = (alerts) => {
    var resources = []

    if (!alerts) {
      return resources
    }

    alerts.items.forEach(resource => {
        var state = resource.status

        if (this.props.severity) {
            if (!this.props.severity.includes(state.rule.spec.severity)){
                return
            }
        }

        if (this.props.status) {
            if (!this.props.status.includes(state.status)) {
                return
            }
        }

        resources.push(resource)
    })

    return resources
  }


  render() {
    const { alerts } = this.props
    const resources = this.filterOnRules(alerts)

    return (
      <>
        <Table
            dataSource={resources}
            columns={MonitoringTable.columns}
            //expandedRowRender={record => <MonitoringAlert rule={record}/> }
        />
      </>
    )
  }
}
