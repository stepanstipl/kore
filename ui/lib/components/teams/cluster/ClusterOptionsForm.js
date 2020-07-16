import * as React from 'react'
import PropTypes from 'prop-types'

import { Typography, Form, Select, Card, Radio, Modal, Input, Collapse } from 'antd'
const { Paragraph, Text, Title } = Typography
const { Option } = Select

import PlanViewer from '../../plans/PlanViewer'
import UsePlanForm from '../../plans/UsePlanForm'
import { patterns } from '../../../utils/validation'
import CostEstimate from '../../costs/CostEstimate'

class ClusterOptionsForm extends React.Component {
  static propTypes = {
    form: PropTypes.any.isRequired,
    team: PropTypes.object.isRequired,
    selectedCloud: PropTypes.string.isRequired,
    selectedProvider: PropTypes.string.isRequired,
    credentials: PropTypes.array.isRequired,
    accountManagement: PropTypes.object,
    plans: PropTypes.array.isRequired,
    teamClusters: PropTypes.array.isRequired,
    onPlanValuesChange: PropTypes.func,
    validationErrors: PropTypes.array
  }

  state = {
    cloudAccountType: 'KORE',
    planValuesForEstimate: null
  }

  componentDidUpdate(prevProps) {
    // Reset the selected plan if the credential/accountManagement changes:
    if (this.props.credentials !== prevProps.credentials || this.props.accountManagement !== prevProps.accountManagement) {
      this.props.form.setFieldsValue({ 'plan': null })
    }
  }

  onPlanChange = e => {
    if (this.props.form.getFieldValue('clusterName')) {
      this.props.form.setFieldsValue({ 'clusterName': this.generateClusterName(e.target.value) })
    }
  }

  generateClusterName = selectedPlan => {
    let clusterName = `${this.props.team.metadata.name}-${selectedPlan}`
    const matchingClusters = this.props.teamClusters.filter(tc => tc.metadata.name.indexOf(clusterName) === 0)
    if (matchingClusters.length) {
      clusterName = `${clusterName}-${matchingClusters.length + 1}`
    }
    return clusterName
  }

  showPlanDetails = planName => {
    return () => {
      const selectedPlan = this.props.plans.find(p => p.metadata.name === planName)
      Modal.info({
        title: (<><Title level={4}>{selectedPlan.spec.description}</Title><Text>{selectedPlan.spec.summary}</Text></>),
        content: <PlanViewer
          plan={selectedPlan}
          resourceType="cluster"
        />,
        width: 700,
        onOk() {}
      })
    }
  }

  onPlanValuesChange = paramValues => {
    if (this.props.onPlanValuesChange) {
      this.props.onPlanValuesChange(paramValues)
    }
    this.setState({ planValuesForEstimate: paramValues })
  }

  availablePlans = () => {
    if (!this.props.accountManagement) {
      return this.props.plans
    }
    // TODO: need to be able to filter down the plans to ones which appear in the automation rules
    // currently not possible due to only being able to access the Allocation and not the AccountManagement CRD itself
    // return this.props.plans.filter(plan => rulePlans.includes(plan.metadata.name))
    return this.props.plans
  }

  render() {
    const { getFieldDecorator, getFieldValue } = this.props.form
    const { credentials, accountManagement, selectedCloud, selectedProvider } = this.props
    const selectedPlan = getFieldValue('plan')

    const checkForDuplicateName = (rule, value) => {
      const matchingCluster = this.props.teamClusters.find(tc => tc.metadata.name === value)
      if (!matchingCluster) {
        return Promise.resolve()
      }
      return Promise.reject('This name is already used!')
    }

    const cloudAccountName = { 'GCP': 'Project', 'AWS': 'Account', 'Azure': 'Subscription' }[selectedCloud]

    return (
      <Card title="Cluster options">

        {accountManagement && credentials.length >= 1 ? (
          <Form.Item label={cloudAccountName} style={{ marginBottom: '5px' }}>
            <Radio.Group onChange={(e) => this.setState({ cloudAccountType: e.target.value })} value={this.state.cloudAccountType}>
              <Radio value={'KORE'} style={{ marginRight: '20px' }}>
                <Text strong>Kore managed {cloudAccountName.toLowerCase()}<Text type="secondary"> (recommended)</Text></Text>
                <Paragraph style={{ marginLeft: '24px', marginBottom: '0' }}>Kore will create the required {cloudAccountName.toLowerCase()}</Paragraph>
              </Radio>
              <Radio value={'EXISTING'}>
                <Text strong>Use existing {cloudAccountName.toLowerCase()}</Text>
                <Paragraph style={{ marginLeft: '24px', marginBottom: '0' }}>Specify an existing {cloudAccountName.toLowerCase()} the team has access to</Paragraph>
              </Radio>
            </Radio.Group>
          </Form.Item>
        ) : null}

        {accountManagement && credentials.length === 0 ? (
          <Form.Item label={cloudAccountName}>
            <Text>Kore managed {cloudAccountName}</Text>
          </Form.Item>
        ) : null}

        {!accountManagement || this.state.cloudAccountType === 'EXISTING' ? (
          <Form.Item label={cloudAccountName}>
            {getFieldDecorator('credential', {
              rules: [{ required: true, message: 'Please select your credential!' }],
              initialValue: credentials.length === 1 ? credentials[0].metadata.name : undefined
            })(
              <Select placeholder={cloudAccountName}>
                {credentials.map(c => <Option key={c.metadata.name} value={c.metadata.name}>{c.spec.name} - {c.spec.summary}</Option>)}
              </Select>
            )}
          </Form.Item>
        ) : null}

        <Form.Item label="Plan">
          {getFieldDecorator('plan', {
            rules: [{ required: true, message: 'Please select your plan!' }],
          })(
            <Radio.Group onChange={this.onPlanChange}>
              {this.availablePlans().map((p, idx) => (
                <Radio.Button key={idx} value={p.metadata.name}>{p.spec.description}</Radio.Button>
              ))}
            </Radio.Group>
          )}
          {selectedPlan ?
            <a style={{ marginLeft: '20px' }} onClick={this.showPlanDetails(selectedPlan)}>View plan details</a> :
            null
          }
        </Form.Item>
        {selectedPlan ? (
          <Form.Item label="Cluster name">
            {getFieldDecorator('clusterName', {
              rules: [
                { required: true, message: 'Please enter cluster name!' },
                { ...patterns.uriCompatible40CharMax },
                { validator: checkForDuplicateName }
              ],
              initialValue: this.generateClusterName(selectedPlan)
            })(
              <Input />
            )}
          </Form.Item>
        ) : null}
        {selectedPlan ? (
          <Collapse>
            <Collapse.Panel header="Cluster running cost estimate">
              <CostEstimate
                planValues={this.state.planValuesForEstimate}
                resourceType="cluster"
                kind={selectedProvider}
              />
            </Collapse.Panel>
            <Collapse.Panel header="Customize cluster parameters" forceRender={true}>
              <UsePlanForm
                team={this.props.team}
                resourceType="cluster"
                kind={selectedProvider}
                plan={selectedPlan}
                validationErrors={this.props.validationErrors}
                onPlanValuesChange={this.onPlanValuesChange}
                mode="create"
              />
            </Collapse.Panel>
          </Collapse>
        ) : null}
      </Card>
    )
  }
}

const WrappedClusterOptionsForm = Form.create({ name: 'cluster_options' })(ClusterOptionsForm)

export default WrappedClusterOptionsForm
