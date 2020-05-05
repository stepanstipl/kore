import React from 'react'
import Link from 'next/link'
import { Typography, Card, Alert, Divider, Tooltip, Radio, List, Icon, Steps, Button, Row, Col, Modal , Drawer, Form, Input, Popover, Result, message } from 'antd'
const { Title, Text, Paragraph } = Typography
const { Step } = Steps

import copy from '../../../../lib/utils/object-copy'
import canonical from '../../../../lib/utils/canonical'
import KoreApi from '../../../../lib/kore-api'
import GKECredentialsList from '../../../../lib/components/credentials/GKECredentialsList'
import PlanViewer from '../../../../lib/components/plans/PlanViewer'
import CloudSelector from '../../../../lib/components/common/CloudSelector'

// prototype components
import GCPOrganizationsList from '../../../../lib/prototype/components/credentials/GCPOrganizationsList'
import AutomatedProjectForm from '../../../../lib/prototype/components/configure/AutomatedProjectForm'

class CloudAccessPage extends React.Component {

  static staticProps = {
    title: 'Setup cloud access',
    hideSider: true,
    adminOnly: true
  }

  static initialProjectList = [{
    code: 'not-production',
    title: 'Not production',
    description: 'To be used for all environments except production',
    prefix: 'kore',
    suffix: 'not-prod',
    plans: ['gke-development']
  }, {
    code: 'production',
    title: 'Production',
    description: 'Project just for the production environment',
    prefix: 'kore',
    suffix: 'prod',
    plans: ['gke-production']
  }]

  stepsKoreManaged = [
    { id: 'CREDS', title: 'Credentials' },
    { id: 'PROJECTS', title: 'Project automation' }
  ]

  stepsExisting = [
    { id: 'CREDS', title: 'Credentials' },
    { id: 'ACCESS', title: 'Project access' }
  ]

  state = {
    selectedCloud: '',
    gcpManagementType: '',
    gcpProjectAutomationType: '',
    gcpProjectList: [],
    currentStepKoreManaged: 0,
    currentStepExisting: 0,
    stepsKoreManagedComplete: false,
    stepsExistingComplete: false,
    addProject: false,
    associatePlanVisible: false
  }

  nextStepKoreManaged() {
    const currentStepKoreManaged = this.state.currentStepKoreManaged + 1
    this.setState({ currentStepKoreManaged })
  }

  prevStepKoreManaged() {
    const currentStepKoreManaged = this.state.currentStepKoreManaged - 1
    this.setState({ currentStepKoreManaged })
  }

  nextStepExisting() {
    const currentStepExisting = this.state.currentStepExisting + 1
    this.setState({ currentStepExisting })
  }

  prevStepExisting() {
    const currentStepExisting = this.state.currentStepExisting - 1
    this.setState({ currentStepExisting })
  }

  stepsKoreManagedComplete = () => this.setState({ stepsKoreManagedComplete: true })
  stepsExistingComplete = () => this.setState({ stepsExistingComplete: true })

  async fetchComponentData() {
    const api = await KoreApi.client()
    const planList = await api.ListPlans()
    return planList.items
  }

  componentDidMount() {
    this.fetchComponentData().then(plans => {
      this.setState({ plans })
    })
  }

  handleSelectCloud = cloud => {
    if (this.state.selectedCloud !== cloud) {
      const state = copy(this.state)
      state.selectedCloud = cloud
      this.setState(state)
    }
  }

  selectGcpManagementType = e => this.setState({ gcpManagementType: e.target.value })
  selectGcpProjectAutomationType = e => this.setState({ gcpProjectAutomationType: e.target.value })

  deleteGcpProject = (code) => {
    return () => {
      this.setState({
        gcpProjectList: this.state.gcpProjectList.filter(p => p.code !== code)
      })
      message.success('GCP automated project removed')
    }
  }

  onChange = (code, property) => {
    return (value) => {
      const gcpProjectList = copy(this.state.gcpProjectList)
      gcpProjectList.find(p => p.code === code)[property] = value
      this.setState({
        gcpProjectList
      })
    }
  }

