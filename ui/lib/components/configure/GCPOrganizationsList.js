import React from 'react'
import PropTypes from 'prop-types'
import { Typography, Card, List, Button, Drawer, message, Icon } from 'antd'
const { Text, Title } = Typography

import { kore } from '../../../config'
import GCPOrganization from '../team/GCPOrganization'
import GCPOrganizationForm from '../forms/GCPOrganizationForm'
import apiRequest from '../../utils/api-request'
import apiPaths from '../../utils/api-paths'
import copy from '../../utils/object-copy'

class GCPOrganizationsList extends React.Component {

  static propTypes = {
    style: PropTypes.object
  }

  constructor(props) {
    super(props)
    this.state = {
      dataLoading: true,
      editOrganization: false,
      addOrganization: false
    }
  }

  async fetchComponentData() {
    const [ allTeams, gcpOrganizations, allAllocations ] = await Promise.all([
      apiRequest(null, 'get', apiPaths.teams),
      apiRequest(null, 'get', apiPaths.team(kore.koreAdminTeamName).gcpOrganizations),
      apiRequest(null, 'get', apiPaths.team(kore.koreAdminTeamName).allocations)
    ])
    allTeams.items = allTeams.items.filter(t => !kore.ignoreTeams.includes(t.metadata.name))
    gcpOrganizations.items.forEach(org => {
      org.allocation = (allAllocations.items || []).find(alloc => alloc.metadata.name === org.metadata.name)
    })
    return { gcpOrganizations, allTeams }
  }

  componentDidMount() {
    return this.fetchComponentData()
      .then(({ gcpOrganizations, allTeams }) => {
        const state = copy(this.state)
        state.gcpOrganizations = gcpOrganizations
        state.allTeams = allTeams
        state.dataLoading = false
        this.setState(state)
      })
  }

  handleStatusUpdated = (updatedResource, done) => {
    const state = copy(this.state)
    const resource = state.gcpOrganizations.items.find(r => r.metadata.name === updatedResource.metadata.name)
    resource.status = updatedResource.status
    this.setState(state, done)
  }

  editOrganization = org => {
    return async () => {
      const state = copy(this.state)
      state.editOrganization = org ? org : false
      this.setState(state)
    }
  }

  handleEditOrganizationSave = updatedOrg => {
    const state = copy(this.state)

    const editedIntegration = state.gcpOrganizations.items.find(c => c.metadata.name === state.editOrganization.metadata.name)
    editedIntegration.spec = updatedOrg.spec
    editedIntegration.allocation = updatedOrg.allocation
    editedIntegration.status.status = 'Pending'

    state.editOrganization = false
    this.setState(state)
    message.success('GCP organization updated successfully')
  }

  addOrganization = enabled => {
    return () => {
      const state = copy(this.state)
      state.addOrganization = enabled
      this.setState(state)
    }
  }

  handleAddOrganizationSave = async createdOrg => {
    const state = copy(this.state)
    state.gcpOrganizations.items.push(createdOrg)
    state.addOrganization = false
    this.setState(state)
    message.success('GCP organization created successfully')
  }

  render() {
    const { gcpOrganizations, allTeams, editOrganization, addOrganization } = this.state

    return (
      <Card style={this.props.style} title="GCP organizations" extra={<Button type="primary" onClick={this.addOrganization(true)}>+ New</Button>}>
        {!gcpOrganizations ? <Icon type="loading" /> : (
          <>
            <List
              dataSource={gcpOrganizations.items}
              renderItem={org =>
                <GCPOrganization
                  organization={org}
                  allTeams={allTeams.items}
                  editOrganization={this.editOrganization}
                  handleUpdate={this.handleStatusUpdated}
                  refreshMs={2000}
                  stateResourceDataKey="organization"
                  resourceApiPath={`/teams/${kore.koreAdminTeamName}/organizations/${org.metadata.name}`}
                />
              }
            >
            </List>
            {editOrganization ? (
              <Drawer
                title={
                  editOrganization.allocation ? (
                    <div>
                      <Title level={4}>{editOrganization.allocation.spec.name}</Title>
                      <Text>{editOrganization.allocation.spec.summary}</Text>
                    </div>
                  ) : (
                    <Title level={4}>{editOrganization.metadata.name}</Title>
                  )
                }
                visible={!!editOrganization}
                onClose={this.editOrganization(false)}
                width={700}
              >
                <GCPOrganizationForm
                  team={kore.koreAdminTeamName}
                  allTeams={allTeams}
                  data={editOrganization}
                  handleSubmit={this.handleEditOrganizationSave}
                />
              </Drawer>
            ) : null}
            {addOrganization ? (
              <Drawer
                title={<Title level={4}>New GCP organization</Title>}
                visible={addOrganization}
                onClose={this.addOrganization(false)}
                width={700}
              >
                <GCPOrganizationForm
                  team={kore.koreAdminTeamName}
                  allTeams={allTeams}
                  handleSubmit={this.handleAddOrganizationSave}
                />
              </Drawer>
            ) : null}
          </>
        )}
      </Card>
    )
  }
}

export default GCPOrganizationsList
