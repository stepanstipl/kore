import * as React from 'react'
import PropTypes from 'prop-types'
import copy from '../../utils/object-copy'
import canonical from '../../utils/canonical'
import GKECredentials from '../../crd/GKECredentials'
import Generic from '../../crd/Generic'
import Allocation from '../../crd/Allocation'
import apiRequest from '../../utils/api-request'
import apiPaths from '../../utils/api-paths'
import { Button, Form, Input, Alert, Select, message, Typography } from 'antd'

const { Paragraph, Text } = Typography

class GCPOrganizationForm extends React.Component {
  static propTypes = {
    form: PropTypes.any.isRequired,
    team: PropTypes.string.isRequired,
    allTeams: PropTypes.object,
    data: PropTypes.object,
    handleSubmit: PropTypes.func.isRequired,
    saveButtonText: PropTypes.string
  }

  constructor(props) {
    super(props)
    let allocations = []
    if (props.data && props.data.allocation) {
      allocations = props.data.allocation.spec.teams.filter(a => a !== '*')
    }
    this.state = {
      submitting: false,
      formErrorMessage: false,
      allocations
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

  onAllocationsChange = value => {
    const state = copy(this.state)
    state.allocations = value
    this.setState(state)
  }

  handleSubmit = e => {
    e.preventDefault()

    const state = copy(this.state)
    state.submitting = true
    state.formErrorMessage = false
    this.setState(state)

    return this.props.form.validateFields(async (err, values) => {
      if (!err) {
        try {
          const data = this.props.data
          const canonicalName = (data && data.metadata && data.metadata.name) || canonical(values.name)

          const gcpServiceAccountSecret = Generic({
            apiVersion: 'config.kore.appvia.io',
            kind: 'Secret',
            name: canonicalName,
            spec: {
              type: 'credentials',
              description: `GCP admin project Service Account for ${values.parentID}`,
              data: {
                key: values.account
              }
            }
          })

          const gcpOrgResource = Generic({
            apiVersion: 'gcp.compute.kore.appvia.io/v1alpha1',
            kind: 'Organization',
            name: canonicalName,
            spec: {
              parentType: 'organization',
              parentID: values.parentID,
              billingAccount: values.billingAccount,
              serviceAccount: values.serviceAccount,
              credentialsRef: {
                name: canonicalName,
                namespace: this.props.team
              }
            }
          })

          console.log('Creating secret', this.props.team, gcpServiceAccountSecret)

          await apiRequest(null, 'put', `${apiPaths.team(this.props.team).secrets}/${canonicalName}`, gcpServiceAccountSecret)
          const gcpOrgResult = await apiRequest(null, 'put', `${apiPaths.team(this.props.team).gcpOrganizations}/${canonicalName}`, gcpOrgResource)

          const allocationSpec = {
            name: values.name,
            summary: values.summary,
            resource: {
              group: 'gcp.compute.kore.appvia.io',
              version: 'v1alpha1',
              kind: 'Organization',
              namespace: this.props.team,
              name: canonicalName
            },
            teams: this.state.allocations.length > 0 ? this.state.allocations : ['*']
          }
          const allocationResult = await apiRequest(null, 'put', `${apiPaths.team(this.props.team).allocations}/${canonicalName}`, Allocation(canonicalName, allocationSpec))
          gcpOrgResult.allocation = allocationResult
          await this.props.handleSubmit(gcpOrgResult)

        } catch (err) {
          console.error('Error submitting form', err)
          const state = copy(this.state)
          state.submitting = false
          state.formErrorMessage = 'An error occurred saving the organization, please try again'
          this.setState(state)
        }
      }
    })
  }

  render() {
    const { form, data, allTeams, saveButtonText } = this.props
    const { getFieldDecorator, getFieldsError, getFieldError, isFieldTouched } = form
    const { formErrorMessage, allocations, submitting } = this.state
    const formConfig = {
      layout: 'horizontal',
      labelAlign: 'left',
      hideRequiredMark: true,
      labelCol: {
        sm: { span: 24 },
        md: { span: 8 },
        lg: { span: 6 }
      },
      wrapperCol: {
        sm: { span: 24 },
        md: { span: 16 },
        lg: { span: 18 }
      }
    }

    // Only show error after a field is touched.
    const fieldError = fieldKey => isFieldTouched(fieldKey) && getFieldError(fieldKey)
    const allocationMissing = Boolean(data && !data.allocation)

    const FormErrorMessage = () => {
      if (formErrorMessage) {
        return (
          <Alert
            message={formErrorMessage}
            type="error"
            showIcon
            closable
            style={{ marginBottom: '20px'}}
          />
        )
      }
      return null
    }

    return (
      <div>
        {allocationMissing ? (
          <Alert
            message="These credentials are not allocated to any teams"
            description="Enter Allocated team(s), Name and Description below and click Save to allocate these Credentials."
            type="warning"
            showIcon
            style={{ marginBottom: '20px', marginTop: '5px' }}
          />
        ) : null}

        {allTeams.items.length ? (
          <Form>
            <Form.Item label="Allocate team(s)" extra="If nothing selected then this integration will be available to ALL teams">
              {getFieldDecorator('allocations', { initialValue: allocations })(
                <Select
                  mode="multiple"
                  style={{ width: '100%' }}
                  placeholder={allocationMissing ? 'No teams' : 'All teams'}
                  onChange={this.onAllocationsChange}
                >
                  {allTeams.items.map(t => (
                    <Select.Option key={t.metadata.name} value={t.metadata.name}>{t.spec.summary}</Select.Option>
                  ))}
                </Select>
              )}
            </Form.Item>
          </Form>
        ) : (
          <Alert
            message="These credentials will be allocated to all teams"
            description="No teams exist in Kore yet, therefore currently these credentials will be available to all teams created in the future. If you wish to restrict this please return here and allocate to teams once they have been created."
            type="info"
            showIcon
            style={{ marginBottom: '20px', marginTop: '5px' }}
          />
        )}

        <Form {...formConfig} onSubmit={this.handleSubmit}>
          <FormErrorMessage />
          <Form.Item label="Name" validateStatus={fieldError('name') ? 'error' : ''} help={fieldError('name') || 'The name for your credentials eg. MyOrg GKE'}>
            {getFieldDecorator('name', {
              rules: [{ required: true, message: 'Please enter the name!' }],
              initialValue: data && data.allocation && data.allocation.spec.name
            })(
              <Input placeholder="Name" />,
            )}
          </Form.Item>
          <Form.Item label="Description" validateStatus={fieldError('summary') ? 'error' : ''} help={fieldError('summary') || 'A description of your credentials to help when choosing between them'}>
            {getFieldDecorator('summary', {
              rules: [{ required: true, message: 'Please enter the description!' }],
              initialValue: data && data.allocation && data.allocation.spec.summary
            })(
              <Input placeholder="Description" />,
            )}
          </Form.Item>
          <Form.Item label="Parent ID" validateStatus={fieldError('parentID') ? 'error' : ''} help={fieldError('parentID') || 'The GCP organization ID'}>
            {getFieldDecorator('parentID', {
              rules: [{ required: true, message: 'Please enter the parent ID!' }],
              initialValue: data && data.spec.parentID
            })(
              <Input placeholder="Parent ID" />,
            )}
          </Form.Item>
          <Form.Item label="Billing account" validateStatus={fieldError('billingAccount') ? 'error' : ''} help={fieldError('billingAccount') || 'The billing account'}>
            {getFieldDecorator('billingAccount', {
              rules: [{ required: true, message: 'Please enter your billing account!' }],
              initialValue: data && data.spec.billingAccount
            })(
              <Input placeholder="Billing account" />,
            )}
          </Form.Item>
          <Form.Item label="Service account name" validateStatus={fieldError('serviceAccount') ? 'error' : ''} help={fieldError('serviceAccount') || 'The name of the Service Account that will be created'}>
            {getFieldDecorator('serviceAccount', {
              rules: [{ required: true, message: 'Please enter the Service Account name!' }],
              initialValue: data && data.spec.serviceAccount
            })(
              <Input placeholder="Service account name" />,
            )}
          </Form.Item>
          <Form.Item label="Service Account JSON" validateStatus={fieldError('account') ? 'error' : ''} help={fieldError('account') || 'The Service Account JSON with admin permissions on the GCP admin project'}>
            {!data ? (
              getFieldDecorator('account', {
                rules: [{ required: true, message: 'Please enter your service account JSON!' }]
              })(
                <Input.TextArea autoSize={{ minRows: 4, maxRows: 10  }} placeholder="Account JSON" />,
              )
            ) : (
              <Text type="warning">This cannot be seen once it's been created</Text>
            )}
          </Form.Item>
          <Form.Item style={{ marginBottom: '0'}}>
            <Button type="primary" htmlType="submit" loading={submitting} disabled={this.disableButton(getFieldsError())}>{saveButtonText || 'Save'}</Button>
          </Form.Item>
        </Form>
      </div>
    )
  }
}

const WrappedGCPOrganizationForm = Form.create({ name: 'gcp_organization' })(GCPOrganizationForm)

export default WrappedGCPOrganizationForm