  showPlanDetails = (plan) => {
    return () => {
      Modal.info({
        title: (<><Title level={4}>{plan.spec.description}</Title><Text>{plan.spec.summary}</Text></>),
        content: <PlanViewer
          plan={plan}
          resourceType="cluster"
        />,
        width: 700,
        onOk() {}
      })
    }
  }

  handleAssociatePlanVisibleChange = (projectCode) => () => this.setState({ associatePlanVisible: projectCode })

  associatePlan = (projectCode, plan) => {
    return () => {
      const gcpProjectList = copy(this.state.gcpProjectList)
      gcpProjectList.find(p => p.code === projectCode).plans.push(plan)
      this.setState({ gcpProjectList })
      message.success('Plan associated')
    }
  }

  closeAssociatePlans = () => this.setState({ associatePlanVisible: false })

  associatePlanContent = (projectCode) => {
    const project = this.state.gcpProjectList.find(p => p.code === projectCode)
    const cloudPlans = this.state.plans.filter(p => p.spec.kind === this.state.selectedCloud)
    const unassociatedPlans = cloudPlans.filter(p => !project.plans.includes(p.metadata.name))
    if (unassociatedPlans.length === 0) {
      return (
        <>
          <Alert style={{ marginBottom: '20px' }} message="All existing plans are already associated with this project type." />
          <Button type="primary" onClick={this.closeAssociatePlans}>Close</Button>
        </>
      )
    }
    return (
      <>
        <Alert style={{ marginBottom: '20px' }} message="Plans available to be associated with this project type." />
        <List
          dataSource={unassociatedPlans}
          renderItem={plan => <Paragraph>{plan.spec.description} <a style={{ textDecoration: 'underline' }} onClick={this.associatePlan(project.code, plan.metadata.name)}>Associate</a></Paragraph>}
        />
        <Button style={{ marginTop: '10px' }} type="primary" onClick={this.closeAssociatePlans}>Close</Button>
      </>
    )
  }

  unassociatePlan = (projectCode, plan) => {
    return () => {
      const gcpProjectList = copy(this.state.gcpProjectList)
      const project = gcpProjectList.find(p => p.code === projectCode)
      project.plans = project.plans.filter(p => p !== plan)
      this.setState({ gcpProjectList })
      message.success('Plan unassociated')
    }
  }

  setGcpProjectsToDefault = () => this.setState({ gcpProjectList: CloudAccessPage.initialProjectList })

  addProject = (enabled) => () => this.setState({ addProject: enabled })

  handleProjectAdded = (project) => {
    const code = canonical(project.title)
    this.setState({
      gcpProjectList: this.state.gcpProjectList.concat([{ code, plans: [], ...project }]),
      addProject: false
    })
    message.success('GCP automated project added')
  }

  IconTooltip = ({ icon, text }) => (
    <Tooltip title={text}>
      <Icon type={icon} theme="twoTone" />
    </Tooltip>
  )

  IconTooltipButton = ({ icon, text, onClick }) => (
    <Tooltip title={text}>
      <a style={{ marginLeft: '5px' }} onClick={onClick}><Icon type={icon} theme="twoTone" /></a>
    </Tooltip>
  )

