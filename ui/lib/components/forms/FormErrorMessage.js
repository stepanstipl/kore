import PropTypes from 'prop-types'
import { Alert } from 'antd'

const FormErrorMessage = ({ message }) => {
  if (!message) {
    return null
  }
  return (
    <Alert
      message={message}
      type="error"
      showIcon
      closable
      style={{ marginBottom: '20px' }}
    />
  )
}

FormErrorMessage.propTypes = {
  message: PropTypes.oneOfType([
    PropTypes.bool,
    PropTypes.string,
    PropTypes.element
  ])
}

export default FormErrorMessage
