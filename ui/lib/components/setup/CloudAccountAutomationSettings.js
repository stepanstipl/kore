import React from 'react'
import PropTypes from 'prop-types'
import { Alert, Button, Icon, Modal, Typography } from 'antd'
const { Paragraph } = Typography
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()
import Link from 'next/link'
import { titleize, pluralize } from 'inflect'

import KoreManagedCloudAccountsCustom from './KoreManagedCloudAccountsCustom'
import RequestCredentialAccessForm from './forms/RequestCredentialAccessForm'
import AllocationHelpers from '../../utils/allocation-helpers'
import FormErrorMessage from '../forms/FormErrorMessage'
import KoreApi from '../../kore-api'
import canonical from '../../utils/canonical'
import copy from '../../utils/object-copy'
import asyncForEach from '../../utils/async-foreach'
import { errorMessage, successMessage } from '../../utils/message'
import KoreTeamCloudIntegration from './radio-groups/KoreTeamCloudIntegration'
import CloudAccountAutomationType from './radio-groups/CloudAccountAutomationType'
import { filterCloudAccountList } from '../../utils/cloud'

class CloudAccountAutomationSettings extends React.Component {
  static propTypes = {
    provider: PropTypes.oneOf(['GKE', 'EKS']).isRequired,
    cloudOrgsApiMethod: PropTypes.string.isRequired,
    cloud: PropTypes.oneOf(['GCP', 'AWS']).isRequired,
    accountNoun: PropTypes.string.isRequired,
    cloudAccountAutomationType: PropTypes.oneOf(['CLUSTER', 'CUSTOM'])
  }

  state = {
    dataLoading: true,
    submitting: false,
    errorMessage: false,
    email: undefined,
    emailValid: false
  }

  async fetchComponentData() {
    const api = await KoreApi.client()
    const [ plansList, accountManagementList, cloudOrgs, cloudConfig ] = await Promise.all([
      api.ListPlans(this.props.provider),
      api.ListAccounts(),
      api[this.props.cloudOrgsApiMethod](publicRuntimeConfig.koreAdminTeamName),
      api.GetConfig(this.props.cloud)
    ])

    const plans = plansList.items
    const cloudOrgList = cloudOrgs.items
    const accountManagement = accountManagementList.items.find(a => a.spec.provider === this.props.provider)
    const cloudManagementType = accountManagement ? 'KORE' : 'EXISTING'

    let cloudAccountAutomationType = this.props.cloudAccountAutomationType || false
    let cloudAccountList = []

    if (accountManagement) {
      cloudAccountAutomationType = (accountManagement.spec.rules || []).length === 0 ? 'CLUSTER' : 'CUSTOM'
      cloudAccountList = (accountManagement.spec.rules || []).map(rule => ({ code: canonical(rule.name), ...rule }))
    }
    const email = cloudConfig && cloudConfig.spec.values.requestAccessEmail
    const emailValid = email ? true : false
    return { plans, accountManagement, cloudManagementType, cloudAccountList, cloudAccountAutomationType, cloudOrgList, email, emailValid }
  }

  componentDidMount() {
    return this.fetchComponentData()
      .then(data => this.setState({ ...data, dataLoading: false }))
  }

  selectCloudManagementType = e => this.setState({ cloudManagementType: e.target.value })

  selectCloudAccountAutomationType = e => this.setState({ cloudAccountAutomationType: e.target.value })

  disabledSave = () => {
    if (this.state.submitting || !this.state.cloudManagementType) {
      return true
    }
    if (this.state.cloudManagementType === 'EXISTING') {
      return this.state.emailValid ? false : true
    }
    if (this.state.cloudManagementType === 'KORE' && !this.state.cloudAccountAutomationType) {
      return true
    }
    if (this.state.cloudAccountAutomationType === 'CUSTOM' && this.state.cloudAccountList.length === 0) {
      return true
    }
    return false
  }

