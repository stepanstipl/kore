import React from 'react'
import PropTypes from 'prop-types'
import { Badge } from 'antd'

const TextWithCount = ({ title, count, icon }) => (
  <span>
    {title}
    {count !== undefined && count !== -1 && <Badge showZero={true} style={{ marginLeft: '10px', backgroundColor: '#1890ff' }} count={count} />}
    {icon}
  </span>
)

TextWithCount.propTypes = {
  title: PropTypes.string.isRequired,
  count: PropTypes.number,
  icon: PropTypes.node
}

export default TextWithCount
