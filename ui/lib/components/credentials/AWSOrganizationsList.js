import PropTypes from 'prop-types'
import { Typography, List, Button, Drawer, Icon, Alert, Modal } from 'antd'
const { Title } = Typography
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

import AWSOrganization from './AWSOrganization'
import ResourceList from '../resources/ResourceList'
import AWSOrganizationForm from './AWSOrganizationForm'
import KoreApi from '../../kore-api'
import AllocationHelpers from '../../utils/allocation-helpers'
import { errorMessage, loadingMessage, successMessage } from '../../utils/message'

class AWSOrganizationsList extends ResourceList {

  static propTypes = {
    user: PropTypes.object.isRequired
  }

  createdMessage = 'AWS organization created successfully'
  updatedMessage = 'AWS organization updated successfully'
  deletedMessage = 'AWS organization deleted successfully'
  deleteFailedMessage = 'Error deleting AWS organization'

  async fetchComponentData() {
    const api = await KoreApi.client()
    const [ allTeams, awsOrganizations, allAllocations, accountList ] = await Promise.all([
      api.ListTeams(),
      api.ListAWSOrganizations(publicRuntimeConfig.koreAdminTeamName),
      api.ListAllocations(publicRuntimeConfig.koreAdminTeamName),
      api.ListAccounts()
    ])
    allTeams.items = allTeams.items.filter(t => !publicRuntimeConfig.ignoreTeams.includes(t.metadata.name))
    awsOrganizations.items.forEach(org => {
      org.allocation = AllocationHelpers.findAllocationForResource(allAllocations, org)
      org.accountManagement = accountList.items.find(a => a.metadata.name === `am-${org.metadata.name}`)
    })
    return { resources: awsOrganizations, allTeams }
  }

  delete = (awsOrg) => () => {
    Modal.confirm({
      title: `Are you sure you want to delete the AWS Organization "${awsOrg.spec.ouName}"?`,
      content: 'This cannot be undone',
      okText: 'Yes',
      okType: 'danger',
      cancelText: 'No',
      onOk: async () => {
        const key = loadingMessage('Deleting', { duration: 0 })
        try {
          if (awsOrg.accountManagement) {
            loadingMessage('Deleting allocations for organization account management', { key })
            await AllocationHelpers.removeAllocation(awsOrg.accountManagement)
            loadingMessage('Deleting account management', { key })
            await (await KoreApi.client()).RemoveAccount(awsOrg.accountManagement.metadata.name)
          }

          loadingMessage('Deleting allocations for organization', { key })
          await AllocationHelpers.removeAllocation(awsOrg)
          loadingMessage('Deleting organization', { key })
          await (await KoreApi.client()).DeleteAWSOrganization(publicRuntimeConfig.koreAdminTeamName, awsOrg.metadata.name)
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
          message="Give Kore access to your Amazon Web Services organization"
          description="This will allow Kore to manage the organization for you. This includes managing the creation of Accounts and credentials giving Kore teams the ability to create clusters with ease."
          type="info"
          showIcon
          style={{ marginBottom: '20px' }}
        />
        {!resources ? <Icon type="loading" /> : (
          <>
            {resources.items.length === 0 && <Button type="primary" onClick={this.add(true)} style={{ display: 'block', marginBottom: '20px' }} className="new-aws-organization">Configure</Button>}
            <List
              dataSource={resources.items}
              renderItem={org =>
                <AWSOrganization
                  organization={org}
                  allTeams={allTeams.items}
                  editOrganization={this.edit}
                  deleteOrganization={this.delete}
                  handleUpdate={this.handleStatusUpdated}
                  handleDelete={() => {}}
                  refreshMs={2000}
                  propsResourceDataKey="organization"
                  resourceApiRequest={async () => await (await KoreApi.client()).GetAWSOrganization(publicRuntimeConfig.koreAdminTeamName, org.metadata.name)}
                />
              }
            >
            </List>
            {edit ? (
              <Drawer
                title={<Title level={4}>AWS Organization: {edit.spec.parentID}</Title>}
                visible={Boolean(edit)}
                onClose={this.edit(false)}
                width={900}
              >
                <AWSOrganizationForm
                  user={this.props.user}
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
                title={<Title level={4}>New AWS organization</Title>}
                visible={add}
                onClose={this.add(false)}
                width={900}
              >
                <AWSOrganizationForm
                  user={this.props.user}
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

export default AWSOrganizationsList
