import React from 'react'
import { Icon } from 'antd'

import Breadcrumb from '../../../lib/components/layout/Breadcrumb'
import MonitoringOverview from '../../../lib/components/monitoring/MonitoringOverview'
import KoreApi from '../../../lib/kore-api'

export default class MonitoringOverviewPage extends React.Component {

  static staticProps = {
    title: 'Monitoring Overview',
    adminOnly: true
  }

  static REFRESH_MS = 30000

  state = {
    dataLoading: true
  }

  async fetchComponentData() {
    const alerts = await (await KoreApi.client()).monitoring.ListLatestAlerts()
    return { alerts }
  }

  refreshData = async () => {
    const data = await this.fetchComponentData()
    this.setState({ ...data })
  }

  componentDidMount() {
    this.fetchComponentData().then(data => {
      this.setState({ ...data, dataLoading: false })
      this.interval = setInterval(this.refreshData, MonitoringOverviewPage.REFRESH_MS)
    })
  }

  componentWillUnmount() {
    clearInterval(this.interval)
  }

  render() {
    return (
      <>
        <Breadcrumb items={[{ text: 'Monitoring' }, { text: 'Overview' }]} />
        {this.state.dataLoading ? <Icon type="loading" /> : <MonitoringOverview alerts={this.state.alerts} refreshData={this.refreshData} />}
      </>
    )
  }
}
