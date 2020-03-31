import PropTypes from 'prop-types'
import moment from 'moment'
import { List, Avatar, Icon, Typography, message, Tooltip } from 'antd'
const { Text, Title } = Typography

import ResourceVerificationStatus from '../ResourceVerificationStatus'
import AutoRefreshComponent from './AutoRefreshComponent'

class Credentials extends AutoRefreshComponent {
  static propTypes = {
    gke: PropTypes.object.isRequired,
    allTeams: PropTypes.array.isRequired,
    editGKECredential: PropTypes.func.isRequired
  }

  componentDidUpdate(prevProps) {
    const prevStatus = prevProps.gke.status && prevProps.gke.status.status
    const newStatus = this.props.gke.status && this.props.gke.status.status
    if (prevStatus !== newStatus) {
      this.startRefreshing()
    }
  }

  finalStateReached() {
    const { gke } = this.props
    const { status } = gke
    if (status.status === 'Success') {
      return message.success(`GCP Service Account for project "${gke.spec.project}" verified successfully`)
    }
    if (status.status === 'Failure') {
      return message.error(`GCP Service Account for project "${gke.spec.project}" could not be verified`)
    }
  }

  render() {
    const { gke, editGKECredential, allTeams } = this.props
    const created = moment(gke.metadata.creationTimestamp).fromNow()

    const displayAllocations = () => {
      if (!gke.allocation) {
        return <Text>No teams <Tooltip title="This project is not allocated to any teams, click edit to fix this."><Icon type="warning" theme="twoTone" twoToneColor="orange" /></Tooltip> </Text>
      }
      const allocatedTeams = allTeams.filter(team => gke.allocation.spec.teams.includes(team.metadata.name)).map(team => team.spec.summary)
      return allocatedTeams.length > 0 ? allocatedTeams.join(', ') : 'All teams'
    }

    return (
      <List.Item key={gke.metadata.name} actions={[
        <ResourceVerificationStatus key="verification_status" resourceStatus={gke.status} />,
        <Text key="show_creds"><a onClick={editGKECredential(gke)}><Icon type="eye" theme="filled"/> Edit</a></Text>
      ]}>
        <List.Item.Meta
          avatar={<Avatar icon="project" />}
          title={
            <>
              <Title level={4} style={{ display: 'inline', marginRight: '15px' }}>{gke.spec.project}</Title>
              <Text style={{ marginRight: '5px' }}>{gke.allocation.spec.name}</Text>
              <Tooltip title={gke.allocation.spec.summary}>
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

export default Credentials