  saveSettings = async () => {
    this.setState({ submitting: true, errorMessage: false })
    if (this.state.cloudManagementType === 'EXISTING') {
      try {
        const api = await KoreApi.client()
        const config = { requestAccessEmail: this.state.email }
        await api.UpdateConfig(this.props.cloud, KoreApi.resources().generateConfigResource(this.props.cloud, config))

        if (this.state.accountManagement) {
          // disable account management by deleting the CRD and Allocation
          await api.RemoveAccount(this.state.accountManagement.metadata.name)
          await AllocationHelpers.removeAllocation(this.state.accountManagement)
          this.setState({
            submitting: false,
            cloudAccountAutomationType: false,
            accountManagement: false,
            cloudAccountList: []
          })
          return successMessage(`${titleize(this.props.accountNoun)} automation settings saved`)
        } else {
          this.setState({ submitting: false })
          return successMessage(`${titleize(this.props.accountNoun)} automation settings saved`)
        }
      } catch (err) {
        console.error(`Error saving ${this.props.accountNoun} automation settings`, err)
        errorMessage(`Failed to save ${this.props.accountNoun} automation settings`)
        return this.setState({
          submitting: false,
          errorMessage: 'A problem occurred trying to save, please try again later.'
        })
      }
    }
    try {
      const api = await KoreApi.client()
      // create AccountManagement and Allocation CRDs, for each cloud Org
      // each cloud org will use the same settings
      await asyncForEach(this.state.cloudOrgList, async (cloudOrg) => {
        let cloudAccountList = this.state.cloudAccountAutomationType === 'CUSTOM' ? copy(this.state.cloudAccountList) : false
        cloudAccountList = filterCloudAccountList(cloudAccountList, this.state.plans)
        const resourceVersion = this.state.accountManagement && this.state.accountManagement.metadata.resourceVersion
        const resourceName = `am-${this.props.cloud.toLowerCase()}`
        const accountMgtResource = KoreApi.resources().generateAccountManagementResource(resourceName, this.props.provider, cloudOrg, cloudAccountList, resourceVersion)
        await api.UpdateAccount(resourceName, accountMgtResource)
        await AllocationHelpers.storeAllocation({ resourceToAllocate: accountMgtResource })
        this.setState({ submitting: false, accountManagement: accountMgtResource, cloudAccountList })
      })
      successMessage(`${titleize(this.props.accountNoun)} automation settings saved`)
    } catch (error) {
      console.error(`Error saving ${this.props.accountNoun} automation settings`, error)
      errorMessage(`Failed to save ${this.props.accountNoun} automation settings`)
      this.setState({ submitting: false, errorMessage: 'A problem occurred trying to save, please try again later.' })
    }
  }

  handleResetCloudAccountList = (cloudAccountList) => this.setState({ cloudAccountList })

  handleCloudAccountAdded = (cloudAccount) => {
    const code = canonical(cloudAccount.name)
    this.setState({
      cloudAccountList: this.state.cloudAccountList.concat([{ code, plans: [], ...cloudAccount }]),
    })
    successMessage(`${this.props.cloud} automated ${this.props.accountNoun} added`)
  }

  handleCloudAccountDeleted = (code) => {
    return () => {
      this.setState({
        cloudAccountList: this.state.cloudAccountList.filter(p => p.code !== code)
      })
      successMessage(`${this.props.cloud} automated ${this.props.accountNoun} removed`)
    }
  }

  handleCloudAccountChange = (code, property) => {
    return (value) => {
      const cloudAccountList = copy(this.state.cloudAccountList)
      cloudAccountList.find(p => p.code === code)[property] = value
      this.setState({ cloudAccountList })
    }
  }

  handleCloudAccountEdited = (code) => {
    return (cloudAccount) => {
      cloudAccount.code = canonical(cloudAccount.name)
      const cloudAccountList = copy(this.state.cloudAccountList)
      const updatedCloudAccount = cloudAccountList.find(p => p.code === code)
      Object.keys(cloudAccount).forEach(k => updatedCloudAccount[k] = cloudAccount[k])
      this.setState({ cloudAccountList })
    }
  }

