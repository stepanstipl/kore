import PropTypes from 'prop-types'
import moment from 'moment'
import { Avatar, List, Alert, Icon, Drawer, Typography, Button, Tooltip, Modal } from 'antd'
const { Title, Text, Paragraph } = Typography
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

import ResourceList from '../resources/ResourceList'
import KoreApi from '../../kore-api'
import Policy from './Policy'
import PolicyForm from './PolicyForm'
import AllocationHelpers from '../../utils/allocation-helpers'
import { isReadOnlyCRD } from '../../utils/crd-helpers'
import { successMessage, warningMessage } from '../../utils/message'

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
        successMessage(`Policy ${policy.spec.description} deleted`)
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
        <Text key="view_policy"><a id={`policy_view_${policy.metadata.name}`} onClick={this.view(policy)}><Icon type="eye" theme="filled"/> View</a></Text>,
        <Text key="edit_policy">
          <Tooltip title="Edit this policy">
            <a id={`policy_edit_${policy.metadata.name}`} onClick={readonly ? () => warningMessage('Read Only', { description: 'This policy is read-only. Create a new policy to further restrict or allow changes.' }) : this.edit(policy)} style={{ color: readonly ? 'lightgray' : null }}><Icon type="edit" theme="filled"/> Edit</a>
          </Tooltip>
        </Text>,
        <Text key="delete_policy">
          <Tooltip title="Delete this policy">
            <a id={`policy_delete_${policy.metadata.name}`} onClick={readonly ? () => warningMessage('Read Only', { description: 'This policy is read-only and cannot be deleted. Create a new policy to further restrict or allow changes.' }) : this.delete(policy)} style={{ color: readonly ? 'lightgray' : null }}><Icon type="delete" theme="filled"/> Delete</a>
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
    const drawerVisible = Boolean(view || edit || add)
    let drawerTitle = null
    let drawerClose = () => {}
    if (view) {
      drawerTitle = <><Title level={4}>{view.spec.summary}</Title><Text>{view.spec.description}</Text></>
      drawerClose = this.view(false)
    } else if (edit) {
      drawerTitle = <><Title level={4}>{edit.spec.summary}</Title><Text>{edit.spec.description}</Text></>
      drawerClose = this.edit(false)
    } else if (add) {
      drawerTitle = <Title level={4}>New {this.props.kind} policy</Title>
      drawerClose = this.add(false)
    }

    return (
      <>
        <Alert
          message="Manage plan policies"
          description="These policies define what aspects of a plan can be edited by teams when they are creating or updating clusters. This allows you to control what aspects of clusters teams can manage for themselves."
          type="info"
          showIcon
          style={{ marginBottom: '20px' }}
        />
        <Button id="add" type="primary" onClick={this.add(true)} style={{ display: 'block', marginBottom: '20px' }}>+ New</Button>

        <Drawer
          title={drawerTitle}
          visible={drawerVisible}
          onClose={drawerClose}
          width={900}>
          {!view ? null : 
            <>
              <Paragraph style={{ textAlign: 'center' }}>{this.describeAllocation(view.allocation)}</Paragraph>
              <Policy policy={view} mode="view" />
            </>
          }
          {!edit ? null :
            <PolicyForm
              kind={this.props.kind}
              policy={edit}
              handleSubmit={this.handleEditSave}
            />
          }
          {!add ? null : 
            <PolicyForm
              kind={this.props.kind}
              policy={{ spec: { kind: this.props.kind, properties: [] } }}
              handleSubmit={this.handleAddSave}
            />
          }
        </Drawer>

        {!resources ? <Icon type="loading" /> : (
          <List
            id="policy_list"
            dataSource={resources.items}
            renderItem={policy => this.policyItem(policy)}
          />
        )}
      </>
    )
  }
}

export default PolicyList
