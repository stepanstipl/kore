import React from 'react'
import PropTypes from 'prop-types'
import axios from 'axios'
import Link from 'next/link'
import Router from 'next/router'
import Error from 'next/error'
import { Typography, Card, List, Tag, Button, Avatar, Popconfirm, Select, Drawer, Badge, Alert, Icon, Modal, Dropdown, Menu, Tabs } from 'antd'
const { Paragraph, Text } = Typography
const { Option } = Select
const { TabPane } = Tabs
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

import Breadcrumb from '../../../lib/components/layout/Breadcrumb'
import InviteLink from '../../../lib/components/teams/InviteLink'
import NamespaceClaimForm from '../../../lib/components/teams/namespace/NamespaceClaimForm'
import ResourceStatusTag from '../../../lib/components/resources/ResourceStatusTag'
import apiRequest from '../../../lib/utils/api-request'
import copy from '../../../lib/utils/object-copy'
import asyncForEach from '../../../lib/utils/async-foreach'
import apiPaths from '../../../lib/utils/api-paths'
import redirect from '../../../lib/utils/redirect'
import KoreApi from '../../../lib/kore-api'

// prototype imports
import TeamData from '../../../lib/prototype/utils/dummy-team-data'
import Cluster from '../../../lib/prototype/components/teams/cluster/Cluster'
import NamespaceClaim from '../../../lib/prototype/components/teams/namespace/NamespaceClaim'
import ServiceForm from '../../../lib/prototype/components/teams/service/ServiceForm'
import { successMessage, errorMessage, loadingMessage } from '../../../lib/utils/message'

class TeamDashboard extends React.Component {
  static propTypes = {
    invitation: PropTypes.bool,
    team: PropTypes.object.isRequired,
    members: PropTypes.object.isRequired,
    user: PropTypes.object.isRequired,
    clusters: PropTypes.object.isRequired,
    namespaceClaims: PropTypes.object.isRequired,
    available: PropTypes.object.isRequired,
    teamRemoved: PropTypes.func.isRequired
  }

  static staticProps = {
    title: 'Team dashboard'
  }

  static initialCloudServiceList = [{
    name: 'Amazon RDS for PostgreSQL',
    description: 'Amazon Simple Storage Service (Amazon S3) is storage for the Internet. You can use Amazon S3 to store and retrieve any amount of data at any time, from anywhere on the web. You can accomplish these tasks using the simple and intuitive web interface of the AWS Management Console.',
    plan: 'Development',
    imageURL: 'https://s3.amazonaws.com/awsservicebroker/icons/AmazonRDS_LARGE.png',
    cluster: 'demo-notprod',
    namespaceClaim: 'dev'
  }]

  constructor(props) {
    super(props)
    this.state = {
      members: props.members,
      allUsers: [],
      membersToAdd: [],
      clusters: props.clusters,
      createNamespace: false,
      namespaceClaims: props.namespaceClaims,
      createCloudService: false,
      cloudServices: TeamDashboard.initialCloudServiceList
    }
  }

  static async getTeamDetails() {
    const getTeam = () => Promise.resolve(TeamData.team)
    const getTeamMembers = () => Promise.resolve(TeamData.members)
    const getTeamClusters = () => Promise.resolve(TeamData.clusters)
    const getNamespaceClaims = () => Promise.resolve(TeamData.namespaceClaims)
    const getAvailable = () => Promise.resolve(TeamData.allocations)

    return axios.all([getTeam(), getTeamMembers(), getTeamClusters(), getNamespaceClaims(), getAvailable()])
      .then(axios.spread(function (team, members, clusters, namespaceClaims, available) {
        return { team, members, clusters, namespaceClaims, available }
      }))
      .catch(err => {
        throw new Error(err.message)
      })
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
    return teamDetails
  }

  getAllUsers = async () => {
    const users = await apiRequest(null, 'get', apiPaths.users)
    if (users.items) {
      return users.items.map(user => user.spec.username).filter(user => user !== 'admin')
    }
    return []
  }

  componentDidMount() {
    return this.getAllUsers()
      .then(users => {
        const state = copy(this.state)
        state.allUsers = users
        this.setState(state)
      })
  }

