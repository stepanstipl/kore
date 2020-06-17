import React from 'react'
import { Collapse, Typography } from 'antd'
const { Text, Title } = Typography

import Breadcrumb from '../../../lib/components/layout/Breadcrumb'

// prototype imports
import TeamMonthlyCostTable from '../../../lib/prototype/components/costs/TeamMonthlyCostTable'
import CurrentCosts from '../../../lib/prototype/components/costs/CurrentCosts'

class CostReport extends React.Component {

  static staticProps = {
    title: 'Costs reports',
    hideSider: true,
    adminOnly: true
  }

  render() {
    return (
      <>
        <Breadcrumb items={[ { text: 'Reports', href: '/prototype/reports', link: '/prototype/reports' }, { text: 'Costs report' } ]} />
        <Title level={2}>Costs report</Title>
        <CurrentCosts currentCost={1254.62} predictedCost={1734.72} predictedCostChangePercent={3.1} />
        <Collapse bordered={false} defaultActiveKey={['current']}>
          <Collapse.Panel key="current" className="enlarged-header" header="June 2020 (current)" extra={<Text>£1254.62</Text>}>
            <Collapse>
              <Collapse.Panel className="enlarged-header" header="Team A" extra={<Text>£554.43</Text>}>
                <TeamMonthlyCostTable />
              </Collapse.Panel>
              <Collapse.Panel className="enlarged-header" header="Team B" extra={<Text>£700.19</Text>}>
                <TeamMonthlyCostTable />
              </Collapse.Panel>
            </Collapse>
          </Collapse.Panel>
          <Collapse.Panel className="enlarged-header" header="May 2020" extra={<Text>£765.43</Text>}>
          </Collapse.Panel>
          <Collapse.Panel className="enlarged-header" header="April 2020" extra={<Text>£734.14</Text>}>
          </Collapse.Panel>
          <Collapse.Panel className="enlarged-header" header="March 2020" extra={<Text>£695.97</Text>}>
          </Collapse.Panel>
        </Collapse>
      </>
    )
  }
}

export default CostReport
