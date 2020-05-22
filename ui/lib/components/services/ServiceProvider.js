import * as React from 'react'
import PropTypes from 'prop-types'
import KoreApi from '../../kore-api'
import { Typography, List, Avatar, Alert, Tooltip, Switch, Icon } from 'antd'
const { Title, Text } = Typography

export default class ServiceProvider extends React.Component {
  static propTypes = {
    provider: PropTypes.object
  }

  state = {
    kinds: [],
    loading: true
  }

  loadKindsForServiceProvider = async () => {
    const { provider } = this.props
    this.setState({ loading: true })
    if (!provider) {
      this.setState({
        loading: false, kinds: []
      })
    }
    const api = await KoreApi.client()
    const allKinds = await api.ListServiceKinds()
    this.setState({
      loading: false, 
      kinds: allKinds.items.filter((k) => provider.status.supportedKinds.indexOf(k.metadata.name) > -1)
    })
  }

  componentDidMountComplete = null
  componentDidMount = () => {
    this.componentDidMountComplete = Promise.resolve().then(async () => {
      await this.loadKindsForServiceProvider()
    })
  }

  toggleKindEnabled = async (kind) => {
    const api = await KoreApi.client()
    await api.UpdateServiceKind(kind.metadata.name, { ...kind, spec: { ...kind.spec, enabled: !kind.spec.enabled } })
    await this.loadKindsForServiceProvider()
  }

  renderKind = (kind) => {
    return (
      <List.Item key={kind.metadata.name} actions={[
        <Text key="enable">
          <Tooltip title={kind.spec.enabled ? 'Disable' : 'Enable'}>
            <a onClick={() => this.toggleKindEnabled(kind)}>
              <Switch checked={kind.spec.enabled} />
            </a>
          </Tooltip>
        </Text>
      ]}>
        <List.Item.Meta 
          avatar={<Avatar src={kind.spec.imageURL} />} 
          title={kind.spec.displayName} 
          description={kind.spec.description} />
      </List.Item>
    )
  }

  render() {
    const { kinds, loading } = this.state
    return (
      <>
        {loading ? <Icon type="loading" /> : (
          <>
            <Title level={4}>Enabled Services</Title>
            <Alert type="info" message="These services are available to your teams"/>
            <List dataSource={kinds.filter((k) => k.spec.enabled)} renderItem={(kind) => this.renderKind(kind)} />
            <Title level={4} style={{ marginTop: '30px' }}>Disabled Services</Title>
            <Alert type="warning" message="These services need to be enabled before teams can consume them"/>
            <List dataSource={kinds.filter((k) => !k.spec.enabled)} renderItem={(kind) => this.renderKind(kind)} />
          </>
        )}
      </>
    )
  }
}