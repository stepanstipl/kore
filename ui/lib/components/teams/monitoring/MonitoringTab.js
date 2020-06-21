import React from 'react'
import PropTypes from 'prop-types'
import { Icon } from 'antd'

import KoreApi from '../../../kore-api'
import MonitoringOverview from '../../monitoring/MonitoringOverview'

export default class MonitoringTab extends React.Component {
  static propTypes = {
    team: PropTypes.object.isRequired,
    getOverviewStatus: PropTypes.func
  }

  state = {
    dataLoading: true,
  }

  async fetchComponentData() {
    const api = await KoreApi.client()
    const alerts = await api.monitoring.ListTeamAlerts(this.props.team.metadata.name)

    return { alerts }
  }

  componentDidMount() {
    return this.fetchComponentData().then(data => {
      this.setState({ ...data, dataLoading: false })
    })
  }

  componentDidUpdate(prevProps) {
    if (prevProps.team.metadata.name !== this.props.team.metadata.name) {
      this.setState({ dataLoading: true })
      return this.fetchComponentData().then(data => this.setState({ ...data, dataLoading: false }))
    }
  }

  render() {
    if (this.state.dataLoading) {
      return <Icon type="loading" />
    }

    return (
      <MonitoringOverview alerts={this.state.alerts} />
    )
  }
}
