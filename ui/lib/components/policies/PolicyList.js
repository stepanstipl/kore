import PropTypes from 'prop-types'
import moment from 'moment'
import { message, Avatar, List, Alert, Icon, Drawer, Typography, Button, Tooltip, Modal } from 'antd'
const { Title, Text, Paragraph } = Typography
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

import ResourceList from '../resources/ResourceList'
import KoreApi from '../../kore-api'
import Policy from './Policy'
import PolicyForm from './PolicyForm'
import AllocationHelpers from '../../utils/allocation-helpers'
import { isReadOnlyCRD } from '../../utils/crd-helpers'

class PolicyList extends ResourceList {
  static propTypes = {
    kind: PropTypes.string,
    style: PropTypes.object
  }

  async fetchComponentData() {
    const api = await KoreApi.client()
    const [ policyList, allAllocations ] = await Promise.all([
      api.ListPlanPolicies(this.props.kind),
      api.ListAllocations(publicRuntimeConfig.koreAdminTeamName)
    ])
    policyList.items.forEach((p) => {
      p.allocation = AllocationHelpers.findAllocationForResource(allAllocations, p)
    })
    return { resources: policyList }
  }

  delete = (policy) => () => {
    Modal.confirm({
      title: `Are you sure you want to delete the policy ${policy.spec.description}?`,
      content: 'This cannot be undone',
      okText: 'Yes',
      okType: 'danger',
      cancelText: 'No',
      onOk: async () => {
        await AllocationHelpers.removeAllocation(policy)
        await (await KoreApi.client()).RemovePlanPolicy(policy.metadata.name)
        message.success(`Policy ${policy.spec.description} deleted`)
        await this.refresh()
      }
    })
  }

  updatedMessage = 'Policy saved successfully'
  createdMessage = 'Policy created successfully'

  describeAllocation = (allocation) => {
    if (!allocation) {
      return <><span style={{ fontWeight: 'bold' }}>No teams allocated</span> <Tooltip title="This policy is not allocated to any teams, edit the policy to fix this."><Icon type="warning" theme="twoTone" twoToneColor="orange" /></Tooltip></>
    }
    if (AllocationHelpers.isAllTeams(allocation)) {
      return <>This policy applies to <span style={{ fontWeight: 'bold' }}>all teams</span></>
    }
    return <>This policy applies to: <span style={{ fontWeight: 'bold' }}>{allocation.spec.teams.join(', ')}</span></>
  }

  policyItem = (policy) => {
    const created = moment(policy.metadata.creationTimestamp).fromNow()
    const readonly = isReadOnlyCRD(policy)
    return (
      <List.Item key={policy.metadata.name} actions={[
        <Text key="view_policy"><a onClick={this.view(policy)}><Icon type="eye" theme="filled"/> View</a></Text>,
        <Text key="edit_policy">
          <Tooltip title={readonly ? 'Read-only' : 'Edit this policy'}>
            <a onClick={readonly ? () => {} : this.edit(policy)} style={{ color: readonly ? 'lightgray' : null }}><Icon type="edit" theme="filled"/> Edit</a>
          </Tooltip>
        </Text>,
        <Text key="delete_policy">
          <Tooltip title={readonly ? 'Read-only' : 'Delete this policy'}>
            <a onClick={readonly ? () => {} : this.delete(policy)} style={{ color: readonly ? 'lightgray' : null }}><Icon type="delete" theme="filled"/> Delete</a>
          </Tooltip>
        </Text>
      ]}>
        <List.Item.Meta
          avatar={<Avatar icon="lock" />}
          title={policy.spec.summary}
          description={(<Text>{policy.spec.description}. {this.describeAllocation(policy.allocation)}</Text>)}
        />
        <Text type='secondary'>Created {created}</Text>
      </List.Item>
    )
  }

  render() {
    const { resources, view, edit, add } = this.state

    return (
      <>
        <Alert
          message="Manage plan policies"
          description="These policies define what aspects of a plan can be edited by teams when they are creating or updating clusters. This allows you to control what aspects of clusters teams can manage for themselves."
          type="info"
          showIcon
          style={{ marginBottom: '20px' }}
        />
        <Button type="primary" onClick={this.add(true)} style={{ display: 'block', marginBottom: '20px' }}>+ New</Button>
        {!resources ? <Icon type="loading" /> : (
          <>
            <List
              dataSource={resources.items}
              renderItem={policy => this.policyItem(policy)}
            >
            </List>

            {view ? (
              <Drawer
                title={<><Title level={4}>{view.spec.summary}</Title><Text>{view.spec.description}</Text></>}
                visible={Boolean(view)}
                onClose={this.view(false)}
                width={900}
              >
                <Paragraph style={{ textAlign: 'center' }}>{this.describeAllocation(view.allocation)}</Paragraph>
                <Policy policy={view} mode="view" />
              </Drawer>
            ) : null}

            {edit ? (
              <Drawer
                title={<><Title level={4}>{edit.spec.summary}</Title><Text>{edit.spec.description}</Text></>}
                visible={Boolean(edit)}
                onClose={this.edit(false)}
                width={900}
              >
                <PolicyForm
                  kind={this.props.kind}
                  policy={edit}
                  handleSubmit={this.handleEditSave}
                />
              </Drawer>
            ) : null}

            {add ? (
              <Drawer
                title={<Title level={4}>New {this.props.kind} policy</Title>}
                visible={add}
                onClose={this.add(false)}
                width={900}
              >
                <PolicyForm
                  kind={this.props.kind}
                  policy={{ spec: { kind: this.props.kind, properties: [] } }}
                  handleSubmit={this.handleAddSave}
                />
              </Drawer>
            ) : null}            
          </>
        )}
      </>
    )
  }
}

export default PolicyList
