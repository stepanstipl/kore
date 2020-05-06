import React from 'react'

import { Alert, Tabs } from 'antd'

import Breadcrumb from '../../lib/components/layout/Breadcrumb'
import GKECredentialsList from '../../lib/components/credentials/GKECredentialsList'
import GCPOrganizationsList from '../../lib/components/credentials/GCPOrganizationsList'
import EKSCredentialsList from '../../lib/components/credentials/EKSCredentialsList'
import PlanList from '../../lib/components/plans/PlanList'
import PolicyList from '../../lib/components/policies/PolicyList'
import GCPProjectAutomationSettings from '../../lib/components/setup/GCPProjectAutomationSettings'
import CloudTabs from '../../lib/components/common/CloudTabs'

class ConfigureCloudPage extends React.Component {

  state = {
    selectedCloud: 'GCP',
    activeKey: 'orgs'
  }

  handleSelectCloud = cloud => {
    this.setState({ selectedCloud: cloud })
  }

  render() {
    const { selectedCloud, activeKey } = this.state

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
          <Tabs activeKey={activeKey} onChange={(key) => this.setState({ activeKey: key })} tabPosition="left" style={{ marginTop: '20px' }}>
            <Tabs.TabPane tab="Organization credentials" key="orgs">
              <GCPOrganizationsList />
            </Tabs.TabPane>
            <Tabs.TabPane tab="Project credentials" key="projects">
              <GKECredentialsList />
            </Tabs.TabPane>
            <Tabs.TabPane tab="Plans" key="plans">
              <PlanList kind="GKE" />
            </Tabs.TabPane>
            <Tabs.TabPane tab="Policies" key="policies">
              <PolicyList kind="GKE" />
            </Tabs.TabPane>
            <Tabs.TabPane tab="Project automation" key="project_automation">
              <GCPProjectAutomationSettings tabActiveKey={activeKey} setTabActiveKey={(key) => this.setState({ activeKey: key })} />
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
            <Tabs.TabPane tab="Policies" key="policies">
              <PolicyList kind="EKS" />
            </Tabs.TabPane>
          </Tabs>
        ) : null}
      </>
    )
  }
}

export default ConfigureCloudPage
