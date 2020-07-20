import React from 'react'
import PropTypes from 'prop-types'
import Link from 'next/link'
import { Button, Card, Divider, Icon, Result, Steps, Typography } from 'antd'
const { Paragraph } = Typography
const { Step } = Steps
import { pluralize, titleize } from 'inflect'
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

import RequestCredentialAccessForm from './forms/RequestCredentialAccessForm'
import CredentialsList from '../credentials/CredentialsList'
import FormErrorMessage from '../forms/FormErrorMessage'
import KoreApi from '../../kore-api'
import { errorMessage } from '../../utils/message'

class ExistingCloudAccounts extends React.Component {

  static propTypes = {
    cloud: PropTypes.oneOf(['GCP', 'AWS']),
    accountNoun: PropTypes.string.isRequired,
    setupComplete: PropTypes.bool.isRequired,
    handleSetupComplete: PropTypes.func.isRequired
  }

  steps = [
    { id: 'CREDS', title: 'Credentials', contentFn: 'stepsContentCreds', completeFn: 'stepsCompleteCreds'  },
    { id: 'ACCESS', title: `${titleize(this.props.accountNoun)} access`, contentFn: 'stepsContentAccess', completeFn: 'stepsCompleteAccess'  }
  ]

  state = {
    currentStep: 0,
    credsCount: 0,
    email: undefined,
    emailValid: false,
    submitting: false,
    errorMessage: false
  }

  async fetchComponentData() {
    const cloudConfig = await (await KoreApi.client()).GetConfig(this.props.cloud)
    const email = cloudConfig.spec && cloudConfig.spec.values.requestAccessEmail
    const emailValid = email ? true : false
    return { email, emailValid }
  }

  componentDidMount() {
    this.fetchComponentData().then(data => this.setState({ ...data }))
  }

  stepsContentCreds = () => (
    <>
      <Paragraph style={{ fontSize: '16px', fontWeight: '600' }}>Add one or more {this.props.cloud} {this.props.accountNoun} credentials</Paragraph>
      <CredentialsList provider={publicRuntimeConfig.clusterProviderMap[this.props.cloud]} getResourceItemList={this.setCredsCount} />
    </>
  )

  setupCompleteContent = () => (
    <Card>
      <Result
        status="success"
        title="Setup complete!"
        subTitle={`Kore will use existing ${this.props.cloud} ${pluralize(this.props.accountNoun)} that it's given access to`}
        extra={<Link href="/setup/kore/complete"><Button type="primary" key="continue">Continue</Button></Link>}
      >
        <Paragraph><Icon type="check-circle" theme="twoTone" twoToneColor="#52c41a" /> {this.props.cloud} {this.props.accountNoun} credentials</Paragraph>
        <Paragraph style={{ marginBottom: '0' }}><Icon type="check-circle" theme="twoTone" twoToneColor="#52c41a" /> {titleize(this.props.accountNoun)} access guidance</Paragraph>
      </Result>
    </Card>
  )

  nextStep() {
    const currentStep = this.state.currentStep + 1
    this.setState({ currentStep })
  }

  prevStep() {
    const currentStep = this.state.currentStep - 1
    this.setState({ currentStep })
  }

  setupComplete = async () => {
    this.setState({ submitting: true })
    try {
      const config = { requestAccessEmail: this.state.email }
      await (await KoreApi.client()).UpdateConfig(this.props.cloud, KoreApi.resources().generateConfigResource(this.props.cloud, config))
      this.setState({ submitting: false })
      this.props.handleSetupComplete()
    } catch (err) {
      console.error(`Error saving ${this.props.cloud} existing ${this.props.accountNoun} settings`, err)
      errorMessage(`Failed to save ${this.props.cloud} existing ${this.props.accountNoun} settings`)
      this.setState({ submitting: false, errorMessage: 'A problem occurred trying to save, please try again later.' })
    }
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

  stepsCompleteCreds = () => this.state.credsCount >= 1

  stepsCompleteAccess = () => this.state.emailValid

  stepContent = () => (
    <div className="steps-content" style={{ marginTop: '20px', marginBottom: '20px' }}>
      {this[this.steps[this.state.currentStep].contentFn]()}
    </div>
  )

  setCredsCount = (list) => this.setState({ credsCount: list.length })

  stepsContentAccess = () => {
    return (
      <RequestCredentialAccessForm
        cloud={this.props.cloud}
        data={{ email: this.state.email }}
        onChange={(email, errors) => this.setState({ email, emailValid: Boolean(!errors) })}
      />
    )
  }

  render() {
    if (this.props.setupComplete) {
      return this.setupCompleteContent()
    }

    return (
      <Card>
        {this.stepsHeader()}
        <Divider />
        <FormErrorMessage message={this.state.errorMessage} />
        {this.stepContent()}
        <Divider />
        {this.stepsActions()}
      </Card>
    )
  }
}

export default ExistingCloudAccounts