  render() {
    const {
      selectedCloud,
      gcpManagementType,
      gcpProjectAutomationType,
      gcpProjectList,
      plans,
      currentStepKoreManaged,
      currentStepExisting,
      stepsKoreManagedComplete,
      stepsExistingComplete,
      addProject,
      associatePlanVisible
    } = this.state

    return (
      <div>
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
        { selectedCloud === 'GKE' ? (
          <Card>

            <div style={{ marginBottom: '15px' }}>
              <Paragraph style={{ fontSize: '16px', fontWeight: '600' }}>How do you want Kore teams to integrate with GCP projects?</Paragraph>
              <Radio.Group onChange={this.selectGcpManagementType} value={gcpManagementType} disabled={stepsKoreManagedComplete || stepsExistingComplete}>
                <Radio value={'KORE'} style={{ marginRight: '20px' }}>
                  <Text style={{ fontSize: '16px', fontWeight: '600' }}>Kore managed projects <Text type="secondary">(recommended)</Text></Text>
                  <Paragraph style={{ marginLeft: '24px', marginBottom: '0' }}>Kore will manage the GCP projects required for teams</Paragraph>
                </Radio>
                <Radio value={'EXTERNAL'}>
                  <Text style={{ fontSize: '16px', fontWeight: '600' }}>Use existing projects</Text>
                  <Paragraph style={{ marginLeft: '24px', marginBottom: '0' }}>Kore teams will use existing GCP projects that it&apos;s given access to</Paragraph>
                </Radio>
              </Radio.Group>
            </div>

            {gcpManagementType === 'KORE' && stepsKoreManagedComplete ? (
              <Card>
                <Result
                  status="success"
                  title="Kore managed setup complete!"
                  subTitle="Kore will manage the lifecycle of GCP projects along with Kore teams"
                  extra={<Link href="/prototype/setup/kore/complete"><Button type="primary" key="continue">Continue</Button></Link>}
                >
                  <Paragraph><Icon type="check-circle" theme="twoTone" twoToneColor="#52c41a" /> GCP organization credentials</Paragraph>
                  <Paragraph><Icon type="check-circle" theme="twoTone" twoToneColor="#52c41a" /> Project automation</Paragraph>
                </Result>
              </Card>
            ) : null}

            {gcpManagementType === 'KORE' && !stepsKoreManagedComplete ? (
              <>
                <Card>
                  <Steps current={currentStepKoreManaged}>
                    {this.stepsKoreManaged.map(item => (
                      <Step key={item.title} title={item.title} />
                    ))}
                  </Steps>
                  <Divider />
                  <div className="steps-content" style={{ marginTop: '20px', marginBottom: '20px' }}>
                    {this.stepsKoreManaged[currentStepKoreManaged].id === 'CREDS' ? (
                      <>
                        <Paragraph style={{ fontSize: '16px', fontWeight: '600' }}>Add one or more GCP organization credentials</Paragraph>
                        <GCPOrganizationsList />
                      </>
                    ) : null}
                    {this.stepsKoreManaged[currentStepKoreManaged].id === 'PROJECTS' ? (
                      <>
                        <Paragraph style={{ fontSize: '16px', fontWeight: '600' }}>Configure GCP project automation</Paragraph>
                        <Alert
                          message="This defines how Kore will automate GCP projects for teams"
                          description="When a team is created in Kore and a cluster is requested, Kore will ensure the GCP project is also created and the cluster placed inside it. You must also specify the plans available for each type of project, this is to ensure the correct cluster specification is being used."
                          type="info"
                          showIcon
                          style={{ marginBottom: '20px' }}
                        />

                        <div style={{ marginBottom: '15px' }}>
                          <Paragraph style={{ fontSize: '16px', fontWeight: '600' }}>How do you want Kore to automate GCP projects for teams?</Paragraph>
                          <Radio.Group onChange={this.selectGcpProjectAutomationType} value={gcpProjectAutomationType}>
                            <Radio value={'CLUSTER'} style={{ marginRight: '20px' }}>
                              <Text style={{ fontSize: '16px', fontWeight: '600' }}>One per cluster</Text>
                              <Paragraph style={{ marginLeft: '24px', marginBottom: '0' }}>Kore will create a GCP project for each cluster a team provisions</Paragraph>
                            </Radio>
                            <Radio value={'CUSTOM'}>
                              <Text style={{ fontSize: '16px', fontWeight: '600' }}>Custom</Text>
                              <Paragraph style={{ marginLeft: '24px', marginBottom: '0' }}>Configure how Kore will create GCP projects for teams</Paragraph>
                            </Radio>
                          </Radio.Group>
                        </div>

                        {gcpProjectAutomationType === 'CLUSTER' ? (
                          <Alert
                            message="For every cluster a team creates Kore will also create a GCP project and provision the cluster inside it. The GCP project will share the name given to the cluster."
                            type="info"
                            showIcon
                            style={{ marginBottom: '20px' }}
                          />
                        ) : null}

                        {gcpProjectAutomationType === 'CUSTOM' ? (

                          <>
                            <Alert
                              message="When a team is created in Kore and a cluster is requested, Kore will ensure the associated GCP project is also created and the cluster placed inside it. You must also specify the plans available for each type of project, this is to ensure the correct cluster specification is being used."
                              type="info"
                              showIcon
                              style={{ marginBottom: '20px' }}
                            />
                            <div style={{ display: 'block', marginBottom: '20px' }}>
                              <Button type="primary" onClick={this.addProject(true)}>+ New</Button>
                              <Button style={{ marginLeft: '10px' }} onClick={this.setGcpProjectsToDefault}>Set to Kore defaults</Button>
                            </div>
                            {gcpProjectList.length === 0 ? (
                              <Paragraph>No automated projects configured, you can &apos;Set to Kore defaults&apos; and/or add new ones. </Paragraph>
                            ) : (
                              <List
                                itemLayout="vertical"
                                bordered={true}
                                dataSource={gcpProjectList}
                                renderItem={project => (
                                  <List.Item actions={[<a key="delete" onClick={this.deleteGcpProject(project.code)}><Icon type="delete" /> Remove</a>]}>
                                    <List.Item.Meta
                                      title={<Text editable={{ onChange: this.onChange(project.code, 'title') }} style={{ fontSize: '16px' }}>{project.title}</Text>}
                                      description={<Text editable={{ onChange: this.onChange(project.code, 'description') }}>{project.description}</Text>}
                                    />

                                    <Row gutter={16}>
                                      <Col span={8}>
                                        <Card
                                          title="Naming"
                                          size="small"
                                          bordered={false}
                                        >
                                          <Paragraph>The project will be named using the team name, with the prefix and suffix below</Paragraph>
                                          <Row style={{ padding: '5px 0' }}>
                                            <Col span={8}>
                                              Prefix
                                            </Col>
                                            <Col span={16}>
                                              <Text editable={{ onChange: this.onChange(project.code, 'prefix') }}>{project.prefix}</Text>
                                            </Col>
                                          </Row>
                                          <Row style={{ padding: '5px 0' }}>
                                            <Col span={8}>
                                              Suffix
                                            </Col>
                                            <Col span={16}>
                                              <Text editable={{ onChange: this.onChange(project.code, 'suffix') }}>{project.suffix}</Text>
                                            </Col>
                                          </Row>
                                          <Row style={{ paddingTop: '15px' }}>
                                            <Col span={8}>
                                              Example
                                            </Col>
                                            <Col span={16}>
                                              <Text>{project.prefix}-<span style={{ fontStyle: 'italic' }}>team-name</span>-{project.suffix}</Text>
                                            </Col>
                                          </Row>
                                        </Card>
                                      </Col>
                                      <Col span={8}>
                                        <Card
                                          title="Cluster plans"
                                          size="small"
                                          bordered={false}
                                        >
                                          <Paragraph>The cluster plans associated with this project.</Paragraph>
                                          {project.plans.length === 0 ? <div style={{ padding: '5px 0' }}>No plans</div> : null}
                                          {(plans || []).filter(p => project.plans.includes(p.metadata.name)).map((plan, i) => (
                                            <div key={i} style={{ padding: '5px 0' }}>
                                              <Text style={{ marginRight: '10px' }}>{plan.spec.description}</Text>
                                              <this.IconTooltip icon="info-circle" text={plan.spec.summary} />
                                              <this.IconTooltipButton icon="eye" text="View plan" onClick={this.showPlanDetails(plan)} />
                                              <this.IconTooltipButton icon="delete" text="Unassociate plan" onClick={this.unassociatePlan(project.code, plan.metadata.name)} />
                                            </div>
                                          ))}
                                          <div style={{ padding: '5px 0' }}>

                                            <Popover
                                              content={this.associatePlanContent(project.code)}
                                              title={`${project.title}: Associate plans`}
                                              trigger="click"
                                              visible={associatePlanVisible === project.code}
                                              onVisibleChange={this.handleAssociatePlanVisibleChange(project.code)}
                                            >
                                              <a>+ Associate plan</a>
                                            </Popover>
                                          </div>
                                        </Card>
                                      </Col>
                                    </Row>

                                  </List.Item>
                                )}
                              />
                            )}
                            {addProject ? (
                              <Drawer
                                title={<Title level={4}>New project</Title>}
                                visible={addProject}
                                onClose={this.addProject(false)}
                                width={700}
                              >
                                <AutomatedProjectForm handleSubmit={this.handleProjectAdded} handleCancel={this.addProject(false)} />
                              </Drawer>
                            ) : null}
                          </>

                        ) : null}

                      </>
                    ) : null}
                  </div>
                  <Divider />
                  <div className="steps-action">
                    {currentStepKoreManaged < this.stepsKoreManaged.length - 1 && (
                      <Button type="primary" onClick={() => this.nextStepKoreManaged()}>
                        Next
                      </Button>
                    )}
                    {currentStepKoreManaged === this.stepsKoreManaged.length - 1 && (
                      <Button type="primary" onClick={this.stepsKoreManagedComplete}>
                        Done
                      </Button>
                    )}
                    {currentStepKoreManaged > 0 && (
                      <Button style={{ marginLeft: 8 }} onClick={() => this.prevStepKoreManaged()}>
                        Previous
                      </Button>
                    )}
                  </div>
                </Card>

              </>
            ) : null}

            {gcpManagementType === 'EXTERNAL' && stepsExistingComplete ? (
              <Card>
                <Result
                  status="success"
                  title="Use existing setup complete!"
                  subTitle="Kore will use existing GCP projects that it's given access to"
                  extra={<Link href="/prototype/setup/kore/complete"><Button type="primary" key="continue">Continue</Button></Link>}
                >
                  <Paragraph><Icon type="check-circle" theme="twoTone" twoToneColor="#52c41a" /> GCP project credentials</Paragraph>
                  <Paragraph><Icon type="check-circle" theme="twoTone" twoToneColor="#52c41a" /> Project access guidance</Paragraph>
                </Result>
              </Card>
            ) : null}

            {gcpManagementType === 'EXTERNAL' && !stepsExistingComplete ? (
              <>
                <Card>

                  <Steps current={currentStepExisting}>
                    {this.stepsExisting.map(item => (
                      <Step key={item.title} title={item.title} />
                    ))}
                  </Steps>
                  <Divider />

                  {this.stepsExisting[currentStepExisting].id === 'CREDS' ? (
                    <>
                      <Paragraph style={{ fontSize: '16px', fontWeight: '600' }}>Add one or more GCP project credentials</Paragraph>
                      <GKECredentialsList />
                    </>
                  ) : null}
                  {this.stepsExisting[currentStepExisting].id === 'ACCESS' ? (
                    <>
                      <Paragraph style={{ fontSize: '16px', fontWeight: '600' }}>Project credential access for teams</Paragraph>
                      <Alert
                        message="Team access"
                        description="When using Kore with existing GCP projects, you must allocate the project credentials to teams in order for them to provision clusters within those projects. When a new team is created they may not have access to any project credentials, here you can provide an email address which will be displayed to a team in this situation, in order to request access to a GCP project through Kore."
                        type="info"
                        showIcon
                        style={{ marginBottom: '20px' }}
                      />
                      <Form.Item labelAlign="left" labelCol={{ span: 2 }} wrapperCol={{ span: 6 }} label="Email" help="Email for teams who need access to project credentails">
                        <Input placeholder="Title" />
                      </Form.Item>
                    </>
                  ) : null}

                  <Divider />
                  <div className="steps-action">
                    {currentStepExisting < this.stepsExisting.length - 1 && (
                      <Button type="primary" onClick={() => this.nextStepExisting()}>
                        Next
                      </Button>
                    )}
                    {currentStepExisting === this.stepsExisting.length - 1 && (
                      <Button type="primary" onClick={this.stepsExistingComplete}>
                        Done
                      </Button>
                    )}
                    {currentStepExisting > 0 && (
                      <Button style={{ marginLeft: 8 }} onClick={() => this.prevStepExisting()}>
                        Previous
                      </Button>
                    )}
                  </div>
                  
                </Card>
              </>
            ) : null}

          </Card>
        ) : null }

        { selectedCloud === 'EKS' ? (
          <Card>
            <Alert
              message="Amazon Web Services"
              description="TODO"
              type="info"
              showIcon
              style={{ marginBottom: '20px' }}
            />
          </Card>
        ) : null}
      </div>
    )
  }
}

export default CloudAccessPage
