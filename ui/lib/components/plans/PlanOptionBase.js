import * as React from 'react'
import PropTypes from 'prop-types'
import { Alert } from 'antd'
import { startCase } from 'lodash'

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
    id: PropTypes.string,
    disableCustom: PropTypes.bool, // set to true to force not using a custom control
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
    Object.keys(property.items.properties).forEach((p) => {
      const prop = property.items.properties[p]
      newItem[p] = prop.const !== undefined && prop.const !== null ? prop.const : prop.default
    })
    values.push(newItem)
    return values
  }

  removeFromArray = (values, indToRemove) => {
    values.splice(indToRemove, 1)
    return values
  }

  validationErrors = (name, exactMatch) => {
    const ve = this.getValidationErrors(name, exactMatch)

    if (ve.length === 0) {
      return null
    }

    return ve.map((ve, i) => <Alert key={`${name}.valError.${i}`} type="error" message={ve.message} style={{ marginTop: '10px' }} />)
  }

  hasValidationErrors = (name, exactMatch) => {
    return this.getValidationErrors(name, exactMatch).length > 0
  }

  getValidationErrors = (name, exactMatch) => {
    if (!this.props.validationErrors) {
      return []
    }

    const dotName = name.replace(/\[([0-9+])\]/g, '.$1')
    const valErrors = this.props.validationErrors.filter(v => {
      const f = v.field.replace(/^spec\.configuration\./, '')
      if (exactMatch === true) {
        return f === dotName
      } else {
        return f.indexOf(dotName) === 0
      }
    })
    return valErrors
  }

  prepCommonProps = (props, defaultIfNoDefault = undefined) => {
    const { property, value, name } = props
    const onChange = props.onChange || (() => {})
    const displayName = props.displayName || property.title || startCase(name)
    const help = props.help || this.describe(property)
    const defaultValue = (property.const !== undefined && property.const !== null ? property.const : property.default) || defaultIfNoDefault
    const valueOrDefault = value !== undefined && value !== null ? value : defaultValue
    const id = props.id || `plan_input_${name}`
    return { onChange, displayName, help, defaultValue, valueOrDefault, id }
  }
}
