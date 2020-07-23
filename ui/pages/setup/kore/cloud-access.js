import PropTypes from 'prop-types'
import React from 'react'
import { Alert, Typography } from 'antd'
const { Title } = Typography
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

import CloudSelector from '../../../lib/components/common/CloudSelector'
import CloudSetup from '../../../lib/components/setup/CloudSetup'
import KoreApi from '../../../lib/kore-api'

class CloudAccessPage extends React.Component {

  static propTypes = {
    gkeCredentialsList: PropTypes.array.isRequired,
    eksCredentialsList: PropTypes.array.isRequired,
    aksCredentialsList: PropTypes.array.isRequired,
    gcpAccountManagement: PropTypes.object,
    awsAccountManagement: PropTypes.object,
  }

  static staticProps = {
    title: 'Setup cloud access',
    hideSider: true,
    adminOnly: true
  }

  static async getPageData(ctx) {
    try {
      const api = await KoreApi.client(ctx)
      const [ gkeCredentialsList, eksCredentialsList, aksCredentialsList, accountManagementList ] = await Promise.all([
        api.ListGKECredentials(publicRuntimeConfig.koreAdminTeamName),
        api.ListEKSCredentials(publicRuntimeConfig.koreAdminTeamName),
        api.ListAKSCredentials(publicRuntimeConfig.koreAdminTeamName),
        api.ListAccounts()
      ])
      const gcpAccountManagement = accountManagementList.items.find(a => a.spec.provider === 'GKE')
      const awsAccountManagement = accountManagementList.items.find(a => a.spec.provider === 'EKS')
      return {
        gkeCredentialsList: gkeCredentialsList.items,
        eksCredentialsList: eksCredentialsList.items,
        aksCredentialsList: aksCredentialsList.items,
        gcpAccountManagement,
        awsAccountManagement
      }
    } catch (err) {
      throw new Error(err.message)
    }
  }

  static getInitialProps = async (ctx) => {
    const data = await CloudAccessPage.getPageData(ctx)
    return data
  }

  state = {
    selectedCloud: false
  }

  handleSelectCloud = (cloud) => this.setState({ selectedCloud: cloud })

  render() {
    const { gcpAccountManagement, gkeCredentialsList, awsAccountManagement, eksCredentialsList, aksCredentialsList } = this.props
    const { selectedCloud } = this.state

    return (
      <>
        <Title>Setup cloud access</Title>
        <Alert
          message="Setup cloud access"
          description="Choose a cloud provider below to configure how Kore uses your cloud accounts."
          type="info"
          showIcon
          style={{ marginBottom: '20px' }}
        />
        <div style={{ marginTop: '20px', marginBottom: '20px' }}>
          <CloudSelector selectedCloud={selectedCloud} handleSelectCloud={this.handleSelectCloud} />
        </div>
        {selectedCloud === 'GCP' && <CloudSetup provider="GKE" cloud="GCP" accountNoun="project" accountManagement={gcpAccountManagement} credentialsList={gkeCredentialsList} />}
        {selectedCloud === 'AWS' && <CloudSetup provider="EKS" cloud="AWS" accountNoun="account" accountManagement={awsAccountManagement} credentialsList={eksCredentialsList} />}
        {selectedCloud === 'Azure' && <CloudSetup provider="AKS" cloud="Azure" accountNoun="subscription" credentialsList={aksCredentialsList} />}
      </>
    )
  }
}

export default CloudAccessPage
