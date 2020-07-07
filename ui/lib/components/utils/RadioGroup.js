import PropTypes from 'prop-types'
import React from 'react'
import { Radio, Typography } from 'antd'
const { Paragraph, Text } = Typography

const RadioGroup = ({ heading, onChange, options, value, disabled, style }) => {
  return (
    <>
      {heading ? <Paragraph style={{ fontSize: '16px', fontWeight: '600' }}>{heading}</Paragraph> : null}
      <Radio.Group onChange={onChange} value={value} disabled={disabled} style={{ ...style }}>
        {options.map(option => (
          <div key={option.value} style={{ display: 'inline-block', marginRight: '20px' }}>
            <Radio className={option.className} value={option.value} style={{ float: 'left' }}>
              <Text style={{ fontSize: '16px', fontWeight: '600' }}>{option.title}</Text>
              {option.description ? <Paragraph style={{ marginLeft: '24px', marginBottom: '0' }}>{option.description}</Paragraph> : null}
            </Radio>
            {option.extra}
          </div>
        ))}
      </Radio.Group>
    </>
  )
}

RadioGroup.propTypes = {
  heading: PropTypes.string,
  onChange: PropTypes.func,
  options: PropTypes.array.isRequired,
  value: PropTypes.string.isRequired,
  disabled: PropTypes.bool,
  style: PropTypes.object
}

export default RadioGroup
