import React from 'react'
import PropTypes from 'prop-types'
import { Alert, Button, Card, Col, Icon, List, message, Modal, Popover, Row, Typography } from 'antd'
const { Paragraph, Text, Title } = Typography

import PlanViewer from '../plans/PlanViewer'
import IconTooltip from '../utils/IconTooltip'
import DataField from '../utils/DataField'

class GCPAutomatedProjectList extends React.Component {

  static propTypes = {
    automatedProjectList: PropTypes.array.isRequired,
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

  handleAssociatePlanVisibleChange = (projectCode) => () => this.setState({ associatePlanVisible: projectCode })

  associatePlan = (projectCode, plan) => () => {
    const project = this.props.automatedProjectList.find(p => p.code === projectCode)
    project.plans.push(plan)
    this.props.handleChange(projectCode, 'plans')(project.plans)
    message.success('Plan associated')
  }

  closeAssociatePlans = () => this.setState({ associatePlanVisible: false })

  associatePlanContent = (projectCode) => {
    const project = this.props.automatedProjectList.find(p => p.code === projectCode)
    const associatedPlans = this.props.automatedProjectList.map(ap => ap.plans).flat()
    const unassociatedPlans = this.props.plans.filter(p => !associatedPlans.includes(p.metadata.name))
    if (unassociatedPlans.length === 0) {
      return (
        <>
          <Alert style={{ marginBottom: '20px' }} message="All cluster plans are already associated with a project type." />
          <Button type="primary" onClick={this.closeAssociatePlans}>Close</Button>
        </>
      )
    }
    return (
      <>
        <Alert style={{ marginBottom: '20px' }} message="Cluster plans available to be associated with this project type." />
        <List
          dataSource={unassociatedPlans}
          renderItem={plan => <Paragraph>{plan.spec.description} <a style={{ textDecoration: 'underline' }} onClick={this.associatePlan(project.code, plan.metadata.name)}>Associate</a></Paragraph>}
        />
        <Button style={{ marginTop: '10px' }} type="primary" onClick={this.closeAssociatePlans}>Close</Button>
      </>
    )
  }

  unassociatePlan = (projectCode, plan) => () => {
    const project = this.props.automatedProjectList.find(p => p.code === projectCode)
    const updatedPlans = project.plans.filter(p => p !== plan)
    this.props.handleChange(projectCode, 'plans')(updatedPlans)
    message.success('Plan unassociated')
  }

  automatedProjectContent = ({ project }) => (
    <Row gutter={16}>
      <Col span={12}>
        {this.automatedProjectNaming({ project })}
      </Col>
      <Col span={12}>
        {this.automatedProjectPlans({ project })}
      </Col>
    </Row>
  )

  automatedProjectNaming = ({ project }) => (
    <Card
      title="Naming"
      size="small"
      bordered={false}
    >
      <Paragraph>The project will be named using the team name, with the prefix and suffix below</Paragraph>
      <DataField label="Prefix" value={project.prefix} />
      <DataField label="Suffix" value={project.suffix} />
      <DataField label="Example" value={<Text>{project.prefix ? `${project.prefix}-` : ''}<span style={{ fontStyle: 'italic' }}>team-name</span>-{project.suffix}</Text>} style={{ paddingTop: '15px' }} />
    </Card>
  )

  automatedProjectPlans = ({ project }) => (
    <Card
      title="Cluster plans"
      size="small"
      bordered={false}
    >
      <Paragraph>The cluster plans associated with this project.</Paragraph>
      {project.plans.length === 0 ? <Text type="warning" style={{ padding: '5px 0' }}>No plans</Text> : null}
      {(this.props.plans || []).filter(p => project.plans.includes(p.metadata.name)).map((plan, i) => (
        <div key={i} style={{ padding: '5px 0' }}>
          <Text style={{ marginRight: '10px' }}>{plan.spec.description}</Text>
          <IconTooltip icon="info-circle" text={plan.spec.summary} />
          <IconTooltip icon="eye" text="View plan" onClick={this.showPlanDetails(plan)} />
          <IconTooltip icon="delete" text="Unassociate plan" onClick={this.unassociatePlan(project.code, plan.metadata.name)} />
        </div>
      ))}
      <div style={{ padding: '5px 0' }}>
        <Popover
          content={this.associatePlanContent(project.code)}
          title={`${project.name}: Associate plans`}
          trigger="click"
          visible={this.state.associatePlanVisible === project.code}
          onVisibleChange={this.handleAssociatePlanVisibleChange(project.code)}
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
        dataSource={this.props.automatedProjectList}
        renderItem={project => (
          <List.Item actions={[
            <a key="delete" onClick={this.props.handleDelete(project.code)}><Icon type="delete" /> Remove</a>,
            <a key="edit" onClick={this.props.handleEdit(project.code)}><Icon type="edit" /> Edit</a>
          ]}>
            <List.Item.Meta
              title={<Text style={{ fontSize: '16px' }}>{project.name}</Text>}
              description={<Text>{project.description}</Text>}
            />
            {this.automatedProjectContent({ project })}
          </List.Item>
        )}
      />
    )
  }

}

export default GCPAutomatedProjectList
