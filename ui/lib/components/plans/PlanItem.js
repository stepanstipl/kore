import React from 'react'
import PropTypes from 'prop-types'
import moment from 'moment'
import { List, Avatar, Icon, Typography } from 'antd'
const { Text } = Typography

import IconTooltip from '../utils/IconTooltip'

class PlanItem extends React.Component {
  static propTypes = {
    plan: PropTypes.object.isRequired,
    viewPlan: PropTypes.func.isRequired,
    editPlan: PropTypes.func.isRequired,
    displayUnassociatedPlanWarning: PropTypes.bool.isRequired
  }

  actions = () => {
    const actions = []
    if (this.props.displayUnassociatedPlanWarning) {
      actions.push(<IconTooltip key="warning" icon="warning" color="orange" text="This plan not associated with any GCP automated projects and will not be available for teams to use. Edit this plan or go to Project automation settings to review this."/>)
    }
    actions.push(<Text key="view_plan"><a onClick={this.props.viewPlan(this.props.plan)}><Icon type="eye" theme="filled"/> View</a></Text>)
    actions.push(<Text key="edit_plan"><a onClick={this.props.editPlan(this.props.plan)}><Icon type="edit" theme="filled"/> Edit</a></Text>,)
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
