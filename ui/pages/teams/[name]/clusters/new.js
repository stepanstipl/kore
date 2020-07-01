import React from 'react'
import PropTypes from 'prop-types'

import TeamHeader from '../../../../lib/components/teams/TeamHeader'
import ClusterBuildForm from '../../../../lib/components/teams/cluster/ClusterBuildForm'
import KoreApi from '../../../../lib/kore-api'

class NewTeamClusterPage extends React.Component {
  static propTypes = {
    user: PropTypes.object.isRequired,
    team: PropTypes.object.isRequired,
    clusters: PropTypes.object.isRequired,
    teamRemoved: PropTypes.func.isRequired
  }

  static staticProps = {
    title: 'New team cluster'
  }

  static async getPageData({ query }) {
    const teamName = query.name
    try {
      const api = await KoreApi.client()
      const [ team, clusters ] = await Promise.all([
        api.GetTeam(teamName),
        api.ListClusters(teamName)
      ])
      return { team, clusters }
    } catch (err) {
      throw new Error(err.message)
    }
  }

  static getInitialProps = async (ctx) => {
    const data = await NewTeamClusterPage.getPageData(ctx)
    return data
  }

  render() {
    const { user, team, teamRemoved, clusters } = this.props

    return (
      <>
        <TeamHeader team={team} breadcrumbExt={[
          { text: 'Clusters', href: '/teams/[name]/[tab]', link: `/teams/${team.metadata.name}/clusters` },
          { text: 'New cluster' }
        ]} teamRemoved={teamRemoved} />

        <ClusterBuildForm
          user={user}
          team={team}
          teamClusters={clusters.items}
          skipButtonText="Cancel"
        />
      </>
    )
  }
}

export default NewTeamClusterPage
