import React from 'react'
import PropTypes from 'prop-types'
import { Typography } from 'antd'
const { Title } = Typography
import moment from 'moment'

import KoreApi from '../../../lib/kore-api'
import Breadcrumb from '../../../lib/components/layout/Breadcrumb'
import SecurityScanResult from '../../../lib/components/security/SecurityScanResult'

export default class SecurityScanPage extends React.Component {
  static propTypes = {
    result: PropTypes.object.isRequired
  }

  constructor(props) {
    super(props)
  }

  static getInitialProps = async ctx => {
    const api = await KoreApi.client(ctx)
    const result = await api.security.GetSecurityScan(ctx.query.id)
    if (!result && ctx.res) {
      /* eslint-disable-next-line require-atomic-updates */
      ctx.res.statusCode = 404
    }
    return { result }
  }

  render = () => {
    const { result } = this.props
    const checked = moment(result.spec.checkedAt).format('D MMM YYYY HH:mm:ss')
    return (
      <div>
        <Breadcrumb
          items={[
            { text: 'Security', href: '/security', link: '/security' },
            { text: `Resource: ${result.spec.resource.kind} ${result.spec.resource.name} as at ${checked}` }
          ]}
        />
        <Title level={3} style={{ marginBottom: '20px' }}>Security scan of {result.spec.resource.kind.toLowerCase()} {result.spec.resource.name} at {checked}</Title>
        <SecurityScanResult result={result} />
      </div>
    )
  }
}
