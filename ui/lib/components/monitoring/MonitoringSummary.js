import React from 'react'
import PropTypes from 'prop-types'
import { Tag } from 'antd'
import { filter } from 'lodash'

export default class MonitoringSummary extends React.Component {
  static propTypes = {
    record: PropTypes.object.isRequired,
  }

  static filtered = [
    "alertname",
    "endpoint",
    "fstype",
    "job",
    "prometheus",
    "service",
    "severity",
  ]

  filterOnLabels = (labels) => {
    var list = {}

    if (!labels) {
      return {}
    }

    for (const [key, value] of Object.entries(labels.spec.labels)) {
        if (MonitoringSummary.filtered.includes(key)) {
            continue
        }
        list[key] = value
    }

    return list
  }

  render() {
    const labels = this.filterOnLabels(this.props.record)

    return (
        <div>
            {this.props.record.spec.summary}<br/>
            <Tag color="green">team={this.props.record.metadata.namespace}</Tag>
            <Tag color="green">kind={this.props.record.status.rule.spec.resource.kind}</Tag>
            <Tag color="green">name={this.props.record.status.rule.spec.resource.name}</Tag>
            <Tag color="green">status={this.props.record.status.status}</Tag>
            {Object.keys(labels).map(key => <Tag color="green">{key}={labels[key]}</Tag>)}
        </div>
    )
  }
}
