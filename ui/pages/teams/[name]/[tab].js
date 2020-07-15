import React from 'react'
import PropTypes from 'prop-types'
import Router from 'next/router'
import Error from 'next/error'
import { Alert, Icon, Tabs } from 'antd'
const { TabPane } = Tabs

import KoreApi from '../../../lib/kore-api'
import ClustersTab from '../../../lib/components/teams/cluster/ClustersTab'
import MembersTab from '../../../lib/components/teams/members/MembersTab'
import SecurityTab from '../../../lib/components/teams/security/SecurityTab'
import MonitoringTab from '../../../lib/components/teams/monitoring/MonitoringTab'
import SecurityStatusIcon from '../../../lib/components/security/SecurityStatusIcon'
import TeamHeader from '../../../lib/components/teams/TeamHeader'
import TextWithCount from '../../../lib/components/utils/TextWithCount'

class TeamDashboardTabPage extends React.Component {
  static propTypes = {
    invitation: PropTypes.bool,
    team: PropTypes.object.isRequired,
    user: PropTypes.object.isRequired,
    teamRemoved: PropTypes.func.isRequired,
    tabActiveKey: PropTypes.string
  }

  static staticProps = {
    title: 'Team dashboard'
  }

  constructor(props) {
    super(props)
    this.state = {
      tabActiveKey: this.props.tabActiveKey,
      memberCount: -1,
      clusterCount: -1,
      securityStatus: false
    }
  }

  static async getTeamDetails(ctx) {
    try {
      const team = await (await KoreApi.client(ctx)).GetTeam(ctx.query.name)
      return { team }
    } catch (err) {
      throw new Error(err.message)
    }
  }

  static getInitialProps = async (ctx) => {
    const teamDetails = await TeamDashboardTabPage.getTeamDetails(ctx)
    if (!teamDetails.team && ctx.res) {
      /* eslint-disable-next-line require-atomic-updates */
      ctx.res.statusCode = 404
    }
    if (ctx.query.invitation === 'true') {
      teamDetails.invitation = true
    }
    teamDetails.tabActiveKey = ctx.query.tab || 'clusters'
    return teamDetails
  }

  componentDidUpdate(prevProps) {
    const prevTeamName = prevProps.team && prevProps.team.metadata && prevProps.team.metadata.name
    if (this.props.team && this.props.team.metadata.name !== prevTeamName) {
      this.setState({ tabActiveKey: 'clusters' })
    }
    if (this.state.tabActiveKey !== this.props.tabActiveKey) {
      this.setState({ tabActiveKey: this.props.tabActiveKey })
    }
  }

  handleTabChange = (key) => {
    Router.push('/teams/[name]/[tab]', `/teams/${this.props.team.metadata.name}/${key}`)
  }

  render() {
    const { team, invitation, teamRemoved } = this.props

    if (!team) {
      return <Error statusCode={404} />
    }

    return (
      <>
        <TeamHeader team={team} teamRemoved={teamRemoved} />

        {invitation ? (
          <Alert
            message="You have joined this team from an invitation"
            type="info"
            showIcon
            style={{ marginBottom: '10px' }}
          />
        ) : null}

        <Tabs activeKey={this.state.tabActiveKey} onChange={(key) => this.handleTabChange(key)} tabBarStyle={{ marginBottom: '20px' }}>
          <TabPane key="clusters" tab={<TextWithCount title="Clusters" count={this.state.clusterCount} />} forceRender={true}>
            <ClustersTab user={this.props.user} team={this.props.team} getClusterCount={(count) => this.setState({ clusterCount: count })} />
          </TabPane>

          <TabPane key="members" tab={<TextWithCount title="Members" count={this.state.memberCount} />} forceRender={true}>
            <MembersTab user={this.props.user} team={this.props.team} getMemberCount={(count) => this.setState({ memberCount: count })} />
          </TabPane>

          <TabPane key="security" tab={<TextWithCount title="Security" icon={<SecurityStatusIcon status="Compliant" size="small" style={{ verticalAlign: 'middle' }} />} />} forceRender={true}>
            <SecurityTab team={this.props.team} getOverviewStatus={(status) => this.setState({ securityStatus: status })} />
          </TabPane>

          <TabPane key="monitoring" tab={<TextWithCount title="Monitoring" icon={<Icon type="monitor"/>} />} forceRender={true}>
            <MonitoringTab user={this.props.user} team={this.props.team} />
          </TabPane>
        </Tabs>

      </>
    )
  }
}

export default TeamDashboardTabPage
