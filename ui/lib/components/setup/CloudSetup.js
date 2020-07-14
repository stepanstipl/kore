import React from 'react'
import PropTypes from 'prop-types'
import { Card } from 'antd'

import KoreTeamCloudIntegration from './radio-groups/KoreTeamCloudIntegration'
import KoreManagedCloudAccounts from './KoreManagedCloudAccounts'
import ExistingCloudAccounts from './ExistingCloudAccounts'

class CloudSetup extends React.Component {

  static propTypes = {
    provider: PropTypes.oneOf(['GKE', 'EKS']),
    cloud: PropTypes.oneOf(['GCP', 'AWS']),
    accountNoun: PropTypes.string.isRequired,
    accountManagement: PropTypes.object,
    credentialsList: PropTypes.array.isRequired
  }

  state = {
    cloudManagementType: this.props.accountManagement ? 'KORE' : (this.props.credentialsList.length >= 1 ? 'EXISTING' : ''),
    setupComplete: false
  }

  selectCloudManagementType = (e) => this.setState({ cloudManagementType: e.target.value })
  setupComplete = () => this.setState({ setupComplete: true })

  render() {
    const { cloudManagementType, setupComplete } = this.state

    return (
      <Card>
        <KoreTeamCloudIntegration
          cloud={this.props.cloud}
          accountNoun={this.props.accountNoun}
          onChange={this.selectCloudManagementType}
          value={cloudManagementType}
          disabled={setupComplete}
        />
        {cloudManagementType === 'KORE' && <KoreManagedCloudAccounts provider={this.props.provider} cloud={this.props.cloud} accountNoun={this.props.accountNoun} accountManagement={this.props.accountManagement} setupComplete={setupComplete} handleSetupComplete={this.setupComplete} />}
        {cloudManagementType === 'EXISTING' && <ExistingCloudAccounts cloud={this.props.cloud} accountNoun={this.props.accountNoun} setupComplete={setupComplete} handleSetupComplete={this.setupComplete} />}
      </Card>
    )
  }
}

export default CloudSetup
