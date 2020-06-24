import React from 'react'
import PropTypes from 'prop-types'
import { Card, Icon, Statistic } from 'antd'

export default class MonitoringStatistic extends React.Component {
  static propTypes = {
    alerts: PropTypes.object.isRequired,
    color: PropTypes.string,
    description: PropTypes.string.isRequired,
    severity: PropTypes.string,
    status: PropTypes.string.isRequired,
  }

  filterByStatus = () => {
    if (!this.props.alerts) {
      return 0
    }

    const filtered = this.props.alerts.items
      .filter(alert => (!this.props.severity || alert.status.rule.spec.severity === this.props.severity) && alert.status.status === this.props.status)
    return filtered.length
  }

  render() {
    const { color, description } = this.props
    const count = this.filterByStatus()

    return (
      <Card>
        <Statistic
          title={description}
          value={count}
          valueStyle={{ color: (count > 0 ? color : 'green') }}
          prefix={<Icon type="alert" />}
          suffix=""
        />
      </Card>
    )
  }
}
