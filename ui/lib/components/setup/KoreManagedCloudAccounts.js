import React from 'react'
import PropTypes from 'prop-types'
import Link from 'next/link'
import { Alert, Button, Card, Divider, Icon, Result, Steps, Typography } from 'antd'
const { Paragraph } = Typography
const { Step } = Steps
import { pluralize, titleize } from 'inflect'

import KoreApi from '../../kore-api'
import copy from '../../utils/object-copy'
import canonical from '../../utils/canonical'
import asyncForEach from '../../utils/async-foreach'
import FormErrorMessage from '../forms/FormErrorMessage'
import AllocationHelpers from '../../utils/allocation-helpers'
import { successMessage } from '../../utils/message'
import KoreManagedCloudAccountsConfigure from './KoreManagedCloudAccountsConfigure'
import GCPOrganizationsList from '../credentials/GCPOrganizationsList'
import AWSOrganizationsList from '../credentials/AWSOrganizationsList'

class KoreManagedCloudAccounts extends React.Component {

  static propTypes = {
    user: PropTypes.object.isRequired,
    provider: PropTypes.oneOf(['GKE', 'EKS']).isRequired,
    cloud: PropTypes.oneOf(['GCP', 'AWS']),
    accountNoun: PropTypes.string.isRequired,
    accountManagement: PropTypes.object,
    setupComplete: PropTypes.bool.isRequired,
    handleSetupComplete: PropTypes.func.isRequired
  }

  steps = [
    { id: 'CREDS', title: 'Credentials', contentFn: 'stepsContentCreds', completeFn: 'stepsCompleteCreds' },
    { id: 'ACCOUNTS', title: `${titleize(this.props.accountNoun)} automation`, contentFn: 'stepsContentAccounts', completeFn: 'stepsCompleteAccounts' }
  ]

  constructor(props) {
    super(props)
    let cloudAccountList = []
    if (props.accountManagement) {
      cloudAccountList = (props.accountManagement.spec.rules || []).map(rule => ({ code: canonical(rule.name), ...rule }))
    }
    this.state = {
      cloudOrgList: [],
      currentStep: 0,
      cloudAccountList,
      plans: [],
      submitting: false,
      errorMessage: false
    }
  }

  async fetchComponentData() {
    const planList = await (await KoreApi.client()).ListPlans(this.props.provider)
    return planList.items
  }

  componentDidMount() {
    this.fetchComponentData().then(plans => this.setState({ plans }))
  }

  handleResetCloudAccountList = (cloudAccountList) => this.setState({ cloudAccountList })

