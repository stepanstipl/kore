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

  filterByStatus = (alerts) => {
    let count = 0

    if (!alerts) {
      return 0
    }

    alerts.items.forEach(resource => {
      if (this.props.severity) {
        if (this.props.severity !== resource.status.rule.spec.severity) {
          return
        }
      }
      if (resource.status.status === this.props.status) {
        count++
      }
    })

    return count
  }

  render() {
    const { alerts, color, description } = this.props
    const count = this.filterByStatus(alerts)

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
