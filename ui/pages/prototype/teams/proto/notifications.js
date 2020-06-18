import React from 'react'
import moment from 'moment'
import Link from 'next/link'
import { Col, Icon, List, Row, Typography } from 'antd'
const { Title } = Typography

import Breadcrumb from '../../../../lib/components/layout/Breadcrumb'
import TeamNotificationData from '../../../../lib/prototype/utils/dummy-notification-data'
import EventNotificationsTab from '../../../../lib/prototype/components/teams/settings/EventNotificationsTab'
import IconTooltip from '../../../../lib/components/utils/IconTooltip'

class TeamNotifications extends React.Component {

  state = {
    dataLoading: true
  }

  componentDidMount() {
    this.setState({
      dataLoading: false,
      notifications: TeamNotificationData.notifications.items
    })
  }

  acknowledge = (index) => () => {
    this.setState(state => ({
      notifications: state.notifications.map((n, i) => i === index ? { ...n, acknowledged: true } : n)
    }))
  }

  renderItem = ({ event, detail, creationTimestamp, acknowledged }, index) => {
    const actions = ! acknowledged ? [<IconTooltip key="seen" icon="check-circle" text="Mark as seen" onClick={this.acknowledge(index)} />] : []
    return (
      <List.Item
        key={index}
        style={{ paddingLeft: '15px', paddingRight: '15px', backgroundColor: !acknowledged ? '#efefef' : '' }}
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
    if (this.state.dataLoading) {
      return <Icon type="loading" />
    }

    return (
      <>
        <Breadcrumb items={[{ text: 'Proto', link: '/prototype/teams/proto', href: '/prototype/teams/proto' }, { text: 'Notifications' }]}/>
        <Row>
          <Col span={18}>
            <Title level={3}>Team notifications</Title>
          </Col>
          <Col style={{ textAlign: 'right' }} span={6}>
            <Link href="/prototype/teams/proto/settings/notifications">
              <a style={{ textDecoration: 'underline' }}>Notification settings</a>
            </Link>
          </Col>
        </Row>
        <List className="notifications">
          {this.state.notifications.map((n, i) => this.renderItem(n, i))}
        </List>
      </>
    )
  }
}

export default TeamNotifications
