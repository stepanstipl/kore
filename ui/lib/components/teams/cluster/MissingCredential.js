import React from 'react'
import Link from 'next/link'
import PropTypes from 'prop-types'
import { Alert, Button, Icon, Typography } from 'antd'
const { Paragraph, Text } = Typography
import { pluralize } from 'inflect'

import KoreApi from '../../../kore-api'
import { getCloudInfo } from '../../../utils/cloud'

class MissingCredential extends React.Component {
  static propTypes = {
    team: PropTypes.string.isRequired,
    cloud: PropTypes.oneOf(['GCP', 'AWS', 'Azure']).isRequired,
  }

  state = {
    dataLoading: true,
    cloudInfo: getCloudInfo(this.props.cloud)
  }

  async fetchComponentData() {
    const cloudConfig = await (await KoreApi.client()).GetConfig(this.props.cloud)
    const email = cloudConfig && cloudConfig.spec.values.requestAccessEmail
    return { email }
  }

  componentDidMount() {
    this.fetchComponentData()
      .then(data => this.setState({ ...data, dataLoading: false }))
  }

  componentDidUpdate(prevProps) {
    if (prevProps.cloud !== this.props.cloud) {
      this.setState({ dataLoading: true, cloudInfo: getCloudInfo(this.props.cloud) })
      this.fetchComponentData()
        .then(data => this.setState({ ...data, dataLoading: false }))
    }
  }

  render() {
    const { dataLoading, email, cloudInfo } = this.state

    if (dataLoading) {
      return <Icon type="loading" />
    }

    return (
      <>
        <Alert
          message={`${cloudInfo.cloud} ${cloudInfo.accountNoun} access not found`}
          description={
            <>
              <Paragraph>This team does not have access to create clusters in any {cloudInfo.cloud} {pluralize(cloudInfo.accountNoun)}. Please use the contact below to grant this team access to a {cloudInfo.cloud} {cloudInfo.accountNoun}.</Paragraph>
              <Text strong>
                {email ? <a style={{ textDecoration: 'underline' }} href={`mailto: ${this.state.email}`}>{this.state.email}</a> : 'Kore administrator'}
              </Text>
            </>
          }
          type="warning"
          showIcon
          style={{ marginBottom: '30px' }}
        />
        <Button type="primary">
          <Link href="/teams/[name]" as={`/teams/${this.props.team}`}>
            <a>Team dashboard</a>
          </Link>
        </Button>
      </>
    )
  }
}

export default MissingCredential
