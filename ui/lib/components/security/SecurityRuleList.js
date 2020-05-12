import React from 'react'
import PropTypes from 'prop-types'
import { List, Avatar, Tooltip, Icon } from 'antd'
import inflect from 'inflect'
import Link from 'next/link'

export default class SecurityRuleList extends React.Component {
  static propTypes = {
    rules: PropTypes.object.isRequired
  }

  appliesTo = (rule) => {
    return rule.spec.appliesTo.map((a) => inflect.pluralize(a.toLowerCase())).join(', ')
  }

  render() {
    const { rules } = this.props
    return (
      <>
        <List>
          {rules.items.map((rule) => (
            <List.Item key={rule.spec.code} actions={[<Link key="view" href="/security/rules/[code]" as={`/security/rules/${rule.spec.code}`}><a><Tooltip placement="left" title="View rule details"><Icon type="info-circle" /> View rule details</Tooltip></a></Link>]}>
              <List.Item.Meta
                title={<Link href="/security/rules/[code]" as={`/security/rules/${rule.spec.code}`}><a><Tooltip title="View rule details">{rule.spec.code}: {rule.spec.name}</Tooltip></a></Link>}
                description={`Applies to ${this.appliesTo(rule)}`}
                avatar={<Avatar icon="schedule" />}
              />
            </List.Item>
          ))}
        </List>
      </>
    )
  }
}