import * as React from 'react'
import axios from 'axios'
import PropTypes from 'prop-types'
import { Typography, Steps, Button, Alert, Form, Input, Card } from 'antd'
import { auth } from '../../../../config'
import CopyTextWithLabel from '../../../../lib/components/CopyTextWithLabel'
import copy from '../../../../lib/utils/object-copy'
import redirect from '../../../../lib/utils/redirect'

const { Step } = Steps
const { Title, Paragraph, Text } = Typography

class ConfigureHubForm extends React.Component {
  static propTypes = {
    form: PropTypes.any.isRequired
  }

  render() {
    const { getFieldDecorator } = this.props.form
    return (
      <Form>
        <Form.Item label="Client ID" extra="The OAuth client ID generated by Github.">
          {getFieldDecorator('clientID', {
            rules: [{ required: true, message: 'Please enter the client ID' }],
          })(
            <Input placeholder="Client ID" />
          )}
        </Form.Item>
        <Form.Item label="Client secret" extra="The OAuth client secret generated by Github.">
          {getFieldDecorator('clientSecret', {
            rules: [{ required: true, message: 'Please enter the client secret' }],
          })(
            <Input placeholder="Client secret" />
          )}
        </Form.Item>
      </Form>
    )
  }
}

const WrappedConfigureHubForm = Form.create({ name: 'configure_hub' })(ConfigureHubForm)

class GithubSetupPage extends React.Component {
  static staticProps = {
    title: 'Configure authentication provider',
    hideSider: true,
    authUrl: auth.url,
    authCallbackUrl: auth.callbackUrl,
    adminOnly: true
  }

  static propTypes = {
    authUrl: PropTypes.string.isRequired,
    authCallbackUrl: PropTypes.string.isRequired
  }

  constructor(props) {
    super(props)
    this.state = {
      current: 0,
      configureHubErrorMessage: false,
      configureHubSubmitting: false
    }
  }

  next() {
    this.steps[this.state.current].process()
  }

  prev() {
    const state = copy(this.state)
    state.current = this.state.current - 1
    this.setState(state)
  }

  async configure({ clientID, clientSecret }) {
    const state = copy(this.state)
    state.configureHubSubmitting = true
    this.setState(state)
    try {
      const body = {
        displayName: 'GitHub',
        name: 'github',
        config: {
          clientID,
          clientSecret
        }
      }
      await axios.post(`${window.location.origin}/login/auth/configure`, body)
      return redirect(null, '/setup/authentication/github/complete')
    } catch (err) {
      console.error('Error submitting form', err)
      const state = copy(this.state)
      state.configureHubSubmitting = false
      state.configureHubErrorMessage = 'There was a problem saving the configuration, please try again.'
      this.setState(state)
    }
  }

  steps = [
    {
      title: 'Configure Github',
      content: () => (
        <div>
          <Alert
            message="Firstly, you need to complete the setup of your Github organisation to allow it to be used as authentication for Kore. The following instructions will guide you through this."
            type="info"
          />
          <Card style={{ marginTop: '20px' }}>
            <Alert
              showIcon
              message="You must be an administrator of your GitHub organisation to complete this"
              type="warning"
              style={{ marginBottom: '20px' }}
            />
            <Paragraph>You need to create an OAuth App within your GitHub organisation.</Paragraph>
            <Paragraph>Github&apos;s instructions can be found at <a href="https://developer.github.com/apps/building-oauth-apps/creating-an-oauth-app/">https://developer.github.com/apps/building-oauth-apps/creating-an-oauth-app/</a></Paragraph>
            <Paragraph>Navigate to the <Text strong>Settings</Text> tab under your GitHub organisation.</Paragraph>
            <Paragraph>Choose <Text strong>OAuth Apps</Text> under <Text strong>Developer settings</Text> on the left side menu.</Paragraph>
            <Paragraph>Click the <Text strong>New OAuth App</Text> button and complete the form using the following to help.</Paragraph>
            <Paragraph>
              <CopyTextWithLabel label="Application name" text="Appvia Kore" />
              <CopyTextWithLabel label="Homepage URL" text={this.props.authUrl} />
              <CopyTextWithLabel label="Authorization callback URL" text={this.props.authCallbackUrl} />
            </Paragraph>
            <Paragraph>Click <Text strong>Register application</Text> once complete to reveal the client credentials.</Paragraph>
          </Card>
        </div>
      ),
      process: () => {
        const state = copy(this.state)
        state.current = this.state.current + 1
        this.setState(state)
      }
    },
    {
      title: 'Configure Appvia Kore',
      content: () => (
        <div>
          <Alert
            message="Enter the credentials generated by Github in the previous step below. This will enable Kore to use Github to authenticate users in your Organisation."
            type="info"
          />
          <Card style={{ marginTop: '20px' }}>
            {this.state.configureHubErrorMessage ? (
              <Alert
                message={this.state.configureHubErrorMessage}
                type="error"
                showIcon
                closable
                style={{ marginBottom: '20px'}}
              />
            ) : null}
            <WrappedConfigureHubForm wrappedComponentRef={(inst) => this.configureHubFormRef = inst} />
          </Card>
        </div>
      ),
      process: () => {
        this.configureHubFormRef.props.form.validateFields(async (err, values) => {
          if (!err) {
            await this.configure(values)
          }
        })
      }
    }
  ]

  render() {
    const { current } = this.state
    return (
      <div style={{ padding: '10px 50px' }}>
        <Title>Configuring Github OpenID Connect authentication</Title>
        <div>
          <Steps current={current}>
            {this.steps.map(item => (
              <Step key={item.title} title={item.title} />
            ))}
          </Steps>
          <div className="steps-content" style={{ padding: '30px 0' }}>{this.steps[current].content()}</div>
          <div className="steps-action">
            {current < this.steps.length - 1 && (
              <Button type="primary" onClick={() => this.next()}>Next</Button>
            )}
            {current === this.steps.length - 1 && (
              <Button type="primary" loading={this.state.configureHubSubmitting} onClick={() => this.next()}>Save</Button>
            )}
            {current > 0 && (
              <Button style={{ marginLeft: 8 }} onClick={() => this.prev()}>Previous</Button>
            )}
          </div>
        </div>
      </div>
    )
  }
}
export default GithubSetupPage
