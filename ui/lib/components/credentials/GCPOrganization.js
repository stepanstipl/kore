import PropTypes from 'prop-types'
import moment from 'moment'
import { List, Avatar, Icon, Typography, Tooltip } from 'antd'
const { Text } = Typography

import ResourceVerificationStatus from '../resources/ResourceVerificationStatus'
import AutoRefreshComponent from '../teams/AutoRefreshComponent'
import { successMessage, errorMessage } from '../../utils/message'

class GCPOrganization extends AutoRefreshComponent {
  static propTypes = {
    organization: PropTypes.object.isRequired,
    allTeams: PropTypes.array.isRequired,
    editOrganization: PropTypes.func.isRequired
  }

  componentDidUpdate(prevProps) {
    const prevStatus = prevProps.organization.status && prevProps.organization.status.status
    const newStatus = this.props.organization.status && this.props.organization.status.status
    if (prevStatus !== newStatus) {
      this.startRefreshing()
    }
  }

  finalStateReached({ state, deleted }) {
    const { organization } = this.props
    if (deleted) {
      return successMessage(`GCP organization "${organization.allocation.spec.name}" deleted`)
    }
    if (state === AutoRefreshComponent.FINAL_STATES.SUCCESS) {
      return successMessage(`GCP organization "${organization.allocation.spec.name}" created successfully`)
    }
    if (state === AutoRefreshComponent.FINAL_STATES.FAILURE) {
      return errorMessage(`GCP organization "${organization.allocation.spec.name}" failed to be created`)
    }
  }

  render() {
    const { organization, editOrganization, allTeams } = this.props
    const created = moment(organization.metadata.creationTimestamp).fromNow()

    const displayAllocations = () => {
      if (!organization.allocation) {
        return <Text>No teams <Tooltip title="This organization is not allocated to any teams, click edit to fix this."><Icon type="warning" theme="twoTone" twoToneColor="orange" /></Tooltip> </Text>
      }
      const allocatedTeams = allTeams.filter(team => organization.allocation.spec.teams.includes(team.metadata.name)).map(team => team.spec.summary)
      return allocatedTeams.length > 0 ? allocatedTeams.join(', ') : 'All teams'
    }

    return (
      <List.Item key={organization.metadata.name} actions={[
        <ResourceVerificationStatus key="verification_status" resourceStatus={organization.status} />,
        <Text key="edit"><a onClick={editOrganization(organization)}><Icon type="edit" theme="filled"/> Edit</a></Text>
      ]}>
        <List.Item.Meta
          avatar={<Avatar icon="cloud" />}
          title={
            <>
              <Text style={{ display: 'inline', marginRight: '15px', fontSize: '20px', fontWeight: '600' }}>{organization.spec.parentID}</Text>
              <Text style={{ marginRight: '5px' }}>{organization.allocation ? organization.allocation.spec.name : organization.metadata.name}</Text>
              <Tooltip title={organization.allocation ? organization.allocation.spec.summary : organization.spec.summary}>
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

export default GCPOrganization
