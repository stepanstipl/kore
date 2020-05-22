import * as React from 'react'
import { Form, Select, Table, Icon } from 'antd'
import { startCase, debounce } from 'lodash'
import PlanOptionBase from '../PlanOptionBase'
import KoreApi from '../../../kore-api'

export default class PlanOptionClusterUsers extends PlanOptionBase {
  constructor(props) {
    super(props)
    this.lastSearchId = 0
    this.searchUsers = debounce(this.searchUsers, 250)
  }

  state = {
    loadingUsers: false,
    allUsers: null,
    searchUsers: []
  }

  searchUsers = async (value) => {
    this.lastSearchId += 1
    const searchId = this.lastSearchId

    // Load the user list once (we're doing filtering client-side at the moment)
    let allUsers = this.state.allUsers
    if (!allUsers) {
      this.setState({ loadingUsers: true })
      if (this.props.manage) {
        allUsers = await (await KoreApi.client()).ListUsers()
        // all users and team members have different shaped return objects, so map them both to a plain array of usernames
        allUsers = allUsers.items.map((u) => u.spec.username)
      } else {
        allUsers = await (await KoreApi.client()).ListTeamMembers(this.props.team.metadata.name)
        // all users and team members have different shaped return objects, so map them both to a plain array of usernames
        allUsers = allUsers.items
      }
      this.setState({ allUsers, loadingUsers: false })
    }

    if (searchId !== this.lastSearchId) {
      // if this has been superceded by another search, just bin it it.
      return
    }

    const searchUsers = allUsers.filter((u) => {
      // remove users already in the list
      if (this.props.value && this.props.value.find((existingUser) => existingUser.username === u)) {
        return false
      }
      // remove users who don't match the entered value
      if (u.indexOf(value) === -1) {
        return false
      }
      return true
    })

    this.setState({
      loadingUsers: false,
      searchUsers: searchUsers
    })
  }

  addUser = (username) => {
    if (!this.props.editable || !this.props.onChange) {
      return
    }

    // No-op if already present:
    if (this.props.value && this.props.value.find((u) => u.username === username)) {
      return
    }

    // Need to handle the value being undefined in the case where this is a new plan or no
    // users are defined yet.
    let newValue
    if (this.props.value) {
      newValue = [...this.props.value, { username, roles: [] }]
    } else {
      newValue = [{ username, roles: [] }]
    }

    this.setState({ searchUsers: [] })

    this.props.onChange(this.props.name, newValue)
  }

  removeUser = (username) => {
    if (!this.props.editable || !this.props.onChange) {
      return
    }

    this.setState({ searchUsers: [] })

    this.props.onChange(
      this.props.name, 
      this.props.value.filter((u) => u.username !== username)
    )
  }

  setUserRoles = (username, roles) => {
    if (!this.props.editable || !this.props.onChange) {
      return
    }

    this.props.onChange(
      this.props.name, 
      this.props.value.map((u) => u.username !== username ? u : { ...u, roles: roles })
    )
  }

  render() {
    const { name, value, editable, property } = this.props

    const valueOrDefault = value || property.default || []
    const displayName = this.props.displayName || startCase(name)
    const description = this.props.manage ? 'Set default users to be added to every cluster created from this plan' : 'Control which team members have access to this cluster'
    const roles = ['cluster-admin', 'admin', 'edit', 'view']
    const columns = [
      { title: 'User', dataIndex: 'username', key: 'username', width: '45%' },
      { title: 'Roles', dataIndex: 'roles', key: 'tags', width: '45%', render: function renderRoles(userRoles, r) { 
        return (
          <Select mode="multiple" value={userRoles}  onChange={(selectedRoles) => this.setUserRoles(r.username, selectedRoles)}>
            {!editable ? null : roles.map((role) => userRoles.indexOf(role) === -1 ? <Select.Option key={role} value={role}>{role}</Select.Option> : null)}
          </Select>
        )
      }.bind(this) },
      { key: 'action', width: '10%', render: function renderAction(_, r) {
        if (!editable) {
          return null
        }
        return <><div style={{ textAlign: 'right' }}><a onClick={() => this.removeUser(r.username)}><Icon type="delete" title="Delete" /></a></div></>
      }.bind(this) },
    ]

    return (
      <Form.Item label={displayName} help={description}>
        <Table 
          size="small" 
          pagination={false} 
          dataSource={valueOrDefault} 
          columns={columns} 
          rowKey={r => r.username}
          footer={!editable ? null : () => (
            <Select
              mode="single" showSearch={true}
              placeholder="Start typing username to find users to add"
              onSearch={this.searchUsers}
              onChange={(v) => this.addUser(v)}
              value={undefined}
            >
              {this.state.searchUsers.map((user, idx) => (
                <Select.Option key={idx} value={user}>{user}</Select.Option>
              ))}
            </Select>
          )}
        >
        </Table>
        {this.validationErrors(name)}
      </Form.Item>
    )
  }
}
