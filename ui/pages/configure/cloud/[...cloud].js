import React from 'react'
import PropTypes from 'prop-types'
import { Tabs } from 'antd'
import Router from 'next/router' 

import Breadcrumb from '../../../lib/components/layout/Breadcrumb'
import CredentialsList from '../../../lib/components/credentials/CredentialsList'
import GCPOrganizationsList from '../../../lib/components/credentials/GCPOrganizationsList'
import AWSOrganizationsList from '../../../lib/components/credentials/AWSOrganizationsList'
import PlanList from '../../../lib/components/plans/PlanList'
import PolicyList from '../../../lib/components/policies/PolicyList'
import CloudAccountAutomationSettings from '../../../lib/components/setup/CloudAccountAutomationSettings'
import CloudTabs from '../../../lib/components/common/CloudTabs'
import CloudServiceAdmin from '../../../lib/components/services/CloudServiceAdmin'
import { featureEnabled, KoreFeatures } from '../../../lib/utils/features'

export default class ConfigureCloudPage extends React.Component {
  static propTypes = {
    user: PropTypes.object.isRequired,
    selectedCloud: PropTypes.string.isRequired,
    activeKeys: PropTypes.object.isRequired
  }

  static defaultTabs = {
    selectedCloud: 'GCP',
    activeKeys: {
      'GCP': 'orgs',
      'AWS': 'orgs',
      'Azure': 'subscriptions'
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

  handleSelectCloud = (cloud) => {
    Router.push('/configure/cloud/[...cloud]', `/configure/cloud/${cloud}/${this.props.activeKeys[cloud]}`)
  }

  handleSelectKey = (cloud, key) => {
    Router.push('/configure/cloud/[...cloud]', `/configure/cloud/${cloud}/${key}`)
  }

  render() {
    const { selectedCloud, activeKeys, user } = this.props

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
                <CredentialsList provider="GKE" />
              </Tabs.TabPane>
              <Tabs.TabPane tab="Project automation" key="project_automation">
                <CloudAccountAutomationSettings provider="GKE" cloudOrgsApiMethod="ListGCPOrganizations" cloud="GCP" accountNoun="project" />
              </Tabs.TabPane>
              <Tabs.TabPane tab="Cluster plans" key="plans">
                <PlanList kind="GKE" />
              </Tabs.TabPane>
              <Tabs.TabPane tab="Cluster policies" key="policies">
                <PolicyList kind="GKE"/>
              </Tabs.TabPane>
            </Tabs>
          ) : null}
          {selectedCloud === 'AWS' ? (
            <Tabs activeKey={activeKeys['AWS']} onChange={(key) => this.handleSelectKey('AWS', key)} destroyInactiveTabPane={true} tabPosition="left" style={{ marginTop: '20px' }}>
              <Tabs.TabPane tab="Organization credentials" key="orgs">
                <AWSOrganizationsList autoAllocateToAllTeams={true} user={user} />
              </Tabs.TabPane>
              <Tabs.TabPane tab="Account credentials" key="accounts">
                <CredentialsList provider="EKS" />
              </Tabs.TabPane>
              <Tabs.TabPane tab="Account automation" key="account-automation">
                <CloudAccountAutomationSettings provider="EKS" cloudOrgsApiMethod="ListAWSOrganizations" cloud="AWS" accountNoun="account" />
              </Tabs.TabPane>
              <Tabs.TabPane tab="Cluster plans" key="plans">
                <PlanList kind="EKS" />
              </Tabs.TabPane>
              <Tabs.TabPane tab="Cluster policies" key="policies">
                <PolicyList kind="EKS" />
              </Tabs.TabPane>
              {!featureEnabled(KoreFeatures.SERVICES) ? null :
                <Tabs.TabPane tab="Cloud services" key="services">
                  <CloudServiceAdmin cloud="AWS" />
                </Tabs.TabPane>
              }
            </Tabs>
          ) : null}
          {selectedCloud === 'Azure' ? (
            <Tabs activeKey={activeKeys['Azure']} onChange={(key) => this.handleSelectKey('Azure', key)} destroyInactiveTabPane={true} tabPosition="left" style={{ marginTop: '20px' }}>
              <Tabs.TabPane tab="Subscription credentials" key="subscriptions">
                <CredentialsList provider="AKS" />
              </Tabs.TabPane>
              <Tabs.TabPane tab="Cluster plans" key="plans">
                <PlanList kind="AKS" />
              </Tabs.TabPane>
              <Tabs.TabPane tab="Cluster policies" key="policies">
                <PolicyList kind="AKS" />
              </Tabs.TabPane>
            </Tabs>
          ) : null}
        </div>
      </>
    )
  }
}