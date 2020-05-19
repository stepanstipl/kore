import React from 'react'
import PropTypes from 'prop-types'
import { Icon } from 'antd'

export default class SecurityStatusIcon extends React.Component {
  static propTypes = {
    status: PropTypes.string.isRequired,
    inactive: PropTypes.bool,
    text: PropTypes.string,
    size: PropTypes.oneOf(['small', 'default']),
    style: PropTypes.object
  }

  static statusIcons = {
    'Compliant': { 'icon': 'safety-certificate', 'color': 'LimeGreen' },
    'Warning': { 'icon': 'warning', 'color': 'DarkOrange' },
    'Failure': { 'icon': 'stop', 'color': 'Red' }
  }

  render() {
    const { status, inactive, text, size, style } = this.props
    const color = inactive ? 'Silver' : SecurityStatusIcon.statusIcons[status].color
    const fontSize = size === 'small' ? '20px' : '45px'

    if (!text) {
      return (
        <Icon type={SecurityStatusIcon.statusIcons[status].icon} style={{ fontSize, paddingLeft: '10px', paddingRight: '10px', ...style }} theme="twoTone" twoToneColor={color} />
      )
    }
    return (
      <>
        <Icon type={SecurityStatusIcon.statusIcons[status].icon} style={{ marginRight: 8, ...style }} theme="twoTone" twoToneColor={color} />
        <span style={{ color: inactive ? 'rgba(0, 0, 0, 0.65)' : 'black' }}>{text}</span>
      </>
    )
  }
}
