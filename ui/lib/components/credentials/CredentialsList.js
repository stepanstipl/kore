import PropTypes from 'prop-types'
import { Typography, List, Button, Drawer, Alert, Icon, Modal } from 'antd'
const { Title } = Typography
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()
import { pluralize } from 'inflect'

import Credentials from './Credentials'
import ResourceList from '../resources/ResourceList'
import GKECredentialsForm from './GKECredentialsForm'
import EKSCredentialsForm from './EKSCredentialsForm'
import AKSCredentialsForm from './AKSCredentialsForm'
import KoreApi from '../../kore-api'
import AllocationHelpers from '../../utils/allocation-helpers'
import { errorMessage, successMessage, loadingMessage } from '../../utils/message'
import { getProviderCloudInfo } from '../../utils/cloud'

class CredentialsList extends ResourceList {

  static propTypes = {
    provider: PropTypes.oneOf(['GKE', 'EKS', 'AKS']),
    style: PropTypes.object
  }

  cloudInfo = getProviderCloudInfo(this.props.provider)

  createdMessage = `${this.cloudInfo.cloud} ${this.cloudInfo.accountNoun} credentials created successfully`
  updatedMessage = `${this.cloudInfo.cloud} ${this.cloudInfo.accountNoun} credentials updated successfully`
  deletedMessage = `${this.cloudInfo.cloud} ${this.cloudInfo.accountNoun} credentials deleted successfully`
  deleteFailedMessage = `Error deleting ${this.cloudInfo.cloud} ${this.cloudInfo.accountNoun} credentials`

  async fetchComponentData() {
    const api = await KoreApi.client()
    const [ allTeams, credentials, allAllocations ] = await Promise.all([
      api.ListTeams(),
      api[`List${this.props.provider}Credentials`](publicRuntimeConfig.koreAdminTeamName),
      api.ListAllocations(publicRuntimeConfig.koreAdminTeamName)
    ])
    allTeams.items = allTeams.items.filter(t => !publicRuntimeConfig.ignoreTeams.includes(t.metadata.name))
    credentials.items.forEach(credential => {
      credential.allocation = AllocationHelpers.findAllocationForResource(allAllocations, credential)
    })
    return { resources: credentials, allTeams }
  }

  delete = (cred) => () => {
    Modal.confirm({
      title: `Are you sure you want to delete the credentials: ${cred.spec[this.cloudInfo.credentialsIdentifierKey]}?`,
      content: 'This cannot be undone',
      okText: 'Yes',
      okType: 'danger',
      cancelText: 'No',
      onOk: async () => {
        const key = loadingMessage('Deleting allocations for credential', { duration: 0 })
        try {
          await AllocationHelpers.removeAllocation(cred)
          loadingMessage('Deleting credential', { key, duration: 0 })
          await (await KoreApi.client())[`Delete${this.props.provider}Credentials`](publicRuntimeConfig.koreAdminTeamName, cred.metadata.name)
          successMessage(this.deletedMessage, { key })
        } catch (err) {
          console.error(err)
          errorMessage(this.deleteFailedMessage, { key })
        }
        await this.refresh()
      }
    })
  }

  drawerClose = () => {
    if (this.state.add) {
      this.add(false)()
    }
    if (this.state.edit) {
      this.edit(false)()
    }
  }

  drawerVisible = () => Boolean(this.state.edit || this.state.add)

  drawerTitle = () => {
    if (this.state.add) {
      return <Title level={4}>New {this.cloudInfo.cloud} {this.cloudInfo.accountNoun}</Title>
    }
    if (this.state.edit) {
      return <Title level={4}>{this.cloudInfo.cloud} {this.cloudInfo.accountNoun}: {this.state.edit.spec[this.cloudInfo.credentialsIdentifierKey]}</Title>
    }
  }

  drawerHandleSubmit = () => {
    if (this.state.add) {
      return this.handleAddSave
    }
    if (this.state.edit) {
      return this.handleEditSave
    }
  }

  render() {
    const { resources, allTeams, edit, add } = this.state

    return (
      <>
        <Alert
          message={`Give Kore access to your existing  ${this.cloudInfo.cloudLong} ${pluralize(this.cloudInfo.accountNoun)}`}
          description={`This will enable Kore to build clusters inside a ${this.cloudInfo.cloud} ${this.cloudInfo.accountNoun} that you already manage outside of Kore. You must create a Service Account inside your ${this.cloudInfo.accountNoun} and add the key in JSON format here.`}
          type="info"
          showIcon
          style={{ marginBottom: '20px' }}
        />
        <Button id="add" type="primary" onClick={this.add(true)} style={{ display: 'block', marginBottom: '20px' }}>+ New</Button>
        {!resources ? <Icon type="loading" /> : (
          <>
            <List
              id={`${this.props.provider.toLowerCase()}creds_list`}
              dataSource={resources.items}
              renderItem={eks =>
                <Credentials
                  provider={this.props.provider}
                  identifierKey={this.cloudInfo.credentialsIdentifierKey}
                  credentials={eks}
                  allTeams={allTeams.items}
                  editCredential={this.edit}
                  deleteCredential={this.delete}
                  handleUpdate={this.handleStatusUpdated}
                  handleDelete={() => {}}
                  refreshMs={2000}
                  propsResourceDataKey="credentials"
                  resourceApiRequest={async () => await (await KoreApi.client())[`Get${this.props.provider}Credentials`](publicRuntimeConfig.koreAdminTeamName, eks.metadata.name)}
                />
              }
            >
            </List>
            <Drawer
              title={this.drawerTitle()}
              visible={this.drawerVisible()}
              onClose={this.drawerClose}
              width={700}
            >
              {(edit || add) ? (
                <>
                  {this.props.provider === 'GKE' ? <GKECredentialsForm
                    team={publicRuntimeConfig.koreAdminTeamName}
                    allTeams={allTeams}
                    data={edit || undefined}
                    handleSubmit={this.drawerHandleSubmit()}
                  /> : null}
                  {this.props.provider === 'EKS' ? <EKSCredentialsForm
                    team={publicRuntimeConfig.koreAdminTeamName}
                    allTeams={allTeams}
                    data={edit || undefined}
                    handleSubmit={this.drawerHandleSubmit()}
                  /> : null}
                  {this.props.provider === 'AKS' ? <AKSCredentialsForm
                    team={publicRuntimeConfig.koreAdminTeamName}
                    allTeams={allTeams}
                    data={edit || undefined}
                    handleSubmit={this.drawerHandleSubmit()}
                  /> : null}
                </>
              ) : null}
            </Drawer>
          </>
        )}
      </>
    )
  }
}

export default CredentialsList
