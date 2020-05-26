import PropTypes from 'prop-types'
import moment from 'moment'
import { List, Icon, Typography, Popconfirm, message, Tooltip, Avatar } from 'antd'
const { Text } = Typography

import { inProgressStatusList } from '../../../utils/ui-helpers'
import ResourceStatusTag from '../../resources/ResourceStatusTag'
import AutoRefreshComponent from '../AutoRefreshComponent'

class ServiceCredential extends AutoRefreshComponent {
  static propTypes = {
    team: PropTypes.string.isRequired,
    serviceCredential: PropTypes.object.isRequired,
    deleteServiceCredential: PropTypes.func.isRequired
  }

  finalStateReached() {
    const { serviceCredential } = this.props
    const { status, deleted } = serviceCredential
    if (deleted) {
      return message.success(`Service Credential successfully deleted: ${serviceCredential.metadata.name}`)
    }
    if (status.status === 'Success') {
      return message.success(`Service Credential successfully created: ${serviceCredential.metadata.name}`)
    }
    if (status.status === 'Failure') {
      return message.error(`Service Credential failed to create: ${serviceCredential.metadata.name}`)
    }
  }

  deleteServiceCredential = () => {
    this.props.deleteServiceCredential(this.props.serviceCredential.metadata.name, () => {
      this.startRefreshing()
    })
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
            title={`Are you sure you want to delete ${serviceCredential.metadata.name}?`}
            onConfirm={this.deleteServiceCredential}
            okText="Yes"
            cancelText="No"
          >
            <a><Tooltip title="Delete this service credential"><Icon type="delete" /></Tooltip></a>
          </Popconfirm>
        )
        actions.push(deleteAction)
      }

      actions.push(<ResourceStatusTag resourceStatus={serviceCredential.status} />)
      return actions
    }

    return (
      <List.Item actions={actions()}>
        <List.Item.Meta
          avatar={<Avatar icon="database" />}
          title={<Text>{serviceCredential.spec.kind} <Text style={{ fontFamily: 'monospace', marginLeft: '15px' }}>{serviceCredential.metadata.name}</Text></Text>}
          description={
            <>
              <div>
                <Text>Cluster: <b>{serviceCredential.spec.cluster.name}</b>, namespace: <b>{serviceCredential.spec.clusterNamespace}</b>, secret name: <b>{serviceCredential.spec.secretName}</b></Text>
              </div>
              <div>
                <Text type='secondary'>Created {created}</Text>
                {deleted ? <Text type='secondary'><br/>Deleted {deleted}</Text> : null }
              </div>
            </>
          }
        />
      </List.Item>
    )
  }

}

export default ServiceCredential
