import { notification, Icon } from 'antd'

/**
 * Shows a message to indicate the successful completion of an operation.
 * @param {string} message 
 * @param {messageOpts} opts 
 * @returns {string} key which can be used to update message, e.g. warningMessage('warning!!', { key })
 */
export function successMessage(message, opts) {
  return showMessage('success', message, null, opts)
}

/**
 * Shows a message to indicate an operation is ongoing.
 * @param {string} message 
 * @param {messageOpts} opts 
 * @returns {string} key which can be used to update message, e.g. successMessage('completed!!', { key })
 */
export function loadingMessage(message, opts) {
  return showMessage('open', message, <Icon type="loading" />, opts)
}

/**
 * Shows a message to indicate the failure of an operation.
 * @param {string} message 
 * @param {messageOpts} opts 
 * @returns {string} key which can be used to update message, e.g. successMessage('completed!!', { key })
 */
export function errorMessage(message, opts) {
  return showMessage('error', message, null, opts)
}

/**
 * Shows a message to warn the user of something.
 * @param {string} message 
 * @param {messageOpts} opts 
 * @returns {string} key which can be used to update message, e.g. successMessage('completed!!', { key })
 */
export function warningMessage(message, opts) {
  return showMessage('warning', message, null, opts)
}

export class messageOpts {
  /**
   * Duration to show the message in seconds, 0 for indefinite.
   */
  duration
  /**
   * Key to identify message - use the same key on multiple messages
   * to update a message.
   */
  key
  /**
   * Extended description for the message.
   */
  description
}

function showMessage(type, message, icon, opts) {
  if (!opts) {
    opts = {}
  }
  const key = opts.key || Math.random().toString(36).substr(2, 5)
  const duration = opts.duration || 2.5
  notification[type]({
    icon,
    message,
    duration,
    key,
    placement: 'bottomRight',
    description: opts.description
  })
  return key
}
