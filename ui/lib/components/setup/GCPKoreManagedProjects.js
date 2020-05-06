import React from 'react'
import PropTypes from 'prop-types'
import Link from 'next/link'
import { Alert, Button, Card, Divider, Icon, message, Radio, Result, Steps, Typography } from 'antd'
const { Paragraph, Text } = Typography
const { Step } = Steps

import KoreApi from '../../kore-api'
import copy from '../../utils/object-copy'
import canonical from '../../utils/canonical'
import asyncForEach from '../../utils/async-foreach'
import GCPOrganizationsList from '../credentials/GCPOrganizationsList'
import GCPKoreManagedProjectsCustom from './GCPKoreManagedProjectsCustom'
import FormErrorMessage from '../forms/FormErrorMessage'
import V1beta1AccountManagement from '../../kore-api/model/V1beta1AccountManagement'
import V1beta1AccountManagementSpec from '../../kore-api/model/V1beta1AccountManagementSpec'
import V1ObjectMeta from '../../kore-api/model/V1ObjectMeta'
import V1Ownership from '../../kore-api/model/V1Ownership'
import V1beta1AccountsRule from '../../kore-api/model/V1beta1AccountsRule'
import AllocationHelpers from '../../utils/allocation-helpers'

class GCPKoreManagedProjects extends React.Component {

  static propTypes = {
    accountManagement: PropTypes.object
  }

  static propTypes = {
    setupComplete: PropTypes.bool.isRequired,
    handleSetupComplete: PropTypes.func.isRequired
  }

  steps = [
    { id: 'CREDS', title: 'Credentials', contentFn: 'stepsContentCreds', completeFn: 'stepsCompleteCreds' },
    { id: 'PROJECTS', title: 'Project automation', contentFn: 'stepsContentProjects', completeFn: 'stepsCompleteProjects' }
  ]

  constructor(props) {
    super(props)
    let gcpProjectAutomationType = false
    let gcpProjectList = []
    if (props.accountManagement) {
      gcpProjectAutomationType = (props.accountManagement.spec.rules || []).length === 0 ? 'CLUSTER' : 'CUSTOM'
      gcpProjectList = (props.accountManagement.spec.rules || []).map(rule => ({ code: canonical(rule.name), ...rule }))
    }
    this.state = {
      gcpOrgList: [],
      currentStep: 0,
      gcpProjectAutomationType,
      gcpProjectList,
      plans: [],
      submitting: false,
      errorMessage: false
    }
  }

  async fetchComponentData() {
    const api = await KoreApi.client()
    const planList = await api.ListPlans('GKE')
    return planList.items
  }

  componentDidMount() {
    this.fetchComponentData().then(plans => this.setState({ plans }))
  }

  selectGcpProjectAutomationType = e => this.setState({ gcpProjectAutomationType: e.target.value })

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

  stepsCompleteCreds = () => this.state.gcpOrgList.length >= 1

  stepsCompleteProjects = () => {
    if (this.state.gcpProjectAutomationType === 'CLUSTER') {
      return true
    }
    return this.state.gcpProjectList.length >= 1
  }

