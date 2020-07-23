import React from 'react'
import PropTypes from 'prop-types'
import { Alert, Form, Icon, Input, Modal, Typography } from 'antd'
const { Paragraph } = Typography

import { patterns } from '../../../utils/validation'

class RequestCredentialAccessForm extends React.Component {
  static propTypes = {
    form: PropTypes.object.isRequired,
    cloud: PropTypes.oneOf(['GCP', 'AWS', 'Azure']).isRequired,
    data: PropTypes.object,
    onChange: PropTypes.func,
    helpInModal: PropTypes.bool
  }

  cloudContent = {
    'GCP': {
      accountNoun: 'Project',
      help: (
        <div>
          <p>When using Kore with existing GCP projects, you must allocate the project credentials to teams in order for them to provision clusters within those projects.</p>
          <p style={{ marginBottom: '0' }}>When a new team is created they may not have access to any project credentials, here you can provide an email address which will be displayed to a team in this situation, in order to request access to a GCP project through Kore.</p>
        </div>
      )
    },
    'AWS': {
      accountNoun: 'Account',
      help: (
        <div>
          <p>When using Kore with existing AWS accounts, you must allocate the account credentials to teams in order for them to provision clusters within those accounts.</p>
          <p style={{ marginBottom: '0' }}>When a new team is created they may not have access to any account credentials, here you can provide an email address which will be displayed to a team in this situation, in order to request access to an AWS account through Kore.</p>
        </div>
      )
    },
    'Azure': {
      accountNoun: 'Subscription',
      help: (
        <div>
          <p>When using Kore with existing Azure subscriptions, you must allocate the subscriptions credentials to teams in order for them to provision clusters within those accounts.</p>
          <p style={{ marginBottom: '0' }}>When a new team is created they may not have access to any subscriptions credentials, here you can provide an email address which will be displayed to a team in this situation, in order to request access to an Azure subscription through Kore.</p>
        </div>
      )
    }
  }

  onChange = () => this.props.onChange && this.props.onChange(this.props.form.getFieldValue('email'), this.props.form.getFieldError('email'))

  projectCredentialAccessHelp = () => {
    Modal.info({
      title: 'Team access',
      content: this.cloudContent[this.props.cloud].help,
      onOk() {},
      width: 500
    })
  }

  render() {
    const { form, helpInModal, data } = this.props
    const { getFieldDecorator, isFieldTouched, getFieldError } = form
    const content = this.cloudContent[this.props.cloud]

    const emailError = isFieldTouched('email') && getFieldError('email')

    return (
      <>
        {helpInModal ? (
          <Paragraph style={{ fontSize: '16px', fontWeight: '600' }}>{content.accountNoun} credential access for teams <Icon style={{ marginLeft: '5px' }} type="info-circle" theme="twoTone" onClick={this.projectCredentialAccessHelp}/></Paragraph>
        ) : (
          <>
            <Paragraph style={{ fontSize: '16px', fontWeight: '600' }}>{content.accountNoun} credential access for teams</Paragraph>
            <Alert
              message="Team access"
              description={content.help}
              type="info"
              showIcon
              style={{ marginBottom: '20px' }}
            />
          </>
        )}
        <Form>
          <Form.Item
            labelAlign="left"
            labelCol={{ span: 2 }}
            wrapperCol={{ span: 8 }}
            label="Email"
            validateStatus={emailError ? 'error' : ''}
            help={emailError || `Email for teams who need access to ${content.accountNoun.toLowerCase()} credentials`}
            onChange={this.onChange}
          >
            {getFieldDecorator('email', {
              rules: [{ required: true, message: 'Please enter the email!' }, { ...patterns.email }],
              initialValue: data && data.email
            })(
              <Input type="email" placeholder="Email address" />
            )}
          </Form.Item>
        </Form>
      </>
    )
  }
}
const WrappedRequestCredentialAccessForm = Form.create({ name: 'request_credential_access_form' })(RequestCredentialAccessForm)

export default WrappedRequestCredentialAccessForm
