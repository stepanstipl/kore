import * as React from 'react'
import axios from 'axios'
import PropTypes from 'prop-types'
import Router from 'next/router'
import { Typography, Steps, Button, Alert, Row, Col, Form, Input, Card, List, Icon } from 'antd'
const { Step } = Steps
const { Title, Paragraph, Text } = Typography
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

import CopyTextWithLabel from '../../../../lib/components/utils/CopyTextWithLabel'
import copy from '../../../../lib/utils/object-copy'
import redirect from '../../../../lib/utils/redirect'

class ConfigureKoreForm extends React.Component {
  static propTypes = {
    form: PropTypes.any.isRequired
  }
  render() {
    const { getFieldDecorator } = this.props.form
    return (
      <Form>
        <Form.Item label="Client ID" extra="The OAuth client ID generated by Google.">
          {getFieldDecorator('clientID', {
            rules: [{ required: true, message: 'Please enter the client ID' }],
          })(
            <Input placeholder="Client ID" />
          )}
        </Form.Item>
        <Form.Item label="Client secret" extra="The OAuth client secret generated by Google.">
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

const WrappedConfigureKoreForm = Form.create({ name: 'configure_kore' })(ConfigureKoreForm)

class GoogleSetupPage extends React.Component {
  static staticProps = {
    title: 'Configure authentication provider',
    hideSider: true,
    authUrl: publicRuntimeConfig.authOpenidUrl,
    authCallbackUrl: publicRuntimeConfig.authOpenidCallbackUrl,
    baseUrl: publicRuntimeConfig.koreBaseUrl,
    adminOnly: true
  }

  static propTypes = {
    authUrl: PropTypes.string.isRequired,
    authCallbackUrl: PropTypes.string.isRequired,
    baseUrl: PropTypes.string.isRequired
  }

  constructor(props) {
    super(props)
    this.state = {
      current: 0,
      configureKoreErrorMessage: false,
      configureKoreSubmitting: false,
      setupGoogle: [
        {
          step: 1,
          title: 'Create project',
          content: () => (
            <div>
              <Paragraph>Create a project for Kore on your Google account.</Paragraph>
              <Paragraph>Once created navigate to the <Text strong>APIs & Services</Text> screen for your project.</Paragraph>
            </div>
          ),
          complete: false
        },
        {
          step: 2,
          title: 'Configure consent screen',
          content: () => (
            <div>
              <Paragraph>Google instructions can be found at <a href="https://developers.google.com/identity/protocols/OpenIDConnect#consentpageexperience">https://developers.google.com/identity/protocols/OpenIDConnect#consentpageexperience</a></Paragraph>
              <Paragraph>Choose <Text strong>OAuth consent screen</Text> from the left side menu.</Paragraph>
              <Paragraph>Choose <Text strong>User Type</Text> &quot;Internal&quot; to ensure access is locked down to your Organisation only.</Paragraph>
              <Paragraph>Complete the rest of the OAuth consent screen using the following to help.</Paragraph>
              <Paragraph>
                <CopyTextWithLabel label="Application name" text="Kore" />
                <div style={{ padding: '5px 0' }}>
                  <Row>
                    <Col xs={24} sm={12} md={9} lg={7} xl={5}>
                      <Text strong>Authorised domains</Text>
                    </Col>
                    <Col>
                      <Text>Enter your top level domain</Text>
                    </Col>
                  </Row>
                </div>
                <CopyTextWithLabel label="Application homepage link" text={this.props.baseUrl} />
                <CopyTextWithLabel label="Application Privacy Policy link" text={this.props.baseUrl} />
              </Paragraph>
              <Paragraph>Click <Text strong>Save</Text> once complete.</Paragraph>
            </div>
          ),
          complete: false
        },
        {
          step: 3,
          title: 'Create OAuth client ID credential',
          content: () => (
            <div>
              <Paragraph>Google instructions can be found at <a href="https://developers.google.com/identity/protocols/OpenIDConnect#getcredentials">https://developers.google.com/identity/protocols/OpenIDConnect#getcredentials</a></Paragraph>
              <Paragraph>Choose <Text strong>Credentials</Text> from the left side menu.</Paragraph>
              <Paragraph>Click <Text strong>Create credentials</Text> &gt; <Text strong>OAuth client ID</Text>.</Paragraph>
              <Paragraph>Complete the <Text strong>Create OAuth client ID</Text> form using the following to help.</Paragraph>
              <Paragraph>
                <CopyTextWithLabel label="Name" text="Kore" />
                <CopyTextWithLabel label="Authorised redirect URIs" text={this.props.authCallbackUrl} />
              </Paragraph>
              <Paragraph>Click <Text strong>Create</Text> once completed to reveal the client credentials.</Paragraph>
            </div>
          ),
          complete: false
        }
      ]
    }
  }

  disableNextButton = () => {
    if (this.steps[this.state.current]) {
      return !this.steps[this.state.current].complete()
    }
    return false
  }

  next() {
    this.steps[this.state.current].process()
  }

  prev() {
    const state = copy(this.state)
    state.current = this.state.current - 1
    this.setState(state)
  }

  setStepComplete(stepNumber, value) {
    return () => {
      const state = copy(this.state)
      const step = state.setupGoogle.find(s => s.step === stepNumber)
      if (value === false) {
        state.setupGoogle.filter(s => s.step > step.step).forEach(s => s.complete = false)
      }
      step.complete = value
      this.setState(state)
    }
  }

  async configure({ clientID, clientSecret }) {
    const state = copy(this.state)
    state.configureKoreSubmitting = true
    this.setState(state)
    try {
      const body = {
        displayName: 'Google',
        name: 'google',
        config: {
          clientID,
          clientSecret
        }
      }
      await axios.post(`${window.location.origin}/login/auth/configure`, body)
      return redirect({
        router: Router,
        path: '/setup/authentication/google/complete'
      })
    } catch (err) {
      console.error('Error submitting form', err)
      const state = copy(this.state)
      state.configureKoreSubmitting = false
      state.configureKoreErrorMessage = 'There was a problem saving the configuration, please try again.'
      this.setState(state)
    }
  }

  steps = [
    {
      title: 'Configure Google Account',
      content: () => (
        <div>
          <Alert
            message="Firstly, you need to complete the setup of your Google organisation to allow it to be used as authentication for Kore. The following instructions will guide you through this."
            type="info"
          />
          <Card style={{ marginTop: '20px' }}>
            <Alert
              showIcon
              message="You must be an administrator of your Google organisation to complete this"
              type="warning"
              style={{ marginBottom: '20px' }}
            />
            <List
              dataSource={this.state.setupGoogle.filter(s => s.step === 1 || this.state.setupGoogle.find(i => i.step === s.step-1).complete)}
              renderItem={item => (
                <List.Item>
                  <Row style={{ width: '100%' }}>
                    <Col span={22}>
                      <Typography.Paragraph style={{ fontSize: '16px' }} strong>{item.title}</Typography.Paragraph>
                      {!item.complete ? (
                        <div>
                          {item.content()}
                          <a onClick={this.setStepComplete(item.step, true)}><Icon style={{ marginRight: '5px' }} type="check-circle" /> I&apos;ve done this</a>
                        </div>
                      ) : null}
                    </Col>
                    <Col span={2}>
                      {item.complete? (
                        <div>
                          <Icon style={{ fontSize: '20px', marginRight: '10px' }} type="check-circle" theme="twoTone" twoToneColor="#52c41a" />
                          <a onClick={this.setStepComplete(item.step, false)}>Redo</a>
                        </div>
                      ) : null}
                    </Col>
                  </Row>
                </List.Item>
              )}
            />
          </Card>
        </div>
      ),
      process: () => {
        const state = copy(this.state)
        state.current = this.state.current + 1
        this.setState(state)
      },
      complete: () => {
        const incomplete = this.state.setupGoogle.filter(s => !s.complete)
        return incomplete.length === 0 ? true : false
      }
    },
    {
      title: 'Configure Kore',
      content: () => (
        <div>
          <Alert
            message="Enter the credentials generated by Google in the previous step below. This will enable Kore to use Google to authenticate users in your Organisation."
            type="info"
          />
          <Card style={{ marginTop: '20px' }}>
            <WrappedConfigureKoreForm wrappedComponentRef={(inst) => this.configureKoreFormRef = inst} />
          </Card>
        </div>
      ),
      process: () => {
        this.configureKoreFormRef.props.form.validateFields(async (err, values) => {
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
        <Title>Configuring Google OpenID Connect authentication</Title>
        <div>
          <Steps current={current}>
            {this.steps.map(item => (
              <Step key={item.title} title={item.title} />
            ))}
          </Steps>
          <div className="steps-content" style={{ padding: '30px 0' }}>{this.steps[current].content()}</div>
          <div className="steps-action">
            {current < this.steps.length - 1 && (
              <Button type="primary" disabled={this.disableNextButton(current)} onClick={() => this.next()}>Next</Button>
            )}
            {current === this.steps.length - 1 && (
              <Button type="primary" loading={this.state.configureKoreSubmitting} onClick={() => this.next()}>Save</Button>
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
export default GoogleSetupPage
