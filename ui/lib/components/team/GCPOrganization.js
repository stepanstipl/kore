import PropTypes from 'prop-types'
import moment from 'moment'
import { List, Avatar, Icon, Typography, message, Tooltip } from 'antd'
const { Text, Title } = Typography

import ResourceStatusTag from '../ResourceStatusTag'
import AutoRefreshComponent from './AutoRefreshComponent'

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

  finalStateReached() {
    const { organization } = this.props
    const { allocation, status } = organization
    if (status.status === 'Success') {
      return message.success(`GCP organization "${allocation.spec.name}" created successfully`)
    }
    if (status.status === 'Failure') {
      return message.error(`GCP organization "${allocation.spec.name}" failed to be created`)
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
        <ResourceStatusTag key="status" resourceStatus={organization.status} />,
        <Text key="edit"><a onClick={editOrganization(organization)}><Icon type="eye" theme="filled"/> Edit</a></Text>
      ]}>
        <List.Item.Meta
          avatar={<Avatar icon="cloud" />}
          title={
            <>
              <Title level={4} style={{ display: 'inline', marginRight: '15px' }}>{organization.spec.parentID}</Title>
              <Text style={{ marginRight: '5px' }}>{organization.allocation.spec.name}</Text>
              <Tooltip title={organization.allocation.spec.summary}>
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
