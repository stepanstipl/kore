import PropTypes from 'prop-types'
import moment from 'moment'
import { List, Avatar, Icon, Typography, message, Tooltip } from 'antd'
const { Text } = Typography

import ResourceVerificationStatus from '../ResourceVerificationStatus'
import AutoRefreshComponent from './AutoRefreshComponent'

class EKSCredentials extends AutoRefreshComponent {
  static propTypes = {
    eksCredentials: PropTypes.object.isRequired,
    allTeams: PropTypes.array.isRequired,
    editEKSCredentials: PropTypes.func.isRequired
  }

  componentDidUpdate(prevProps) {
    const prevStatus = prevProps.eksCredentials.status && prevProps.eksCredentials.status.status
    const newStatus = this.props.eksCredentials.status && this.props.eksCredentials.status.status
    if (prevStatus !== newStatus) {
      this.startRefreshing()
    }
  }

  finalStateReached() {
    const { eksCredentials } = this.props
    const { status } = eksCredentials
    if (status.status === 'Success') {
      return message.success(`AWS credentials for account "${eksCredentials.spec.accountID}" verified successfully`)
    }
    if (status.status === 'Failure') {
      return message.error(`AWS credentials for account "${eksCredentials.spec.accountID}" could not be verified`)
    }
  }

  render() {
    const { eksCredentials, editEKSCredentials, allTeams } = this.props
    const created = moment(eksCredentials.metadata.creationTimestamp).fromNow()

    const displayAllocations = () => {
      if (!eksCredentials.allocation) {
        return <Text>No teams <Tooltip title="This account is not allocated to any teams, click edit to fix this."><Icon type="warning" theme="twoTone" twoToneColor="orange" /></Tooltip> </Text>
      }
      const allocatedTeams = allTeams.filter(team => eksCredentials.allocation.spec.teams.includes(team.metadata.name)).map(team => team.spec.summary)
      return allocatedTeams.length > 0 ? allocatedTeams.join(', ') : 'All teams'
    }

    return (
      <List.Item key={eksCredentials.metadata.name} actions={[
        <ResourceVerificationStatus key="verification_status" resourceStatus={eksCredentials.status} />,
        <Text key="show_creds"><a onClick={editEKSCredentials(eksCredentials)}><Icon type="edit" theme="filled"/> Edit</a></Text>
      ]}>
        <List.Item.Meta
          avatar={<Avatar icon="amazon" />}
          title={
            <>
              <Text style={{ display: 'inline', marginRight: '15px', fontSize: '20px', fontWeight: '600' }}>{eksCredentials.spec.accountID}</Text>
              <Text style={{ marginRight: '5px' }}>{eksCredentials.allocation.spec.name}</Text>
              <Tooltip title={eksCredentials.allocation.spec.summary}>
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
