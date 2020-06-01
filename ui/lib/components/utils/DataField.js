import PropTypes from 'prop-types'
import { Col, Row, Typography } from 'antd'
const { Text } = Typography

const DataField = ({ label, value, labelColSpan, valueColSpan, textProps, style }) => (
  <Row style={{ padding: '5px 0', ...style }}>
    <Col span={labelColSpan || 6}>{label}</Col>
    <Col span={valueColSpan || 18}>
      {typeof value === 'string' ? <Text {...(textProps || {})}>{value}</Text> : <>{value}</>}
    </Col>
  </Row>
)

DataField.propTypes = {
  label: PropTypes.string.isRequired,
  value: PropTypes.oneOfType([
    PropTypes.string,
    PropTypes.node
  ]),
  labelColSpan: PropTypes.number,
  valueColSpan: PropTypes.number,
  textProps: PropTypes.object,
  style: PropTypes.object
}

export default DataField
