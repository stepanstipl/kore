import PropTypes from 'prop-types'
import { Typography, List, Button, Drawer, Alert, Icon, Modal } from 'antd'
const { Title } = Typography
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

import ResourceList from '../resources/ResourceList'
import EKSCredentialsForm from './EKSCredentialsForm'
import EKSCredentials from './EKSCredentials'
import KoreApi from '../../kore-api'
import AllocationHelpers from '../../utils/allocation-helpers'
import { loadingMessage, successMessage, errorMessage } from '../../utils/message'

class EKSCredentialsList extends ResourceList {

  static propTypes = {
    style: PropTypes.object
  }

  createdMessage = 'AWS account credentials created successfully'
  updatedMessage = 'AWS account credentials updated successfully'
  deletedMessage = 'AWS account credentials deleted successfully'
  deleteFailedMessage = 'Error deleting AWS account credentials'

  async fetchComponentData() {
    const api = await KoreApi.client()
    const [ allTeams, eksCredentials, allAllocations ] = await Promise.all([
      api.ListTeams(),
      api.ListEKSCredentials(publicRuntimeConfig.koreAdminTeamName),
      api.ListAllocations(publicRuntimeConfig.koreAdminTeamName)
    ])
    allTeams.items = allTeams.items.filter(t => !publicRuntimeConfig.ignoreTeams.includes(t.metadata.name))
    eksCredentials.items.forEach(eks => {
      eks.allocation = AllocationHelpers.findAllocationForResource(allAllocations, eks)
    })
    return { resources: eksCredentials, allTeams }
  }

  delete = (cred) => () => {
    Modal.confirm({
      title: `Are you sure you want to delete the credentials ${cred.spec.accountID}?`,
      content: 'This cannot be undone',
      okText: 'Yes',
      okType: 'danger',
      cancelText: 'No',
      onOk: async () => {
        const key = loadingMessage('Deleting allocations for credential', { duration: 0 })
        try {
          await AllocationHelpers.removeAllocation(cred)
          loadingMessage('Deleting credential', { key, duration: 0 })
          await (await KoreApi.client()).DeleteEKSCredentials(publicRuntimeConfig.koreAdminTeamName, cred.metadata.name)
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
          message="Give Kore access to your existing AWS accounts"
          description="This will enable Kore to build clusters inside an AWS account that you already manage outside of Kore. You must create an access key in the account and provide the details here."
          type="info"
          showIcon
          style={{ marginBottom: '20px' }}
        />
        <Button id="add" type="primary" onClick={this.add(true)} style={{ display: 'block', marginBottom: '20px' }}>+ New</Button>
        {!resources ? <Icon type="loading" /> : (
          <>
            <List
              dataSource={resources.items}
              renderItem={eks =>
                <EKSCredentials
                  eksCredentials={eks}
                  allTeams={allTeams.items}
                  editEKSCredentials={this.edit}
                  deleteEKSCredentials={this.delete}
                  handleUpdate={this.handleStatusUpdated}
                  handleDelete={() => {}}
                  refreshMs={2000}
                  propsResourceDataKey="eksCredentials"
                  resourceApiPath={`/teams/${publicRuntimeConfig.koreAdminTeamName}/ekscredentials/${eks.metadata.name}`}
                />
              }
            >
            </List>

            {edit ? (
              <Drawer
                title={<Title level={4}>AWS account: {edit.spec.accountID}</Title>}
                visible={Boolean(edit)}
                onClose={this.edit(false)}
                width={700}
              >
                <EKSCredentialsForm
                  team={publicRuntimeConfig.koreAdminTeamName}
                  allTeams={allTeams}
                  data={edit}
                  handleSubmit={this.handleEditSave}
                />
              </Drawer>
            ) : null}
            {add ? (
              <Drawer
                title={<Title level={4}>New AWS account</Title>}
                visible={add}
                onClose={this.add(false)}
                width={700}
              >
                <EKSCredentialsForm
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

export default EKSCredentialsList
