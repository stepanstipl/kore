import PropTypes from 'prop-types'
import { List, Alert, Icon, Drawer, Typography, Button, Modal } from 'antd'
const { Title, Text } = Typography

import PlanItem from './PlanItem'
import ManageClusterPlanForm from './ManageClusterPlanForm'
import ResourceList from '../resources/ResourceList'
import PlanViewer from './PlanViewer'
import KoreApi from '../../kore-api'
import AllocationHelpers from '../../utils/allocation-helpers'
import { successMessage, errorMessage, loadingMessage } from '../../utils/message'

class PlanList extends ResourceList {

  static propTypes = {
    kind: PropTypes.string,
    style: PropTypes.object
  }

  createdMessage = `${this.props.kind} plan created successfully`
  updatedMessage = `${this.props.kind} plan updated successfully`

  infoDescription = {
    GKE: 'These plans define the specification of the clusters that can be created using the Google Kubernetes Engine (GKE) on GCP. These help to give teams an easy way to provision clusters which match the requirements of the organization.',
    EKS: 'These plans define the specification of the clusters that can be created using the Elastic Kubernetes Service (EKS) on AWS. These help to give teams an easy way to provision clusters which match the requirements of the organization.'
  }

  async fetchComponentData() {
    const api = await KoreApi.client()
    const [ planList, accountManagementList ] = await Promise.all([
      api.ListPlans(this.props.kind),
      api.ListAccounts()
    ])

    const accountManagement = accountManagementList.items.find(a => a.spec.provider === this.props.kind)
    if (accountManagement) {
      planList.items = planList.items.map(plan => {
        (accountManagement.spec.rules || []).forEach(rule => rule.plans.forEach(rulePlan => {
          if (rulePlan === plan.metadata.name) {
            plan = { ...plan, gcpAutomatedProject: rule }
          }
        }))
        return plan
      })
    }
    return { resources: planList, accountManagement }
  }

  componentDidUpdate(prevProps) {
    // reload data if coming back from another tab
    if (prevProps.kind !== this.props.kind) {
      this.fetchComponentData().then(data => this.setState({ ...data }))
    }
  }

  processPlanCreateEdit = (process) => {
    return (args) => {
      process && process(args)
    }
  }

  unassociatedPlanWarning = (plan) => {
    if (!this.state.accountManagement) {
      return false
    }
    if (!this.state.accountManagement.spec.rules) {
      return false
    }
    if (!plan.gcpAutomatedProject) {
      return true
    }
    return false
  }

  delete = (plan) => () => {
    Modal.confirm({
      title: `Are you sure you want to delete the plan ${plan.spec.description}?`,
      content: 'This cannot be undone',
      okText: 'Yes',
      okType: 'danger',
      cancelText: 'No',
      onOk: async () => {
        const key = loadingMessage(`Deleting allocations for plan ${plan.spec.description}`, { duration: 0 })
        try {
          await AllocationHelpers.removeAllocation(plan)
          loadingMessage(`Deleting plan ${plan.spec.description}`, { key, duration: 0 })
          await (await KoreApi.client()).RemovePlan(plan.metadata.name)
          successMessage(`${plan.spec.description} plan deleted`, { key })
        } catch (err) {
          console.error(err)
          errorMessage(`Error deleting plan ${plan.spec.description}`, { key })
        }
        await this.refresh()
      }
    })
  }

  render() {
    const { resources, view, edit, add } = this.state
    const drawerVisible = Boolean(view || edit || add)
    let drawerTitle = null
    let drawerClose = () => {}
    if (view) {
      drawerTitle = <><Title level={4}>{view.spec.summary}</Title><Text>{view.spec.description}</Text></>
      drawerClose = this.view(false)
    } else if (edit) {
      drawerTitle = <><Title level={4}>{edit.spec.summary}</Title><Text>{edit.spec.description}</Text></>
      drawerClose = this.edit(false)
    } else if (add) {
      drawerTitle = <Title level={4}>New {this.props.kind} plan</Title>
      drawerClose = this.add(false)
    }

    return (
      <>
        <Alert
          message="Manage the cluster plans"
          description={this.infoDescription[this.props.kind]}
          type="info"
          showIcon
          style={{ marginBottom: '20px' }}
        />
        <Button id="add" type="primary" onClick={this.add(true)} style={{ display: 'block', marginBottom: '20px' }}>+ New</Button>

        <Drawer
          title={drawerTitle}
          visible={drawerVisible}
          onClose={drawerClose}
          width={900}>
          {!view ? null : 
            <PlanViewer
              plan={view}
              resourceType="cluster"
              displayUnassociatedPlanWarning={this.unassociatedPlanWarning(view)}
            />
          }
          {!edit ? null :
            <ManageClusterPlanForm
              mode="edit"
              kind={this.props.kind}
              data={edit}
              handleSubmit={(args) => this.handleEditSave(args)}
              displayUnassociatedPlanWarning={this.unassociatedPlanWarning(edit)}
            />
          }
          {!add ? null : 
            <ManageClusterPlanForm
              mode="create"
              kind={this.props.kind}
              handleSubmit={(args) => this.handleAddSave(args)}
            />
          }
        </Drawer>

        {!resources ? <Icon type="loading" /> : (
          <>
            <List
              id="plans_list"
              dataSource={resources.items}
              renderItem={plan => <PlanItem plan={plan} viewPlan={this.view} editPlan={this.edit} deletePlan={this.delete} displayUnassociatedPlanWarning={this.unassociatedPlanWarning(plan)} /> }
            >
            </List>
          </>
        )}
      </>
    )
  }
}

export default PlanList
