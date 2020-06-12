import React from 'react'
import PropTypes from 'prop-types'
import Link from 'next/link'
import Router from 'next/router'
import Error from 'next/error'
import { Typography, Button, Badge, Alert, Icon, Modal, Dropdown, Menu, Tabs } from 'antd'
const { Paragraph, Text } = Typography
const { TabPane } = Tabs

import Breadcrumb from '../../../lib/components/layout/Breadcrumb'
import redirect from '../../../lib/utils/redirect'
import KoreApi from '../../../lib/kore-api'
import SecurityStatusIcon from '../../../lib/components/security/SecurityStatusIcon'
import { successMessage, errorMessage } from '../../../lib/utils/message'

// prototype imports
import TeamData from '../../../lib/prototype/utils/dummy-team-data'
import ClustersTab from '../../../lib/prototype/components/teams/cluster/ClustersTab'

class CostsDemoTeamDashboardPage extends React.Component {
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
      serviceCount: -1,
      securityStatus: false
    }
  }

  static async getTeamDetails() {
    try {
      const team = await Promise.resolve(TeamData.team)
      return { team }
    } catch (err) {
      throw new Error(err.message)
    }
  }

  static getInitialProps = async (ctx) => {
    const teamDetails = await CostsDemoTeamDashboardPage.getTeamDetails()
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

  deleteTeam = async () => {
    try {
      const team = this.props.team.metadata.name
      await (await KoreApi.client()).RemoveTeam(team)
      this.props.teamRemoved(team)
      successMessage(`Team "${team}" deleted`)
      return redirect({ router: Router, path: '/' })
    } catch (err) {
      console.log('Error deleting team', err)
      errorMessage('Team could not be deleted, please try again later')
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
        <Menu.Item key="costs">
          <Link href="/prototype/teams/demo/costs">
            <a>
              <Icon type="pound" style={{ marginRight: '5px' }} />
              Team costs
            </a>
          </Link>
        </Menu.Item>
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

    if (!team) {
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

        <Tabs activeKey={this.state.tabActiveKey} onChange={(key) => this.handleTabChange(key)} tabBarStyle={{ marginBottom: '20px' }}>
          <TabPane key="clusters" tab={this.getTabTitle({ title: 'Clusters', count: this.state.clusterCount })} forceRender={true}>
            <ClustersTab user={this.props.user} team={this.props.team} getClusterCount={(count) => this.setState({ clusterCount: count })} />
          </TabPane>

          <TabPane key="members" tab={this.getTabTitle({ title: 'Members', count: this.state.memberCount })} forceRender={true}>
            {/*<MembersTab user={this.props.user} team={this.props.team} getMemberCount={(count) => this.setState({ memberCount: count })} />*/}
          </TabPane>

          <TabPane key="security" tab={this.getTabTitle({ title: 'Security', icon: <SecurityStatusIcon status="Compliant" size="small" style={{ verticalAlign: 'middle' }} /> })} forceRender={true}>
            {/*<SecurityTab team={this.props.team} getOverviewStatus={(status) => this.setState({ securityStatus: status })} />*/}
          </TabPane>
        </Tabs>

      </div>
    )
  }
}

export default CostsDemoTeamDashboardPage
