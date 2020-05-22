import PropTypes from 'prop-types'
import { List, Alert, Icon, Drawer, Typography, Button } from 'antd'
const { Title, Text } = Typography

import PlanItem from './PlanItem'
import PlanForm from './PlanForm'
import ResourceList from '../resources/ResourceList'
import PlanViewer from './PlanViewer'
import KoreApi from '../../kore-api'

class PlanList extends ResourceList {

  static propTypes = {
    kind: PropTypes.string,
    style: PropTypes.object,
    tabActiveKey: PropTypes.string.isRequired
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
    if (prevProps.tabActiveKey !== this.props.tabActiveKey) {
      this.fetchComponentData().then(data => this.setState({ ...data }))
    }
  }

  handleValidationErrors = validationErrors => {
    this.setState({ validationErrors })
  }

  processAndClearValidationErrors = process => {
    return args => {
      process && process(args)
      this.handleValidationErrors(null)
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

  render() {
    const { resources, view, edit, add, validationErrors } = this.state

    return (
      <>
        <Alert
          message="Manage the cluster plans"
          description={this.infoDescription[this.props.kind]}
          type="info"
          showIcon
          style={{ marginBottom: '20px' }}
        />
        <Button type="primary" onClick={this.add(true)} style={{ display: 'block', marginBottom: '20px' }}>+ New</Button>
        {!resources ? <Icon type="loading" /> : (
          <>
            <List
              dataSource={resources.items}
              renderItem={plan => <PlanItem plan={plan} viewPlan={this.view} editPlan={this.edit} displayUnassociatedPlanWarning={this.unassociatedPlanWarning(plan)} /> }
            >
            </List>

            {view ? (
              <Drawer
                title={<><Title level={4}>{view.spec.summary}</Title><Text>{view.spec.description}</Text></>}
                visible={Boolean(view)}
                onClose={this.view(false)}
                width={900}
              >
                <PlanViewer
                  plan={view}
                  resourceType="cluster"
                  displayUnassociatedPlanWarning={this.unassociatedPlanWarning(view)}
                />
              </Drawer>
            ) : null}

            {edit ? (
              <Drawer
                title={<><Title level={4}>{edit.spec.summary}</Title><Text>{edit.spec.description}</Text></>}
                visible={Boolean(edit)}
                onClose={this.processAndClearValidationErrors(this.edit(false))}
                width={900}
              >
                <PlanForm
                  kind={this.props.kind}
                  data={edit}
                  validationErrors={validationErrors}
                  handleValidationErrors={this.handleValidationErrors}
                  handleSubmit={this.processAndClearValidationErrors(this.handleEditSave)}
                  displayUnassociatedPlanWarning={this.unassociatedPlanWarning(edit)}
                />
              </Drawer>
            ) : null}

            {add ? (
              <Drawer
                title={<Title level={4}>New {this.props.kind} plan</Title>}
                visible={add}
                onClose={this.processAndClearValidationErrors(this.add(false))}
                width={900}
              >
                <PlanForm
                  kind={this.props.kind}
                  validationErrors={validationErrors}
                  handleValidationErrors={this.handleValidationErrors}
                  handleSubmit={this.processAndClearValidationErrors(this.handleAddSave)}
                />
              </Drawer>
            ) : null}
          </>
        )}
      </>
    )
  }
}

export default PlanList
