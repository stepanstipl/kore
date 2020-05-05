import React from 'react'
import PropTypes from 'prop-types'

import Breadcrumb from '../../lib/components/layout/Breadcrumb'
import AuditViewer from '../../lib/components/common/AuditViewer'
import KoreApi from '../../lib/kore-api'

class AuditPage extends React.Component {
  static propTypes = {
    events: PropTypes.array.isRequired,
  }

  state = {
    events: []
  }

  static staticProps = {
    title: 'Audit Viewer',
    adminOnly: true
  }

  static async getPageData(ctx) {
    try {
      const eventList = await (await KoreApi.client(ctx)).ListAuditEvents()
      const events = eventList.items
      return { events }
    } catch (err) {
      throw new Error(err.message)
    }
  }

  static getInitialProps = async ctx => {
    const data = await AuditPage.getPageData(ctx)
    return data
  }

  constructor(props) {
    super(props)
    this.state = { events: props.events }
  }

  render() {
    return (
      <div>
        <Breadcrumb items={[{ text: 'Audit' }, { text: 'Audit Viewer' }]} />
        <AuditViewer items={this.state.events} />
      </div>
    )
  }
}

export default AuditPage
