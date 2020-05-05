import React from 'react'
import PropTypes from 'prop-types'

import Breadcrumb from '../../../lib/components/layout/Breadcrumb'
import SecurityEventViewer from '../../../lib/prototype/components/security/SecurityEventViewer'
import SecurityData from '../../../lib/prototype/utils/dummy-security-data'

class SecurityReviewPage extends React.Component {
  static propTypes = {
    securityEvents: PropTypes.array.isRequired,
  }

  state = {
    securityEvents: []
  }

  static staticProps = {
    title: 'Security | Detailed Review',
    adminOnly: true
  }

  static async getPageData() {
    return await Promise.resolve({ securityEvents: SecurityData.events })
  }

  static getInitialProps = async ctx => {
    const data = await SecurityReviewPage.getPageData(ctx)
    return data
  }

  constructor(props) {
    super(props)
    this.state = { securityEvents: props.securityEvents }
  }

  render() {
    return (
      <div>
        <Breadcrumb
          items={[
            { text: 'Security', href: '/prototype/security', link: '/prototype/security' },
            { text: 'Detailed Review', href: '/prototype/security/review', link: '/prototype/security/review' }
          ]}
        />
        <SecurityEventViewer items={this.state.securityEvents} />
      </div>
    )
  }
}

export default SecurityReviewPage