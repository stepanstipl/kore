import React from 'react'
import PropTypes from 'prop-types'
import { Icon } from 'antd'

import KoreApi from '../../../kore-api'
import MonitoringDashboard from '../../monitoring/MonitoringDashboard'

export default class MonitoringTab extends React.Component {
  static propTypes = {
    team: PropTypes.object.isRequired,
    getOverviewStatus: PropTypes.func
  }

  static REFRESH_MS = 10000

  state = {
    dataLoading: true,
  }

  async fetchComponentData() {
    const alerts = await (await KoreApi.client()).monitoring.ListTeamAlerts(this.props.team.metadata.name)
    return { alerts }
  }

  refreshData = async () => {
    const data = await this.fetchComponentData()
    this.setState({ ...data })
  }

  componentDidMount() {
    return this.fetchComponentData().then(data => {
      this.setState({ ...data, dataLoading: false })
      this.interval = setInterval(this.refreshData, MonitoringTab.REFRESH_MS)
    })
  }

  componentDidUpdate(prevProps) {
    if (prevProps.team.metadata.name !== this.props.team.metadata.name) {
      this.setState({ dataLoading: true })
      return this.fetchComponentData().then(data => this.setState({ ...data, dataLoading: false }))
    }
  }

  componentWillUnmount() {
    clearInterval(this.interval)
  }

  render() {
    if (this.state.dataLoading) {
      return <Icon type="loading" />
    }

    return (
      <MonitoringDashboard alerts={this.state.alerts} refreshData={this.refreshData} />
    )
  }
}
