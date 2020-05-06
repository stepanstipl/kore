import React from 'react'
import PropTypes from 'prop-types'
import { Alert, Button, Icon, message, Modal, Radio, Typography } from 'antd'
const { Paragraph, Text } = Typography
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

import GCPKoreManagedProjectsCustom from './GCPKoreManagedProjectsCustom'
import RequestCredentialAccessForm from './forms/RequestCredentialAccessForm'
import V1beta1AccountManagement from '../../kore-api/model/V1beta1AccountManagement'
import V1beta1AccountManagementSpec from '../../kore-api/model/V1beta1AccountManagementSpec'
import V1ObjectMeta from '../../kore-api/model/V1ObjectMeta'
import V1Ownership from '../../kore-api/model/V1Ownership'
import V1beta1AccountsRule from '../../kore-api/model/V1beta1AccountsRule'
import AllocationHelpers from '../../utils/allocation-helpers'
import FormErrorMessage from '../forms/FormErrorMessage'
import KoreApi from '../../kore-api'
import canonical from '../../utils/canonical'
import copy from '../../utils/object-copy'
import asyncForEach from '../../utils/async-foreach'

class GCPProjectAutomationSettings extends React.Component {

  static propTypes = {
    tabActiveKey: PropTypes.string.isRequired,
    setTabActiveKey: PropTypes.func.isRequired
  }

  state = {
    dataLoading: true,
    submitting: false,
    errorMessage: false
  }

  async fetchComponentData() {
    const api = await KoreApi.client()
    const [ plansList, accountManagementList, gcpOrgs ] = await Promise.all([
      api.ListPlans('GKE'),
      api.ListAccounts(),
      api.ListGCPOrganizations(publicRuntimeConfig.koreAdminTeamName)
    ])

    const plans = plansList.items
    const gcpOrgList = gcpOrgs.items
    const accountManagement = accountManagementList.items.find(a => a.spec.provider === 'GKE')
    const gcpManagementType = accountManagement ? 'KORE' : 'EXISTING'

    let gcpProjectAutomationType = false
    let gcpProjectList = []

    if (accountManagement) {
      gcpProjectAutomationType = (accountManagement.spec.rules || []).length === 0 ? 'CLUSTER' : 'CUSTOM'
      gcpProjectList = (accountManagement.spec.rules || []).map(rule => ({ code: canonical(rule.name), ...rule }))
    }

    return { plans, accountManagement, gcpManagementType, gcpProjectList, gcpProjectAutomationType, gcpOrgList }
  }

  componentDidMount() {
    return this.fetchComponentData()
      .then(data => this.setState({ ...data, dataLoading: false }))
  }

  componentDidUpdate(prevProps) {
    // reload data if coming back from another tab
    if (prevProps.tabActiveKey !== this.props.tabActiveKey) {
      this.fetchComponentData().then(data => this.setState({ ...data }))
    }
  }

  selectGcpManagementType = e => this.setState({ gcpManagementType: e.target.value })

  selectGcpProjectAutomationType = e => this.setState({ gcpProjectAutomationType: e.target.value })

  generateAccountManagementResource = (gcpOrgResource, gcpProjectList) => {
    const resource = new V1beta1AccountManagement()
    resource.setApiVersion('accounts.kore.appvia.io/v1beta1')
    resource.setKind('AccountManagement')

    const meta = new V1ObjectMeta()
    meta.setName(gcpOrgResource.metadata.name)
    meta.setNamespace('kore-admin')
    if (this.state.accountManagement) {
      meta.setResourceVersion(this.state.accountManagement.metadata.resourceVersion)
    }
    resource.setMetadata(meta)

    const spec = new V1beta1AccountManagementSpec()
    spec.setProvider('GKE')

    const owner = new V1Ownership()
    const groupVersion = gcpOrgResource.apiVersion.split('/')
    owner.setGroup(groupVersion[0])
    owner.setVersion(groupVersion[1])
    owner.setKind(gcpOrgResource.kind)
    owner.setName(gcpOrgResource.metadata.name)
    owner.setNamespace(gcpOrgResource.metadata.namespace)
    spec.setOrganization(owner)

    if (gcpProjectList) {
      const rules = gcpProjectList.map(project => {
        const rule = new V1beta1AccountsRule()
        rule.setName(project.name)
        rule.setDescription(project.description)
        rule.setPrefix(project.prefix)
        rule.setSuffix(project.suffix)
        rule.setPlans(project.plans)
        return rule
      })
      spec.setRules(rules)
    }

    resource.setSpec(spec)

    return resource
  }

