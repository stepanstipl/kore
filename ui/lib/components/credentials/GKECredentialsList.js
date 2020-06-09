import PropTypes from 'prop-types'
import { Typography, List, Button, Drawer, Alert, Icon, Modal } from 'antd'
const { Title } = Typography
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

import GKECredentials from './GKECredentials'
import ResourceList from '../resources/ResourceList'
import GKECredentialsForm from './GKECredentialsForm'
import KoreApi from '../../kore-api'
import AllocationHelpers from '../../utils/allocation-helpers'
import { errorMessage, successMessage, loadingMessage } from '../../utils/message'

class GKECredentialsList extends ResourceList {

  static propTypes = {
    style: PropTypes.object
  }

  createdMessage = 'GCP project credentials created successfully'
  updatedMessage = 'GCP project credentials updated successfully'
  deletedMessage = 'GCP project credentials deleted successfully'
  deleteFailedMessage = 'Error deleting GCP project credentials'

  async fetchComponentData() {
    const api = await KoreApi.client()
    const [ allTeams, gkeCredentials, allAllocations ] = await Promise.all([
      api.ListTeams(),
      api.ListGKECredentials(publicRuntimeConfig.koreAdminTeamName),
      api.ListAllocations(publicRuntimeConfig.koreAdminTeamName)
    ])
    allTeams.items = allTeams.items.filter(t => !publicRuntimeConfig.ignoreTeams.includes(t.metadata.name))
    gkeCredentials.items.forEach(gke => {
      gke.allocation = AllocationHelpers.findAllocationForResource(allAllocations, gke)
    })
    return { resources: gkeCredentials, allTeams }
  }

  delete = (cred) => () => {
    Modal.confirm({
      title: `Are you sure you want to delete the credentials ${cred.spec.project}?`,
      content: 'This cannot be undone',
      okText: 'Yes',
      okType: 'danger',
      cancelText: 'No',
      onOk: async () => {
        const key = loadingMessage('Deleting allocations for credential', { duration: 0 })
        try {
          await AllocationHelpers.removeAllocation(cred)
          loadingMessage('Deleting credential', { key, duration: 0 })
          await (await KoreApi.client()).RemoveGKECredential(publicRuntimeConfig.koreAdminTeamName, cred.metadata.name)
          successMessage(this.deletedMessage, { key })
        } catch (err) {
          console.error(err)
          errorMessage(this.deleteFailedMessage, { key })
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
              id="gkecreds_list"
              dataSource={resources.items}
              renderItem={gke =>
                <GKECredentials
                  gkeCredentials={gke}
                  allTeams={allTeams.items}
                  editGKECredential={this.edit}
                  deleteGKECredential={this.delete}
                  handleUpdate={this.handleStatusUpdated}
                  handleDelete={() => {}}
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
