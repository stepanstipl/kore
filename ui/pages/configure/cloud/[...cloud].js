import React from 'react'
import PropTypes from 'prop-types'
import { Tabs } from 'antd'
import Router from 'next/router' 

import Breadcrumb from '../../../lib/components/layout/Breadcrumb'
import GKECredentialsList from '../../../lib/components/credentials/GKECredentialsList'
import GCPOrganizationsList from '../../../lib/components/credentials/GCPOrganizationsList'
import EKSCredentialsList from '../../../lib/components/credentials/EKSCredentialsList'
import PlanList from '../../../lib/components/plans/PlanList'
import PolicyList from '../../../lib/components/policies/PolicyList'
import GCPProjectAutomationSettings from '../../../lib/components/setup/GCPProjectAutomationSettings'
import CloudTabs from '../../../lib/components/common/CloudTabs'
import CloudServiceAdmin from '../../../lib/components/services/CloudServiceAdmin'
import { featureEnabled, KoreFeatures } from '../../../lib/utils/features'

export default class ConfigureCloudPage extends React.Component {
  static propTypes = {
    selectedCloud: PropTypes.string.isRequired,
    activeKeys: PropTypes.object.isRequired
  }

  static defaultTabs = {
    selectedCloud: 'GCP',
    activeKeys: {
      'GCP': 'orgs',
      'AWS': 'accounts'
    }
  }

  static getInitialProps = async (ctx) => {
    const { cloud } = ctx.query
    if (!cloud || cloud.length === 0 || cloud.length === 1 && cloud[0] === 'index') {
      return ConfigureCloudPage.defaultTabs
    }
    if (cloud.length === 1) {
      return {
        ...ConfigureCloudPage.defaultTabs,
        selectedCloud: cloud[0]
      }
    }

    return {
      ...ConfigureCloudPage.defaultTabs,
      selectedCloud: cloud[0],
      activeKeys: {
        ...ConfigureCloudPage.defaultTabs.activeKeys,
        [cloud[0]]: cloud[1]
      }
    }
  }

  handleSelectCloud = cloud => {
    Router.push('/configure/cloud/[...cloud]', `/configure/cloud/${cloud}/${this.props.activeKeys[cloud]}`)
  }

  handleSelectKey = (cloud, key) => {
    Router.push('/configure/cloud/[...cloud]', `/configure/cloud/${cloud}/${key}`)
  }

  render() {
    const { selectedCloud, activeKeys } = this.props

    return (
      <>
        <Breadcrumb items={[{ text: 'Configure' }, { text: 'Cloud' }]} />
        <CloudTabs selectedKey={selectedCloud} handleSelectCloud={this.handleSelectCloud}/>
        <div id="cloud_subtabs">
          {selectedCloud === 'GCP' ? (
            <Tabs activeKey={activeKeys['GCP']} onChange={(key) => this.handleSelectKey('GCP', key)} destroyInactiveTabPane={true} tabPosition="left" style={{ marginTop: '20px' }}>
              <Tabs.TabPane tab="Organization credentials" key="orgs">
                <GCPOrganizationsList autoAllocateToAllTeams={true} />
              </Tabs.TabPane>
              <Tabs.TabPane tab="Project credentials" key="projects">
                <GKECredentialsList />
              </Tabs.TabPane>
              <Tabs.TabPane tab="Project automation" key="project_automation">
                <GCPProjectAutomationSettings />
              </Tabs.TabPane>
              <Tabs.TabPane tab="Cluster Plans" key="plans">
                <PlanList kind="GKE" />
              </Tabs.TabPane>
              <Tabs.TabPane tab="Cluster Policies" key="policies">
                <PolicyList kind="GKE"/>
              </Tabs.TabPane>
            </Tabs>
          ) : null}
          {selectedCloud === 'AWS' ? (
            <Tabs activeKey={activeKeys['AWS']} onChange={(key) => this.handleSelectKey('AWS', key)} destroyInactiveTabPane={true} tabPosition="left" style={{ marginTop: '20px' }}>
              <Tabs.TabPane tab="Account credentials" key="accounts">
                <EKSCredentialsList />
              </Tabs.TabPane>
              <Tabs.TabPane tab="Cluster Plans" key="plans">
                <PlanList kind="EKS" />
              </Tabs.TabPane>
              <Tabs.TabPane tab="Cluster Policies" key="policies">
                <PolicyList kind="EKS" />
              </Tabs.TabPane>
              {!featureEnabled(KoreFeatures.SERVICES) ? null :
                <Tabs.TabPane tab="Cloud Services" key="services">
                  <CloudServiceAdmin cloud="AWS" />
                </Tabs.TabPane>
              }
            </Tabs>
          ) : null}
        </div>
      </>
    )
  }
}