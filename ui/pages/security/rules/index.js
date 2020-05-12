import React from 'react'
import PropTypes from 'prop-types'

import Breadcrumb from '../../../lib/components/layout/Breadcrumb'
import KoreApi from '../../../lib/kore-api'
import SecurityRuleList from '../../../lib/components/security/SecurityRuleList'

export default class SecurityRulesPage extends React.Component {
  static propTypes = {
    rules: PropTypes.array.isRequired,
  }

  static staticProps = {
    title: 'Security Rules'
  }

  static async getPageData(ctx) {
    try {
      const rules = await (await KoreApi.client(ctx)).security.ListSecurityRules()
      return { rules }
    } catch (err) {
      throw new Error(err.message)
    }
  }

  static getInitialProps = async ctx => {
    const data = await SecurityRulesPage.getPageData(ctx)
    return data
  }

  constructor(props) {
    super(props)
  }

  render() {
    const { rules } = this.props
    return (
      <div>
        <Breadcrumb items={[{ text: 'Security' }, { text: 'Security Rules' }]} />
        <SecurityRuleList rules={rules} />
      </div>
    )
  }
}
