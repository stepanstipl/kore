import React from 'react'
import Link from 'next/link'
import PropTypes from 'prop-types'
import { List, Icon, Tooltip, Table, Tag, Button, Card, Timeline } from 'antd'

export default class MonitoringTimeline extends React.Component {
  static propTypes = {
    alerts: PropTypes.object.isRequired,
  }

  render() {
    const { alerts } = this.props

    return (
      <>
        <Timeline>
            {alerts.items.forEach(e => <Timeline.Item color="green">{e.metadata.name}</Timeline.Item>)}
        </Timeline>
      </>
    )
  }
}
