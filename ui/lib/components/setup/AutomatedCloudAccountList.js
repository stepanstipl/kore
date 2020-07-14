import React from 'react'
import PropTypes from 'prop-types'
import { Alert, Button, Card, Col, Icon, List, Modal, Popover, Row, Typography } from 'antd'
const { Paragraph, Text, Title } = Typography

import PlanViewer from '../plans/PlanViewer'
import IconTooltip from '../utils/IconTooltip'
import DataField from '../utils/DataField'
import { successMessage } from '../../utils/message'

class AutomatedCloudAccountList extends React.Component {

  static propTypes = {
    automatedCloudAccountList: PropTypes.array.isRequired,
    plans: PropTypes.array.isRequired,
    handleChange: PropTypes.func.isRequired,
    handleDelete: PropTypes.func.isRequired,
    handleEdit: PropTypes.func.isRequired
  }

  state = {
    associatePlanVisible: false
  }

  showPlanDetails = (plan) => () => {
    Modal.info({
      title: <><Title level={4}>{plan.spec.description}</Title><Text>{plan.spec.summary}</Text></>,
      content: <PlanViewer plan={plan} resourceType="cluster" />,
      width: 700,
      onOk() {}
    })
  }

  handleAssociatePlanVisibleChange = (cloudAccountCode) => () => this.setState({ associatePlanVisible: cloudAccountCode })

  associatePlan = (cloudAccountCode, plan) => () => {
    const cloudAccount = this.props.automatedCloudAccountList.find(p => p.code === cloudAccountCode)
    cloudAccount.plans.push(plan)
    this.props.handleChange(cloudAccountCode, 'plans')(cloudAccount.plans)
    successMessage('Plan associated')
  }

  closeAssociatePlans = () => this.setState({ associatePlanVisible: false })

  associatePlanContent = (cloudAccountCode) => {
    const cloudAccount = this.props.automatedCloudAccountList.find(p => p.code === cloudAccountCode)
    const associatedPlans = this.props.automatedCloudAccountList.map(ap => ap.plans).flat()
    const unassociatedPlans = this.props.plans.filter(p => !associatedPlans.includes(p.metadata.name))
    if (unassociatedPlans.length === 0) {
      return (
        <>
          <Alert style={{ marginBottom: '20px' }} message="All cluster plans are already associated." />
          <Button type="primary" onClick={this.closeAssociatePlans}>Close</Button>
        </>
      )
    }
    return (
      <>
        <Alert style={{ marginBottom: '20px' }} message="Cluster plans available to be associated." />
        <List
          dataSource={unassociatedPlans}
          renderItem={plan => <Paragraph>{plan.spec.description} <a style={{ textDecoration: 'underline' }} onClick={this.associatePlan(cloudAccount.code, plan.metadata.name)}>Associate</a></Paragraph>}
        />
        <Button style={{ marginTop: '10px' }} type="primary" onClick={this.closeAssociatePlans}>Close</Button>
      </>
    )
  }

  unassociatePlan = (cloudAccountCode, plan) => () => {
    const cloudAccount = this.props.automatedCloudAccountList.find(p => p.code === cloudAccountCode)
    const updatedPlans = cloudAccount.plans.filter(p => p !== plan)
    this.props.handleChange(cloudAccountCode, 'plans')(updatedPlans)
    successMessage('Plan unassociated')
  }

  automatedCloudAccountContent = ({ cloudAccount }) => (
    <Row gutter={16}>
      <Col span={12}>
        {this.automatedCloudAccountNaming({ cloudAccount })}
      </Col>
      <Col span={12}>
        {this.automatedCloudAccountPlans({ cloudAccount })}
      </Col>
    </Row>
  )

  automatedCloudAccountNaming = ({ cloudAccount }) => (
    <Card
      title="Naming"
      size="small"
      bordered={false}
    >
      <Paragraph>The project will be named using the team name, with the prefix and suffix below</Paragraph>
      <DataField label="Prefix" value={cloudAccount.prefix} />
      <DataField label="Suffix" value={cloudAccount.suffix} />
      <DataField label="Example" value={<Text>{cloudAccount.prefix ? `${cloudAccount.prefix}-` : ''}<span style={{ fontStyle: 'italic' }}>team-name</span>-{cloudAccount.suffix}</Text>} style={{ paddingTop: '15px' }} />
    </Card>
  )

  automatedCloudAccountPlans = ({ cloudAccount }) => (
    <Card
      title="Cluster plans"
      size="small"
      bordered={false}
    >
      <Paragraph>The cluster plans associated with this project.</Paragraph>
      {cloudAccount.plans.length === 0 ? <Text type="warning" style={{ padding: '5px 0' }}>No plans</Text> : null}
      {(this.props.plans || []).filter(p => cloudAccount.plans.includes(p.metadata.name)).map((plan, i) => (
        <div key={i} style={{ padding: '5px 0' }}>
          <Text style={{ marginRight: '10px' }}>{plan.spec.description}</Text>
          <IconTooltip icon="info-circle" text={plan.spec.summary} />
          <IconTooltip icon="eye" text="View plan" onClick={this.showPlanDetails(plan)} />
          <IconTooltip icon="delete" text="Unassociate plan" onClick={this.unassociatePlan(cloudAccount.code, plan.metadata.name)} />
        </div>
      ))}
      <div style={{ padding: '5px 0' }}>
        <Popover
          content={this.associatePlanContent(cloudAccount.code)}
          title={`${cloudAccount.name}: Associate plans`}
          trigger="click"
          visible={this.state.associatePlanVisible === cloudAccount.code}
          onVisibleChange={this.handleAssociatePlanVisibleChange(cloudAccount.code)}
        >
          <a>+ Associate plan</a>
        </Popover>
      </div>
    </Card>
  )

  render() {
    return (
      <List
        itemLayout="vertical"
        bordered={true}
        dataSource={this.props.automatedCloudAccountList}
        renderItem={cloudAccount => (
          <List.Item actions={[
            <a key="delete" onClick={this.props.handleDelete(cloudAccount.code)}><Icon type="delete" /> Remove</a>,
            <a key="edit" onClick={this.props.handleEdit(cloudAccount.code)}><Icon type="edit" /> Edit</a>
          ]}>
            <List.Item.Meta
              title={<Text style={{ fontSize: '16px' }}>{cloudAccount.name}</Text>}
              description={<Text>{cloudAccount.description}</Text>}
            />
            {this.automatedCloudAccountContent({ cloudAccount })}
          </List.Item>
        )}
      />
    )
  }

}

export default AutomatedCloudAccountList
