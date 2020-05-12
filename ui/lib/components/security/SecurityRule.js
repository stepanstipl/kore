import React from 'react'
import PropTypes from 'prop-types'
import { Descriptions } from 'antd'
import Markdown from 'react-markdown'
import inflect from 'inflect'

export default class SecurityRule extends React.Component {
  static propTypes = {
    rule: PropTypes.object.isRequired
  }

  appliesTo = (rule) => {
    return rule.spec.appliesTo.map((a) => inflect.pluralize(a)).join(', ')
  }

  render() {
    const { rule } = this.props
    return (
      <>
        <Descriptions>
          <Descriptions.Item label="Code">{rule.spec.code}</Descriptions.Item>
        </Descriptions>
        <Descriptions>
          <Descriptions.Item label="Name">{rule.spec.name}</Descriptions.Item>
        </Descriptions>
        <Descriptions>
          <Descriptions.Item label="Applies to">{this.appliesTo(rule)}</Descriptions.Item>
        </Descriptions>
        <Markdown source={rule.spec.description} />
      </>
    )
  }
}