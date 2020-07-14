import React from 'react'
import PropTypes from 'prop-types'
import moment from 'moment'
import { List, Avatar, Icon, Typography, Tooltip } from 'antd'
const { Text } = Typography
import { pluralize, titleize } from 'inflect'

import IconTooltip from '../utils/IconTooltip'
import { isReadOnlyCRD } from '../../utils/crd-helpers'
import { warningMessage } from '../../utils/message'
import { getPlanCloudInfo } from '../../utils/plans'

class PlanItem extends React.Component {
  static propTypes = {
    plan: PropTypes.object.isRequired,
    viewPlan: PropTypes.func.isRequired,
    editPlan: PropTypes.func.isRequired,
    deletePlan: PropTypes.func.isRequired,
    displayUnassociatedPlanWarning: PropTypes.bool.isRequired
  }

  planCloudInfo = getPlanCloudInfo(this.props.plan.spec.kind)

  actions = () => {
    const readonly = isReadOnlyCRD(this.props.plan)
    const actions = []
    if (this.props.displayUnassociatedPlanWarning) {
      actions.push(<IconTooltip key="warning" icon="warning" color="orange" text={`This plan not associated with any ${this.planCloudInfo.cloud} automated ${pluralize(this.planCloudInfo.accountNoun)} and will not be available for teams to use. Edit this plan or go to ${titleize(this.planCloudInfo.accountNoun)} automation settings to review this.`}/>)
    }
    actions.push(<Text key="view_plan"><a id={`plans_view_${this.props.plan.metadata.name}`} onClick={this.props.viewPlan(this.props.plan)}><Icon type="eye" theme="filled"/> View</a></Text>)
    actions.push(
      <Text key="edit_plan">
        <Tooltip title="Edit this plan">
          <a id={`plans_edit_${this.props.plan.metadata.name}`} onClick={readonly ? () => warningMessage('Read Only', { description: 'This plan is read-only. Create a new plan if this built-in plan does not meet your needs.' }) : this.props.editPlan(this.props.plan)} style={{ color: readonly ? 'lightgray' : null }}><Icon type="edit" theme="filled"/> Edit</a>
        </Tooltip>
      </Text>
    )
    actions.push(
      <Text key="delete_plan">
        <Tooltip title="Delete this plan">
          <a id={`plans_delete_${this.props.plan.metadata.name}`} onClick={readonly ? () => warningMessage('Read Only', { description: 'This plan is read-only and cannot be deleted. To prevent teams using this plan, remove the allocation.' }) : this.props.deletePlan(this.props.plan)} style={{ color: readonly ? 'lightgray' : null }}><Icon type="delete" theme="filled"/> Delete</a>
        </Tooltip>
      </Text>
    )
    return actions
  }

  render() {
    const plan = this.props.plan
    const created = moment(plan.metadata.creationTimestamp).fromNow()

    return (
      <List.Item key={plan.metadata.name} actions={this.actions()}>
        <List.Item.Meta
          avatar={<Avatar icon="build" />}
          title={plan.spec.description}
          description={plan.spec.summary}
        />
        <Text type='secondary'>Created {created}</Text>
      </List.Item>
    )
  }

}

export default PlanItem
