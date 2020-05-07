import React from 'react'
import PropTypes from 'prop-types'

import Breadcrumb from '../../lib/components/layout/Breadcrumb'
import SecurityOverview from '../../lib/components/security/SecurityOverview'
import KoreApi from '../../lib/kore-api'

export default class SecurityPage extends React.Component {
  static propTypes = {
    overview: PropTypes.object.isRequired,
  }

  static staticProps = {
    title: 'Security Overview',
    adminOnly: true
  }

  static async getPageData(ctx) {
    try {
      const overview = await (await KoreApi.client(ctx)).security.GetSecurityOverview()
      return { overview }
    } catch (err) {
      throw new Error(err.message)
    }
  }

  static getInitialProps = async ctx => {
    const data = await SecurityPage.getPageData(ctx)
    return data
  }

  constructor(props) {
    super(props)
  }

  render() {
    const { overview } = this.props
    return (
      <div>
        <Breadcrumb items={[{ text: 'Security' }, { text: 'Security Overview' }]} />
        <SecurityOverview overview={overview} />
      </div>
    )
  }
}
