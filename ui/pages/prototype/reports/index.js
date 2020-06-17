import React from 'react'
import Link from 'next/link'
import { Button, List, Typography } from 'antd'
const { Text, Title } = Typography

import Breadcrumb from '../../../lib/components/layout/Breadcrumb'

class ReportsIndex extends React.Component {
  static REPORT_LIST = [{
    title: 'Security',
    description: 'View the security status for cluster plans and all team clusters, across the organisation',
    path: '/prototype/reports/security'
  }, {
    title: 'Costs',
    description: 'View the organisation costs, broken down by team',
    path: '/prototype/reports/costs'
  }]

  render() {
    return (
      <>
        <Breadcrumb items={[ { text: 'Reports' } ]} />
        <Title level={2}>Kore reports</Title>
        <List
          dataSource={ReportsIndex.REPORT_LIST}
          renderItem={report => (
            <List.Item actions={[<Link key="view" href={report.path}><Button type="primary">View</Button></Link>]}>
              <List.Item.Meta
                title={<Text style={{ fontSize: '20px', fontWeight: '600' }}>{report.title}</Text>}
                description={report.description}
              />
            </List.Item>
          )}
        />
      </>
    )
  }
}

export default ReportsIndex
