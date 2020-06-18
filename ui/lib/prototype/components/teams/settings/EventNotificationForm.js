import React from 'react'
import PropTypes from 'prop-types'
import { Alert, Button, Card, Divider, Form, Icon, Input, Radio, Select, Switch, Typography } from 'antd'
const { Paragraph, Text } = Typography

import EventNotificationsTab from './EventNotificationsTab'
import IconTooltip from '../../../../components/utils/IconTooltip'
import { successMessage } from '../../../../utils/message'

class EventNotificationForm extends React.Component {
  static propTypes = {
    form: PropTypes.object.isRequired,
    handleCancel: PropTypes.func.isRequired,
    handleSubmit: PropTypes.func.isRequired
  }

  state = {
    submitting: false,
    slackNotification: false,
    emailNotification: false,
    emailNotificationAudience: null
  }

  componentDidMount() {
    // To disabled submit button at the beginning.
    this.props.form.validateFields()
  }

  disableButton = (fieldsError) => {
    if (this.state.submitting) {
      return true
    }
    return Object.keys(fieldsError).some(field => fieldsError[field])
  }

  switch = (name) => (
    <Switch
      checked={this.state[name]}
      onChange={(checked) => this.setState({ [name]: checked })}
      checkedChildren={<Icon type="check" />}
      unCheckedChildren={<Icon type="close" />}
    />
  )

  handleSubmit = (e) => {
    e.preventDefault()

    this.props.form.validateFields((err, values) => {

      if (err) {
        return
      }

      const notifications = []
      if (!this.state.slackNotification && !this.state.emailNotification) {
        notifications.push({
          metadata: { name: `${values.event}-slack` },
          spec: { event: values.event }
        })
      }

      if (this.state.slackNotification) {
        notifications.push({
          metadata: { name: `${values.event}-slack` },
          spec: {
            event: values.event,
            channel: values.slackChannel,
            integration: {
              group: 'integrations.kore.appvia.io',
              version: 'v1',
              kind: 'SlackIntegration',
              namespace: 'proto',
              name: 'proto-slack-integration'
            }
          }
        })
      }
      if (this.state.emailNotification) {
        notifications.push({
          metadata: { name: `${values.event}-email` },
          spec: {
            event: values.event,
            emailAddressList: values.emailAddressList,
            integration: {
              group: 'integrations.kore.appvia.io',
              version: 'v1',
              kind: 'EmailIntegration',
              namespace: 'proto',
              name: 'proto-email-integration'
            }
          }
        })
      }
      successMessage(`Notifications created for the "${EventNotificationsTab.EVENTS[values.event]}" event`)
      this.props.handleSubmit(notifications)
    })
  }

  fieldError = (field) => this.props.form.isFieldTouched(field) && this.props.form.getFieldError(field)

  render() {
    const { getFieldDecorator, getFieldsError } = this.props.form
    const formConfig = {
      layout: 'horizontal',
      labelAlign: 'left',
      hideRequiredMark: true,
      labelCol: {
        sm: { span: 24 },
        md: { span: 6 },
        lg: { span: 4 }
      },
      wrapperCol: {
        span: 12
      }
    }

    return (
      <Form {...formConfig} onSubmit={this.handleSubmit}>
        <Form.Item label="Event" validateStatus={this.fieldError('event') ? 'error' : ''} help={this.fieldError('event') || ''}>
          {getFieldDecorator('event', {
            rules: [{ required: true, message: 'Please select the event!' }]
          })(
            <Select placeholder="Select event">
              {Object.keys(EventNotificationsTab.EVENTS).map(event => <Select.Option key={event}>{EventNotificationsTab.EVENTS[event]}</Select.Option>)}
            </Select>
          )}
        </Form.Item>
        <Card style={{ marginBottom: '20px' }}>
          <Alert
            message="Choose notification actions"
            description="Notifications will always be sent to Kore, but you can also choose additional ways to be notified."
            type="info"
            showIcon
            style={{ marginBottom: '20px' }}
          />
          <Card title={<Text style={{ fontSize: '16px' }}>Kore notification</Text>} headStyle={{ border: 'none', padding: '0 12px' }} bodyStyle={{ padding: '0 12px' }} bordered={false} extra={<Switch defaultChecked disabled={true} checkedChildren={<Icon type="check" />} unCheckedChildren={<Icon type="close" />} />}>
            <Paragraph type="secondary">Notifications are always available in Kore.</Paragraph>
          </Card>

          <Divider style={{ marginBottom: '10px', marginTop: '10px' }} />

          <Card title={<Text style={{ fontSize: '16px' }}>Slack notification <IconTooltip text="Send a notification to your slack channel when this event occurs" icon="info-circle" /></Text>} headStyle={{ border: 'none', padding: '0 12px' }} bodyStyle={{ padding: '0 12px' }} bordered={false} extra={this.switch('slackNotification')}>
            {this.state.slackNotification ? (
              <>
                <Paragraph type="secondary">Send a notification to your slack channel</Paragraph>
                <Form.Item validateStatus={this.fieldError('slackChannel') ? 'error' : ''} help={this.fieldError('slackChannel') || ''}>
                  {getFieldDecorator('slackChannel', {
                    rules: [{ required: true, message: 'Please enter the slack channel!' }]
                  })(
                    <Input placeholder="Channel" />
                  )}
                </Form.Item>
              </>
            ) : null}
          </Card>

          <Divider style={{ marginBottom: '10px', marginTop: '10px' }} />

          <Card title={<Text style={{ fontSize: '16px' }}>Email notification <IconTooltip text="Send en email when this event occurs, choose to email all team members or specific email addresses" icon="info-circle" /></Text>} headStyle={{ border: 'none', padding: '0 12px' }} bodyStyle={{ padding: '0 12px' }} bordered={false} extra={this.switch('emailNotification')}>
            {this.state.emailNotification ? (
              <>
                <div style={{ marginBottom: '15px' }}>
                  <Radio.Group onChange={(e) => this.setState({ emailNotificationAudience: e.target.value })} value={this.state.emailNotificationAudience}>
                    <Radio value={'TEAM'} style={{ marginRight: '20px' }}>
                      <Text>All team members</Text>
                    </Radio>
                    <Radio value={'SPECIFIC'}>
                      <Text>Specific email addresses</Text>
                    </Radio>
                  </Radio.Group>
                </div>

                {this.state.emailNotificationAudience === 'SPECIFIC' ? (
                  <Form.Item wrapperCol={{ span: 24 }} validateStatus={this.fieldError('emailAddressList') ? 'error' : ''} help={this.fieldError('emailAddressList') || ''}>
                    {getFieldDecorator('emailAddressList', {
                      rules: [{ required: true, message: 'Please enter at least one email address!' }]
                    })(
                      <Select mode="tags" style={{ width: '100%' }} placeholder="Enter email addresses (paste in comma-separated)" tokenSeparators={[',']} />
                    )}
                  </Form.Item>
                ) : null}
              </>
            ) : null}
          </Card>

        </Card>
        <Form.Item>
          <Button type="primary" htmlType="submit" loading={this.state.submitting} disabled={this.disableButton(getFieldsError())}>Save</Button>
          <Button type="link" onClick={this.props.handleCancel}>Cancel</Button>
        </Form.Item>
      </Form>
    )
  }
}

const WrapperEventNotificationForm = Form.create({ name: 'notification_event' })(EventNotificationForm)

export default WrapperEventNotificationForm