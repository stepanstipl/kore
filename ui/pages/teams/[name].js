import React from 'react'
import PropTypes from 'prop-types'
import Link from 'next/link'
import Router from 'next/router'
import Error from 'next/error'
import { Typography, Button, message, Badge, Alert, Icon, Modal, Dropdown, Menu, Tabs } from 'antd'
const { Paragraph, Text } = Typography
const { TabPane } = Tabs
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

import Breadcrumb from '../../lib/components/layout/Breadcrumb'
import redirect from '../../lib/utils/redirect'
import KoreApi from '../../lib/kore-api'
import ClustersTab from '../../lib/components/teams/cluster/ClustersTab'
import ServicesTab from '../../lib/components/teams/service/ServicesTab'
import MembersTab from '../../lib/components/teams/members/MembersTab'
import SecurityTab from '../../lib/components/teams/security/SecurityTab'
import SecurityStatusIcon from '../../lib/components/security/SecurityStatusIcon'

class TeamDashboard extends React.Component {
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
      tabActiveKey: this.props.tabActiveKey || 'clusters',
      memberCount: -1,
      clusterCount: -1,
      serviceCount: -1,
      securityStatus: false
    }
  }

  static async getTeamDetails(ctx) {
    try {
      const name = ctx.query.name
      const api = await KoreApi.client(ctx)
      const team = await api.GetTeam(name)
      return { team }
    } catch (err) {
      throw new Error(err.message)
    }
  }

  static getInitialProps = async ctx => {
    const teamDetails = await TeamDashboard.getTeamDetails(ctx)
    if (Object.keys(teamDetails.team).length === 0 && ctx.res) {
      /* eslint-disable-next-line require-atomic-updates */
      ctx.res.statusCode = 404
    }
    if (ctx.query.invitation === 'true') {
      teamDetails.invitation = true
    }
    if (ctx.query.tab) {
      teamDetails.tabActiveKey = ctx.query.tab
    }
    return teamDetails
  }

  componentDidUpdate(prevProps) {
    const teamFound = Object.keys(this.props.team).length
    const prevTeamName = prevProps.team.metadata && prevProps.team.metadata.name
    if (teamFound && this.props.team.metadata.name !== prevTeamName) {
      this.setState({ tabActiveKey: 'clusters' })
    }
  }

  deleteTeam = async () => {
    try {
      const team = this.props.team.metadata.name
      await (await KoreApi.client()).RemoveTeam(team)
      this.props.teamRemoved(team)
      message.success(`Team "${team}" deleted`)
      return redirect({ router: Router, path: '/' })
    } catch (err) {
      console.log('Error deleting team', err)
      message.error('Team could not be deleted, please try again later')
    }
  }

  deleteTeamConfirm = () => {
    const { clusterCount } = this.state
    if (clusterCount > 0) {
      return Modal.warning({
        title: 'Warning: team cannot be deleted',
        content: 'The clusters must be deleted first',
        onOk() {}
      })
    }

    Modal.confirm({
      title: 'Are you sure you want to delete this team?',
      content: 'This cannot be undone',
      okText: 'Yes',
      okType: 'danger',
      cancelText: 'No',
      onOk: this.deleteTeam
    })
  }

  settingsMenu = ({ team }) => {
    const menu = (
      <Menu>
        <Menu.Item key="audit">
          <Link href="/teams/[name]/audit" as={`/teams/${team.metadata.name}/audit`}>
            <a>
              <Icon type="table" style={{ marginRight: '5px' }} />
              Team audit viewer
            </a>
          </Link>
        </Menu.Item>
        <Menu.Item key="delete" className="ant-btn-danger" onClick={this.deleteTeamConfirm}>
          <Icon type="delete" style={{ marginRight: '5px' }} />
          Delete team
        </Menu.Item>
      </Menu>
    )
    return (
      <Dropdown overlay={menu}>
        <Button>
          <Icon type="setting" style={{ marginRight: '10px' }} />
          <Icon type="down" />
        </Button>
      </Dropdown>
    )
  }

  getTabTitle = ({ title, count, icon }) => (
    <span>
      {title}
      {count !== undefined && count !== -1 && <Badge showZero={true} style={{ marginLeft: '10px', backgroundColor: '#1890ff' }} count={count} />}
      {icon}
    </span>
  )

  render() {
    const { team, invitation } = this.props

    if (Object.keys(team).length === 0) {
      return <Error statusCode={404} />
    }

    return (
      <div>
        <div style={{ display: 'inline-block', width: '100%' }}>
          <div style={{ float: 'left', marginTop: '8px' }}>
            <Breadcrumb items={[{ text: team.spec.summary }]} />
          </div>
          <div style={{ float: 'right' }}>
            {this.settingsMenu({ team })}
          </div>
        </div>
        <Paragraph>
          {team.spec.description ? <Text strong>{team.spec.description}</Text> : <Text style={{ fontStyle: 'italic' }} type="secondary">No description</Text> }
          <Text style={{ float: 'right' }}><Text strong>Team ID: </Text>{team.metadata.name}</Text>
        </Paragraph>

        {invitation ? (
          <Alert
            message="You have joined this team from an invitation"
            type="info"
            showIcon
            style={{ marginBottom: '10px' }}
          />
        ) : null}

        <Tabs activeKey={this.state.tabActiveKey} onChange={(key) => this.setState({ tabActiveKey: key })} tabBarStyle={{ marginBottom: '20px' }}>
          <TabPane key="clusters" tab={this.getTabTitle({ title: 'Clusters', count: this.state.clusterCount })} forceRender={true}>
            <ClustersTab user={this.props.user} team={this.props.team} getClusterCount={(count) => this.setState({ clusterCount: count })} />
          </TabPane>

          {!publicRuntimeConfig.featureGates['services'] ? null : (
            <TabPane key="services" tab={this.getTabTitle({ title: 'Cloud services', count: this.state.serviceCount })} forceRender={true}>
              <ServicesTab user={this.props.user} team={this.props.team} getServiceCount={(count) => this.setState({ serviceCount: count })} />
            </TabPane>
          )}

          <TabPane key="members" tab={this.getTabTitle({ title: 'Members', count: this.state.memberCount })} forceRender={true}>
            <MembersTab user={this.props.user} team={this.props.team} getMemberCount={(count) => this.setState({ memberCount: count })} />
          </TabPane>

          <TabPane key="security" tab={this.getTabTitle({ title: 'Security', icon: <SecurityStatusIcon status="Compliant" size="small" style={{ verticalAlign: 'middle' }} /> })} forceRender={true}>
            <SecurityTab team={this.props.team} getOverviewStatus={(status) => this.setState({ securityStatus: status })} />
          </TabPane>
        </Tabs>

      </div>
    )
  }
}

export default TeamDashboard