  generateAccountManagementResource = (gcpOrgResource, gcpProjectList) => {
    const resource = new V1beta1AccountManagement()
    resource.setApiVersion('accounts.kore.appvia.io/v1beta1')
    resource.setKind('AccountManagement')

    const meta = new V1ObjectMeta()
    meta.setName(`am-${gcpOrgResource.metadata.name}`)
    meta.setNamespace('kore-admin')
    if (this.props.accountManagement) {
      meta.setResourceVersion(this.props.accountManagement.metadata.resourceVersion)
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

  setupComplete = async () => {
    this.setState({ submitting: true })
    try {
      const api = await KoreApi.client()
      // create AccountManagement and Allocation CRDs, for each GCP Org
      // each GCP org will use the same settings
      await asyncForEach(this.state.gcpOrgList, async (gcpOrg) => {
        const gcpProjectList = this.state.gcpProjectAutomationType === 'CUSTOM' ? this.state.gcpProjectList : false
        const accountMgtResource = this.generateAccountManagementResource(gcpOrg, gcpProjectList)
        await api.UpdateAccount(`am-${gcpOrg.metadata.name}`, accountMgtResource)
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
      <Paragraph style={{ fontSize: '16px', fontWeight: '600' }}>Configure your GCP organization credentials</Paragraph>
      <GCPOrganizationsList getResourceItemList={this.setGcpOrgList} autoAllocateToAllTeams={true} />
    </>
  )

  setGcpOrgList = (list) => this.setState({ gcpOrgList: list })

  stepsContentProjects = () => (
    <>
      <Paragraph style={{ fontSize: '16px', fontWeight: '600' }}>Configure GCP project automation</Paragraph>
      <Alert
        message="This defines how Kore will automate GCP projects for teams"
        description="When a team is created in Kore and a cluster is requested, Kore will ensure the GCP project is also created and the cluster placed inside it. You must also specify the plans available for each type of project, this is to ensure the correct cluster specification is being used."
        type="info"
        showIcon
        style={{ marginBottom: '20px' }}
      />

      <Paragraph style={{ fontSize: '16px', fontWeight: '600' }}>How do you want Kore to automate GCP projects for teams?</Paragraph>
      <Radio.Group onChange={this.selectGcpProjectAutomationType} value={this.state.gcpProjectAutomationType}>
        <Radio className="automated-projects-cluster" value={'CLUSTER'} style={{ marginRight: '20px' }}>
          <Text style={{ fontSize: '16px', fontWeight: '600' }}>One per cluster</Text>
          <Paragraph style={{ marginLeft: '24px', marginBottom: '0' }}>Kore will create a GCP project for each cluster a team provisions</Paragraph>
        </Radio>
        <Radio className="automated-projects-custom" value={'CUSTOM'}>
          <Text style={{ fontSize: '16px', fontWeight: '600' }}>Custom</Text>
          <Paragraph style={{ marginLeft: '24px', marginBottom: '0' }}>Configure how Kore will create GCP projects for teams</Paragraph>
        </Radio>
      </Radio.Group>

      {this.state.gcpProjectAutomationType === 'CLUSTER' && (
        <Alert
          message="For every cluster a team creates Kore will also create a GCP project and provision the cluster inside it. The GCP project will be named using the team name and the name of the cluster created."
          type="info"
          showIcon
          style={{ marginBottom: '20px' }}
        />
      )}

      {this.state.gcpProjectAutomationType === 'CUSTOM' && (
        <GCPKoreManagedProjectsCustom
          gcpProjectList={this.state.gcpProjectList}
          plans={this.state.plans}
          handleChange={this.handleGcpProjectChange}
          handleDelete={this.handleGcpProjectDeleted}
          handleEdit={this.handleGcpProjectEdited}
          handleAdd={this.handleGcpProjectAdded}
          handleReset={this.handleResetGcpProjectList}
        />
      )}
    </>
  )

  render() {
    if (this.props.setupComplete) {
      return (
        <Card className="kore-managed-setup-complete">
          <Result
            status="success"
            title="Setup complete!"
            subTitle="Kore will manage the lifecycle of GCP projects along with Kore teams"
            extra={<Link href="/setup/kore/complete"><Button type="primary" key="continue">Continue</Button></Link>}
          >
            <Paragraph><Icon type="check-circle" theme="twoTone" twoToneColor="#52c41a" /> GCP organization credentials</Paragraph>
            <Paragraph style={{ marginBottom: '0' }}><Icon type="check-circle" theme="twoTone" twoToneColor="#52c41a" /> Project automation</Paragraph>
          </Result>
        </Card>
      )
    }

    return (
      <Card>
        <this.stepsHeader />
        <Divider />
        <this.stepContent />
        <Divider />
        <this.stepsActions />
      </Card>
    )
  }
}

export default GCPKoreManagedProjects
