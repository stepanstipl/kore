import PropTypes from 'prop-types'
import { Table } from 'antd'

const Plan = ({ plan }) => {
  const planColumns = [{
    title: 'Property',
    dataIndex: 'property',
    key: 'property',
  }, {
    title: 'Value',
    dataIndex: 'value',
    key: 'value',
  }]

  const planValues = plan ?
    Object.keys(plan.spec.configuration).map(key => {
      let value = plan.spec.configuration[key]
      if (key === 'authorizedMasterNetworks') {
        const complexValue = plan.spec.configuration[key]
        value = `${complexValue[0].name}: ${complexValue[0].cidr}`
      }
      return {
        key,
        property: key,
        value: `${value}`
      }
    }) :
    null

  return (
    <Table
      size="small"
      pagination={false}
      columns={planColumns}
      dataSource={planValues}
    />
  )
}

Plan.propTypes = {
  plan: PropTypes.object.isRequired
}

export default Plan