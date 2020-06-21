import React from 'react'
import PropTypes from 'prop-types'
import { Divider, Icon, Tooltip, Alert, Table, Row, Col, Card } from 'antd'
import inflect from 'inflect'
import Link from 'next/link'
import getConfig from 'next/config'

import MonitoringStatistic from './MonitoringStatistic'
import MonitoringTable from './MonitoringTable'
import MonitoringTimeline from './MonitoringTimeline'

const { publicRuntimeConfig } = getConfig()
const { Meta } = Card;

export default class MonitoringOverview extends React.Component {
  static propTypes = {
    alerts: PropTypes.object.isRequired
  }

  render() {
    const { alerts } = this.props

    if (!alerts) {
      return null
    }

    return (
      <>
        <div>
          <Row gutter={16}>
            <Col span={5}>
              <MonitoringStatistic
                alerts={alerts}
                description="No. OK Alerts"
                color='green'
                status="OK"
              />
            </Col>
            <Col span={5}>
              <MonitoringStatistic
                alerts={alerts}
                description="No. Critical Alerts"
                severity='Critical'
                status="Active"
              />
            </Col>
            <Col span={5}>
              <MonitoringStatistic
                alerts={alerts}
                description="No. Warning Alerts"
                severity="Warning"
                status="Active"
              />
            </Col>
            <Col span={5}>
              <MonitoringStatistic
                alerts={alerts}
                description="No. Silenced Alerts"
                status="Silenced"
              />
            </Col>
          </Row>
        </div>
        <Divider />
        <Card
          title="Alerts"
          style={{ padding: '5px' }}
          size="small"
        >
          <MonitoringTable
            alerts={alerts}
            severity={['Critical', 'Warning']}
            status={['Active', 'Silenced']}
          />
        </Card>
        <Divider />
      </>
    )
  }
}