  disabledSave = () => {
    if (this.state.submitting || !this.state.gcpManagementType || !this.state.gcpProjectAutomationType) {
      return true
    }
    if (this.state.gcpProjectAutomationType === 'CUSTOM' && this.state.gcpProjectList.length === 0) {
      return true
    }
    return false
  }


  saveSettings = async () => {
    this.setState({ submitting: true, errorMessage: false })
    if (this.state.gcpManagementType === 'EXISTING') {
      // Noop - need to implement somewhere to store the request access email address
      this.setState({ submitting: false })
      message.success('Project automation settings saved')
      return
    }
    try {
      const api = await KoreApi.client()
      // create AccountManagement and Allocation CRDs, for each GCP Org
      // each GCP org will use the same settings
      await asyncForEach(this.state.gcpOrgList, async (gcpOrg) => {
        const gcpProjectList = this.state.gcpProjectAutomationType === 'CUSTOM' ? this.state.gcpProjectList : false
        const accountMgtResource = this.generateAccountManagementResource(gcpOrg, gcpProjectList)
        await api.UpdateAccount(gcpOrg.metadata.name, accountMgtResource)
        await AllocationHelpers.storeAllocation({ resourceToAllocate: accountMgtResource })
      })
      this.setState({ submitting: false })
      message.success('Project automation settings saved')
    } catch (error) {
      console.error('Error saving project automation settings', error)
      message.success('Failed to save project automation settings')
      this.setState({ submitting: false, errorMessage: 'A problem occurred trying to save, please try again later.' })
    }
  }

  handleResetGcpProjectList = (gcpProjectList) => this.setState({ gcpProjectList })

  handleGcpProjectAdded = (project) => {
    const code = canonical(project.name)
    this.setState({
      gcpProjectList: this.state.gcpProjectList.concat([{ code, plans: [], ...project }]),
    })
    message.success('GCP automated project added')
  }

  handleGcpProjectDeleted = (code) => {
    return () => {
      this.setState({
        gcpProjectList: this.state.gcpProjectList.filter(p => p.code !== code)
      })
      message.success('GCP automated project removed')
    }
  }

  handleGcpProjectChange = (code, property) => {
    return (value) => {
      const gcpProjectList = copy(this.state.gcpProjectList)
      gcpProjectList.find(p => p.code === code)[property] = value
      this.setState({ gcpProjectList })
    }
  }

  handleGcpProjectEdited = (code) => {
    return (project) => {
      project.code = canonical(project.name)
      const gcpProjectList = copy(this.state.gcpProjectList)
      const updatedGcpProject = gcpProjectList.find(p => p.code === code)
      Object.keys(project).forEach(k => updatedGcpProject[k] = project[k])
      this.setState({ gcpProjectList })
    }
  }

  projectAutomationHelp = () => {
    Modal.info({
      title: 'This defines how Kore will automate GCP projects for teams',
      content: (
        <div>
          <p>When a team is created in Kore and a cluster is requested, Kore will ensure the GCP project is also created and the cluster placed inside it.</p>
          <p>You can specify how the GCP projects should be created for Kore teams.</p>
        </div>
      ),
      onOk() {},
      width: 500
    })
  }

  projectAutomationClusterHelp = () => {
    Modal.info({
      title: 'Project automation: One per cluster',
      content: 'For every cluster a team creates Kore will also create a GCP project and provision the cluster inside it. The GCP project will share the name given to the cluster.',
      onOk() {},
      width: 500
    })
  }

  projectAutomationCustomHelp = () => {
    Modal.info({
      title: 'Project automation: Custom',
      content: (
        <div>
          <p>When a team is created in Kore and a cluster is requested, Kore will ensure the associated GCP project is also created and the cluster placed inside it.</p>
          <p>You must also specify the plans available for each type of project, this is to ensure the correct cluster specification is being used.</p>
        </div>
      ),
      onOk() {},
      width: 500
    })
  }

  gcpOrgRequired = () => (
    <Alert
      message="GCP organization credentials not found"
      description={
        <div>
          <Paragraph style={{ marginTop: '10px' }}>No GCP organization credentials have been configured in Kore. Without these, Kore will be unable to managed GCP projects on your behalf.</Paragraph>
          <Button type="secondary" onClick={() => this.props.setTabActiveKey('orgs')}>Setup organization credentials</Button>
        </div>
      }
      type="warning"
      showIcon
      style={{ marginTop: '10px' }}
    />
  )

