import React from 'react'
import PropTypes from 'prop-types'
import { Avatar, Button, Col, Divider, Icon, List, Popconfirm, Row, Select, Tag, Typography } from 'antd'
const { Text } = Typography
const { Option } = Select

import KoreApi from '../../../kore-api'
import copy from '../../../utils/object-copy'
import asyncForEach from '../../../utils/async-foreach'
import InviteLink from '../InviteLink'
import { successMessage } from '../../../utils/message'

class MembersTab extends React.Component {

  static propTypes = {
    user: PropTypes.object.isRequired,
    team: PropTypes.object.isRequired,
    getMemberCount: PropTypes.func
  }

  state = {
    dataLoading: true,
    members: [],
    users: [],
    membersToAdd: []
  }

  async fetchComponentData () {
    try {
      const api = await KoreApi.client()
      let [members, users] = await Promise.all([
        api.ListTeamMembers(this.props.team.metadata.name),
        api.ListUsers()
      ])
      members = members.items
      this.props.getMemberCount && this.props.getMemberCount(members.length)
      users = users.items.map(user => user.spec.username).filter(user => user !== 'admin')

      return { members, users }
    } catch (err) {
      console.error('Unable to load data for members tab', err)
      return {}
    }
  }

  componentDidMountComplete = null
  componentDidMount() {
    this.componentDidMountComplete = this.fetchComponentData().then(data => {
      this.setState({ ...data, dataLoading: false })
    })
  }

  componentDidUpdate(prevProps) {
    if (prevProps.team.metadata.name !== this.props.team.metadata.name) {
      this.setState({ dataLoading: true })
      return this.fetchComponentData().then(data => this.setState({ ...data, dataLoading: false }))
    }
  }

  addTeamMembersUpdated = (membersToAdd) => this.setState({ membersToAdd })

  addTeamMembers = async () => {
    const members = copy(this.state.members)
    try {
      const api = await KoreApi.client()
      await asyncForEach(this.state.membersToAdd, async member => {
        await api.AddTeamMember(this.props.team.metadata.name, member)
        members.push(member)
        successMessage(`Team member added: ${member}`)
      })
      this.props.getMemberCount && this.props.getMemberCount(members.length)
      this.setState({ members, membersToAdd: [] })
    } catch (err) {
      console.error('Error adding team member', err)
      errorMessage('Error adding team members, please try again.')
    }
  }

  deleteTeamMember = (member) => {
    return async () => {
      let members = copy(this.state.members)
      try {
        const api = await KoreApi.client()
        await api.RemoveTeamMember(this.props.team.metadata.name, member)
        members = members.filter(m => m !== member)
        this.props.getMemberCount && this.props.getMemberCount(members.length)
        this.setState({ members })
        successMessage(`Team member removed: ${member}`)
      } catch (err) {
        console.error('Error removing team member', err)
        errorMessage('Error removing team member, please try again.')
      }
    }
  }

  memberActions = (member) => {
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
    if (member !== this.props.user.id) {
      return [deleteAction]
    }
    return []
  }

  render() {
    const { team } = this.props
    const { dataLoading, members, users, membersToAdd } = this.state

    if (dataLoading) {
      return <Icon type="loading" />
    }

    const membersAvailableToAdd = users.filter(user => !members.includes(user))

    return (
      <>
        <Row gutter={16} style={{ marginBottom: '20px' }}>
          <Col span={12}>
            <Select
              mode="multiple"
              placeholder="Add existing users to this team"
              onChange={this.addTeamMembersUpdated}
              style={{ width: '100%' }}
              value={membersToAdd}
            >
              {membersAvailableToAdd.map((user, idx) => <Option key={idx} value={user}>{user}</Option>)}
            </Select>
          </Col>
          <Col span={3}><Button key="add" type="secondary" onClick={this.addTeamMembers}>Add</Button></Col>
          <Col span={9}><InviteLink team={team.metadata.name} /></Col>
        </Row>

        <Divider />

        <List
          dataSource={members}
          renderItem={m => (
            <List.Item actions={this.memberActions(m)}>
              <List.Item.Meta avatar={<Avatar icon="user" />} title={<Text>{m} {m === this.props.user.id ? <Tag>You</Tag>: null}</Text>} />
            </List.Item>
          )}
        >
        </List>
      </>
    )
  }
}

export default MembersTab
