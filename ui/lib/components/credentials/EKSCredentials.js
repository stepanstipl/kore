import PropTypes from 'prop-types'
import moment from 'moment'
import { List, Avatar, Icon, Typography, Tooltip } from 'antd'
const { Text } = Typography

import ResourceVerificationStatus from '../resources/ResourceVerificationStatus'
import AutoRefreshComponent from '../teams/AutoRefreshComponent'
import { successMessage, errorMessage } from '../../utils/message'

class EKSCredentials extends AutoRefreshComponent {
  static propTypes = {
    eksCredentials: PropTypes.object.isRequired,
    allTeams: PropTypes.array.isRequired,
    editEKSCredentials: PropTypes.func.isRequired,
    deleteEKSCredentials: PropTypes.func.isRequired
  }

  componentDidUpdate(prevProps) {
    const prevStatus = prevProps.eksCredentials.status && prevProps.eksCredentials.status.status
    const newStatus = this.props.eksCredentials.status && this.props.eksCredentials.status.status
    if (prevStatus !== newStatus) {
      this.startRefreshing()
    }
  }

  stableStateReached({ state }) {
    const { eksCredentials } = this.props
    if (state === AutoRefreshComponent.STABLE_STATES.SUCCESS) {
      return successMessage(`AWS credentials for account "${eksCredentials.spec.accountID}" verified successfully`)
    }
    if (state === AutoRefreshComponent.STABLE_STATES.FAILURE) {
      return errorMessage(`AWS credentials for account "${eksCredentials.spec.accountID}" could not be verified`)
    }
  }

  render() {
    const { eksCredentials, editEKSCredentials, deleteEKSCredentials, allTeams } = this.props
    const created = moment(eksCredentials.metadata.creationTimestamp).fromNow()

    const displayAllocations = () => {
      if (!eksCredentials.allocation) {
        return <Text>No teams <Tooltip title="This account is not allocated to any teams, click edit to fix this."><Icon type="warning" theme="twoTone" twoToneColor="orange" /></Tooltip> </Text>
      }
      const allocatedTeams = allTeams.filter(team => eksCredentials.allocation.spec.teams.includes(team.metadata.name)).map(team => team.spec.summary)
      return allocatedTeams.length > 0 ? allocatedTeams.join(', ') : 'All teams'
    }

    return (
      <List.Item id={`ekscreds_${eksCredentials.metadata.name}`} key={eksCredentials.metadata.name} actions={[
        <ResourceVerificationStatus key="verification_status" resourceStatus={eksCredentials.status} />,
        <Text key="delete_creds"><a id={`ekscreds_del_${eksCredentials.metadata.name}`} onClick={deleteEKSCredentials(eksCredentials)}><Icon type="delete" theme="filled"/> Delete</a></Text>,
        <Text key="show_creds"><a id={`ekscreds_edit_${eksCredentials.metadata.name}`} onClick={editEKSCredentials(eksCredentials)}><Icon type="edit" theme="filled"/> Edit</a></Text>
      ]}>
        <List.Item.Meta
          avatar={<Avatar icon="amazon" />}
          title={
            <>
              <Text style={{ display: 'inline', marginRight: '15px', fontSize: '20px', fontWeight: '600' }}>{eksCredentials.spec.accountID}</Text>
              <Text style={{ marginRight: '5px' }}>{eksCredentials.allocation ? eksCredentials.allocation.spec.name : null}</Text>
              <Tooltip title={eksCredentials.allocation ? eksCredentials.allocation.spec.summary : null}>
                <Icon type="info-circle" theme="twoTone" />
              </Tooltip>
            </>
          }
          description={
            <Text>Allocated to: {displayAllocations()}</Text>
          }
        />
        <Text type='secondary'>Created {created}</Text>
      </List.Item>
    )
  }

}

export default EKSCredentials
