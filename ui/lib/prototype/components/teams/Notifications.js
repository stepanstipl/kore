import React from 'react'
import Link from 'next/link'
import { Button, Badge, Icon, List, Dropdown, Typography } from 'antd'
const { Text } = Typography
import moment from 'moment'

import TeamNotificationData from '../../utils/dummy-notification-data'
import IconTooltip from '../../../components/utils/IconTooltip'
import EventNotificationsTab from './settings/EventNotificationsTab'

class Notifications extends React.Component {

  state = {
    notifications: []
  }

  componentDidMount() {
    this.setState({
      notifications: TeamNotificationData.notifications.items
    })
  }

  acknowledgeAll = () => {
    this.setState(state => ({
      notifications: state.notifications.map(n => ({ ...n, acknowledged: true }))
    }))
  }

  acknowledge = (index) => () => {
    this.setState(state => ({
      notifications: state.notifications.map((n, i) => i === index ? { ...n, acknowledged: true } : n)
    }))
  }

  renderItem = ({ event, detail, creationTimestamp, acknowledged }, index) => {
    const actions = ! acknowledged ? [<IconTooltip key="seen" icon="check-circle" text="Mark as seen" onClick={this.acknowledge(index)} />] : []
    const lastItem = index === this.state.notifications.length - 1
    return (
      <List.Item
        key={index}
        style={{ paddingLeft: '15px', paddingRight: '15px', borderBottom: !lastItem ? '1px solid #CCC' : '', backgroundColor: !acknowledged ? '#efefef' : '' }}
        actions={actions}>
        <List.Item.Meta
          title={EventNotificationsTab.EVENTS[event] || event}
          description={detail}
        />
        <div style={{ marginRight: '10px' }}>{moment(creationTimestamp).fromNow()}</div>
      </List.Item>
    )
  }

  render() {
    const notifications = this.state.notifications
    const unseenCount = notifications.filter(n => !n.acknowledged).length

    return (
      <Dropdown trigger={['click']} overlayStyle={{ width: '500px', border: '1px solid #999', padding: 0 }} overlay={
        <List className="notifications">
          <List.Item
            key="heading"
            style={{ paddingLeft: '15px', paddingRight: '15px', backgroundColor: '#3d5b58' }}
            actions={[
              <a style={{ textDecoration: 'underline', color: '#FFF' }} key="seen" onClick={this.acknowledgeAll}>Mark all as seen</a>,
              <Link key="all" href="/prototype/teams/proto/notifications">
                <a style={{ textDecoration: 'underline', color: '#FFF' }}>See all</a>
              </Link>
            ]}>
            <Text strong style={{ color: '#FFF', width: '100%' }}>Notifications</Text>
          </List.Item>
          {notifications.length > 0 ? notifications.map((n, i) => this.renderItem(n, i)) : null}
        </List>
      }>
        <Badge count={unseenCount} overflowCount={9}>
          <Button>
            <Icon type="bell" />
          </Button>
        </Badge>
      </Dropdown>
    )
  }
}

export default Notifications
