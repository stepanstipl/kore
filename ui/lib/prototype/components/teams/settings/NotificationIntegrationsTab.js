import React from 'react'
import { Button, Card, Form, Icon, Input, Switch, Typography } from 'antd'
const { Paragraph, Text } = Typography

import { successMessage } from '../../../../utils/message'

// prototype imports
import TeamNotificationData from '../../../utils/dummy-notification-data'

class NotificationIntegrationsTab extends React.Component {

  state = {
    dataLoading: true
  }

  componentDidMount() {
    this.setState({
      dataLoading: false,
      integrations: TeamNotificationData.integrations.items
    })
  }

  toggleIntegration = (integrationKind, enabled) => {
    this.setState(state => {
      const integrations = state.integrations
      const updated = integrations.find(i => i.kind === integrationKind)
      updated.spec.enabled = enabled
      return { integrations }
    }, () => successMessage(`${integrationKind} ${enabled ? 'enabled' : 'disabled'}`))
  }

  renderSwitch = (integrationKind) => (
    <Switch
      checked={this.state.integrations.find(i => i.kind === integrationKind).spec.enabled}
      onChange={(checked) => this.toggleIntegration(integrationKind, checked)}
      checkedChildren={<Icon type="check" />}
      unCheckedChildren={<Icon type="close" />}
    />
  )

  render() {
    const { dataLoading, integrations } = this.state
    if (dataLoading) {
      return <Icon type="loading" />
    }

    const slackIntegration = integrations.find(i => i.kind === 'SlackIntegration')
    const emailIntegration = integrations.find(i => i.kind === 'EmailIntegration')

    return (
      <>
        <Card style={{ marginBottom: '20px' }} title={<Text style={{ fontSize: '16px', fontWeight: '600' }}>Slack</Text>} headStyle={{ border: 'none' }} bodyStyle={{ padding: slackIntegration.spec.enabled ? '20px' : 0 }} extra={this.renderSwitch('SlackIntegration')}>
          {slackIntegration.spec.enabled ? (
            <>
              <Paragraph type="secondary">
                Visit <a style={{ textDecoration: 'underline' }} href="#">Slack&apos;s Kore Integration page</a> and click &quot;Add to Slack&quot;. Choose a channel, and click the &quot;Add Kore Integration&quot; button. Copy your Webhook URL and click the &quot;Save Integration&quot; button.
              </Paragraph>
              <Paragraph type="secondary">
                Add your Webhook URL below.
              </Paragraph>
              <Form.Item label="Webhook URL" labelAlign="left" labelCol={{ span: 4 }} wrapperCol={{ span: 16 }}>
                <Input value={slackIntegration.spec.webhookURL} />
              </Form.Item>
              <Form.Item label="Name" labelAlign="left" labelCol={{ span: 4 }} wrapperCol={{ span: 16 }}>
                <Input value={slackIntegration.spec.displayName} />
              </Form.Item>
              <Form.Item style={{ marginBottom: 0 }}>
                <Button type="primary" onClick={() => successMessage('Slack integration changes saved')}>Save</Button>
              </Form.Item>
            </>
          ) : null}
        </Card>

        <Card title={<Text style={{ fontSize: '16px', fontWeight: '600' }}>Email</Text>} headStyle={{ border: 'none' }} bodyStyle={{ padding: emailIntegration.spec.enabled ? '20px' : 0 }} extra={this.renderSwitch('EmailIntegration')}>
          {emailIntegration.spec.enabled ? (
            <>
              <Paragraph style={{ marginBottom: 0 }} type="secondary">
                Choose to email all members of the team or specific email addresses, per event.
              </Paragraph>
            </>
          ) : null}
        </Card>
      </>
    )
  }
}

export default NotificationIntegrationsTab
