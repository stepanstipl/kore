import React from 'react'
import PropTypes from 'prop-types'
import { Typography } from 'antd'
const { Title } = Typography

import KoreApi from '../../../../../../../lib/kore-api'
import Breadcrumb from '../../../../../../../lib/components/layout/Breadcrumb'
import SecurityScanResult from '../../../../../../../lib/components/security/SecurityScanResult'
import { errorMessage } from '../../../../../../../lib/utils/message'

export default class SecurityResourcePage extends React.Component {
  static propTypes = {
    result: PropTypes.object.isRequired
  }

  state = {
    history: null,
    historyLoading: false
  }

  constructor(props) {
    super(props)
  }

  static getInitialProps = async ctx => {
    const api = await KoreApi.client(ctx)
    const result = await api.security.GetSecurityScanForResource(ctx.query.group, ctx.query.version, ctx.query.kind, ctx.query.namespace, ctx.query.name)
    if (!result && ctx.res) {
      /* eslint-disable-next-line require-atomic-updates */
      ctx.res.statusCode = 404
    }
    return { result }
  }

  loadHistory = async () => {
    this.setState({ historyLoading: true })
    const resource = this.props.result.spec.resource
    try {
      const api = await KoreApi.client()
      const history = await api.security.ListSecurityScansForResource(resource.group, resource.version, resource.kind, resource.namespace, resource.name)
      this.setState({ historyLoading: false, history: history })
    } catch (err) {
      console.error('Error loading history', err)
      this.setState({ historyLoading: false })
      errorMessage(`Failed to load history for resource ${resource.name}`)
    }
  }

  render = () => {
    const { result } = this.props
    const { history, historyLoading } = this.state
  
    return (
      <div>
        <Breadcrumb
          items={[
            { text: 'Security', href: '/security', link: '/security' },
            { text: `Resource: ${result.spec.resource.kind} ${result.spec.resource.name}` }
          ]}
        />
        <Title level={3} style={{ marginBottom: '20px' }}>Current security status for {result.spec.resource.kind.toLowerCase()} {result.spec.resource.name}</Title>
        <SecurityScanResult result={result} history={history} historyLoading={historyLoading} onHistoryRequest={this.loadHistory} />
      </div>
    )
  }
}
