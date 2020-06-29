import { Typography, List, Button, Drawer, Icon, Alert, Modal } from 'antd'
const { Title } = Typography
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

import GCPOrganization from './GCPOrganization'
import ResourceList from '../resources/ResourceList'
import GCPOrganizationForm from './GCPOrganizationForm'
import KoreApi from '../../kore-api'
import AllocationHelpers from '../../utils/allocation-helpers'
import { errorMessage, loadingMessage, successMessage } from '../../utils/message'

class GCPOrganizationsList extends ResourceList {

  createdMessage = 'GCP organization created successfully'
  updatedMessage = 'GCP organization updated successfully'
  deletedMessage = 'GCP organization deleted successfully'
  deleteFailedMessage = 'Error deleting GCP organization'

  async fetchComponentData() {
    const api = await KoreApi.client()
    const [ allTeams, gcpOrganizations, allAllocations, accountList ] = await Promise.all([
      api.ListTeams(),
      api.ListGCPOrganizations(publicRuntimeConfig.koreAdminTeamName),
      api.ListAllocations(publicRuntimeConfig.koreAdminTeamName),
      api.ListAccounts()
    ])
    allTeams.items = allTeams.items.filter(t => !publicRuntimeConfig.ignoreTeams.includes(t.metadata.name))
    gcpOrganizations.items.forEach(org => {
      org.allocation = AllocationHelpers.findAllocationForResource(allAllocations, org)
      org.accountManagement = accountList.items.find(a => a.metadata.name === `am-${org.metadata.name}`)
    })
    return { resources: gcpOrganizations, allTeams }
  }

  delete = (gcpOrg) => () => {
    Modal.confirm({
      title: `Are you sure you want to delete the GCP Organization "${gcpOrg.spec.parentID}"?`,
      content: 'This cannot be undone',
      okText: 'Yes',
      okType: 'danger',
      cancelText: 'No',
      onOk: async () => {
        const key = loadingMessage('Deleting', { duration: 0 })
        try {
          if (gcpOrg.accountManagement) {
            loadingMessage('Deleting allocations for organization account management', { key })
            await AllocationHelpers.removeAllocation(gcpOrg.accountManagement)
            loadingMessage('Deleting account management', { key })
            await (await KoreApi.client()).RemoveAccount(gcpOrg.accountManagement.metadata.name)
          }

          loadingMessage('Deleting allocations for organization', { key })
          await AllocationHelpers.removeAllocation(gcpOrg)
          loadingMessage('Deleting organization', { key })
          await (await KoreApi.client()).DeleteGCPOrganization(publicRuntimeConfig.koreAdminTeamName, gcpOrg.metadata.name)
          successMessage(this.deletedMessage, { key })
        } catch (err) {
          console.log('Error deleting org', err.statusCode, err.response)
          let msg = this.deleteFailedMessage
          if (err.statusCode === 403) {
            msg += `: ${err.response.body.message}`
          }
          errorMessage(msg, { key, duration: 10 })
        }
        await this.refresh()
      }
    })
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
        {!resources ? <Icon type="loading" /> : (
          <>
            {resources.items.length === 0 && <Button type="primary" onClick={this.add(true)} style={{ display: 'block', marginBottom: '20px' }} className="new-gcp-organization">Configure</Button>}
            <List
              dataSource={resources.items}
              renderItem={org =>
                <GCPOrganization
                  organization={org}
                  allTeams={allTeams.items}
                  editOrganization={this.edit}
                  deleteOrganization={this.delete}
                  handleUpdate={this.handleStatusUpdated}
                  handleDelete={() => {}}
                  refreshMs={2000}
                  propsResourceDataKey="organization"
                  resourceApiRequest={async () => await (await KoreApi.client()).GetGCPOrganization(publicRuntimeConfig.koreAdminTeamName, org.metadata.name)}
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
                  team={publicRuntimeConfig.koreAdminTeamName}
                  allTeams={allTeams}
                  data={edit}
                  handleSubmit={this.handleEditSave}
                  autoAllocateToAllTeams={this.props.autoAllocateToAllTeams}
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
                  team={publicRuntimeConfig.koreAdminTeamName}
                  allTeams={allTeams}
                  handleSubmit={this.handleAddSave}
                  autoAllocateToAllTeams={this.props.autoAllocateToAllTeams}
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
