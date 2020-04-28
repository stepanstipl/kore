import PropTypes from 'prop-types'
import moment from 'moment'
import { List, Avatar, Icon, Typography, message, Tooltip } from 'antd'
const { Text } = Typography

import ResourceVerificationStatus from '../../../components/ResourceVerificationStatus'
import AutoRefreshComponent from '../../../components/team/AutoRefreshComponent'

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
    const { organization, editOrganization } = this.props
    const created = moment(organization.metadata.creationTimestamp).fromNow()

    return (
      <List.Item key={organization.metadata.name} actions={[
        <ResourceVerificationStatus key="verification_status" resourceStatus={organization.status} />,
        <Text key="edit"><a onClick={editOrganization(organization)}><Icon type="edit" theme="filled"/> Edit</a></Text>
      ]}>
        <List.Item.Meta
          avatar={<Avatar icon="cloud" />}
          title={<Text style={{ display: 'inline', marginRight: '15px', fontSize: '20px', fontWeight: '600' }}>{organization.spec.parentID}</Text>}
          description={
            <>
              <Text style={{ marginRight: '5px' }}>{organization.allocation.spec.name}</Text>
              <Tooltip title={organization.allocation.spec.summary}>
                <Icon type="info-circle" theme="twoTone" />
              </Tooltip>
            </>
          }
        />
        <Text type='secondary'>Created {created}</Text>
      </List.Item>
    )
  }

}

export default GCPOrganization
