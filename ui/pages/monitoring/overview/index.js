import React from 'react'

import Breadcrumb from '../../../lib/components/layout/Breadcrumb'
import MonitoringOverview from '../../../lib/components/monitoring/MonitoringOverview'

export default class MonitoringPage extends React.Component {

  static staticProps = {
    title: 'Monitoring Overview',
    adminOnly: true
  }

  render() {
    return (
      <>
        <Breadcrumb items={[{ text: 'Monitoring' }, { text: 'Overview' }]} />
        <MonitoringOverview />
      </>
    )
  }
}
