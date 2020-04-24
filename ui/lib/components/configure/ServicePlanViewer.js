import React from 'react'
import PropTypes from 'prop-types'
import KoreApi from '../../kore-api'
import { Table, Icon, Tag, Typography } from 'antd'
const { Text } = Typography
import { startCase } from 'lodash'

class ServicePlanViewer extends React.Component {

  static propTypes = {
    servicePlan: PropTypes.object.isRequired
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
    const schema = await (await KoreApi.client()).GetServicePlanSchema(this.props.servicePlan.spec.kind)
    this.setState({
      schema,
      dataLoading: false
    })
  }

  render() {
    const { spec } = this.props.servicePlan
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

    const propertyDisplayName = (name, description) => (
      <>
        <Text strong>{startCase(name)}</Text> {description ? <><br/><Text type="secondary">{description}</Text></> : null}
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
        return propertyDisplayValue(schema.items, value)
      }
      case 'object': {
        return value.map((v, i) => (
          <Table key={i} size="small" pagination={false} columns={columns} style={{ paddingBottom: i < value.length-1 ? '10px': '' }}
            dataSource={Object.keys(schema.properties).map(p => ({
              key: p,
              property: propertyDisplayName(p),
              value: v[p]
            }))}
          />
        ))
      }
      default: return `${value}`
      }
    }

    const servicePlanValues = Object.keys(schema.properties).map(name => {
      return {
        key: name,
        property: propertyDisplayName(name, schema.properties[name].description),
        value: propertyDisplayValue(schema.properties[name], spec.configuration[name])
      }
    })

    return (
      <Table
        size="small"
        pagination={false}
        columns={columns}
        dataSource={servicePlanValues}
      />
    )

  }
}

export default ServicePlanViewer