  accountAutomationHelp = () => {
    Modal.info({
      title: `This defines how Kore will automate ${this.props.cloud} ${pluralize(this.props.accountNoun)} for teams`,
      content: (
        <div>
          <p>When a team is created in Kore and a cluster is requested, Kore will ensure the {this.props.cloud} {this.props.accountNoun} is also created and the cluster placed inside it.</p>
          <p>You can specify how the {this.props.cloud} {pluralize(this.props.accountNoun)} should be created for Kore teams.</p>
        </div>
      ),
      onOk() {},
      width: 500
    })
  }

  cloudOrgRequired = () => (
    <Alert
      message={`${this.props.cloud} organization credentials not found`}
      description={
        <div>
          <Paragraph style={{ marginTop: '10px' }}>No {this.props.cloud} organization credentials have been configured in Kore. Without these, Kore will be unable to managed {this.props.cloud} projects on your behalf.</Paragraph>
          <Link href="/configure/cloud/[...cloud]" as={`/configure/cloud/${this.props.cloud}/orgs`}><Button type="secondary">Setup organization credentials</Button></Link>
        </div>
      }
      type="warning"
      showIcon
      style={{ marginTop: '10px' }}
    />
  )

  koreManagedAccountSettings = () => (
    <>
      <Paragraph style={{ fontSize: '16px', fontWeight: '600' }}>Configure {this.props.cloud} {this.props.accountNoun} automation <Icon style={{ marginLeft: '5px' }} type="info-circle" theme="twoTone" onClick={this.accountAutomationHelp}/></Paragraph>

      <CloudAccountAutomationType
        cloud={this.props.cloud}
        accountNoun={this.props.accountNoun}
        value={this.state.cloudAccountAutomationType}
        onChange={this.selectCloudAccountAutomationType}
        inlineHelp={true}
        valuesFilter={this.props.cloudAccountAutomationType ? [this.props.cloudAccountAutomationType] : undefined}
      />

      {this.state.cloudAccountAutomationType === 'CUSTOM' && (
        <KoreManagedCloudAccountsCustom
          cloudAccountList={this.state.cloudAccountList}
          plans={this.state.plans}
          handleChange={this.handleCloudAccountChange}
          handleDelete={this.handleCloudAccountDeleted}
          handleEdit={this.handleCloudAccountEdited}
          handleAdd={this.handleCloudAccountAdded}
          handleReset={this.handleResetCloudAccountList}
          hideGuidance={true}
          cloud={this.props.cloud}
          accountNoun={this.props.accountNoun}
        />
      )}
    </>
  )

  render() {
    const { dataLoading, cloudManagementType, submitting, errorMessage, cloudOrgList, email } = this.state
    if (dataLoading) {
      return <Icon type="loading" />
    }
    const cloudOrgsExist = Boolean(cloudOrgList.length)

    return (
      <>
        <FormErrorMessage message={errorMessage} />

        <KoreTeamCloudIntegration
          cloud={this.props.cloud}
          accountNoun={this.props.accountNoun}
          value={cloudManagementType}
          disabled={submitting}
          onChange={this.selectCloudManagementType}
        />

        {cloudManagementType === 'KORE' && !cloudOrgsExist && this.cloudOrgRequired()}
        {cloudManagementType === 'KORE' && cloudOrgsExist && this.koreManagedAccountSettings()}
        {cloudManagementType === 'EXISTING' ? (
          <RequestCredentialAccessForm
            data={{ email }}
            cloud={this.props.cloud}
            helpInModal={true}
            onChange={(email, errors) => this.setState({ email, emailValid: Boolean(!errors) })}
          />
        ) : null}
        <Button style={{ marginTop: '20px', display: 'block' }} type="primary" loading={submitting} disabled={this.disabledSave()} onClick={this.saveSettings}>Save</Button>
      </>
    )
  }
}

export default CloudAccountAutomationSettings
