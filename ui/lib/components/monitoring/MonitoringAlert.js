import React from 'react'
import PropTypes from 'prop-types'
import { Badge, Descriptions, Divider, Timeline, Table, Card } from 'antd'

export default class MonitoringAlert extends React.Component {
  static propTypes = {
    rule: PropTypes.object.isRequired
  }

  render() {
    const { rule } = this.props

    return (
      <>
        <Card>
          <h3>Alert Details</h3>
          <em>Provides on the alert</em>
          <Divider />
        </Card>
      </>
    )
  }
}
