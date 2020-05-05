import PropTypes from 'prop-types'
import { Typography, List, Button, Drawer, Alert, Icon } from 'antd'
const { Title } = Typography
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

import GKECredentials from './GKECredentials'
import ResourceList from '../resources/ResourceList'
import GKECredentialsForm from './GKECredentialsForm'
import KoreApi from '../../kore-api'

class GKECredentialsList extends ResourceList {

  static propTypes = {
    style: PropTypes.object
  }

  createdMessage = 'GCP project credentials created successfully'
  updatedMessage = 'GCP project credentials updated successfully'

  async fetchComponentData() {
    const api = await KoreApi.client()
    const [ allTeams, gkeCredentials, allAllocations ] = await Promise.all([
      api.ListTeams(),
      api.ListGKECredentials(publicRuntimeConfig.koreAdminTeamName),
      api.ListAllocations(publicRuntimeConfig.koreAdminTeamName)
    ])
    allTeams.items = allTeams.items.filter(t => !publicRuntimeConfig.ignoreTeams.includes(t.metadata.name))
    gkeCredentials.items.forEach(gke => {
      gke.allocation = (allAllocations.items || []).find(alloc => alloc.metadata.name === gke.metadata.name)
    })
    return { resources: gkeCredentials, allTeams }
  }

  render() {
    const { resources, allTeams, edit, add } = this.state

    return (
      <>
        <Alert
          message="Give Kore access to your existing Google Cloud Platform projects"
          description="This will enable Kore to build clusters inside a GCP project that you already manage outside of Kore. You must create a Service Account inside your project and add the key in JSON format here."
          type="info"
          showIcon
          style={{ marginBottom: '20px' }}
        />
        <Button type="primary" onClick={this.add(true)} style={{ display: 'block', marginBottom: '20px' }}>+ New</Button>
        {!resources ? <Icon type="loading" /> : (
          <>
            <List
              dataSource={resources.items}
              renderItem={gke =>
                <GKECredentials
                  gkeCredentials={gke}
                  allTeams={allTeams.items}
                  editGKECredential={this.edit}
                  handleUpdate={this.handleStatusUpdated}
                  refreshMs={2000}
                  propsResourceDataKey="gkeCredentials"
                  resourceApiPath={`/teams/${publicRuntimeConfig.koreAdminTeamName}/gkecredentials/${gke.metadata.name}`}
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
                  team={publicRuntimeConfig.koreAdminTeamName}
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
                  team={publicRuntimeConfig.koreAdminTeamName}
                  allTeams={allTeams}
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

export default GKECredentialsList
