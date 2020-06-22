import React from 'react'
import Link from 'next/link'
import PropTypes from 'prop-types'
import { Icon, Tooltip, Table, Tag } from 'antd'

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
      render: (text) => (
        <Link
          key='view'
          passHref
          href='/docs'
        >
          <a><Tooltip placement='left' title='View the definition of this rule'>
            <Icon type='info-circle' />  {text}
          </Tooltip></a>
        </Link>
      )
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
      render: (text, record) => (
        <Link
          key="view_cluster"
          href='/teams/{text}/[record.spec.resource.kind]/namespaces'
        >
          <a>
            <Tooltip placement="left" title="View the resourc">
              {record.spec.resource.namespace}/{record.spec.resource.kind}
            </Tooltip>
          </a>
        </Link>
      )
    },
  ]

  render() {
    const { rules } = this.props

    return (
      <>
        <Table
          dataSource={rules.items}
          columns={MonitoringRulesTable.columns}
        />
      </>
    )
  }
}
