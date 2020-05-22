import PropTypes from 'prop-types'
import * as React from 'react'
import { Select } from 'antd'

export default class ConstrainedDropdown extends React.Component {
  static propTypes = {
    readOnly: PropTypes.bool,
    value: PropTypes.string,
    allowedValues: PropTypes.array.isRequired,
    onChange: PropTypes.func.isRequired,
  }

  render() {
    const { value, readOnly } = this.props
    let { allowedValues } = this.props
    // Support a list of strings or a list of { value: 'x', display: 'y' } objects:
    if (allowedValues.length > 0 && (typeof allowedValues[0])==='string') {
      allowedValues = allowedValues.map((v) => ({ value: v }))
    }

    return (
      <Select value={value} disabled={readOnly} defaultValue={null} onChange={this.props.onChange} style={{ width: '100%' }}>
        {allowedValues.map((allowedValue) => <Select.Option key={allowedValue.value} value={allowedValue.value}>{allowedValue.display ? allowedValue.display : allowedValue.value}</Select.Option>)}
      </Select>
    )
  }
}