import React from 'react'
import { Button, Card, Divider, Drawer, Form, Icon, List, Popconfirm, Select, Typography } from 'antd'
const { Paragraph, Text } = Typography

import { successMessage } from '../../../../utils/message'

// prototype imports
import TeamNotificationData from '../../../utils/dummy-notification-data'
import EventNotificationForm from './EventNotificationForm'

class EventNotificationsTab extends React.Component {

  static EVENTS = {
    CLUSTER_CREATED: 'Cluster created',
    CLUSTER_DELETED: 'Cluster deleted',
    SERVICE_CREATED: 'Cloud service created',
    SERVICE_DELETED: 'Cloud service deleted'
  }

  static INITIAL_FILTERS = {
    events: [],
    integrations: []
  }

  state = {
    dataLoading: true,
    filters: EventNotificationsTab.INITIAL_FILTERS,
    newNotification: false
  }

  componentDidMount() {
    this.setState({
      dataLoading: false,
      notifications: TeamNotificationData.notifications.items,
      integrations: TeamNotificationData.integrations.items
    })
  }

  filterChanged = (name) => (value) => {
    this.setState(state => ({
      filters: { ...state.filters, [name]: value }
    }))
  }

  notificationsCreated = (notificationList) => {
    this.setState(state => ({
      notifications: [ ...state.notifications, ...notificationList ],
      newNotification: false
    }))
  }

  renderNotificationAction = (notification) => {
    const kore = (
      <div style={{ display: 'block' }}>
        <img src="/static/images/appvia-colour.svg" width={14} style={{ marginRight: '8px' }} />
        <Text>Kore notification</Text>
      </div>
    )

    if (!notification.spec.integration) {
      return kore
    }

    const integration = this.state.integrations.find(i => i.metadata.name === notification.spec.integration.name)
    if (notification.spec.integration.kind === 'SlackIntegration') {
      return (
        <>
          {kore}
          <div style={{ display: 'block' }}>
            <Icon type="slack" style={{ marginRight: '8px' }} />
            <Text>Slack notification</Text>
            <Divider type="vertical" />
            <Text>{integration.spec.displayName}</Text>
            <Divider type="vertical" />
            <Text>#{notification.spec.channel}</Text>
          </div>
        </>
      )
    }
    if (notification.spec.integration.kind === 'EmailIntegration') {
      return (
        <>
          {kore}
          <div style={{ display: 'block' }}>
            <Icon type="mail" style={{ marginRight: '8px' }} />
            <Text>Email notification</Text>
            <Divider type="vertical" />
            <Text>{notification.spec.emailAddressList ? notification.spec.emailAddressList.join(', ') : 'All team members'}</Text>
          </div>
        </>
      )
    }
    return null
  }

  deleteNotification = (notification) => () => {
    this.setState(state => ({
      notifications: state.notifications.filter(n => n.metadata.name !== notification.metadata.name)
    }), () => successMessage('Notification deleted'))
  }

  notificationActions = (notification) => {
    return [(
      <Popconfirm
        key="delete"
        title="Are you sure you want to delete this notification?"
        onConfirm={this.deleteNotification(notification)}
        okText="Yes"
        cancelText="No"
      >
        <Icon type="delete" />
      </Popconfirm>
    )]
  }

  renderFilters = () => {
    return (
      <Card size="small" style={{ marginBottom: '20px' }}>
        <Form.Item labelAlign="left" labelCol={{ span: 4 }} wrapperCol={{ span: 20 }} label={<Text strong>Filter by events</Text>} style={{ marginBottom: '10px' }}>
          <Select
            value={this.state.filters.events}
            mode="multiple"
            style={{ width: '100%' }}
            placeholder="Showing notifications for all events"
            onChange={this.filterChanged('events')}
          >
            {Object.keys(EventNotificationsTab.EVENTS).map(event => <Select.Option key={event}>{EventNotificationsTab.EVENTS[event]}</Select.Option>)}
          </Select>
        </Form.Item>
        <Form.Item labelAlign="left" labelCol={{ span: 4 }} wrapperCol={{ span: 20 }} label={<Text strong>Filter by integrations</Text>} style={{ marginBottom: '10px' }}>
          <Select
            value={this.state.filters.integrations}
            mode="multiple"
            style={{ width: '100%' }}
            placeholder="Showing notifications for all integrations"
            onChange={this.filterChanged('integrations')}
          >
            {this.state.integrations.map(integration => <Select.Option key={integration.metadata.name}>{integration.kind}</Select.Option>)}
          </Select>
        </Form.Item>
        <a style={{ display: 'block', marginTop: '10px', marginBottom: '5px', textDecoration: 'underline' }} onClick={() => this.setState({ filters: EventNotificationsTab.INITIAL_FILTERS })}>Clear filters</a>
      </Card>
    )
  }

  render() {
    const { dataLoading, notifications, newNotification, filters } = this.state
    if (dataLoading) {
      return <Icon type="loading" />
    }
    const filteredNotifications = notifications.filter(n => {
      const eventMatch = (filters.events.length === 0 ? [n.spec.event] : filters.events).includes(n.spec.event)
      const integrationMatch = filters.integrations.length === 0 || (n.spec.integration && (filters.integrations.length === 0 ? [n.spec.integration.name] : filters.integrations).includes(n.spec.integration.name) )
      return eventMatch && integrationMatch
    })

    return (
      <>
        <Button onClick={() => this.setState({ newNotification: true })} type="primary">New notification</Button>
        <Divider />
        {this.renderFilters()}
        <Paragraph strong>Showing {filteredNotifications.length} of {notifications.length} notifications</Paragraph>
        <List
          dataSource={filteredNotifications}
          renderItem={item => (
            <List.Item actions={this.notificationActions(item)}>
              <List.Item.Meta
                title={<Text style={{ fontSize: '16px' }}>{EventNotificationsTab.EVENTS[item.spec.event]}</Text>}
                description={this.renderNotificationAction(item)}
              />
            </List.Item>
          )}
        />
        <Drawer
          title="New event notification"
          visible={newNotification}
          onClose={() => this.setState({ newNotification: false })}
          width={700}
        >
          {newNotification ? (
            <EventNotificationForm
              handleCancel={() => this.setState({ newNotification: false })}
              handleSubmit={this.notificationsCreated}
            />
          ) : null}
        </Drawer>
      </>
    )
  }
}

export default EventNotificationsTab