  koreManagedProjectsSettings = () => (
    <>
      <Paragraph style={{ fontSize: '16px', fontWeight: '600' }}>Configure GCP project automation <Icon style={{ marginLeft: '5px' }} type="info-circle" theme="twoTone" onClick={this.projectAutomationHelp}/></Paragraph>

      <Paragraph style={{ fontSize: '16px', fontWeight: '600' }}>How do you want Kore to automate GCP projects for teams?</Paragraph>
      <Radio.Group onChange={this.selectGcpProjectAutomationType} value={this.state.gcpProjectAutomationType}>
        <div style={{ display: 'inline-block', marginRight: '20px' }}>
          <Radio value={'CLUSTER'} style={{ float: 'left' }}>
            <Text style={{ fontSize: '16px', fontWeight: '600' }}>One per cluster</Text>
            <Paragraph style={{ marginLeft: '24px', marginBottom: '0' }}>Kore will create a GCP project for each cluster a team provisions</Paragraph>
          </Radio>
          <Icon style={{ marginTop: '28px' }} type="info-circle" theme="twoTone" onClick={this.projectAutomationClusterHelp}/>
        </div>
        <div style={{ display: 'inline-block' }}>
          <Radio value={'CUSTOM'} style={{ float: 'left' }}>
            <Text style={{ fontSize: '16px', fontWeight: '600' }}>Custom</Text>
            <Paragraph style={{ marginLeft: '24px', marginBottom: '0' }}>Configure how Kore will create GCP projects for teams</Paragraph>
          </Radio>
          <Icon style={{ marginTop: '28px' }} type="info-circle" theme="twoTone" onClick={this.projectAutomationCustomHelp}/>
        </div>
      </Radio.Group>

      {this.state.gcpProjectAutomationType === 'CUSTOM' && (
        <GCPKoreManagedProjectsCustom
          gcpProjectList={this.state.gcpProjectList}
          plans={this.state.plans}
          handleChange={this.handleGcpProjectChange}
          handleDelete={this.handleGcpProjectDeleted}
          handleEdit={this.handleGcpProjectEdited}
          handleAdd={this.handleGcpProjectAdded}
          handleReset={this.handleResetGcpProjectList}
          hideGuidance={true}
        />
      )}
    </>
  )

  render() {
    const { dataLoading, gcpManagementType, submitting, errorMessage, gcpOrgList } = this.state
    if (dataLoading) {
      return <Icon type="loading" />
    }
    const gcpOrgsExist = Boolean(gcpOrgList.length)

    return (
      <>
        <FormErrorMessage message={errorMessage} />
        <div style={{ marginBottom: '15px' }}>
          <Paragraph style={{ fontSize: '16px', fontWeight: '600' }}>How do you want Kore teams to integrate with GCP projects?</Paragraph>
          <Radio.Group onChange={this.selectGcpManagementType} value={gcpManagementType} disabled={submitting}>
            <Radio value={'KORE'} style={{ marginRight: '20px' }}>
              <Text style={{ fontSize: '16px', fontWeight: '600' }}>Kore managed projects <Text type="secondary">(recommended)</Text></Text>
              <Paragraph style={{ marginLeft: '24px', marginBottom: '0' }}>Kore will manage the GCP projects required for teams</Paragraph>
            </Radio>
            <Radio value={'EXISTING'}>
              <Text style={{ fontSize: '16px', fontWeight: '600' }}>Use existing projects</Text>
              <Paragraph style={{ marginLeft: '24px', marginBottom: '0' }}>Kore teams will use existing GCP projects that it&apos;s given access to</Paragraph>
            </Radio>
          </Radio.Group>
        </div>
        {gcpManagementType === 'KORE' && !gcpOrgsExist && <this.gcpOrgRequired />}
        {gcpManagementType === 'KORE' && gcpOrgsExist && <this.koreManagedProjectsSettings />}
        {gcpManagementType === 'EXISTING' && <RequestCredentialAccessForm cloud="GKE" helpInModal={true} />}
        <Button style={{ marginTop: '20px', display: 'block' }} type="primary" loading={submitting} disabled={this.disabledSave()} onClick={this.saveSettings}>Save</Button>
      </>
    )
  }
}

export default GCPProjectAutomationSettings
