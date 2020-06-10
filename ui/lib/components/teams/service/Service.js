import PropTypes from 'prop-types'
import moment from 'moment'
import { List, Icon, Typography, Popconfirm, Tooltip, Avatar, Tag } from 'antd'
const { Text } = Typography

import { inProgressStatusList } from '../../../utils/ui-helpers'
import ResourceStatusTag from '../../resources/ResourceStatusTag'
import AutoRefreshComponent from '../AutoRefreshComponent'
import Link from 'next/link'
import { getKoreLabel, isReadOnlyCRD } from '../../../utils/crd-helpers'
import { successMessage } from '../../../utils/message'

class Service extends AutoRefreshComponent {
  static propTypes = {
    team: PropTypes.string.isRequired,
    cluster: PropTypes.object.isRequired,
    service: PropTypes.object.isRequired,
    serviceKind: PropTypes.object,
    deleteService: PropTypes.func.isRequired,
    style: PropTypes.object
  }

  finalStateReached() {
    const { service } = this.props
    const { status, deleted } = service
    if (deleted) {
      return successMessage(`Service successfully deleted: ${service.metadata.name}`)
    }
    if (status.status === 'Success') {
      return successMessage(`Service successfully created: ${service.metadata.name}`)
    }
    if (status.status === 'Failure') {
      return errorMessage(`Service failed to create: ${service.metadata.name}`)
    }
  }

  isApplicationService = () => {
    return getKoreLabel(this.props.serviceKind, 'platform') === 'Kubernetes'
  }

  deleteService = () => {
    this.props.deleteService(this.props.service.metadata.name, () => {
      this.startRefreshing()
    })
  }

  render() {
    const { service, serviceKind, team, cluster, style } = this.props

    if (service.deleted) {
      return null
    }

    const styleOverrides = style || {}

    const created = moment(service.metadata.creationTimestamp).fromNow()
    const deleted = service.metadata.deletionTimestamp ? moment(service.metadata.deletionTimestamp).fromNow() : false

    const actions = () => {
      const actions = []
      const status = service.status.status || 'Pending'
      const readonly = isReadOnlyCRD(service)

      actions.push((
        <Link key="view" href="/teams/[name]/clusters/[cluster]/services/[service]" as={`/teams/${team}/clusters/${cluster.metadata.name}/services/${service.metadata.name}`}><a><Tooltip title="Service details"><Icon type="info-circle" /></Tooltip></a></Link>
      ))

      if (!readonly && !inProgressStatusList.includes(status)) {
        const deleteAction = (
          <Popconfirm
            key="delete"
            title="Are you sure you want to delete this service?"
            onConfirm={this.deleteService}
            okText="Yes"
            cancelText="No"
          >
            <a><Tooltip title="Delete this service"><Icon type="delete" /></Tooltip></a>
          </Popconfirm>
        )
        actions.push(deleteAction)
      }

      actions.push(<ResourceStatusTag resourceStatus={service.status} />)
      return actions
    }

    return (
      <List.Item actions={actions()} style={{ ...styleOverrides }}>
        {this.isApplicationService() ? (
          <List.Item.Meta
            avatar={serviceKind && serviceKind.spec.imageURL ? <Avatar src={serviceKind.spec.imageURL} /> : <Avatar icon="cloud-server" />}
            title={<><Text style={{ marginRight: '15px' }}>{service.metadata.name}</Text><Tag style={{ margin: 0 }}>{serviceKind.spec.displayName}</Tag></>}
            description={<Text>Namespace: <Text strong>{service.spec.clusterNamespace}</Text></Text>} />
        ) : (
          <List.Item.Meta
            avatar={serviceKind && serviceKind.spec.imageURL ? <Avatar src={serviceKind.spec.imageURL} /> : <Avatar icon="cloud-server" />}
            title={<><Link href="/teams/[name]/clusters/[cluster]/services/[service]" as={`/teams/${team}/clusters/${cluster.metadata.name}/services/${service.metadata.name}`}><a><Text style={{ marginRight: '15px', fontSize: '16px', textDecoration: 'underline' }}>{service.metadata.name}</Text></a></Link><Tag style={{ margin: 0 }}>{serviceKind.spec.displayName}</Tag></>}
          />
        )}
        <div>
          <Text type='secondary'>Created {created}</Text>
          {deleted ? <Text type='secondary'><br/>Deleted {deleted}</Text> : null }
        </div>
      </List.Item>
    )
  }

}

export default Service
