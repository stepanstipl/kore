import PropTypes from 'prop-types'
import moment from 'moment'
import { List, Icon, Typography, Popconfirm } from 'antd'
const { Text } = Typography

import ResourceStatusTag from '../../resources/ResourceStatusTag'
import AutoRefreshComponent from '../AutoRefreshComponent'
import { successMessage, errorMessage } from '../../../utils/message'

class NamespaceClaim extends AutoRefreshComponent {
  static propTypes = {
    team: PropTypes.string.isRequired,
    namespaceClaim: PropTypes.object.isRequired,
    deleteNamespace: PropTypes.func.isRequired
  }

  stableStateReached({ state, deleted }) {
    const { namespaceClaim } = this.props
    if (deleted) {
      return successMessage(`Namespace "${namespaceClaim.spec.name}" deleted`)
    }
    if (state === AutoRefreshComponent.STABLE_STATES.SUCCESS) {
      return successMessage(`Namespace "${namespaceClaim.spec.name}" created`)
    }
    if (state === AutoRefreshComponent.STABLE_STATES.FAILURE) {
      return errorMessage(`Namespace "${namespaceClaim.spec.name}" failed be to created`)
    }
  }

  deleteNamespace = () => {
    this.props.deleteNamespace(this.props.namespaceClaim.metadata.name, () => {
      this.startRefreshing()
    })
  }

  render() {
    const { namespaceClaim } = this.props

    if (namespaceClaim.deleted) {
      return null
    }

    const created = moment(namespaceClaim.metadata.creationTimestamp).fromNow()
    const deleted = namespaceClaim.metadata.deletionTimestamp ? moment(namespaceClaim.metadata.deletionTimestamp).fromNow() : false

    const actions = () => {
      const actions = []
      const status = namespaceClaim.status.status || 'Pending'
      if (status === 'Success') {
        const deleteAction = (
          <Popconfirm
            key="delete"
            title="Are you sure you want to delete this namespace?"
            onConfirm={this.deleteNamespace}
            okText="Yes"
            cancelText="No"
          >
            <a id={`namespace_delete_${namespaceClaim.spec.name}`}><Icon type="delete" /></a>
          </Popconfirm>
        )
        actions.push(deleteAction)
      }
      actions.push(<ResourceStatusTag id={`namespace_status_${namespaceClaim.spec.name}`} resourceStatus={namespaceClaim.status} />)
      return actions
    }

    return (
      <List.Item className="namespace" id={`namespace_${namespaceClaim.spec.name}`} style={{ paddingTop: 0 }} actions={actions()}>
        <List.Item.Meta style={{ marginLeft: '5px' }} title={<span>{namespaceClaim.spec.name}</span>} />
        <div>
          <Text type='secondary'>Created {created}</Text>
          {deleted && <Text type='secondary'><br/>Deleted {deleted}</Text>}
        </div>
      </List.Item>
    )
  }

}

export default NamespaceClaim
