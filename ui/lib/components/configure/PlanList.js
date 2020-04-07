import PropTypes from 'prop-types'
import { List, Alert, Icon, Drawer, Typography } from 'antd'

const { Title, Text } = Typography

import PlanItem from '../team/PlanItem'
import ResourceList from '../configure/ResourceList'
import Plan from '../configure/Plan'
import KoreApi from '../../kore-api'

class PlanList extends ResourceList {

  static propTypes = {
    kind: PropTypes.string,
    style: PropTypes.object
  }

  async fetchComponentData() {
    const api = await KoreApi.client()
    const planList = await api.ListPlans(this.props.kind)
    return { resources: planList }
  }

  render() {
    const { resources, view } = this.state

    return (
      <>
        <Alert
          message="Manage the cluster plans"
          description="These plans define the specification of the clusters that can be created using the Google Kubernetes Engine on GCP. These help to give teams an easy way to provision clusters which match the requirements of the organization."
          type="info"
          showIcon
          style={{ marginBottom: '20px' }}
        />
        {!resources ? <Icon type="loading" /> : (
          <>
            <List
              dataSource={resources.items}
              renderItem={plan => <PlanItem plan={plan} viewPlan={this.view} /> }
            >
            </List>

            {view ? (
              <Drawer
                title={<><Title level={4}>{view.spec.description}</Title><Text>{view.spec.summary}</Text></>}
                visible={Boolean(view)}
                onClose={this.view(false)}
                width={700}
              >
                <Plan plan={view} />
              </Drawer>
            ) : null}
          </>
        )}
      </>
    )
  }
}

export default PlanList
