import * as React from 'react'
import PropTypes from 'prop-types'
import { Tabs } from 'antd'

class CloudTabs extends React.Component {
  static propTypes = {
    handleSelectCloud: PropTypes.func.isRequired,
    selectedKey: PropTypes.string.isRequired
  }

  render() {
    const { handleSelectCloud, selectedKey } = this.props

    return (
      <Tabs activeKey={selectedKey} onChange={handleSelectCloud}>
        <Tabs.TabPane tab={
          <span id="tab-gcp">
            <img src="/static/images/GCP.png" height="40px" style={{ marginRight: '10px' }}/>
            Google Cloud Platform
          </span>
        } key="GCP" />
        <Tabs.TabPane tab={
          <span id="tab-aws">
            <img src="/static/images/AWS.png" height="40px" style={{ marginRight: '15px' }} />
            Amazon Web Services
          </span>
        } key="AWS" />
        <Tabs.TabPane tab={
          <span id="tab-azure">
            <img src="/static/images/Azure.svg" height="25px" style={{ marginRight: '15px' }} />
            Microsoft Azure (coming soon)
          </span>
        } disabled key="Azure" />
      </Tabs>
    )
  }
}

export default CloudTabs
