import React from 'react'
import PropTypes from 'prop-types'

import Breadcrumb from '../../../lib/components/layout/Breadcrumb'
import MonitoringOverview from '../../../lib/components/monitoring/MonitoringOverview'
import KoreApi from '../../../lib/kore-api'

export default class MonitoringPage extends React.Component {
  static propTypes = {
    alerts: PropTypes.object.isRequired,
  }

  static staticProps = {
    title: 'Monitoring Overview',
    adminOnly: true
  }

  static async getPageData(ctx) {
    try {
      const alerts = await (await KoreApi.client(ctx)).monitoring.ListLatestAlerts()
      return { alerts }
    } catch (err) {
      throw new Error(err.message)
    }
  }

  static getInitialProps = async ctx => {
    const data = await MonitoringPage.getPageData(ctx)

    return data
  }

  constructor(props) {
    super(props)
  }

  render() {
    const { alerts } = this.props
    return (
      <div>
        <Breadcrumb items={[{ text: 'Monitoring' }, { text: 'Overview' }]} />
        <MonitoringOverview alerts={alerts} />
      </div>
    )
  }
}
