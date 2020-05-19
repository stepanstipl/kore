import React from 'react'
import PropTypes from 'prop-types'
import { Icon } from 'antd'

import KoreApi from '../../../kore-api'
import SecurityOverview from '../../security/SecurityOverview'

export default class SecurityTab extends React.Component {
  static propTypes = {
    team: PropTypes.object.isRequired,
    getOverviewStatus: PropTypes.func
  }

  state = {
    dataLoading: true,
  }

  async fetchComponentData() {
    const api = await KoreApi.client()
    const overview = await api.GetTeamSecurityOverview(this.props.team.metadata.name)
    let status = 'Compliant'
    if (overview.spec.openIssueCounts && overview.spec.openIssueCounts.Warning) {
      status = 'Warning'
    }
    if (overview.spec.openIssueCounts && overview.spec.openIssueCounts.Failure) {
      status = 'Failure'
    }
    this.props.getOverviewStatus && this.props.getOverviewStatus(status)
    return { overview }
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
      <SecurityOverview overview={this.state.overview} />
    )
  }
}
