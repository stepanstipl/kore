import React from 'react'
import PropTypes from 'prop-types'
import { Table } from 'antd'

class AuditViewer extends React.Component {
  static propTypes = {
    items: PropTypes.arrayOf(PropTypes.object).isRequired
  }

  static columns = [
    {
      title: 'Time',
      dataIndex: 'spec.createdAt',
      defaultSortOrder: 'descend',
      sortDirections: ['descend','ascend'],
      sorter: (a, b) => { return a.spec.createdAt.localeCompare(b.spec.createdAt)},
    },
    {
      title: 'Resource',
      dataIndex: 'spec.resource'
    },
    {
      title: 'URI',
      dataIndex: 'spec.resourceURI'
    },
    {
      title: 'Operation',
      dataIndex: 'spec.operation'
    },
    {
      title: 'User',
      dataIndex: 'spec.user'
    },
    {
      title: 'Result',
      dataIndex: 'spec.responseCode'
    }
  ];  

  render() {
    return (
      <Table dataSource={this.props.items} columns={AuditViewer.columns} rowKey={r => r.spec.id} />
    )
  }
}

export default AuditViewer