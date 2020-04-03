import PropTypes from 'prop-types'
import { Typography, List, Button, Drawer, Icon, Alert } from 'antd'
const { Title } = Typography

import { kore } from '../../../config'
import GCPOrganization from '../team/GCPOrganization'
import ResourceList from '../configure/ResourceList'
import GCPOrganizationForm from '../forms/GCPOrganizationForm'
import KoreApi from '../../kore-api'

class GCPOrganizationsList extends ResourceList {

  static propTypes = {
    style: PropTypes.object
  }

  createdMessage = 'GCP organization created successfully'
  updatedMessage = 'GCP organization updated successfully'

  async fetchComponentData() {
    const api = await KoreApi.client()
    const [ allTeams, gcpOrganizations, allAllocations ] = await Promise.all([
      api.ListTeams(),
      api.ListGCPOrganizations(kore.koreAdminTeamName),
      api.ListAllocations(kore.koreAdminTeamName)
    ])
    allTeams.items = allTeams.items.filter(t => !kore.ignoreTeams.includes(t.metadata.name))
    gcpOrganizations.items.forEach(org => {
      org.allocation = (allAllocations.items || []).find(alloc => alloc.metadata.name === org.metadata.name)
    })
    return { resources: gcpOrganizations, allTeams }
  }

  render() {
    const { resources, allTeams, edit, add } = this.state

    return (
      <>
        <Alert
          message="Give Kore access to your Google Cloud Platform organization"
          description="This will allow Kore to manage the organization for you. This includes managing the creation of projects and Service Accounts giving Kore teams the ability to create clusters with ease."
          type="info"
          showIcon
          style={{ marginBottom: '20px' }}
        />
        <Button type="primary" onClick={this.add(true)} style={{ display: 'block', marginBottom: '20px' }}>+ New</Button>
        {!resources ? <Icon type="loading" /> : (
          <>
            <List
              dataSource={resources.items}
              renderItem={org =>
                <GCPOrganization
                  organization={org}
                  allTeams={allTeams.items}
                  editOrganization={this.edit}
                  handleUpdate={this.handleStatusUpdated}
                  refreshMs={2000}
                  stateResourceDataKey="organization"
                  resourceApiPath={`/teams/${kore.koreAdminTeamName}/organizations/${org.metadata.name}`}
                />
              }
            >
            </List>
            {edit ? (
              <Drawer
                title={<Title level={4}>GCP Organization: {edit.spec.parentID}</Title>}
                visible={Boolean(edit)}
                onClose={this.edit(false)}
                width={700}
              >
                <GCPOrganizationForm
                  team={kore.koreAdminTeamName}
                  allTeams={allTeams}
                  data={edit}
                  handleSubmit={this.handleEditSave}
                />
              </Drawer>
            ) : null}
            {add ? (
              <Drawer
                title={<Title level={4}>New GCP organization</Title>}
                visible={add}
                onClose={this.add(false)}
                width={700}
              >
                <GCPOrganizationForm
                  team={kore.koreAdminTeamName}
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

export default GCPOrganizationsList
