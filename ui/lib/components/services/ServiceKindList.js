import * as React from 'react'
import PropTypes from 'prop-types'
import { sortBy } from 'lodash'
import KoreApi from '../../kore-api'
import { Alert, Avatar, Col, Icon, List, Row, Switch, Tooltip, Typography } from 'antd'
import Link from 'next/link'
import { successMessage, errorMessage } from '../../utils/message'
const { Text, Title } = Typography
import { featureEnabled, KoreFeatures } from '../../utils/features'
import { isReadOnlyCRD } from '../../utils/crd-helpers'

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
    this.setState({ loading: true })
    let kinds = await (await KoreApi.client()).ListServiceKinds()
    // We have to filter out service kinds which are disabled and read-only (as you can't enable them)
    const exclude = (k) => !k.spec.enabled && isReadOnlyCRD(k)
    kinds.items = kinds.items.filter(k => !exclude(k))

    if (!featureEnabled(KoreFeatures.APPLICATION_SERVICES)) {
      kinds.items = kinds.items.filter(k => !k.metadata.labels || k.metadata.labels['kore.appvia.io/platform'] !== 'Kubernetes')
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
      successMessage(`${enabled ? 'Enabled' : 'Disabled'} service "${kind.spec.displayName}"`)
    } catch (error) {
      errorMessage(`Failed to ${enabled ? 'enable' : 'disable'} service "${kind.spec.displayName}", please try again.`)
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
    let filteredKinds = kinds
    if (this.props.filter) {
      filteredKinds = filteredKinds.filter(this.props.filter)
    }
    if (!filteredKinds || filteredKinds.length === 0) {
      return <Alert type="warning" message="No matching services can be found." />
    }
    return (
      <Row type="flex" gutter={[24,24]}>
        <Col span={24} xl={12}>
          <Title level={4}>Enabled Services</Title>
          <Alert style={{ marginBottom: '10px' }} type="info" message="These services are currently available to teams."/>
          <List dataSource={filteredKinds.filter((k) => k.spec.enabled)} renderItem={(kind) => this.renderKind(kind)} />
        </Col>
        <Col span={24} xl={12}>
          <Title level={4}>Disabled Services</Title>
          <Alert style={{ marginBottom: '10px' }} type="warning" message="These services need to be enabled before teams can consume them."/>
          <List dataSource={filteredKinds.filter((k) => !k.spec.enabled)} renderItem={(kind) => this.renderKind(kind)} />
        </Col>
      </Row>
    )
  }
}