  componentDidUpdate(prevProps) {
    const teamFound = Object.keys(this.props.team).length
    const prevTeamName = prevProps.team.metadata && prevProps.team.metadata.name
    if (teamFound && this.props.team.metadata.name !== prevTeamName) {
      const state = copy(this.state)
      state.members = this.props.members
      state.clusters = this.props.clusters
      state.namespaceClaims = this.props.namespaceClaims
      this.getAllUsers()
        .then(users => {
          state.allUsers = users
          this.setState(state)
        })
    }
  }

  addTeamMembersUpdated = membersToAdd => {
    const state = copy(this.state)
    state.membersToAdd = membersToAdd
    this.setState(state)
  }

  addTeamMembers = async () => {
    const state = copy(this.state)
    const members = state.members

    await asyncForEach(this.state.membersToAdd, async member => {
      await apiRequest(null, 'put', `${apiPaths.team(this.props.team.metadata.name).members}/${member}`)
      successMessage(`Team member added: ${member}`)
      members.items.push(member)
    })

    state.membersToAdd = []
    this.setState(state)
  }

  deleteTeamMember = member => {
    return async () => {
      const team = this.props.team.metadata.name
      try {
        await apiRequest(null, 'delete', `${apiPaths.team(team).members}/${member}`)
        const state = copy(this.state)
        const members = state.members
        members.items = members.items.filter(m => m !== member)
        this.setState(state)
        successMessage(`Team member removed: ${member}`)
      } catch (err) {
        console.error('Error removing team member', err)
        errorMessage('Error removing team member, please try again.')
      }
    }
  }

  handleResourceUpdated = resourceType => {
    return (updatedResource, done) => {
      const state = copy(this.state)
      const resource = state[resourceType].items.find(r => r.metadata.name === updatedResource.metadata.name)
      resource.status = updatedResource.status
      this.setState(state, done)
    }
  }

  handleResourceDeleted = resourceType => {
    return (name, done) => {
      const state = copy(this.state)
      const resource = state[resourceType].items.find(r => r.metadata.name === name)
      resource.deleted = true
      this.setState(state, done)
    }
  }

  deleteCluster = async (name, done) => {
    const team = this.props.team.metadata.name
    try {
      const state = copy(this.state)
      const cluster = state.clusters.items.find(c => c.metadata.name === name)
      await apiRequest(null, 'delete', `${apiPaths.team(team).clusters}/${cluster.metadata.name}`)
      cluster.status.status = 'Deleting'
      cluster.metadata.deletionTimestamp = new Date()
      this.setState(state, done)
      loadingMessage(`Cluster deletion requested: ${cluster.metadata.name}`)
    } catch (err) {
      console.error('Error deleting cluster', err)
      errorMessage('Error deleting cluster, please try again.')
    }
  }

  createNamespace = value => {
    return () => {
      const state = copy(this.state)
      state.createNamespace = value
      this.setState(state)
    }
  }

  handleNamespaceCreated = namespaceClaim => {
    const state = copy(this.state)
    state.createNamespace = false
    state.namespaceClaims.items.push(namespaceClaim)
    this.setState(state)
    loadingMessage(`Namespace "${namespaceClaim.spec.name}" requested on cluster "${namespaceClaim.spec.cluster.name}"`)
  }

  deleteNamespace = async (name, done) => {
    const team = this.props.team.metadata.name
    try {
      const state = copy(this.state)
      const namespaceClaim = state.namespaceClaims.items.find(nc => nc.metadata.name === name)
      await apiRequest(null, 'delete', `${apiPaths.team(team).namespaceClaims}/${name}`)
      namespaceClaim.status.status = 'Deleting'
      namespaceClaim.metadata.deletionTimestamp = new Date()
      this.setState(state, done)
      loadingMessage(`Namespace deletion requested: ${namespaceClaim.spec.name}`)
    } catch (err) {
      console.error('Error deleting namespace', err)
      errorMessage('Error deleting namespace, please try again.')
    }
  }

