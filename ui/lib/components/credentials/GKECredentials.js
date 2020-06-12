import PropTypes from 'prop-types'
import moment from 'moment'
import { List, Avatar, Icon, Typography, Tooltip } from 'antd'
const { Text } = Typography

import ResourceVerificationStatus from '../resources/ResourceVerificationStatus'
import AutoRefreshComponent from '../teams/AutoRefreshComponent'
import { successMessage, errorMessage } from '../../utils/message'

class GKECredentials extends AutoRefreshComponent {
  static propTypes = {
    gkeCredentials: PropTypes.object.isRequired,
    allTeams: PropTypes.array.isRequired,
    editGKECredential: PropTypes.func.isRequired
  }

  componentDidUpdate(prevProps) {
    const prevStatus = prevProps.gkeCredentials.status && prevProps.gkeCredentials.status.status
    const newStatus = this.props.gkeCredentials.status && this.props.gkeCredentials.status.status
    if (prevStatus !== newStatus) {
      this.startRefreshing()
    }
  }

  finalStateReached({ state }) {
    const { gkeCredentials } = this.props
    if (state === AutoRefreshComponent.FINAL_STATES.SUCCESS) {
      return successMessage(`GCP credentials for project "${gkeCredentials.spec.project}" verified successfully`)
    }
    if (state === AutoRefreshComponent.FINAL_STATES.FAILURE) {
      return errorMessage(`GCP credentials for project "${gkeCredentials.spec.project}" could not be verified`)
    }
  }

  render() {
    const { gkeCredentials, editGKECredential, deleteGKECredential, allTeams } = this.props
    const created = moment(gkeCredentials.metadata.creationTimestamp).fromNow()

    const displayAllocations = () => {
      if (!gkeCredentials.allocation) {
        return <Text>No teams <Tooltip title="This project is not allocated to any teams, click edit to fix this."><Icon type="warning" theme="twoTone" twoToneColor="orange" /></Tooltip> </Text>
      }
      const allocatedTeams = allTeams.filter(team => gkeCredentials.allocation.spec.teams.includes(team.metadata.name)).map(team => team.spec.summary)
      return allocatedTeams.length > 0 ? allocatedTeams.join(', ') : 'All teams'
    }

    return (
      <List.Item id={`gkecreds_${gkeCredentials.metadata.name}`} key={gkeCredentials.metadata.name} actions={[
        <ResourceVerificationStatus key="verification_status" resourceStatus={gkeCredentials.status} />,
        <Text key="delete_creds"><a id={`gkecreds_del_${gkeCredentials.metadata.name}`} onClick={deleteGKECredential(gkeCredentials)}><Icon type="delete" theme="filled"/> Delete</a></Text>,
        <Text key="show_creds"><a id={`gkecreds_edit_${gkeCredentials.metadata.name}`} onClick={editGKECredential(gkeCredentials)}><Icon type="edit" theme="filled"/> Edit</a></Text>
      ]}>
        <List.Item.Meta
          avatar={<Avatar icon="project" />}
          title={
            <>
              <Text style={{ display: 'inline', marginRight: '15px', fontSize: '20px', fontWeight: '600' }}>{gkeCredentials.spec.project}</Text>
              {gkeCredentials.allocation ? (
                <>
                  <Text style={{ marginRight: '5px' }}>{gkeCredentials.allocation.spec.name}</Text>
                  <Tooltip title={gkeCredentials.allocation.spec.summary}>
                    <Icon type="info-circle" theme="twoTone" />
                  </Tooltip>
                </>
              ) : (
                <Text style={{ marginRight: '5px' }}>Not allocated</Text>
              )}
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

export default GKECredentials
