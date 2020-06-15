import React from 'react'
import PropTypes from 'prop-types'
import { Typography } from 'antd'
const { Title } = Typography

import Breadcrumb from '../../../../lib/components/layout/Breadcrumb'
import ClusterBuildForm from '../../../../lib/components/teams/cluster/ClusterBuildForm'
import KoreApi from '../../../../lib/kore-api'

class NewTeamClusterPage extends React.Component {
  static propTypes = {
    user: PropTypes.object.isRequired,
    team: PropTypes.object.isRequired,
    clusters: PropTypes.object.isRequired
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
    const teamName = this.props.team.metadata.name
    const teamClusters = this.props.clusters.items

    return (
      <div>
        <Breadcrumb
          items={[
            { text: this.props.team.spec.summary, href: '/teams/[name]', link: `/teams/${teamName}` },
            { text: 'Clusters', href: '/teams/[name]/[tab]', link: `/teams/${teamName}/clusters` },
            { text: 'New cluster' }
          ]}
        />
        <Title style={{ marginBottom: '40px' }}>New Cluster for {this.props.team.spec.summary}</Title>
        <ClusterBuildForm
          user={this.props.user}
          team={this.props.team}
          teamClusters={teamClusters}
          skipButtonText="Cancel"
        />
      </div>
    )
  }
}

export default NewTeamClusterPage
