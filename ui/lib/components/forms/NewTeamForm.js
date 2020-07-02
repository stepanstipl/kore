import React from 'react'
import PropTypes from 'prop-types'
import canonical from '../../utils/canonical'
import copy from '../../utils/object-copy'
import { Button, Form, Input, Alert, Typography } from 'antd'
const { Paragraph, Text } = Typography

import KoreApi from '../../kore-api'
import { successMessage } from '../../utils/message'
import FormErrorMessage from './FormErrorMessage'

class NewTeamForm extends React.Component {
  static propTypes = {
    form: PropTypes.any.isRequired,
    handleTeamCreated: PropTypes.func.isRequired,
    user: PropTypes.object.isRequired,
    team: PropTypes.any
  }

  state = {
    submitting: false,
    formErrorMessage: false
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

  handleSubmit = e => {
    e.preventDefault()

    this.setState({
      submitting: true,
      formErrorMessage: false
    })

    return this.props.form.validateFields(async (err, values) => {
      if (err) {
        return this.setState({
          submitting: false,
          formErrorMessage: 'Validation failed'
        })
      }

      const canonicalTeamName = canonical(values.teamName)
      try {
        const api = await KoreApi.client()
        const checkTeam = await api.GetTeam(canonicalTeamName)
        if (!checkTeam) {
          const team = await api.UpdateTeam(canonicalTeamName, KoreApi.resources().generateTeamResource(values))
          await this.props.handleTeamCreated(team)
          this.setState({ submitting: false })
          successMessage('Team created')
        } else {
          const state = copy(this.state)
          state.submitting = false
          state.formErrorMessage = `A team with the name "${values.teamName}" already exists`
          this.setState(state)
        }
      } catch (err) {
        console.error('Error submitting form', err)
        this.setState({
          submitting: false,
          formErrorMessage: 'An error occurred creating the team, please try again'
        })
      }
    })
  }

  fieldError = (name) => this.props.form.isFieldTouched(name) && this.props.form.getFieldError(name)

  render() {
    const { getFieldDecorator, getFieldsError, getFieldValue } = this.props.form
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
      <Form {...formConfig} onSubmit={this.handleSubmit} style={{ marginBottom: '20px' }}>
        <FormErrorMessage message={this.state.formErrorMessage} />
        <Form.Item label="Team name" validateStatus={this.fieldError('teamName') ? 'error' : ''} help={this.fieldError('teamName') || ''}>
          {getFieldDecorator('teamName', {
            rules: [{ required: true, message: 'Please enter your team name!' }],
          })(
            <Input placeholder="Team name" disabled={!!this.props.team} />,
          )}
        </Form.Item>
        <Form.Item label="Team description" validateStatus={this.fieldError('teamDescription') ? 'error' : ''} help={this.fieldError('teamDescription') || ''}>
          {getFieldDecorator('teamDescription', {
            rules: [{ required: true, message: 'Please enter your team description!' }],
          })(
            <Input placeholder="Team description" disabled={!!this.props.team} />,
          )}
        </Form.Item>

        <Alert
          message={
            <div>
              <Paragraph>The team ID is: <Text id="team_id" strong>{canonical(getFieldValue('teamName') || '')}</Text></Paragraph>
              <Paragraph style={{ marginBottom: '0' }}>This is how your team will appear when using the Kore CLI.</Paragraph>
            </div>
          }
          type="info"
        />
        {!this.props.team ? (
          <Form.Item style={{ marginTop: '20px' }}>
            <Button id="save" type="primary" htmlType="submit" loading={this.state.submitting} disabled={this.disableButton(getFieldsError())}>Save</Button>
          </Form.Item>
        ) : null}
      </Form>
    )
  }
}

const WrappedNewTeamForm = Form.create({ name: 'new_team' })(NewTeamForm)

export default WrappedNewTeamForm
