import React from 'react'
import PropTypes from 'prop-types'
import KoreApi from '../../kore-api'
import { Table, Icon, Typography, Tooltip } from 'antd'
const { Text } = Typography
import { startCase } from 'lodash'

export default class Policy extends React.Component {
  static propTypes = {
    policy: PropTypes.object.isRequired,
    mode: PropTypes.oneOf(['view','edit']),
    onPolicyUpdate: PropTypes.func
  }

  state = {
    dataLoading: true,
    schema: null
  }

  componentDidMountComplete = null
  componentDidMount() {
    this.componentDidMountComplete = this.fetchComponentData()
  }

  async fetchComponentData() {
    const schema = await (await KoreApi.client()).GetPlanSchema(this.props.policy.spec.kind)
    this.setState({
      schema,
      dataLoading: false
    })
  }

  componentDidUpdate(prevProps) {
    // If we've changed kind, reload the plan schema.
    if (this.props.policy.spec.kind !== prevProps.policy.spec.kind) {
      this.fetchComponentData()
    }
  }

  handleClick = (name, rule, newValue) => {
    if (!this.props.onPolicyUpdate || this.props.mode === 'view') {
      return
    }
    let newProperties = [...this.props.policy.spec.properties]
    const ind = newProperties.findIndex((p) => p.name === name)
    if (ind === -1) {
      newProperties.push({
        name: name,
        [rule]: newValue
      })
    } else {
      newProperties[ind] = {
        ...newProperties[ind],
        [rule]: newValue
      }
      // Remove rule if no longer used.
      if (!newProperties[ind].disallowUpdate && !newProperties[ind].allowUpdate) {
        newProperties.splice(ind, 1)
      }
    }
    this.props.onPolicyUpdate({
      ...this.props.policy,
      spec: {
        ...this.props.policy.spec,
        properties: newProperties
      }
    })
  }

  propertyDisplayName = (name, description) => (
    <>
      <Text strong>{startCase(name)}</Text> {description ? <><br/><Text type="secondary">{description}</Text></> : null}
    </>
  )

  propertyAllowUpdate = (name) => {
    const policyProperty = this.props.policy.spec.properties.find((p) => p.name === name)
    if (policyProperty && policyProperty.allowUpdate) {
      return <Icon id={`policy_${name}_allow`} style={{ fontSize: '1.5em' }} type="check-circle" theme="twoTone" twoToneColor="#52c41a" onClick={() => this.handleClick(name, 'allowUpdate', false)}/>
    } else {
      return <Icon id={`policy_${name}_allow`} style={{ fontSize: '1.5em' }} type="question-circle" theme="twoTone" twoToneColor="lightgray" onClick={() => this.handleClick(name, 'allowUpdate', true)} />      
    }
  }

  propertyDisallowUpdate = (name) => {
    const policyProperty = this.props.policy.spec.properties.find((p) => p.name === name)
    if (policyProperty && policyProperty.disallowUpdate) {
      return <Icon id={`policy_${name}_disallow`} style={{ fontSize: '1.5em' }} type="close-circle" theme="twoTone" twoToneColor="red" onClick={() => this.handleClick(name, 'disallowUpdate', false)} />
    } else {
      return <Icon id={`policy_${name}_disallow`} style={{ fontSize: '1.5em' }} type="question-circle" theme="twoTone" twoToneColor="lightgray" onClick={() => this.handleClick(name, 'disallowUpdate', true)}/>      
    }
  }

  propertyResult = (name) => {
    const policyProperty = this.props.policy.spec.properties.find((p) => p.name === name)
    const defaultAllow = false
    if (policyProperty) {
      if (policyProperty.disallowUpdate) {
        return <Tooltip id={`policy_${name}_result_tt`} placement="left" title="Changes explicitly denied by this policy, this cannot be changed by another policy"><Icon id={`policy_${name}_result`} style={{ fontSize: '1.5em' }} type="close-circle" theme="twoTone" twoToneColor="red" /></Tooltip>
      }
      if (policyProperty.allowUpdate) {
        return <Tooltip id={`policy_${name}_result_tt`} placement="left" title="Changes explicitly allowed by this policy, but another policy could still disallow edits"><Icon id={`policy_${name}_result`} style={{ fontSize: '1.5em' }} type="check-circle" theme="twoTone" twoToneColor="#52c41a" /></Tooltip>
      }
    }
    if (defaultAllow) {
      return <Tooltip id={`policy_${name}_result_tt`} placement="left" title="Changes will be allowed by default, but another policy could disallow edits"><Icon id={`policy_${name}_result`} style={{ fontSize: '1.5em' }} type="check-square" theme="twoTone" twoToneColor="silver" /></Tooltip>
    }
    return <Tooltip id={`policy_${name}_result_tt`} placement="left" title="Changes will be denied by default, but another policy could allow edits"><Icon id={`policy_${name}_result`} style={{ fontSize: '1.5em' }} type="close-square" theme="twoTone" twoToneColor="silver" /></Tooltip>
  }

  render() {
    const { schema, dataLoading } = this.state

    if (dataLoading) {
      return null
    }

    const columns = [{
      title: 'Property',
      dataIndex: 'property',
      key: 'property',
    }, {
      title: 'Allow Update',
      dataIndex: 'allowUpdate',
      key: 'allowUpdate',
    }, {
      title: 'Disallow Update',
      dataIndex: 'disallowUpdate',
      key: 'disallowUpdate',
    }, {
      title: 'Result',
      dataIndex: 'result',
      key: 'result',
    }]

    const planValues = Object.keys(schema.properties).map(name => {
      return {
        key: name,
        property: this.propertyDisplayName(name, schema.properties[name].description),
        allowUpdate: this.propertyAllowUpdate(name),
        disallowUpdate: this.propertyDisallowUpdate(name),
        result: this.propertyResult(name)
      }
    })

    return (
      <Table
        size="small"
        pagination={false}
        columns={columns}
        dataSource={planValues}
      />
    )
  }
}
