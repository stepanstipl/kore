import * as React from 'react'
import PropTypes from 'prop-types'
import { Tabs } from 'antd'

class CloudTabs extends React.Component {
  static propTypes = {
    handleSelectCloud: PropTypes.func.isRequired,
    defaultSelectedKey: PropTypes.string.isRequired
  }

  render() {
    const { handleSelectCloud, defaultSelectedKey } = this.props

    return (
      <Tabs defaultActiveKey={defaultSelectedKey} onChange={handleSelectCloud}>
        <Tabs.TabPane tab={
          <span>
            <img src="/static/images/GCP.png" height="40px" style={{ marginRight: '10px' }}/>
            Google Cloud Platform
          </span>
        } key="GCP" />
        <Tabs.TabPane tab={
          <span>
            <img src="/static/images/AWS.png" height="40px" style={{ marginRight: '15px' }} />
            Amazon Web Services (coming soon)
          </span>
        } disabled key="AWS" />
        <Tabs.TabPane tab={
          <span>
            <img src="/static/images/Azure.svg" height="25px" style={{ marginRight: '15px' }} />
            Microsoft Azure (coming soon)
          </span>
        } disabled key="Azure" />
      </Tabs>
    )
  }
}

export default CloudTabs
