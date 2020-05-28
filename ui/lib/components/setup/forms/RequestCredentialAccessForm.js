import React from 'react'
import PropTypes from 'prop-types'
import { Alert, Form, Icon, Input, Modal, Typography } from 'antd'
const { Paragraph } = Typography

import { patterns } from '../../../utils/validation'

class RequestCredentialAccessForm extends React.Component {
  static ENABLED = false

  static propTypes = {
    form: PropTypes.object.isRequired,
    cloud: PropTypes.oneOf(['GCP', 'AWS']).isRequired,
    onChange: PropTypes.func,
    helpInModal: PropTypes.bool
  }

  cloudContent = {
    'GCP': {
      accountNoun: 'Project',
      help: RequestCredentialAccessForm.ENABLED ? (
        <div>
          <p>When using Kore with existing GCP projects, you must allocate the project credentials to teams in order for them to provision clusters within those projects.</p>
          <p style={{ marginBottom: '0' }}>When a new team is created they may not have access to any project credentials, here you can provide an email address which will be displayed to a team in this situation, in order to request access to a GCP project through Kore.</p>
        </div>
      ) : (
        <div>
          <p>When using Kore with existing GCP projects, you must allocate the project credentials to teams in order for them to provision clusters within those projects.</p>
          <p style={{ marginBottom: '0' }}>When a new team is created they may not have access to any project credentials. In this case, the user will be asked to contact the Kore administrator to setup the required access.</p>
        </div>
      )
    },
    'AWS': {
      accountNoun: 'Account',
      help: RequestCredentialAccessForm.ENABLED ? (
        <div>
          <p>When using Kore with existing AWS accounts, you must allocate the account credentials to teams in order for them to provision clusters within those accounts.</p>
          <p style={{ marginBottom: '0' }}>When a new team is created they may not have access to any account credentials, here you can provide an email address which will be displayed to a team in this situation, in order to request access to an AWS account through Kore.</p>
        </div>
      ) : (
        <div>
          <p>When using Kore with existing AWS accounts, you must allocate the account credentials to teams in order for them to provision clusters within those accounts.</p>
          <p style={{ marginBottom: '0' }}>When a new team is created they may not have access to any account credentials. In this case, the user will be asked to contact the Kore administrator to setup the required access.</p>
        </div>
      )
    }
  }

  onChange = () => this.props.onChange && this.props.onChange(this.props.form.getFieldError('email'))

  projectCredentialAccessHelp = () => {
    Modal.info({
      title: 'Team access',
      content: this.cloudContent[this.props.cloud].help,
      onOk() {},
      width: 500
    })
  }

  render() {
    const { form, helpInModal } = this.props
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
        {RequestCredentialAccessForm.ENABLED && (
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
              {getFieldDecorator('email', { rules: [{ required: true, message: 'Please enter the email!' }, { ...patterns.email }] })(
                <Input type="email" placeholder="Email address" />
              )}
            </Form.Item>
          </Form>
        )}
      </>
    )
  }
}
const WrappedRequestCredentialAccessForm = Form.create({ name: 'request_credential_access_form' })(RequestCredentialAccessForm)

export default WrappedRequestCredentialAccessForm
