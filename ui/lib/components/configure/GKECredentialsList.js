import React from 'react'
import { Typography, Card, List, Button, Drawer, message, Icon } from 'antd'
const { Text, Title } = Typography

import { kore } from '../../../config'
import Credentials from '../team/Credentials'
import GKECredentialsForm from '../forms/GKECredentialsForm'
import apiRequest from '../../utils/api-request'
import apiPaths from '../../utils/api-paths'
import copy from '../../utils/object-copy'

class GKECredentialsList extends React.Component {

  constructor(props) {
    super(props)
    this.state = {
      dataLoading: true,
      editCredential: false,
      addCredential: false
    }
  }

  async fetchComponentData() {
    const [ allTeams, gkeCredentials, allAllocations ] = await Promise.all([
      apiRequest(null, 'get', apiPaths.teams),
      apiRequest(null, 'get', apiPaths.team(kore.koreAdminTeamName).gkeCredentials),
      apiRequest(null, 'get',apiPaths.team(kore.koreAdminTeamName).allocations)
    ])
    allTeams.items = allTeams.items.filter(t => !kore.ignoreTeams.includes(t.metadata.name))
    gkeCredentials.items.forEach(gke => {
      gke.allocation = (allAllocations.items || []).find(alloc => alloc.metadata.name === gke.metadata.name)
    })
    return { gkeCredentials, allTeams }
  }

  componentDidMount() {
    return this.fetchComponentData()
      .then(({ gkeCredentials, allTeams }) => {
        const state = copy(this.state)
        state.gkeCredentials = gkeCredentials
        state.allTeams = allTeams
        state.dataLoading = false
        this.setState(state)
      })
  }

  handleStatusUpdated = resourceType => {
    return (updatedResource, done) => {
      const state = copy(this.state)
      const resource = state[resourceType].items.find(r => r.metadata.name === updatedResource.metadata.name)
      resource.status = updatedResource.status
      this.setState(state, done)
    }
  }

  editCredential = gkeCredentials => {
    return async () => {
      const state = copy(this.state)
      state.editCredential = gkeCredentials ? { type: 'GKE', integration: gkeCredentials } : false
      this.setState(state)
    }
  }

  handleEditCredentialSave = updatedIntegration => {
    const state = copy(this.state)

    const editedIntegration = state.gkeCredentials.items.find(c => c.metadata.name === state.editCredential.integration.metadata.name)
    editedIntegration.spec = updatedIntegration.spec
    editedIntegration.allocation = updatedIntegration.allocation
    editedIntegration.status.status = 'Pending'

    state.editCredential = false
    this.setState(state)
    message.success('GKE credentials updated successfully')
  }

  addCredential = type => {
    return () => {
      const state = copy(this.state)
      state.addCredential = type
      this.setState(state)
    }
  }

  handleAddCredentialSave = async createdIntegration => {
    const state = copy(this.state)
    state.gkeCredentials.items.push(createdIntegration)
    state.addCredential = false
    this.setState(state)
    message.success('GKE credentials created successfully')
  }

  render() {
    const { gkeCredentials, allTeams, editCredential, addCredential } = this.state

    return (
      <Card title="GKE credentials" extra={<Button type="primary" onClick={this.addCredential('GKE')}>+ New</Button>}>
        {!gkeCredentials ? <Icon type="loading" /> : (
          <>
            <List
              dataSource={gkeCredentials.items}
              renderItem={gke =>
                <Credentials
                  gke={gke}
                  allTeams={allTeams.items}
                  editGKECredential={this.editCredential}
                  handleUpdate={this.handleStatusUpdated('gkeCredentials')}
                  refreshMs={2000}
                  stateResourceDataKey="gke"
                  resourceApiPath={`/teams/${kore.koreAdminTeamName}/gkecredentials/${gke.metadata.name}`}
                />
              }
            >
            </List>
            {editCredential ? (
              <Drawer
                title={
                  editCredential.integration.allocation ? (
                    <div>
                      <Title level={4}>{editCredential.integration.allocation.spec.name}</Title>
                      <Text>{editCredential.integration.allocation.spec.summary}</Text>
                    </div>
                  ) : (
                    <Title level={4}>{editCredential.integration.metadata.name}</Title>
                  )
                }
                visible={!!editCredential}
                onClose={this.clearEditIntegration}
                width={700}
              >
                {editCredential.type === 'GKE' ?
                  <GKECredentialsForm
                    team={kore.koreAdminTeamName}
                    allTeams={allTeams}
                    data={editCredential.integration}
                    handleSubmit={this.handleEditCredentialSave}
                  /> :
                  null
                }
              </Drawer>
            ) : null}
            {addCredential ? (
              <Drawer
                title={<Title level={4}>New {addCredential} credentials</Title>}
                visible={!!addCredential}
                onClose={this.addCredential(false)}
                width={700}
              >
                {addCredential === 'GKE' ?
                  <GKECredentialsForm
                    team={kore.koreAdminTeamName}
                    allTeams={allTeams}
                    handleSubmit={this.handleAddCredentialSave}
                  /> :
                  null
                }
              </Drawer>
            ) : null}
          </>
        )}
      </Card>
    )
  }
}

export default GKECredentialsList
