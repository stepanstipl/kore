import React from 'react'
import PropTypes from 'prop-types'
import { Typography, List, Avatar, Tag, message } from 'antd'
const { Text } = Typography
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

import KoreApi from '../../lib/kore-api'
import Breadcrumb from '../../lib/components/layout/Breadcrumb'

class ConfigureUsersPage extends React.Component {
  static propTypes = {
    user: PropTypes.object.isRequired,
    users: PropTypes.array.isRequired,
    admins: PropTypes.array.isRequired
  }

  state = {
    users: this.props.users,
    admins: this.props.admins
  }

  static staticProps = {
    title: 'Configure users',
    adminOnly: true
  }

  static getInitialProps = async () => {
    const api = await KoreApi.client()
    let [ users, admins ] = await Promise.all([
      api.ListUsers(),
      api.ListTeamMembers(publicRuntimeConfig.koreAdminTeamName)
    ])
    users = users.items
    admins = admins.items
    return { users, admins }
  }

  makeAdmin = (username) => {
    return async () => {
      try {
        await (await KoreApi.client()).AddTeamMember(publicRuntimeConfig.koreAdminTeamName, username)
        this.setState(state => ({
          admins: [ ...state.admins, username ]
        }))
        message.success(`${username} is now admin`)
      } catch (err) {
        console.error('Error trying to make admin')
        message.error(`Failed to make ${username} admin`)
      }
    }
  }

  revokeAdmin = (username) => {
    return async () => {
      try {
        await (await KoreApi.client()).RemoveTeamMember(publicRuntimeConfig.koreAdminTeamName, username)
        this.setState(state => ({
          admins: state.admins.filter(m => m !== username)
        }))
        message.success(`${username} is no longer admin`)
      } catch (err) {
        console.error('Error trying to revoke admin', err)
        message.error(`Failed to revoke admin from user ${username}`)
      }
    }
  }

  render() {
    return (
      <div>
        <Breadcrumb items={[{ text: 'Configure' }, { text: 'Users' }]} />
        <List
          dataSource={this.state.users}
          renderItem={user => {
            const isUser = user.spec.username === this.props.user.id
            const isAdmin = this.state.admins.includes(user.spec.username)
            const actions = []
            if (isAdmin && !isUser) {
              actions.push(<Text key="revoke_admin"><a onClick={this.revokeAdmin(user.spec.username)}>Revoke admin</a></Text>)
            }
            if (!isAdmin) {
              actions.push(<Text key="make_admin"><a onClick={this.makeAdmin(user.spec.username)}>Make admin</a></Text>)
            }
            return (
              <List.Item
                key={user.spec.username}
                actions={actions}
              >
                <List.Item.Meta
                  avatar={<Avatar icon="user" />}
                  title={<Text>{user.spec.username} {isAdmin ? <Tag color="green">admin</Tag> : null}{isUser ? <Tag>You</Tag>: null}</Text>}
                  description={user.spec.email}
                />
              </List.Item>
            )
          }}
        >
        </List>
      </div>
    )
  }
}

export default ConfigureUsersPage
