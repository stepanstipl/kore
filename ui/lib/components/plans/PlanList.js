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
  }

  createdMessage = `${this.props.kind} plan created successfully`
  updatedMessage = `${this.props.kind}  plan updated successfully`

  infoDescription = {
    GKE: 'These plans define the specification of the clusters that can be created using the Google Kubernetes Engine (GKE) on GCP. These help to give teams an easy way to provision clusters which match the requirements of the organization.',
    EKS: 'These plans define the specification of the clusters that can be created using the Elastic Kubernetes Service (EKS) on AWS. These help to give teams an easy way to provision clusters which match the requirements of the organization.'
  }

  async fetchComponentData() {
    const api = await KoreApi.client()
    const planList = await api.ListPlans(this.props.kind)
    return { resources: planList }
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
              renderItem={plan => <PlanItem plan={plan} viewPlan={this.view} editPlan={this.edit} /> }
            >
            </List>

            {view ? (
              <Drawer
                title={<><Title level={4}>{view.spec.description}</Title><Text>{view.spec.summary}</Text></>}
                visible={Boolean(view)}
                onClose={this.view(false)}
                width={900}
              >
                <PlanViewer
                  plan={view}
                  resourceType="cluster"
                />
              </Drawer>
            ) : null}

            {edit ? (
              <Drawer
                title={<><Title level={4}>{edit.spec.description}</Title><Text>{edit.spec.summary}</Text></>}
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
