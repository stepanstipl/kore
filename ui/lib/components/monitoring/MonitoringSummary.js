import React from 'react'
import PropTypes from 'prop-types'
import { Tag } from 'antd'

export default class MonitoringSummary extends React.Component {
  static propTypes = {
    record: PropTypes.object.isRequired,
  }

  static filtered = [
    'alertname',
    'endpoint',
    'fstype',
    'job',
    'prometheus',
    'service',
    'severity',
    'rule_group'
  ]

  filterOnLabels = () => {
    if (!this.props.record) {
      return {}
    }

    const labels = this.props.record.spec.labels
    const processedLabels = {}

    Object.keys(labels)
      .filter(labelKey => !MonitoringSummary.filtered.includes(labelKey))
      .forEach(labelKey => processedLabels[labelKey] = labels[labelKey])

    return processedLabels
  }

  render() {
    const { record } = this.props
    const labels = this.filterOnLabels()
    console.log(record)

    return (
      <div>
        {record.spec.summary}<br/>
        <Tag color="green">team={record.metadata.namespace}</Tag>
        <Tag color="green">kind={record.status.rule.spec.resource.kind}</Tag>
        <Tag color="green">name={record.status.rule.spec.resource.name}</Tag>
        {Object.keys(labels).map(key => <Tag key={key} color="green">{key}={labels[key]}</Tag>)}
      </div>
    )
  }
}
