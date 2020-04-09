import PropTypes from 'prop-types'
import { Typography, List, Button, Drawer, Alert, Icon } from 'antd'
const { Title } = Typography

import { kore } from '../../../config'
import ResourceList from '../configure/ResourceList'
import EKSCredentialsForm from '../forms/EKSCredentialsForm'
import EKSCredentials from '../team/EKSCredentials'
import KoreApi from '../../kore-api'

class EKSCredentialsList extends ResourceList {

  static propTypes = {
    style: PropTypes.object
  }

  createdMessage = 'AWS account credentials created successfully'
  updatedMessage = 'AWS account credentials updated successfully'

  async fetchComponentData() {
    const api = await KoreApi.client()
    const [ allTeams, eksCredentials, allAllocations ] = await Promise.all([
      api.ListTeams(),
      api.ListEKSCredentials(kore.koreAdminTeamName),
      api.ListAllocations(kore.koreAdminTeamName)
    ])
    allTeams.items = allTeams.items.filter(t => !kore.ignoreTeams.includes(t.metadata.name))
    eksCredentials.items.forEach(eks => {
      eks.allocation = (allAllocations.items || []).find(alloc => alloc.metadata.name === eks.metadata.name)
    })
    return { resources: eksCredentials, allTeams }
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
        <Button type="primary" onClick={this.add(true)} style={{ display: 'block', marginBottom: '20px' }}>+ New</Button>
        {!resources ? <Icon type="loading" /> : (
          <>
            <List
              dataSource={resources.items}
              renderItem={eks =>
                <EKSCredentials
                  eksCredentials={eks}
                  allTeams={allTeams.items}
                  editEKSCredentials={this.edit}
                  handleUpdate={this.handleStatusUpdated}
                  refreshMs={2000}
                  propsResourceDataKey="eksCredentials"
                  resourceApiPath={`/teams/${kore.koreAdminTeamName}/ekscredentials/${eks.metadata.name}`}
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
                  team={kore.koreAdminTeamName}
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

export default EKSCredentialsList
