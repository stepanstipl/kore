import React from 'react'
import { Typography } from 'antd'
const { Paragraph, Title } = Typography

import Breadcrumb from '../../../../../lib/components/layout/Breadcrumb'
import TeamMonthlyCostTable from '../../../../../lib/prototype/components/costs/TeamMonthlyCostTable'
import CurrentCosts from '../../../../../lib/prototype/components/costs/CurrentCosts'

class TeamCosts extends React.Component {

  render() {
    return (
      <>
        <Breadcrumb items={[{ text: 'Proto', link: '/prototype/teams/proto', href: '/prototype/teams/proto' }, { text: 'Team costs' }]}/>
        <CurrentCosts team="proto" currentCost={254.62} predictedCost={734.72} predictedCostChangePercent={5.5} />
        <Title level={3}>Costs breakdown for June 2020</Title>
        <Paragraph style={{ marginBottom: '20px' }} type="secondary">This shows the breakdown of the incurred costs for this month, the predicted total expenditure is not broken down here.</Paragraph>
        <TeamMonthlyCostTable />
      </>
    )
  }
}

export default TeamCosts
