import React from 'react'
import PropTypes from 'prop-types'
import { Table, Tag } from 'antd'

export default class MonitoringRulesTable extends React.Component {
  static propTypes = {
    rules: PropTypes.object.isRequired,
  }

  static columns = [
    {
      title: 'Severity',
      dataIndex: 'spec.severity',
      key: 'severity',
      render: (text) => (
        <>
          <Tag key='middle' color={text === 'Critical' ? 'red' : 'orange'}>
            {text}
          </Tag>
        </>
      ),
    },
    {
      title: 'Rule',
      dataIndex: 'metadata.name',
      key: 'name',
    },
    {
      title: 'Summary',
      dataIndex: 'spec.summary',
      key: 'summary',
    },
    {
      title: 'Team/Resource',
      dataIndex: 'spec.resource.kind',
      key: 'resource.kind',
    },
  ]

  render() {
    const { rules } = this.props

    return (
      <Table
        dataSource={rules.items}
        columns={MonitoringRulesTable.columns}
      />
    )
  }
}
