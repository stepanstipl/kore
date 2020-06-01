import * as React from 'react'
import PropTypes from 'prop-types'
import { sortBy } from 'lodash'
import KoreApi from '../../kore-api'
import { Typography, List, Avatar, Alert, Tooltip, Switch, Icon, message } from 'antd'
import Link from 'next/link'
const { Title, Text } = Typography

export default class ServiceKindList extends React.Component {
  static propTypes = {
    /**
     * Pass a function of type (kind) => bool (e.g. `(kind) => kind.metadata.labels['cloud'] === 'AWS'`) to filter the list.
     */
    filter: PropTypes.func
  }

  state = {
    loading: true,
    kinds: []
  }

  componentDidMountComplete = null
  componentDidMount = () => {
    this.componentDidMountComplete = this.loadKinds()
  }

  loadKinds = async () => {
    const { filter } = this.props
    this.setState({ loading: true })
    let kinds = await (await KoreApi.client()).ListServiceKinds()
    if (filter) { 
      kinds.items = kinds.items.filter(filter)
    }
    this.setState({
      loading: false, 
      kinds: kinds.items
    })
  }

  toggleKindEnabled = async (kind, enabled) => {
    try {
      const api = await KoreApi.client()
      const serviceKindResult = await api.UpdateServiceKind(kind.metadata.name, { ...kind, spec: { ...kind.spec, enabled } })
      this.setState({
        kinds: sortBy(this.state.kinds.filter(k => k.metadata.name !== kind.metadata.name).concat([ serviceKindResult ]), k => k.spec.displayName.toLowerCase())
      })
      message.success(`${enabled ? 'Enabled' : 'Disabled'} service "${kind.spec.displayName}"`)
    } catch (error) {
      message.error(`Failed to ${enabled ? 'enable' : 'disable'} service "${kind.spec.displayName}", please try again.`)
    }
  }

  renderKind = (kind) => {
    const actions = []
    
    if (kind.spec.enabled) {
      actions.push(
        <Text key="manage">
          <Link href="/configure/services/[kind]" as={`/configure/services/${kind.metadata.name}`}>
            <Tooltip title={`Manage plans for ${kind.spec.displayName}`}>
              <a>
                <Icon type="setting" /> Manage
              </a>
            </Tooltip>
          </Link>
        </Text>
      )
    }

    actions.push(
      <Text key="enable">
        <Switch onChange={(enabled) => this.toggleKindEnabled(kind, enabled)} checked={kind.spec.enabled} checkedChildren={<Icon type="check" />} unCheckedChildren={<Icon type="close" />} />
      </Text>
    )

    const avatar = kind.spec.imageURL ? <Avatar src={kind.spec.imageURL} /> : <Avatar icon="cloud-server" />
    return (
      <List.Item key={kind.metadata.name} actions={actions}>
        <List.Item.Meta 
          avatar={avatar} 
          title={kind.spec.displayName} 
          description={kind.spec.description} />
      </List.Item>
    )
  }

  render() {
    const { kinds, loading } = this.state
    if (loading) {
      return <Icon type="loading" />
    }
    if (!kinds || kinds.length === 0) {
      return <Alert type="warning" message="No services are available to provision" />
    }
    return (
      <>
        <Title level={4}>Enabled Services</Title>
        <List dataSource={kinds.filter((k) => k.spec.enabled)} renderItem={(kind) => this.renderKind(kind)} />
        <Title level={4} style={{ marginTop: '30px' }}>Disabled Services</Title>
        <Alert type="warning" message="These services need to be enabled before teams can consume them"/>
        <List dataSource={kinds.filter((k) => !k.spec.enabled)} renderItem={(kind) => this.renderKind(kind)} />
      </>
    )
  }
}
