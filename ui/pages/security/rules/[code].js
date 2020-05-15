import React from 'react'
import PropTypes from 'prop-types'
import { Typography } from 'antd'
const { Title } = Typography

import KoreApi from '../../../lib/kore-api'
import Breadcrumb from '../../../lib/components/layout/Breadcrumb'
import SecurityRule from '../../../lib/components/security/SecurityRule'

export default class SecurityRulePage extends React.Component {
  static propTypes = {
    rule: PropTypes.object.isRequired
  }

  constructor(props) {
    super(props)
  }

  static getInitialProps = async ctx => {
    const api = await KoreApi.client(ctx)
    const rule = await api.security.GetSecurityRule(ctx.query.code)
    if (!rule && ctx.res) {
      /* eslint-disable-next-line require-atomic-updates */
      ctx.res.statusCode = 404
    }
    return { rule }
  }

  render = () => {
    const { rule } = this.props
    if (!rule) {
      return null
    }
    return (
      <div>
        <Breadcrumb
          items={[
            { text: 'Security', href: '/security', link: '/security' },
            { text: 'Security Rules', href: '/security/rules', link: '/security/rules' },
            { text: `Rule: ${rule.spec.code} ${rule.spec.name}` }
          ]}
        />
        <Title level={3} style={{ marginBottom: '20px' }}>{rule.spec.code}: {rule.spec.name}</Title>
        <SecurityRule rule={rule} />
      </div>
    )
  }
}
