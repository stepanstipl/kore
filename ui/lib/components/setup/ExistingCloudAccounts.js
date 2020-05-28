import React from 'react'
import PropTypes from 'prop-types'
import { Button, Card, Divider, Steps } from 'antd'
const { Step } = Steps

import RequestCredentialAccessForm from './forms/RequestCredentialAccessForm'

class ExistingCloudAccounts extends React.Component {

  static propTypes = {
    setupComplete: PropTypes.bool.isRequired,
    handleSetupComplete: PropTypes.func.isRequired
  }

  steps = [
    { id: 'CREDS', title: 'Credentials', contentFn: 'stepsContentCreds', completeFn: 'stepsCompleteCreds'  },
    { id: 'ACCESS', title: 'Account access', contentFn: 'stepsContentAccess', completeFn: 'stepsCompleteAccess'  }
  ]

  state = {
    currentStep: 0,
    credsCount: 0,
    emailValid: RequestCredentialAccessForm.ENABLED ? false : true
  }

  stepsContentCreds = () => {
    throw new Error('stepsContentCreds must be implemented')
  }

  setupCompleteContent = () => {
    throw new Error('setupCompleteContent must be implemented')
  }

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
      <RequestCredentialAccessForm cloud="AWS" onChange={(errors) => this.setState({ emailValid: Boolean(!errors) })} />
    )
  }

  render() {
    if (this.props.setupComplete) {
      return <this.setupCompleteContent />
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

export default ExistingCloudAccounts
