import React from 'react'
import PropTypes from 'prop-types'
import { Alert, Table, Icon, Tag, Tooltip, Typography } from 'antd'
const { Paragraph, Text } = Typography
import { startCase } from 'lodash'

import KoreApi from '../../kore-api'
import { featureEnabled, KoreFeatures } from '../../utils/features'

class PlanViewer extends React.Component {

  static propTypes = {
    plan: PropTypes.object.isRequired,
    resourceType: PropTypes.oneOf(['cluster', 'service']).isRequired,
    displayUnassociatedPlanWarning: PropTypes.bool
  }

  state = {
    dataLoading: true,
    schema: null,
    costEstimate: null
  }

  componentDidMountComplete = null
  componentDidMount() {
    this.componentDidMountComplete = this.fetchComponentData()
  }

  async fetchComponentData() {
    let schema, costEstimate
    const api = await KoreApi.client()
    switch (this.props.resourceType) {
    case 'cluster':
      schema = await api.GetPlanSchema(this.props.plan.spec.kind)
      if (featureEnabled(KoreFeatures.COSTS)) {
        costEstimate = await api.costestimates.EstimateClusterPlanCost(this.props.plan)
      }
      break
    case 'service':
      schema = (await (await KoreApi.client()).GetServicePlanDetails(this.props.plan.metadata.name)).schema
      break
    }
    this.setState({
      schema: schema || { properties:[] },
      costEstimate,
      dataLoading: false
    })
  }

  render() {
    const { plan, displayUnassociatedPlanWarning } = this.props
    const { schema, dataLoading } = this.state

    if (dataLoading) {
      return null
    }

    const columns = [{
      title: 'Property',
      dataIndex: 'property',
      key: 'property',
    }, {
      title: 'Value',
      dataIndex: 'value',
      key: 'value',
    }]

    const propertyDisplayName = (name, description, deprecated) => (
      <>
        <Text strong>{!deprecated ? null : <Icon type="warning" twoToneColor="orange" theme="twoTone" />} {startCase(name)}</Text> {description ? <><br/><Text type="secondary">{description}</Text></> : null}
      </>
    )

    const propertyDisplayValue = (schema, value) => {
      if (!value && schema.type !== 'boolean') {
        return ''
      }
      switch (schema.type) {
      case 'string': return value
      case 'boolean': return value ? <Icon type="check-circle" theme="twoTone" twoToneColor="#52c41a" /> : <Icon type="close-circle" theme="twoTone" twoToneColor="red" />
      case 'array': {
        if (schema.items.type !== 'array' && schema.items.type !== 'object') {
          return value.map((v, i) => <Tag key={i}>{v}</Tag>)
        }
        return value.map((v) => propertyDisplayValue(schema.items, v))
      }
      case 'object': {
        if (schema.properties) {
          return <Table size="small" pagination={false} columns={columns} style={{ paddingTop: '5px', paddingBottom: '5px' }}
            dataSource={Object.keys(schema.properties).map(p => ({
              key: p,
              property: propertyDisplayName(p),
              value: propertyDisplayValue(schema.properties[p], value[p])
            }))}
          />
        }
        if (schema.additionalProperties && schema.additionalProperties.type === 'string') {
          const keys = value ? Object.keys(value) : []
          return <Table size="small" pagination={false} columns={columns} style={{ paddingTop: '5px', paddingBottom: '5px' }}
            dataSource={keys.map(p => ({
              key: p,
              property: propertyDisplayName(p),
              value: value[p]
            }))}
          />
        }
        return `${value}`
      }
      default: return `${value}`
      }
    }

    let hasDeprecated = false
    let planValues = Object.keys(plan.spec.configuration).map(name => {
      if (schema.properties[name].deprecated && plan.spec.configuration[name] !== undefined) {
        hasDeprecated = true
      }
      return {
        key: name,
        property: propertyDisplayName(name, schema.properties[name].description, schema.properties[name].deprecated),
        value: propertyDisplayValue(schema.properties[name], plan.spec.configuration[name]),
        deprecated: schema.properties[name].deprecated
      }
    })

    if (!hasDeprecated) {
      planValues = planValues.filter((v) => !v.deprecated)
    }

    return (
      <>
        {plan.gcpAutomatedProject && (
          <Paragraph>GCP project automation: <Tooltip overlay="When using Kore managed GCP projects, clusters using this plan will provisioned inside this project type."><Tag style={{ marginLeft: '10px' }}>{plan.gcpAutomatedProject.name}</Tag></Tooltip></Paragraph>
        )}
        {displayUnassociatedPlanWarning && (
          <Alert
            message="This plan not associated with any GCP automated projects and will not be available for teams to use. Edit this plan or go to Project automation settings to review this."
            type="warning"
            showIcon
            style={{ marginBottom: '20px' }}
          />
        )}
        {!hasDeprecated ? null : (
          <Alert
            message="This plan has values set on deprecated fields"
            type="warning"
            showIcon
            style={{ marginBottom: '20px' }}
          />
        )}
        <Table
          size="small"
          pagination={false}
          columns={columns}
          dataSource={planValues}
        />
      </>
    )

  }
}

export default PlanViewer
