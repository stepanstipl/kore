import React from 'react'
import PropTypes from 'prop-types'

import Breadcrumb from '../../../lib/components/layout/Breadcrumb'
import KoreApi from '../../../lib/kore-api'
import MonitoringRulesTable from '../../../lib/components/monitoring/MonitoringRulesTable'

export default class MonitoringPage extends React.Component {
  static propTypes = {
    rules: PropTypes.object.isRequired,
  }

  static staticProps = {
    title: 'Monitoring Overview',
    adminOnly: true
  }

  static async getPageData(ctx) {
    try {
      const rules = await (await KoreApi.client(ctx)).monitoring.ListRules()
      return { rules }
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
    const { rules } = this.props
    return (
      <div>
        <Breadcrumb items={[{ text: 'Monitoring' }, { text: 'Overview' }]} />
        <MonitoringRulesTable rules={rules} />
      </div>
    )
  }
}
