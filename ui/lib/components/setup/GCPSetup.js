import PropTypes from 'prop-types'
import React from 'react'
import { Card, Radio, Typography } from 'antd'
const { Paragraph, Text } = Typography

import GCPKoreManagedProjects from './GCPKoreManagedProjects'
import GCPExistingProjects from './GCPExistingProjects'

class GCPSetup extends React.Component {

  static propTypes = {
    accountManagement: PropTypes.object,
    gkeCredentialsList: PropTypes.array.isRequired
  }

  state = {
    gcpManagementType: this.props.accountManagement ? 'KORE' : (this.props.gkeCredentialsList.length >= 1 ? 'EXISTING': false),
    setupComplete: false
  }

  selectGcpManagementType = e => this.setState({ gcpManagementType: e.target.value })
  setupComplete = () => this.setState({ setupComplete: true })

  render() {
    const { gcpManagementType, setupComplete } = this.state

    return (
      <Card>
        <div style={{ marginBottom: '15px' }}>
          <Paragraph style={{ fontSize: '16px', fontWeight: '600' }}>How do you want Kore teams to integrate with GCP projects?</Paragraph>
          <Radio.Group onChange={this.selectGcpManagementType} value={gcpManagementType} disabled={setupComplete}>
            <Radio className="use-kore-managed-projects" value={'KORE'} style={{ marginRight: '20px' }}>
              <Text style={{ fontSize: '16px', fontWeight: '600' }}>Kore managed projects <Text type="secondary"> (recommended)</Text></Text>
              <Paragraph style={{ marginLeft: '24px', marginBottom: '0' }}>Kore will manage the GCP projects required for teams</Paragraph>
            </Radio>
            <Radio className="use-existing-projects" value={'EXISTING'}>
              <Text style={{ fontSize: '16px', fontWeight: '600' }}>Use existing projects</Text>
              <Paragraph style={{ marginLeft: '24px', marginBottom: '0' }}>Kore teams will use existing GCP projects that it&apos;s given access to</Paragraph>
            </Radio>
          </Radio.Group>
        </div>
        {gcpManagementType === 'KORE' && <GCPKoreManagedProjects accountManagement={this.props.accountManagement} setupComplete={setupComplete} handleSetupComplete={this.setupComplete} />}
        {gcpManagementType === 'EXISTING' && <GCPExistingProjects setupComplete={setupComplete} handleSetupComplete={this.setupComplete} />}
      </Card>
    )
  }
}

export default GCPSetup
