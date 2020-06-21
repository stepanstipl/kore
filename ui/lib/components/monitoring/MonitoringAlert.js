import React from 'react'
import PropTypes from 'prop-types'
import { Badge, Descriptions, Divider, Timeline, Table, Card } from 'antd'

export default class MonitoringAlert extends React.Component {
  static propTypes = {
    rule: PropTypes.object.isRequired
  }

  static columns = [
    {
        title: 'Severity',
        dataIndex: 'spec.severity',
        key: 'severity',
    },
    {
        title: 'Triggered',
        dataIndex: 'spec.alerts[0].metadata.creationTimestamp',
        key: 'triggered',
    },
    {
        title: 'Rule',
        dataIndex: 'metadata.name',
        key: 'name',
    },
    {
        title: 'Summary',
        dataIndex: 'spec.alerts[0].spec.summary',
        key: 'summary',
    },
    {
        title: 'Team',
        dataIndex: 'spec.resource.namespace',
        key: 'resource.team',
    },
    {
        title: 'Resource',
        dataIndex: 'spec.resource.kind',
        key: 'resource.kind',
    },
  ]

  render() {
    const { rule } = this.props

    return (
      <>
        <Card>
          <h3>Alert Timeline</h3>
          <em>Provides a timeline of the alerting events</em>
          <Divider />
          <Timeline>
            <Timeline.Item color="green">Create a services site 2015-09-01</Timeline.Item>
            <Timeline.Item color="green">Create a services site 2015-09-01</Timeline.Item>
            <Timeline.Item color="red">
              <p>Solve initial network problems 1</p>
              <p>Solve initial network problems 2</p>
              <p>Solve initial network problems 3 2015-09-01</p>
            </Timeline.Item>
          </Timeline>
        </Card>
      </>
    )
  }
}

/*
import { Descriptions, Badge } from 'antd';

        <Card>
          <Descriptions title="Details" span={3}>
              <Badge status="processing" text={rule.severity} />
            <Descriptions.Item label="Severity">{rule.spec.severity}</Descriptions.Item>
            <Descriptions.Item label="Source">{rule.spec.resource}</Descriptions.Item>
          </Descriptions>
        </Card>
ReactDOM.render(
  mountNode,
);
*/