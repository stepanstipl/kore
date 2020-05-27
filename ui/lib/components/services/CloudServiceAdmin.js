import * as React from 'react'
import PropTypes from 'prop-types'
import KoreApi from '../../kore-api'
import { Alert, Icon } from 'antd'
import ServiceKindList from './ServiceKindList'

export default class CloudServiceAdmin extends React.Component {
  static propTypes = {
    cloud: PropTypes.string
  }

  state = {
    loading: true,
    provider: null
  }

  static BROKERS = {
    'AWS': 'aws-servicebroker'
  }

  async getServiceProvider(type) {
    const api = await KoreApi.client()
    const serviceProviders = await api.ListServiceProviders()
    return serviceProviders.items.find((p) => p.spec.type === type)
  }

  componentDidMountComplete = null  
  componentDidMount = () => {
    this.componentDidMountComplete = Promise.resolve().then(async () => {
      let provider = null
      if (CloudServiceAdmin.BROKERS[this.props.cloud]) {
        provider = await this.getServiceProvider(CloudServiceAdmin.BROKERS[this.props.cloud])
      }
      this.setState({
        loading: false, provider
      })
    })
  }

  filterKindsForProvider = (kind) => {
    const { provider } = this.state
    return provider.status.supportedKinds.indexOf(kind.metadata.name) > -1
  }

  render() {
    const { cloud } = this.props
    const { loading, provider } = this.state
    return (
      <>
        {loading ? <Icon type="loading" /> : (
          <>{!provider ? <Alert type="warning" message={`No provider available to provision ${cloud} cloud services`}/> : <ServiceKindList filter={this.filterKindsForProvider} />}</>
        )}
      </>
    )
  }
}