import PropTypes from 'prop-types'
import moment from 'moment'
import Link from 'next/link'
import { List, Icon, Typography, Modal, Popconfirm, Tag, Tooltip } from 'antd'
const { Text, Paragraph } = Typography

import { clusterProviderIconSrcMap, inProgressStatusList } from '../../../utils/ui-helpers'
import ResourceStatusTag from '../../resources/ResourceStatusTag'
import AutoRefreshComponent from '../AutoRefreshComponent'
import { isReadOnlyCRD } from '../../../utils/crd-helpers'
import { successMessage, errorMessage } from '../../../utils/message'

class Cluster extends AutoRefreshComponent {
  static propTypes = {
    team: PropTypes.string.isRequired,
    cluster: PropTypes.object.isRequired,
    plan: PropTypes.object.isRequired,
    namespaceClaims: PropTypes.array.isRequired,
    deleteCluster: PropTypes.func.isRequired
  }

  finalStateReached() {
    const { cluster } = this.props
    const { status, deleted } = cluster
    if (deleted) {
      return successMessage(`Cluster successfully deleted: ${cluster.metadata.name}`)
    }
    if (status.status === 'Success') {
      return successMessage(`Cluster successfully created: ${cluster.metadata.name}`)
    }
    if (status.status === 'Failure') {
      return errorMessage(`Cluster failed to create: ${cluster.metadata.name}`)
    }
  }

  deleteCluster = () => {
    const { namespaceClaims } = this.props
    if (namespaceClaims.length > 0) {
      return Modal.warning({
        title: 'Warning: cluster cannot be deleted',
        content: (
          <div>
            <Paragraph strong>The cluster namespaces must be deleted first</Paragraph>
            <List
              size="small"
              dataSource={namespaceClaims}
              renderItem={ns => <List.Item>{ns.spec.name}</List.Item>}
            />
          </div>
        ),
        onOk() {}
      })
    }

    this.props.deleteCluster(this.props.cluster.metadata.name, () => {
      this.startRefreshing()
    })
  }

  render() {
    const { cluster, team, plan } = this.props

    if (cluster.deleted) {
      return null
    }

    const created = moment(cluster.metadata.creationTimestamp).fromNow()
    const deleted = cluster.metadata.deletionTimestamp ? moment(cluster.metadata.deletionTimestamp).fromNow() : false

    const actions = () => {
      const actions = []
      const status = cluster.status.status || 'Pending'

      actions.push((
        <Link key="view" href="/teams/[name]/clusters/[cluster]" as={`/teams/${team}/clusters/${cluster.metadata.name}`}><a><Tooltip title="Cluster details"><Icon type="info-circle" /></Tooltip></a></Link>
      ))

      if (!isReadOnlyCRD(cluster) && !inProgressStatusList.includes(status)) {
        const deleteAction = (
          <Popconfirm
            key="delete"
            title="Are you sure you want to delete this cluster?"
            onConfirm={this.deleteCluster}
            okText="Yes"
            cancelText="No"
          >
            <a><Tooltip title="Delete this cluster"><Icon type="delete" /></Tooltip></a>
          </Popconfirm>
        )
        actions.push(deleteAction)
      }

      actions.push(<ResourceStatusTag resourceStatus={cluster.status} />)
      return actions
    }

    return (
      <List.Item key={cluster.metadata.name} actions={actions()} style={{ paddingTop: 0, paddingBottom: '5px' }}>
        <List.Item.Meta
          avatar={<img src={clusterProviderIconSrcMap[cluster.spec.kind]} height="32px" />}
          title={<><Link href="/teams/[name]/clusters/[cluster]/[tab]" as={`/teams/${team}/clusters/${cluster.metadata.name}/namespaces`}><a><Text style={{ marginRight: '15px', fontSize: '16px', textDecoration: 'underline' }}>{cluster.metadata.name}</Text></a></Link>{ plan && <Tag style={{ margin: 0 }}>{plan.spec.description}</Tag> }</>}
        />
        <div>
          <Text type='secondary'>Created {created}</Text>
          {deleted ? <Text type='secondary'><br/>Deleted {deleted}</Text> : null }
        </div>
      </List.Item>
    )
  }

}

export default Cluster
