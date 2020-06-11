import React from 'react'
import Link from 'next/link'
import { Collapse, Typography } from 'antd'
const { Paragraph, Text, Title } = Typography

import Breadcrumb from '../../../../../lib/components/layout/Breadcrumb'
import MonthlyCostTables from '../../../../../lib/prototype/components/costs/MonthlyCostTables'

class TeamCostsHistory extends React.Component {

  render() {
    return (
      <>
        <Breadcrumb items={[{ text: 'Demo' }, { text: 'Team costs history' }]}/>

        <Title level={3}>Historical team costs</Title>
        <Paragraph>
          <Link href="/prototype/teams/demo/costs">
            <a style={{ fontSize: '14px', textDecoration: 'underline' }}>See current cost</a>
          </Link>
        </Paragraph>

        <Collapse bordered={false}>
          <Collapse.Panel className="enlarged-header" header="May 2020" extra={<Text>£765.43</Text>}>
            <MonthlyCostTables />
          </Collapse.Panel>
          <Collapse.Panel className="enlarged-header" header="April 2020" extra={<Text>£734.14</Text>}>
            <MonthlyCostTables />
          </Collapse.Panel>
          <Collapse.Panel className="enlarged-header" header="March 2020" extra={<Text>£695.97</Text>}>
            <MonthlyCostTables />
          </Collapse.Panel>
        </Collapse>
      </>
    )
  }
}

export default TeamCostsHistory
