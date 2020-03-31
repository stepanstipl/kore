import PropTypes from 'prop-types'
import { Typography, Card, List, Button, Drawer, Alert, Icon } from 'antd'
const { Title } = Typography

import { kore } from '../../../config'
import Credentials from '../team/Credentials'
import ResourceList from '../configure/ResourceList'
import GKECredentialsForm from '../forms/GKECredentialsForm'
import apiRequest from '../../utils/api-request'
import apiPaths from '../../utils/api-paths'

class GKECredentialsList extends ResourceList {

  static propTypes = {
    style: PropTypes.object
  }

  createdMessage = 'GCP project created successfully'
  updatedMessage = 'GCP project updated successfully'

  async fetchComponentData() {
    const [ allTeams, gkeCredentials, allAllocations ] = await Promise.all([
      apiRequest(null, 'get', apiPaths.teams),
      apiRequest(null, 'get', apiPaths.team(kore.koreAdminTeamName).gkeCredentials),
      apiRequest(null, 'get', apiPaths.team(kore.koreAdminTeamName).allocations)
    ])
    allTeams.items = allTeams.items.filter(t => !kore.ignoreTeams.includes(t.metadata.name))
    gkeCredentials.items.forEach(gke => {
      gke.allocation = (allAllocations.items || []).find(alloc => alloc.metadata.name === gke.metadata.name)
    })
    return { resources: gkeCredentials, allTeams }
  }

  render() {
    const { resources, allTeams, edit, add } = this.state

    return (
      <Card
        title="Projects"
        extra={<Button type="primary" onClick={this.add(true)}>+ New</Button>}
        style={this.props.style}
      >
        <Alert
          message="Give Kore access to your existing Google Cloud Platform projects"
          description="This will enable Kore to build clusters inside a GCP project that you already manage outside of Kore. You must create a Service Account inside your project and add the key in JSON format here."
          type="info"
          showIcon
          style={{ marginBottom: '20px' }}
        />
        {!resources ? <Icon type="loading" /> : (
          <>
            <List
              dataSource={resources.items}
              renderItem={gke =>
                <Credentials
                  gke={gke}
                  allTeams={allTeams.items}
                  editGKECredential={this.edit}
                  handleUpdate={this.handleStatusUpdated}
                  refreshMs={2000}
                  stateResourceDataKey="gke"
                  resourceApiPath={`/teams/${kore.koreAdminTeamName}/gkecredentials/${gke.metadata.name}`}
                />
              }
            >
            </List>
            {edit ? (
              <Drawer
                title={<Title level={4}>GCP project: {edit.spec.project}</Title>}
                visible={Boolean(edit)}
                onClose={this.edit(false)}
                width={700}
              >
                <GKECredentialsForm
                  team={kore.koreAdminTeamName}
                  allTeams={allTeams}
                  data={edit}
                  handleSubmit={this.handleEditSave}
                />
              </Drawer>
            ) : null}
            {add ? (
              <Drawer
                title={<Title level={4}>New GCP project</Title>}
                visible={add}
                onClose={this.add(false)}
                width={700}
              >
                <GKECredentialsForm
                  team={kore.koreAdminTeamName}
                  allTeams={allTeams}
                  handleSubmit={this.handleAddSave}
                />
              </Drawer>
            ) : null}
          </>
        )}
      </Card>
    )
  }
}

export default GKECredentialsList
