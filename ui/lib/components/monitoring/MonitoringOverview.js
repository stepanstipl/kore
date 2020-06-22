import React from 'react'
import { Card, Col, Divider, Icon, Row } from 'antd'

import MonitoringStatistic from './MonitoringStatistic'
import MonitoringTable from './MonitoringTable'
import KoreApi from '../../kore-api'

export default class MonitoringOverview extends React.Component {

  static REFRESH_MS = 30000

  state = {
    dataLoading: true
  }

  async fetchComponentData() {
    const alerts = await (await KoreApi.client()).monitoring.ListLatestAlerts()
    return { alerts }
  }

  refreshData = async () => {
    const data = await this.fetchComponentData()
    this.setState({ ...data })
  }

  componentDidMount() {
    this.fetchComponentData().then(data => {
      this.setState({ ...data, dataLoading: false })
      this.interval = setInterval(this.refreshData, MonitoringOverview.REFRESH_MS)
    })
  }

  componentWillUnmount() {
    clearInterval(this.interval)
  }

  render() {
    const { alerts, dataLoading } = this.state

    if (dataLoading) {
      return <Icon type="loading" />
    }

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
            refreshData={this.refreshData}
          />
        </Card>
        <Divider />
      </>
    )
  }
}
