import React from 'react'
import PropTypes from 'prop-types'
import moment from 'moment'
import { List, Avatar, Icon, Typography } from 'antd'
const { Text } = Typography

class PlanItem extends React.Component {
  static propTypes = {
    plan: PropTypes.object.isRequired,
    viewPlan: PropTypes.func.isRequired
  }

  render() {
    const { plan, viewPlan } = this.props
    const created = moment(plan.metadata.creationTimestamp).fromNow()

    return (
      <List.Item key={plan.metadata.name} actions={[
        <Text key="view_plan"><a onClick={viewPlan(plan)}><Icon type="eye" theme="filled"/> View</a></Text>
      ]}>
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
