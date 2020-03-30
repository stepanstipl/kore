import PropTypes from 'prop-types'
import { Typography, Card, List, Button, Drawer, Icon } from 'antd'
const { Text, Title } = Typography

import { kore } from '../../../config'
import GCPOrganization from '../team/GCPOrganization'
import ResourceList from '../configure/ResourceList'
import GCPOrganizationForm from '../forms/GCPOrganizationForm'
import apiRequest from '../../utils/api-request'
import apiPaths from '../../utils/api-paths'

class GCPOrganizationsList extends ResourceList {

  static propTypes = {
    style: PropTypes.object
  }

  createdMessage = 'GCP organization created successfully'
  updatedMessage = 'GCP organization updated successfully'

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
    return { resources: gcpOrganizations, allTeams }
  }

  render() {
    const { resources, allTeams, edit, add } = this.state

    return (
      <Card style={this.props.style} title="GCP organizations" extra={<Button type="primary" onClick={this.add(true)}>+ New</Button>}>
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
                title={
                  edit.allocation ? (
                    <div>
                      <Title level={4}>{edit.allocation.spec.name}</Title>
                      <Text>{edit.allocation.spec.summary}</Text>
                    </div>
                  ) : (
                    <Title level={4}>{edit.metadata.name}</Title>
                  )
                }
                visible={!!edit}
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
      </Card>
    )
  }
}

export default GCPOrganizationsList