  createCloudService = value => {
    return () => {
      const state = copy(this.state)
      state.createCloudService = value
      this.setState(state)
    }
  }

  handleCloudServiceCreated = cloudService => {
    const state = copy(this.state)
    state.createCloudService = false
    state.cloudServices.push(cloudService)
    this.setState(state)
    loadingMessage(`Cloud service "${cloudService.name}" requested.`)
  }

  deleteCloudService = (name, plan) => () => {
    successMessage(`Cloud service deleted: ${name}`)
    this.setState({ cloudServices: copy(this.state.cloudServices).filter(s => `${s.name}_${s.plan}` !== `${name}_${plan}`) })
  }

  clusterAccess = async () => {
    const apiUrl = new URL(publicRuntimeConfig.koreApiPublicUrl)

    const profileConfigureCommand = `kore profile configure ${apiUrl.hostname}`
    const loginCommand = 'kore login'
    const kubeconfigCommand = `kore kubeconfig -t ${this.props.team.metadata.name}`

    const InfoItem = ({ num, title }) => (
      <div style={{ marginBottom: '10px' }}>
        <Badge style={{ backgroundColor: '#1890ff', marginRight: '10px' }} count={num} />
        <Text strong>{title}</Text>
      </div>
    )
    Modal.info({
      title: 'Cluster access',
      content: (
        <div style={{ marginTop: '20px' }}>
          <InfoItem num="1" title="Download" />
          <Paragraph>If you haven&apos;t already, download the CLI from <a href="https://github.com/appvia/kore/releases">https://github.com/appvia/kore/releases</a></Paragraph>

          <InfoItem num="2" title="Setup profile" />
          <Paragraph>Create a profile</Paragraph>
          <Paragraph className="copy-command" style={{ marginRight: '40px' }} copyable>{profileConfigureCommand}</Paragraph>
          <Paragraph>Enter the Kore API URL as follows</Paragraph>
          <Paragraph className="copy-command" style={{ marginRight: '40px' }} copyable>{apiUrl.origin}</Paragraph>

          <InfoItem num="3" title="Login" />
          <Paragraph>Login to the CLI</Paragraph>
          <Paragraph className="copy-command" style={{ marginRight: '40px' }} copyable>{loginCommand}</Paragraph>

          <InfoItem num="4" title="Setup access" />
          <Paragraph>Then, you can use the Kore CLI to setup access to your team&apos;s clusters</Paragraph>
          <Paragraph className="copy-command" style={{ marginRight: '40px' }} copyable>{kubeconfigCommand}</Paragraph>
          <Paragraph>This will add local kubernetes configuration to allow you to use <Text
            style={{ fontFamily: 'monospace' }}>kubectl</Text> to talk to the provisioned cluster(s).</Paragraph>
          <Paragraph>See examples: <a href="https://kubernetes.io/docs/reference/kubectl/overview/" target="_blank" rel="noopener noreferrer">https://kubernetes.io/docs/reference/kubectl/overview/</a></Paragraph>
        </div>
      ),
      width: 700,
      onOk() {}
    })
  }

  deleteTeam = async () => {
    try {
      const team = this.props.team.metadata.name
      const api = await KoreApi.client()
      await api.RemoveTeam(team)
      this.props.teamRemoved(team)
      successMessage(`Team "${team}" deleted`)
      return redirect({ router: Router, path: '/' })
    } catch (err) {
      console.log('Error deleting team', err)
      errorMessage('Team could not be deleted, please try again later')
    }
  }

