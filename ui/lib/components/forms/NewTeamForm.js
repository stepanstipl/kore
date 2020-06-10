import * as React from 'react'
import PropTypes from 'prop-types'
import canonical from '../../utils/canonical'
import copy from '../../utils/object-copy'
import Team from '../../crd/Team'
import { Button, Form, Input, Alert, Typography } from 'antd'
import KoreApi from '../../kore-api'
import { successMessage } from '../../utils/message'

const { Paragraph, Text } = Typography

class NewTeamForm extends React.Component {
  static propTypes = {
    form: PropTypes.any.isRequired,
    handleTeamCreated: PropTypes.func.isRequired,
    user: PropTypes.object.isRequired,
    team: PropTypes.any
  }

  constructor(props) {
    super(props)
    this.state = {
      submitting: false,
      formErrorMessage: false,
    }
  }

  componentDidMount() {
    // To disabled submit button at the beginning.
    this.props.form.validateFields()
  }

  disableButton = fieldsError => {
    if (this.state.submitting) {
      return true
    }
    return Object.keys(fieldsError).some(field => fieldsError[field])
  }

  handleSubmit = e => {
    e.preventDefault()

    this.setState({
      ...this.state,
      submitting: true,
      formErrorMessage: false
    })

    return this.props.form.validateFields(async (err, values) => {
      if (err) {
        this.setState({
          ...this.state,
          submitting: false,
          formErrorMessage: 'Validation failed'
        })
        return
      }
      const canonicalTeamName = canonical(values.teamName)
      const spec = {
        summary: values.teamName.trim(),
        description: values.teamDescription.trim()
      }
      try {
        const api = await KoreApi.client()
        const checkTeam = await api.GetTeam(canonicalTeamName)
        if (!checkTeam) {
          const team = await api.UpdateTeam(canonicalTeamName, Team(canonicalTeamName, spec))
          await this.props.handleTeamCreated(team)
          const state = copy(this.state)
          state.submitting = false
          this.setState(state)
          successMessage('Team created')
        } else {
          const state = copy(this.state)
          state.submitting = false
          state.formErrorMessage = `A team with the name "${values.teamName}" already exists`
          this.setState(state)
        }
      } catch (err) {
        //console.error('Error submitting form', err)
        const state = copy(this.state)
        state.submitting = false
        state.formErrorMessage = 'An error occurred creating the team, please try again'
        this.setState(state)
      }
    })
  }

  render() {
    const { getFieldDecorator, getFieldsError, getFieldError, isFieldTouched, getFieldValue } = this.props.form
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

    // Only show error after a field is touched.
    const teamNameError = isFieldTouched('teamName') && getFieldError('teamName')
    const teamDescriptionError = isFieldTouched('teamDescription') && getFieldError('teamDescription')

    const FormErrorMessage = () => {
      if (this.state.formErrorMessage) {
        return (
          <Alert
            message={this.state.formErrorMessage}
            type="error"
            showIcon
            closable
            style={{ marginBottom: '20px' }}
          />
        )
      }
      return null
    }

    return (
      <Form {...formConfig} onSubmit={this.handleSubmit} style={{ marginBottom: '20px' }}>
        <FormErrorMessage />
        <Form.Item label="Team name" validateStatus={teamNameError ? 'error' : ''} help={teamNameError || ''}>
          {getFieldDecorator('teamName', {
            rules: [{ required: true, message: 'Please enter your team name!' }],
          })(
            <Input placeholder="Team name" disabled={!!this.props.team} />,
          )}
        </Form.Item>
        <Form.Item label="Team description" validateStatus={teamDescriptionError ? 'error' : ''} help={teamDescriptionError || ''}>
          {getFieldDecorator('teamDescription', {
            rules: [{ required: true, message: 'Please enter your team description!' }],
          })(
            <Input placeholder="Team description" disabled={!!this.props.team} />,
          )}
        </Form.Item>

        <Alert
          message={
            <div>
              <Paragraph>The team ID is: <Text strong>{canonical(getFieldValue('teamName') || '')}</Text></Paragraph>
              <Paragraph style={{ marginBottom: '0' }}>This is how your team will appear when using the Kore CLI.</Paragraph>
            </div>
          }
          type="info"
        />
        {!this.props.team ? (
          <Form.Item style={{ marginTop: '20px' }}>
            <Button type="primary" htmlType="submit" loading={this.state.submitting} disabled={this.disableButton(getFieldsError())}>Save</Button>
          </Form.Item>
        ) : null}
      </Form>
    )
  }
}

const WrappedNewTeamForm = Form.create({ name: 'new_team' })(NewTeamForm)

export default WrappedNewTeamForm
