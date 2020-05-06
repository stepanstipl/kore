import React from 'react'

import AWSExistingAccounts from './AWSExistingAccounts'

class AWSSetup extends React.Component {

  state = {
    setupComplete: false
  }

  setupComplete = () => this.setState({ setupComplete: true })

  render() {
    return <AWSExistingAccounts setupComplete={this.state.setupComplete} handleSetupComplete={this.setupComplete} />
  }
}

export default AWSSetup
