import React from 'react'

import { Alert, Tabs } from 'antd'

import Breadcrumb from '../../lib/components/Breadcrumb'
import GKECredentialsList from '../../lib/components/configure/GKECredentialsList'
import GCPOrganizationsList from '../../lib/components/configure/GCPOrganizationsList'
import EKSCredentialsList from '../../lib/components/configure/EKSCredentialsList'
import PlanList from '../../lib/components/configure/PlanList'
import CloudTabs from '../../lib/components/configure/CloudTabs'

class ConfigureCloudPage extends React.Component {

  state = {
    selectedCloud: 'GCP'
  }

  handleSelectCloud = cloud => {
    this.setState({
      ...this.state,
      selectedCloud: cloud
    })
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
          <Tabs defaultActiveKey={'orgs'} tabPosition="left" style={{ marginTop: '20px' }}>
            <Tabs.TabPane tab="Organization credentials" key="orgs">
              <GCPOrganizationsList />
            </Tabs.TabPane>
            <Tabs.TabPane tab="Project credentials" key="projects">
              <GKECredentialsList />
            </Tabs.TabPane>
            <Tabs.TabPane tab="Plans" key="plans">
              <PlanList kind="GKE" />
            </Tabs.TabPane>
          </Tabs>
        ) : null}
        {selectedCloud === 'AWS' ? (
          <Tabs defaultActiveKey={'accounts'} tabPosition="left" style={{ marginTop: '20px' }}>
            <Tabs.TabPane tab="Account credentials" key="accounts">
              <EKSCredentialsList />
            </Tabs.TabPane>
            <Tabs.TabPane tab="Plans" key="plans">
              <PlanList kind="EKS" />
            </Tabs.TabPane>
          </Tabs>
        ) : null}
      </>
    )
  }
}

export default ConfigureCloudPage
