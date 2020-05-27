import React from 'react'
import PropTypes from 'prop-types'
import KoreApi from '../../../lib/kore-api'
import { Card, Typography, Icon, Switch, Tooltip, List, Avatar, Button, Drawer, Modal } from 'antd'
import ManageServicePlanForm from '../../../lib/components/plans/ManageServicePlanForm'
import { isReadOnlyCRD } from '../../utils/crd-helpers'
const { Paragraph, Text, Title } = Typography

export default class ServiceKindManage extends React.Component {
  static propTypes = {
    kind: PropTypes.object.isRequired
  }

  static initialState = {
    loading: false,
    kind: null,
    plans: [],
    selectedPlan: null,
    selectedPlanMode: 'view'
  }

  constructor(props) {
    super(props)
    this.state = {
      ...ServiceKindManage.initialState,
      kind: props.kind
    }
  }

  componentDidMountComplete = null
  componentDidMount() {
    this.componentDidMountComplete = this.fetchComponentData()
  }

  componentDidUpdateComplete = null
  componentDidUpdate(prevProps) {
    if (this.props.kind !== prevProps.kind) {
      this.setState({ ...ServiceKindManage.initialState })
      this.componentDidUpdateComplete = this.fetchComponentData()
    }
  }

  fetchComponentData = async () => {
    await this.refreshPlans()
  }

  toggleEnabled = async () => {
    const { kind } = this.state
    const api = await KoreApi.client()
    await api.UpdateServiceKind(kind.metadata.name, { ...kind, spec: { ...kind.spec, enabled: !kind.spec.enabled } })
    await this.refreshKind()
  }

  refreshKind = async () => {
    this.setState({ loading: true })
    const api = await KoreApi.client()
    const kind = await api.GetServiceKind(this.state.kind.metadata.name)
    this.setState({ loading: false, kind })
  }

  refreshPlans = async () => {
    this.setState({ loading: true })
    const api = await KoreApi.client()
    const plans = await api.ListServicePlans(this.state.kind.metadata.name)
    this.setState({ loading: false, plans: plans })
  }

  renderPlan = (plan) => {
    const readonly = isReadOnlyCRD(plan)
    const displayName = plan.spec.displayName || plan.spec.summary 
    const actions = [
      <Text key="view"><a onClick={() => this.viewPlan(plan)}><Icon type="eye" theme="filled"/> View</a></Text>,
      <Text key="edit">
        <Tooltip title={readonly ? 'Read-only' : 'Edit this plan'}>
          <a onClick={readonly ? () => {} : () => this.editPlan(plan)} style={{ color: readonly ? 'lightgray' : null }}><Icon type="edit" theme="filled"/> Edit</a>
        </Tooltip>
      </Text>,
      <Text key="delete">
        <Tooltip title={readonly ? 'Read-only' : 'Delete this plan'}>
          <a onClick={readonly ? () => {} : () => this.deletePlan(plan)} style={{ color: readonly ? 'lightgray' : null }}><Icon type="delete" theme="filled"/> Delete</a>
        </Tooltip>
      </Text>
    ]

    return (
      <List.Item key={plan.metadata.name} actions={actions}>
        <List.Item.Meta 
          avatar={<Avatar icon="build" />} 
          title={displayName} 
          description={plan.spec.description} />
      </List.Item>
    )
  }

  createPlan = () => {
    this.setState({ selectedPlan: { metadata: {}, spec: { summary: 'New plan', description: 'New plan' } }, selectedPlanMode: 'create' })
  }

  viewPlan = (plan) => {
    this.setState({ selectedPlan: plan, selectedPlanMode: 'view' })
  }

  editPlan = (plan) => {
    this.setState({ selectedPlan: plan, selectedPlanMode: 'edit' })
  }

  deletePlan = (plan) => {
    Modal.confirm({
      title: `Are you sure you want to delete the plan ${plan.metadata.name}?`,
      content: 'This cannot be undone',
      okText: 'Yes',
      okType: 'danger',
      cancelText: 'No',
      onOk: async () => {
        const api = await KoreApi.client()
        await api.DeleteServicePlan(plan.metadata.name)
        await this.refreshPlans()
      }
    })
  }

  closePlan = () => {
    this.setState({ selectedPlan: null, selectedPlanMode: 'view' })
  }

  planUpdated = async () => {
    this.closePlan()
    await this.refreshPlans()
  }

  render() {
    const { kind, plans, loading, selectedPlan, selectedPlanMode } = this.state
    const displayName = kind.spec.displayName || kind.metadata.name
    return (
      <>
        <Drawer 
          visible={Boolean(selectedPlan)} 
          onClose={() => this.closePlan()} 
          title={!selectedPlan ? '' : <><Title level={4}>Plan: {selectedPlan.spec.summary}</Title></>}
          width={900}>
          {!selectedPlan ? null : (
            <ManageServicePlanForm
              mode={selectedPlanMode}
              resourceType="service"
              kind={kind.metadata.name}
              data={selectedPlan}
              handleSubmit={() => this.planUpdated()}
            />
          )}
        </Drawer>
        <Card 
          title={displayName} 
          style={{ marginBottom: '30px' }}
          extra={(
            <Tooltip placement="left" title={kind.spec.enabled ? 'This is available for teams to consume' : 'Teams cannot use this cloud service'}>
              <Switch loading={loading} onChange={() => this.toggleEnabled()} checked={kind.spec.enabled} checkedChildren="Enabled" unCheckedChildren="Disabled" />
            </Tooltip>
          )}>
          <Paragraph className="logo">
            { kind.spec.imageURL ? (
              <img src={kind.spec.imageURL} height="80px" />
            ) : (
              <Icon type="cloud-server" style={{ fontSize: '80px' }} theme="outlined" />
            ) }
          </Paragraph>
          <Paragraph className="name" strong>{displayName}</Paragraph>
          <Paragraph type="secondary">{kind.spec.description}</Paragraph>
        </Card>
        <Card 
          title={`Plans for ${displayName}`} 
          extra={kind.spec.schema ? <Button type="primary" onClick={this.createPlan}>+ New</Button> : null}>
          <List dataSource={plans.items} renderItem={(plan) => this.renderPlan(plan)} />
        </Card>
      </>
    )
  }
}