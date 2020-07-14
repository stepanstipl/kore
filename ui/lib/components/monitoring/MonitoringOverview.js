import React from 'react'
import PropTypes from 'prop-types'
import { Card, Col, Divider, Row } from 'antd'

import MonitoringStatistic from './MonitoringStatistic'
import MonitoringTable from './MonitoringTable'

export default class MonitoringOverview extends React.Component {

  static propTypes = {
    alerts: PropTypes.object.isRequired,
    refreshData: PropTypes.func
  }

  render() {
    const { alerts, refreshData } = this.props

    return (
      <>
        <Row gutter={[16, 16]}>
          <Col lg={12} xl={6}>
            <MonitoringStatistic
              alerts={alerts}
              description="No. OK Alerts"
              color='green'
              status="OK"
            />
          </Col>
          <Col lg={12} xl={6}>
            <MonitoringStatistic
              alerts={alerts}
              description="No. Critical Alerts"
              severity='Critical'
              status="Active"
            />
          </Col>
          <Col lg={12} xl={6}>
            <MonitoringStatistic
              alerts={alerts}
              description="No. Warning Alerts"
              severity="Warning"
              status="Active"
            />
          </Col>
          <Col lg={12} xl={6}>
            <MonitoringStatistic
              alerts={alerts}
              description="No. Silenced Alerts"
              status="Silenced"
            />
          </Col>
        </Row>
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
            refreshData={refreshData}
          />
        </Card>
        <Divider />
      </>
    )
  }
}
