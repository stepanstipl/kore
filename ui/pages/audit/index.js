import React from 'react'
import PropTypes from 'prop-types'

import apiRequest from '../../lib/utils/api-request'
import apiPaths from '../../lib/utils/api-paths'
import Breadcrumb from '../../lib/components/Breadcrumb'
import AuditViewer from '../../lib/components/AuditViewer'

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

  static async getPageData({ req, res }) {
    const getAuditEvents = () => apiRequest({ req, res }, 'get', apiPaths.audit)

    return getAuditEvents()
      .then((eventList) => {
        var events = eventList.items
        return { events }
      })
      .catch(err => {
        throw new Error(err.message)
      })
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
