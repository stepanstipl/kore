import React from 'react'
import PropTypes from 'prop-types'
import { Alert } from 'antd'

import Breadcrumb from '../../../lib/components/Breadcrumb'
import SecurityData from '../../../lib/utils/dummy-security-data'
import SecurityOverview from '../../../lib/components/security/SecurityOverview'

class SecurityPage extends React.Component {
  static propTypes = {
    overview: PropTypes.object.isRequired,
  }

  state = {
    overview: []
  }

  static staticProps = {
    title: 'Security | Overview',
    adminOnly: true
  }

  static async getPageData() {
    return await Promise.resolve({ overview: SecurityData.overview })
  }

  static getInitialProps = async ctx => {
    const data = await SecurityPage.getPageData(ctx)
    return data
  }

  constructor(props) {
    super(props)
    this.state = { overview: props.overview }
  }

  render() {
    return (
      <div>
        <Breadcrumb
          items={[
            { text: 'Security', href: '/prototype/security', link: '/prototype/security' },
            { text: 'Overview' }
          ]}
        />
        <Alert 
          message="See your current security situation"
          description="This page shows, at a glance, an overview of the current security posture across all teams, clusters and plans managed by Kore."
          type="info"
          showIcon
          style={{ marginBottom: '20px' }}
        />
        <SecurityOverview overview={this.state.overview} />
      </div>
    )
  }
}

export default SecurityPage