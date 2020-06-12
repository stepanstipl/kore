import PropTypes from 'prop-types'
import Link from 'next/link'
import moment from 'moment'
import { Divider, List, Icon, Typography, Popconfirm, Tooltip, Tag } from 'antd'
const { Text } = Typography

import ServiceCredentialSnippet from './ServiceCredentialSnippet'
import { inProgressStatusList } from '../../../utils/ui-helpers'
import ResourceStatusTag from '../../resources/ResourceStatusTag'
import AutoRefreshComponent from '../AutoRefreshComponent'
import { successMessage, errorMessage } from '../../../utils/message'

class ServiceCredential extends AutoRefreshComponent {
  static propTypes = {
    viewPerspective: PropTypes.oneOf(['cluster', 'service']),
    hideClusterInfo: PropTypes.bool,
    team: PropTypes.string.isRequired,
    serviceCredential: PropTypes.object.isRequired,
    serviceKind: PropTypes.object.isRequired,
    deleteServiceCredential: PropTypes.func.isRequired
  }

  finalStateReached() {
    const { serviceCredential } = this.props
    const { status, deleted } = serviceCredential
    if (deleted) {
      return successMessage(`Service access successfully deleted for service "${serviceCredential.spec.service.name}"`)
    }
    if (status.status === 'Success') {
      return successMessage(`Service access successfully created for service "${serviceCredential.spec.service.name}"`)
    }
    if (status.status === 'Failure') {
      return errorMessage(`Service access failed to create for service "${serviceCredential.spec.service.name}"`)
    }
  }

  deleteServiceCredential = () => {
    this.props.deleteServiceCredential(this.props.serviceCredential.metadata.name, () => {
      this.startRefreshing()
    })
  }

  title = () => {
    const { serviceCredential, serviceKind } = this.props
    if (this.props.viewPerspective === 'cluster') {
      return (
        <>
          {serviceKind && <><Tag style={{ margin: 0 }}>{serviceKind.spec.displayName}</Tag><Divider type="vertical" /></>}
          <Link href="/teams/[name]/clusters/[cluster]/services/[service]" as={`/teams/${this.props.team}/clusters/${this.props.cluster.metadata.name}/services/${serviceCredential.spec.service.name}`}><a style={{ textDecoration: 'underline' }}>{serviceCredential.spec.service.name}</a></Link>
          <Divider type="vertical" />
          <Text>Secret name: <Text copyable style={{ fontWeight: 'normal', fontStyle: 'italic' }}>{serviceCredential.spec.secretName}</Text><ServiceCredentialSnippet serviceCredential={serviceCredential} /></Text>
        </>
      )
    }
    if (this.props.viewPerspective === 'service') {
      return (
        <>
          <Text>Cluster: <Link href="/teams/[name]/clusters/[cluster]/[tab]" as={`/teams/${this.props.team}/clusters/${serviceCredential.spec.cluster.name}/namespaces`}><a style={{ textDecoration: 'underline' }}>{serviceCredential.spec.cluster.name}</a></Link></Text>
          <Divider type="vertical" />
          <Text>namespace: <Text strong>{serviceCredential.spec.clusterNamespace}</Text></Text>
          <Divider type="vertical" />
          <Text>Secret name: <Text copyable style={{ fontWeight: 'normal', fontStyle: 'italic' }}>{serviceCredential.spec.secretName}</Text><ServiceCredentialSnippet serviceCredential={serviceCredential} /></Text>
        </>
      )
    }
  }

  render() {
    const { serviceCredential } = this.props

    if (serviceCredential.deleted) {
      return null
    }

    const created = moment(serviceCredential.metadata.creationTimestamp).fromNow()
    const deleted = serviceCredential.metadata.deletionTimestamp ? moment(serviceCredential.metadata.deletionTimestamp).fromNow() : false

    const actions = () => {
      const actions = []
      const status = serviceCredential.status.status || 'Pending'

      if (!inProgressStatusList.includes(status)) {
        const deleteAction = (
          <Popconfirm
            key="delete"
            title={`Are you sure you want to delete access with secret name ${serviceCredential.spec.secretName}?`}
            onConfirm={this.deleteServiceCredential}
            okText="Yes"
            cancelText="No"
          >
            <a><Tooltip title="Delete service access"><Icon type="delete" /></Tooltip></a>
          </Popconfirm>
        )
        actions.push(deleteAction)
      }

      actions.push(<ResourceStatusTag resourceStatus={serviceCredential.status} />)
      return actions
    }

    return (
      <List.Item actions={actions()}>
        <List.Item.Meta title={this.title()} />
        <div>
          <Text type='secondary'>Created {created}</Text>
          {deleted ? <Text type='secondary'><br/>Deleted {deleted}</Text> : null }
        </div>
      </List.Item>
    )
  }

}

export default ServiceCredential
