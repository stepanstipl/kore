import React from 'react'
import PropTypes from 'prop-types'

import Breadcrumb from '../../../lib/components/layout/Breadcrumb'
import SecurityOverview from '../../../lib/components/security/SecurityOverview'
import KoreApi from '../../../lib/kore-api'

export default class TeamSecurityPage extends React.Component {
  static propTypes = {
    overview: PropTypes.object.isRequired,
    team: PropTypes.object.isRequired
  }

  static staticProps = {
    title: 'Team Security Overview',
  }

  static async getPageData(ctx) {
    const teamName = ctx.query.name
    try {
      const api = await KoreApi.client(ctx)
      const [ overview, team ] = await Promise.all([
        api.GetTeamSecurityOverview(teamName), 
        api.GetTeam(teamName)
      ])
      return { overview, team }
    } catch (err) {
      throw new Error(err.message)
    }
  }

  static getInitialProps = async ctx => {
    const data = await TeamSecurityPage.getPageData(ctx)
    return data
  }

  constructor(props) {
    super(props)
  }

  render() {
    const { overview, team } = this.props

    return (
      <div>
        <Breadcrumb
          items={[
            { text: team.spec.summary, href: '/teams/[name]', link: `/teams/${team.metadata.name}` },
            { text: 'Team Security Overview' }
          ]}
        />
        <SecurityOverview overview={overview} />
      </div>
    )
  }
}
