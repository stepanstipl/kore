import React from 'react'

import { Alert } from 'antd'

import Breadcrumb from '../../lib/components/Breadcrumb'
import GKECredentialsList from '../../lib/components/configure/GKECredentialsList'
import GCPOrganizationsList from '../../lib/components/configure/GCPOrganizationsList'
import CloudTabs from '../../lib/components/configure/CloudTabs'
import copy from '../../lib/utils/object-copy'

class ConfigureCloudPage extends React.Component {

  state = {
    selectedCloud: 'GCP'
  }

  handleSelectCloud = cloud => {
    if (this.state.selectedCloud !== cloud) {
      const state = copy(this.state)
      state.selectedCloud = cloud
      this.setState(state)
    }
  }

  render() {
    const { selectedCloud } = this.state

    return (
      <>
        <Breadcrumb items={[{ text: 'Configure' }, { text: 'Cloud' }]} />
        <Alert
          message="Select the cloud provider to configure the settings"
          type="info"
          style={{ marginBottom: '20px' }}
        />
        <CloudTabs defaultSelectedKey={selectedCloud} handleSelectCloud={this.handleSelectCloud}/>
        {selectedCloud === 'GCP' ? (
          <>
            <GCPOrganizationsList style={{ marginBottom: '20px' }} />
            <GKECredentialsList />
          </>
        ) : null}
      </>
    )
  }
}

export default ConfigureCloudPage
