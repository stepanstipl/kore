import React from 'react'
import axios from 'axios'
import PropTypes from 'prop-types'

import AuditViewer from '../../../lib/components/common/AuditViewer'
import TeamHeader from '../../../lib/components/teams/TeamHeader'
import KoreApi from '../../../lib/kore-api'

class TeamAuditPage extends React.Component {
  static propTypes = {
    team: PropTypes.object.isRequired,
    events: PropTypes.array.isRequired,
    teamRemoved: PropTypes.func.isRequired
  }

  state = {
    events: []
  }

  static staticProps = {
    title: 'Team Audit Viewer',
    adminOnly: false
  }

  static async getPageData(ctx) {
    const name = ctx.query.name
    const api = await KoreApi.client(ctx)

    return axios.all([api.GetTeam(name), api.ListTeamAudit(name)])
      .then(axios.spread(function (team, eventList) {
        return { team, events: eventList.items }
      }))
      .catch(err => {
        throw new Error(err.message)
      })
  }

  static getInitialProps = async ctx => {
    const data = await TeamAuditPage.getPageData(ctx)
    return data
  }

  constructor(props) {
    super(props)
    this.state = { events: props.events }
  }

  render() {
    const { team, teamRemoved } = this.props

    return (
      <>
        <TeamHeader team={team} breadcrumbExt={[{ text: 'Team Audit Viewer' }]} teamRemoved={teamRemoved} />
        <AuditViewer items={this.state.events} />
      </>
    )
  }
}

export default TeamAuditPage
