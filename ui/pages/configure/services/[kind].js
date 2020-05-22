import React from 'react'
import PropTypes from 'prop-types'
import KoreApi from '../../../lib/kore-api'
import ServiceKindManage from '../../../lib/components/services/ServiceKindManage'
import Breadcrumb from '../../../lib/components/layout/Breadcrumb'

export default class ServiceKindPage extends React.Component {
  static propTypes = {
    details: PropTypes.object.isRequired
  }

  static getInitialProps = async ctx => {
    const { kind } = ctx.query
    const api = await KoreApi.client(ctx)
    const details = await api.GetServiceKind(kind)
    if ((!details) && ctx.res) {
      /* eslint-disable-next-line require-atomic-updates */
      ctx.res.statusCode = 404
    }
    return { details }
  }

  render() {
    const { details } = this.props
    const displayName = details.spec.displayName || details.metadata.name
    return (
      <>
        <Breadcrumb
          items={[
            { text: 'Services', href: '/configure/services', link: '/configure/services' },
            { text: `Service: ${displayName}` }
          ]}
        />
        <ServiceKindManage kind={details} />
      </>
    )
  }
}