  handleCloudAccountAdded = (project) => {
    const code = canonical(project.name)
    this.setState({
      cloudAccountList: this.state.cloudAccountList.concat([{ code, plans: [], ...project }]),
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
    return (project) => {
      project.code = canonical(project.name)
      const cloudAccountList = copy(this.state.cloudAccountList)
      const updatedCloudAccount = cloudAccountList.find(p => p.code === code)
      Object.keys(project).forEach(k => updatedCloudAccount[k] = project[k])
      this.setState({ cloudAccountList })
    }
  }

  nextStep() {
    const currentStep = this.state.currentStep + 1
    this.setState({ currentStep })
  }

  prevStep() {
    const currentStep = this.state.currentStep - 1
    this.setState({ currentStep })
  }

  stepsHeader = () => (
    <Steps current={this.state.currentStep}>
      {this.steps.map(item => <Step key={item.title} title={item.title} />)}
    </Steps>
  )

  stepsActions = () => (
    <div className="steps-action">
      {this.state.currentStep < this.steps.length - 1 && <Button type="primary" disabled={!this[this.steps[this.state.currentStep].completeFn]()} onClick={() => this.nextStep()}>Next</Button>}
      {this.state.currentStep === this.steps.length - 1 && <Button type="primary" loading={this.state.submitting} disabled={!this[this.steps[this.state.currentStep].completeFn]()} onClick={this.setupComplete}>Save</Button>}
      {this.state.currentStep > 0 && <Button style={{ marginLeft: 8 }} onClick={() => this.prevStep()}>Previous</Button>}
    </div>
  )

  stepsCompleteCreds = () => this.state.cloudOrgList.length >= 1

  stepsCompleteAccounts = () => {
    return this.state.cloudAccountList.length >= 1
  }

  setupComplete = async () => {
    this.setState({ submitting: true })
    try {
      const api = await KoreApi.client()
      // create AccountManagement and Allocation CRDs, for each cloud Org
      // each cloud org will use the same settings
      await asyncForEach(this.state.cloudOrgList, async (cloudOrg) => {
        const cloudAccountList = this.state.cloudAccountList
        const resourceVersion = this.props.accountManagement && this.props.accountManagement.metadata.resourceVersion
        const resourceName = `am-${this.props.cloud.toLowerCase()}`
        const accountMgtResource = KoreApi.resources().generateAccountManagementResource(resourceName, this.props.provider, cloudOrg, cloudAccountList, resourceVersion)
        await api.UpdateAccount(resourceName, accountMgtResource)
        await AllocationHelpers.storeAllocation({ resourceToAllocate: accountMgtResource })
      })
      this.setState({ submitting: false })
      this.props.handleSetupComplete()
    } catch (error) {
      console.error('Error submitting AccountManagement', error)
      this.setState({ submitting: false, errorMessage: 'A problem occurred trying to save, please try again later.' })
    }
  }

  stepContent = () => (
    <div className="steps-content" style={{ marginTop: '20px', marginBottom: '20px' }}>
      <FormErrorMessage message={this.state.errorMessage} />
      {this[this.steps[this.state.currentStep].contentFn]()}
    </div>
  )

  stepsContentCreds = () => (
    <>
      <Paragraph style={{ fontSize: '16px', fontWeight: '600' }}>Configure your {this.props.cloud} organization credentials</Paragraph>
      {this.props.cloud === 'GCP' && <GCPOrganizationsList getResourceItemList={(list) => this.setState({ cloudOrgList: list })} autoAllocateToAllTeams={true} />}
      {this.props.cloud === 'AWS' && <AWSOrganizationsList getResourceItemList={(list) => this.setState({ cloudOrgList: list })} autoAllocateToAllTeams={true} user={this.props.user} />}
    </>
  )

  stepsContentAccounts = () => (
    <>
      <Paragraph style={{ fontSize: '16px', fontWeight: '600' }}>Configure {this.props.cloud} {this.props.accountNoun} automation</Paragraph>
      <Alert
        message={`This defines how Kore will automate ${this.props.cloud} ${pluralize(this.props.accountNoun)} for teams`}
        description={`When a team is created in Kore and a cluster is requested, Kore will ensure the associated ${this.props.cloud} ${this.props.accountNoun} is also created and the cluster placed inside it. You must also specify the plans available for each type of ${this.props.accountNoun}, this is to ensure the correct cluster specification is being used.`}
        type="info"
        showIcon
        style={{ marginBottom: '20px' }}
      />

      <KoreManagedCloudAccountsConfigure
        cloudAccountList={this.state.cloudAccountList}
        plans={this.state.plans}
        handleChange={this.handleCloudAccountChange}
        handleDelete={this.handleCloudAccountDeleted}
        handleEdit={this.handleCloudAccountEdited}
        handleAdd={this.handleCloudAccountAdded}
        handleReset={this.handleResetCloudAccountList}
        cloud={this.props.cloud}
        accountNoun={this.props.accountNoun}
      />
    </>
  )

  renderSetupComplete = () => (
    <Card className="kore-managed-setup-complete">
      <Result
        status="success"
        title="Setup complete!"
        subTitle={`Kore will manage the lifecycle of ${this.props.cloud} ${pluralize(this.props.accountNoun)} along with Kore teams`}
        extra={<Link href="/setup/kore/complete"><Button type="primary" key="continue">Continue</Button></Link>}
      >
        <Paragraph><Icon type="check-circle" theme="twoTone" twoToneColor="#52c41a" /> {this.props.cloud} organization credentials</Paragraph>
        <Paragraph style={{ marginBottom: '0' }}><Icon type="check-circle" theme="twoTone" twoToneColor="#52c41a" /> {titleize(this.props.accountNoun)} automation</Paragraph>
      </Result>
    </Card>
  )

  render() {
    if (this.props.setupComplete) {
      return this.renderSetupComplete()
    }

    return (
      <Card>
        {this.stepsHeader()}
        <Divider />
        {this.stepContent()}
        <Divider />
        {this.stepsActions()}
      </Card>
    )
  }
}

export default KoreManagedCloudAccounts
