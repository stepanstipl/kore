import React from 'react'
import PropTypes from 'prop-types'

import Breadcrumb from '../../../lib/components/layout/Breadcrumb'
import SecurityOverview from '../../../lib/components/security/SecurityOverview'

// prototype imports
import SecurityData from '../../../lib/prototype/utils/dummy-security-data'

export default class SecurityReport extends React.Component {
  static propTypes = {
    overview: PropTypes.object.isRequired,
  }

  static staticProps = {
    title: 'Security Overview',
    hideSider: true,
    adminOnly: true
  }

  static async getPageData() {
    return await Promise.resolve({ overview: SecurityData.overview })
  }

  static getInitialProps = async ctx => {
    const data = await SecurityReport.getPageData(ctx)
    return data
  }

  constructor(props) {
    super(props)
  }

  render() {
    const { overview } = this.props
    return (
      <>
        <Breadcrumb items={[ { text: 'Reports', href: '/prototype/reports', link: '/prototype/reports' }, { text: 'Security report' } ]} />
        <SecurityOverview overview={overview} />
      </>
    )
  }
}