  deleteTeamConfirm = () => {
    const { clusters } = this.state
    if (clusters.items.length > 0) {
      return Modal.warning({
        title: 'Warning: team cannot be deleted',
        content: (
          <>
            <Paragraph strong>The clusters must be deleted first</Paragraph>
            <List
              size="small"
              dataSource={clusters.items}
              renderItem={c => <List.Item>{c.spec.kind} <Text style={{ fontFamily: 'monospace', marginLeft: '15px' }}>{c.metadata.name}</Text></List.Item>}
            />
          </>
        ),
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

  render() {
    const { team, user, invitation } = this.props

    if (Object.keys(team).length === 0) {
      return <Error statusCode={404} />
    }

    const { members, namespaceClaims, allUsers, membersToAdd, createNamespace, createCloudService, clusters, cloudServices } = this.state
    const teamMembers = ['ADD_USER', ...members.items]

    const memberActions = member => {
      const deleteAction = (
        <Popconfirm
          key="delete"
          title="Are you sure you want to remove this user?"
          onConfirm={this.deleteTeamMember(member)}
          okText="Yes"
          cancelText="No"
        >
          <a>Remove</a>
        </Popconfirm>
      )
      if (member !== user.id) {
        return [deleteAction]
      }
      return []
    }

    const membersAvailableToAdd = allUsers.filter(user => !members.items.includes(user))
    const hasActiveClusters = Boolean(clusters.items.filter(c => c.status && c.status.status === 'Success').length)

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
          <Text strong>{team.spec.description}</Text>
          <Text style={{ float: 'right' }}><Text strong>Team ID: </Text>{team.metadata.name}</Text>
        </Paragraph>
        {invitation ? (
          <Alert
            message="You have joined this team from an invitation"
            type="info"
            showIcon
            style={{ marginBottom: '20px' }}
          />
        ) : null}

        <Tabs defaultActiveKey="clusters">
          <TabPane tab={
            <span>
              Clusters
              <Badge showZero={true} style={{ marginLeft: '10px', backgroundColor: '#1890ff' }} count={clusters.items.filter(c => !c.deleted).length} />
            </span>
          } key="clusters">
            <Card
              title={<div><Text style={{ marginRight: '10px' }}>Clusters</Text><Badge style={{ backgroundColor: '#1890ff' }} count={clusters.items.filter(c => !c.deleted).length} /></div>}
              style={{ marginBottom: '20px' }}
              extra={
                <div>
                  {hasActiveClusters ?
                    <Text style={{ marginRight: '20px' }}><a onClick={this.clusterAccess}><Icon type="eye" theme="twoTone" /> Access</a></Text> :
                    null
                  }
                  <Button type="primary">
                    <Link href="/teams/[name]/clusters/new" as={`/teams/${team.metadata.name}/clusters/new`}>
                      <a>+ New</a>
                    </Link>
                  </Button>
                </div>
              }
            >
              <List
                dataSource={clusters.items}
                renderItem={cluster => {
                  const namespaceClaims = (this.state.namespaceClaims.items || []).filter(nc => nc.spec.cluster.name === cluster.metadata.name && !nc.deleted)
                  return (
                    <Cluster
                      team={team.metadata.name}
                      cluster={cluster}
                      namespaceClaims={namespaceClaims}
                      deleteCluster={this.deleteCluster}
                      handleUpdate={this.handleResourceUpdated('clusters')}
                      handleDelete={this.handleResourceDeleted('clusters')}
                      refreshMs={10000}
                      propsResourceDataKey="cluster"
                      resourceApiPath={`${apiPaths.team(team.metadata.name).clusters}/${cluster.metadata.name}`}
                    />
                  )
                }}
              >
              </List>
            </Card>

            <Card
              title={<div><Text style={{ marginRight: '10px' }}>Namespaces</Text><Badge style={{ backgroundColor: '#1890ff' }} count={namespaceClaims.items.filter(c => !c.deleted).length} /></div>}
              style={{ marginBottom: '20px' }}
              extra={clusters.items.length > 0 ? <Button type="primary" onClick={this.createNamespace(true)}>+ New</Button> : null}
            >
              <List
                dataSource={namespaceClaims.items}
                renderItem={namespaceClaim =>
                  <NamespaceClaim
                    team={team.metadata.name}
                    namespaceClaim={namespaceClaim}
                    deleteNamespace={this.deleteNamespace}
                    handleUpdate={this.handleResourceUpdated('namespaceClaims')}
                    handleDelete={this.handleResourceDeleted('namespaceClaims')}
                    refreshMs={15000}
                    propsResourceDataKey="namespaceClaim"
                    resourceApiPath={`${apiPaths.team(team.metadata.name).namespaceClaims}/${namespaceClaim.metadata.name}`}
                  />
                }
              >
              </List>
            </Card>

            <Drawer
              title="Create namespace"
              placement="right"
              closable={false}
              onClose={this.createNamespace(false)}
              visible={createNamespace}
              width={700}
            >
              <NamespaceClaimForm team={team.metadata.name} clusters={clusters} handleSubmit={this.handleNamespaceCreated} handleCancel={this.createNamespace(false)}/>
            </Drawer>
          </TabPane>
          <TabPane tab={
            <span>
              Cloud services
              <Badge showZero={true} style={{ marginLeft: '10px', backgroundColor: '#1890ff' }} count={cloudServices.filter(c => !c.deleted).length} />
            </span>
          } key="cloud_services">
            <Card
              title={<div><Text style={{ marginRight: '10px' }}>Cloud services</Text><Badge style={{ backgroundColor: '#1890ff' }} count={cloudServices.filter(c => !c.deleted).length} /></div>}
              style={{ marginBottom: '20px' }}
              extra={
                <div>
                  <Button type="primary" onClick={this.createCloudService(true)}>+ New</Button>
                </div>
              }
            >
              <List
                dataSource={cloudServices}
                renderItem={service => (
                  <List.Item actions={[
                    <Popconfirm
                      key="delete"
                      title="Are you sure you want to delete this cloud service?"
                      onConfirm={this.deleteCloudService(service.name, service.plan)}
                      okText="Yes"
                      cancelText="No"
                    >
                      <a><Icon type="delete" /></a>
                    </Popconfirm>,
                    <ResourceStatusTag key="status" resourceStatus={{ status: 'Success' }} />
                  ]}>
                    <List.Item.Meta
                      avatar={<Avatar src={service.imageURL} />}
                      title={<p>{service.name} - {service.plan}</p>}
                      description={
                        <Text style={{ fontFamily: 'monospace' }}>{service.cluster} / {service.namespaceClaim}</Text>
                      }
                    />
                  </List.Item>
                )}
              >
              </List>
            </Card>
            <Drawer
              title="Create cloud service"
              placement="right"
              closable={false}
              onClose={this.createCloudService(false)}
              visible={createCloudService}
              width={900}
            >
              <ServiceForm clusters={clusters} namespaceClaims={namespaceClaims} handleSubmit={this.handleCloudServiceCreated} handleCancel={this.createCloudService(false)}/>
            </Drawer>
          </TabPane>
          <TabPane tab={
            <span>
              Members
              <Badge showZero={true} style={{ marginLeft: '10px', backgroundColor: '#1890ff' }} count={members.items.length} />
            </span>
          } key="members">
            <Card
              title={<div><Text style={{ marginRight: '10px' }}>Team members</Text><Badge style={{ backgroundColor: '#1890ff' }} count={members.items.length} /></div>}
              style={{ marginBottom: '16px' }}
              className="team-members"
              extra={<InviteLink team={team.metadata.name} />}
            >
              <List
                dataSource={teamMembers}
                renderItem={m => {
                  if (m === 'ADD_USER') {
                    return <List.Item style={{ paddingTop: '0' }} actions={[<Button key="add" type="secondary" onClick={this.addTeamMembers}>Add</Button>]}>
                      <List.Item.Meta
                        title={
                          <Select
                            mode="multiple"
                            placeholder="Add existing users to this team"
                            onChange={this.addTeamMembersUpdated}
                            style={{ width: '100%' }}
                            value={membersToAdd}
                          >
                            {membersAvailableToAdd.map((user, idx) => (
                              <Option key={idx} value={user}>{user}</Option>
                            ))}
                          </Select>
                        }
                      />
                    </List.Item>
                  } else {
                    return <List.Item actions={memberActions(m)}>
                      <List.Item.Meta avatar={<Avatar icon="user" />} title={<Text>{m} {m === user.id ? <Tag>You</Tag>: null}</Text>} />
                    </List.Item>
                  }
                }}
              >
              </List>
            </Card>
          </TabPane>
        </Tabs>

      </div>
    )
  }
}

export default TeamDashboard
