import React from 'react'
import PropTypes from 'prop-types'
import Link from 'next/link'
import { Button, Card, Divider, Icon, Result, Steps, Typography } from 'antd'
const { Paragraph } = Typography
const { Step } = Steps
import { pluralize, titleize } from 'inflect'

import RequestCredentialAccessForm from './forms/RequestCredentialAccessForm'
import GKECredentialsList from '../credentials/GKECredentialsList'
import EKSCredentialsList from '../credentials/EKSCredentialsList'

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
    emailValid: RequestCredentialAccessForm.ENABLED ? false : true
  }

  stepsContentCreds = () => (
    <>
      <Paragraph style={{ fontSize: '16px', fontWeight: '600' }}>Add one or more {this.props.cloud} {this.props.accountNoun} credentials</Paragraph>
      {this.props.cloud === 'GCP' && <GKECredentialsList getResourceItemList={this.setCredsCount} />}
      {this.props.cloud === 'AWS' && <EKSCredentialsList getResourceItemList={this.setCredsCount} />}
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

  setupComplete = () => {
    this.setState({ setupComplete: true })
    this.props.handleSetupComplete()
  }

  stepsHeader = () => (
    <Steps current={this.state.currentStep}>
      {this.steps.map(item => <Step key={item.title} title={item.title} />)}
    </Steps>
  )

  stepsActions = () => (
    <div className="steps-action">
      {this.state.currentStep < this.steps.length - 1 && <Button type="primary" disabled={!this[this.steps[this.state.currentStep].completeFn]()} onClick={() => this.nextStep()}>Next</Button>}
      {this.state.currentStep === this.steps.length - 1 && <Button type="primary" disabled={!this[this.steps[this.state.currentStep].completeFn]()} onClick={this.props.handleSetupComplete}>Save</Button>}
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
      <RequestCredentialAccessForm cloud={this.props.cloud} onChange={(errors) => this.setState({ emailValid: Boolean(!errors) })} />
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
        {this.stepContent()}
        <Divider />
        {this.stepsActions()}
      </Card>
    )
  }
}

export default ExistingCloudAccounts
