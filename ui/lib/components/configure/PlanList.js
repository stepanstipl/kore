import PropTypes from 'prop-types'
import { List, Alert, Icon, Drawer, Typography, Button } from 'antd'
const { Title, Text } = Typography

import PlanItem from '../team/PlanItem'
import PlanForm from '../plans/PlanForm'
import ResourceList from '../configure/ResourceList'
import Plan from '../configure/Plan'
import KoreApi from '../../kore-api'

class PlanList extends ResourceList {

  static propTypes = {
    kind: PropTypes.string,
    style: PropTypes.object,
  }

  createdMessage = 'GKE plan created successfully'
  updatedMessage = 'GKE plan updated successfully'

  async fetchComponentData() {
    const api = await KoreApi.client()
    const planList = await api.ListPlans(this.props.kind)
    return { resources: planList }
  }

  handleValidationErrors = validationErrors => {
    this.setState({ validationErrors })
  }

  render() {
    const { resources, view, edit, add, validationErrors } = this.state

    return (
      <>
        <Alert
          message="Manage the cluster plans"
          description="These plans define the specification of the clusters that can be created using the Google Kubernetes Engine on GCP. These help to give teams an easy way to provision clusters which match the requirements of the organization."
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
                <Plan plan={view} />
              </Drawer>
            ) : null}

            {edit ? (
              <Drawer
                title={<><Title level={4}>{edit.spec.description}</Title><Text>{edit.spec.summary}</Text></>}
                visible={Boolean(edit)}
                onClose={this.edit(false)}
                width={900}
              >
                <PlanForm
                  kind={this.props.kind}
                  data={edit}
                  validationErrors={validationErrors}
                  handleValidationErrors={this.handleValidationErrors}
                  handleSubmit={this.handleEditSave}
                />
              </Drawer>
            ) : null}

            {add ? (
              <Drawer
                title={<Title level={4}>New {this.props.kind} plan</Title>}
                visible={add}
                onClose={this.add(false)}
                width={900}
              >
                <PlanForm
                  kind={this.props.kind}
                  validationErrors={validationErrors}
                  handleValidationErrors={this.handleValidationErrors}
                  handleSubmit={this.handleAddSave}
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
