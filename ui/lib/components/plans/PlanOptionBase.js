import * as React from 'react'
import PropTypes from 'prop-types'
import { Alert } from 'antd'

export default class PlanOptionBase extends React.Component {
  static propTypes = {
    resourceType: PropTypes.string.isRequired,
    kind: PropTypes.string.isRequired,
    name: PropTypes.string.isRequired,
    plan: PropTypes.object.isRequired,
    property: PropTypes.object.isRequired,
    value: PropTypes.any,
    editable: PropTypes.bool,
    hideNonEditable: PropTypes.bool,
    onChange: PropTypes.func,
    displayName: PropTypes.string,
    validationErrors: PropTypes.array,
    manage: PropTypes.bool, // manage means we're editing a PLAN, false/unspecified means we're USING a plan e.g. to make/edit a cluster
    mode: PropTypes.oneOf(['create','view','edit']),
    team: PropTypes.object, // may be optionally used by custom plan option components to give richer interface when manage=false.
  }

  describe = (property) => {
    let descriptionPieces = []
    if (property.description) {
      descriptionPieces.push(property.description)
    }
    if (property.examples) {
      descriptionPieces.push(`Examples: ${property.examples.join(', ')}`)
    }
    if (property.format) {
      descriptionPieces.push(`Format: ${property.format}`)
    } else if (property.items && property.items.format) {
      descriptionPieces.push(`Format: ${property.items.format}`)
    }
    if (property.immutable) {
      descriptionPieces.push('Cannot be edited after cluster is created.')
    }
    if (descriptionPieces.length === 0) {
      return null
    }
    return descriptionPieces.join(' | ')
  }

  addComplexItemToArray = (property, values) => {
    if (property.items.type === 'array') {
      values.push([])
      return values
    }
    let newItem = {}
    Object.keys(property.items.properties).forEach((p) => newItem[p] = null)
    values.push(newItem)
    return values
  }

  removeFromArray = (values, indToRemove) => {
    values.splice(indToRemove, 1)
    return values
  }

  validationErrors = (name) => {
    if (!this.props.validationErrors) {
      return null
    }
    const dotName = name.replace(/\[([0-9+])\]/g, '.$1')
    const valErrors = this.props.validationErrors.filter(v => v.field.indexOf(dotName) === 0)
    if (valErrors.length === 0) {
      return null
    }
    return valErrors.map((ve, i) => <Alert key={`${name}.valError.${i}`} type="error" message={ve.message} style={{ marginTop: '10px' }} />)
  }

  hasValidationErrors = (name) => {
    const dotName = name.replace(/\[([0-9+])\]/g, '.$1')
    return this.props.validationErrors && this.props.validationErrors.some(v => v.field.indexOf(dotName) === 0)
  }
  
}
