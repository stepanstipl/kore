import React from 'react'
import PropTypes from 'prop-types'
import KoreApi from '../../../lib/kore-api'
import { Card, Divider, Typography, Icon, Switch, Tooltip, List, Avatar, Button, Drawer, Modal } from 'antd'
import ManageServicePlanForm from '../../../lib/components/plans/ManageServicePlanForm'
import { isReadOnlyCRD } from '../../utils/crd-helpers'
import { errorMessage } from '../../utils/message'
const { Text, Title, Paragraph } = Typography

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
          title={displayName}
          description={plan.spec.description} />
      </List.Item>
    )
  }

  createPlan = () => {
    this.setState({ selectedPlan: null, selectedPlanMode: 'create' })
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
        try {
          const api = await KoreApi.client()
          await api.DeleteServicePlan(plan.metadata.name)
          await this.refreshPlans()
        } catch (err) {
          if (err.statusCode === 409 && err.dependents) {
            return Modal.warning({
              title: 'The service plan can not be deleted',
              content: (
                <div>
                  <Paragraph strong>Error: {err.message}</Paragraph>
                  <List
                    size="small"
                    dataSource={err.dependents}
                    renderItem={d => <List.Item>{d.kind}: {d.name}</List.Item>}
                  />
                </div>
              ),
              onOk() {}
            })
          }
          errorMessage('Error deleting service plan, please try again.')
        }
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
          visible={Boolean(selectedPlan) || selectedPlanMode === 'create'}
          onClose={() => this.closePlan()}
          title={!selectedPlan ? '' : <><Title level={4}>Plan: {selectedPlan.spec.summary}</Title></>}
          width={900}>
          <ManageServicePlanForm
            mode={selectedPlanMode}
            resourceType="service"
            kind={kind.metadata.name}
            data={selectedPlan}
            handleSubmit={() => this.planUpdated()}
          />
        </Drawer>

        <List.Item>
          <List.Item.Meta
            className="large-list-item"
            avatar={kind && kind.spec.imageURL ? <Avatar src={kind.spec.imageURL} /> : <Avatar size={60} icon="cloud-server" />}
            title={displayName}
            description={kind.spec.description || 'No description'}
          />
          <div style={{ marginLeft: '30px' }}>
            <Tooltip placement="left" title={kind.spec.enabled ? 'This is available for teams to consume' : 'Teams cannot use this cloud service'}>
              <Switch loading={loading} onChange={() => this.toggleEnabled()} checked={kind.spec.enabled} checkedChildren="Enabled" unCheckedChildren="Disabled" />
            </Tooltip>
          </div>
        </List.Item>

        <Divider />

        <Card
          title={`Plans for ${displayName}`}
          extra={kind.spec.schema ? <Button type="primary" onClick={this.createPlan}>+ New plan</Button> : null}>
          <List dataSource={plans.items} renderItem={(plan) => this.renderPlan(plan)} />
        </Card>
      </>
    )
  }
}
