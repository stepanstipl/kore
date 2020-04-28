import * as React from 'react'
import PropTypes from 'prop-types'

import { Typography, Form, Select, Card, Radio, Modal, Input, Collapse } from 'antd'
const { Text, Title } = Typography
const { Option } = Select

import PlanViewer from '../configure/PlanViewer'
import PlanOptionsForm from '../plans/PlanOptionsForm'
import { patterns } from '../../utils/validation'
import KoreApi from '../../kore-api'

class ClusterOptionsForm extends React.Component {
  static propTypes = {
    form: PropTypes.any.isRequired,
    team: PropTypes.object.isRequired,
    selectedCloud: PropTypes.string.isRequired,
    credentials: PropTypes.array.isRequired,
    plans: PropTypes.array.isRequired,
    teamClusters: PropTypes.array.isRequired,
    onPlanOverridden: PropTypes.func,
    validationErrors: PropTypes.array
  }

  componentDidUpdate(prevProps) {
    // Reset the selected plan if the credential changes:
    if (this.props.credentials !== prevProps.credentials) {
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
          getPlanSchema={async () => await (await KoreApi.client()).GetPlanSchema(selectedPlan.spec.kind)}
        />,
        width: 700,
        onOk() {}
      })
    }
  }

  onPlanOverridden = paramValues => {
    if (this.props.onPlanOverridden) {
      this.props.onPlanOverridden(paramValues)
    }
  }

  render() {
    const { getFieldDecorator, getFieldValue } = this.props.form
    const { credentials, plans } = this.props
    const selectedPlan = getFieldValue('plan')

    const checkForDuplicateName = (rule, value) => {
      const matchingCluster = this.props.teamClusters.find(tc => tc.metadata.name === value)
      if (!matchingCluster) {
        return Promise.resolve()
      }
      return Promise.reject('This name is already used!')
    }

    return (
      <Card title="Cluster options">
        <Form.Item label="Credential">
          {getFieldDecorator('credential', {
            rules: [{ required: true, message: 'Please select your credential!' }],
            initialValue: credentials.length === 1 ? credentials[0].metadata.name : undefined
          })(
            <Select placeholder="Credential">
              {credentials.map(c => <Option key={c.metadata.name} value={c.metadata.name}>{c.spec.name} - {c.spec.summary}</Option>)}
            </Select>
          )}
        </Form.Item>
        <Form.Item label="Plan">
          {getFieldDecorator('plan', {
            rules: [{ required: true, message: 'Please select your plan!' }],
          })(
            <Radio.Group onChange={this.onPlanChange}>
              {plans.map((p, idx) => (
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
            <Collapse.Panel header="Customize cluster parameters">
              <PlanOptionsForm
                team={this.props.team}
                resourceType="cluster"
                plan={selectedPlan}
                validationErrors={this.props.validationErrors}
                onPlanChange={this.onPlanOverridden}
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
