import * as React from 'react'
import PropTypes from 'prop-types'
import { Form, Input, Icon, Checkbox, InputNumber, Select, Button, Card, Alert } from 'antd'
import { startCase } from 'lodash'

class PlanOption extends React.Component {
  static propTypes = {
    name: PropTypes.string.isRequired,
    property: PropTypes.object.isRequired,
    value: PropTypes.any,
    editable: PropTypes.bool,
    hideNonEditable: PropTypes.bool,
    onChange: PropTypes.func,
    displayName: PropTypes.string,
    validationErrors: PropTypes.array
  }

  describe = (property) => {
    let description = ''
    if (property.description) {
      description += property.description
    }
    if (property.format) {
      description += ` Format: ${property.format}`
    } else if (property.items && property.items.format) {
      description += ` Format: ${property.items.format}`
    }
    if (description.length === 0) {
      return null
    }
    return description
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
    const dotName = name.replace(/\[([0-9+])\]/g,'.$1')
    const valErrors = this.props.validationErrors.filter((v) => v.field.indexOf(dotName)===0)
    if (valErrors.length === 0) {
      return null
    }
    return (
      <>
        {valErrors.map((ve, i) => 
          <Alert key={`${name}.valError.${i}`} type="error" message={ve.message} style={{ marginTop: '10px' }} />
        )}
      </>
    )
  }

  render() {
    const { name, property, value, editable, hideNonEditable, onChange } = this.props
    if (!editable && hideNonEditable) {
      return null
    }

    const displayName = this.props.displayName || name
    const locked = !editable ? <Icon type="lock" /> : null

    // Special handling for object types - represent as a card with a plan option for each property:
    if (property.type === 'object') {
      const keys = property.properties ? Object.keys(property.properties) : []
      return (
        <Card size="small" title={startCase(displayName)}>
          {keys.map((key) =>
            <PlanOption 
              key={`${name}.${key}`} 
              name={`${name}.${key}`} 
              displayName={key} 
              property={property.properties[key]} 
              value={value[key]} 
              editable={editable} 
              onChange={onChange} 
              validationErrors={this.props.validationErrors} />
          )}
        </Card>
      )
    }

    // Handle all other types:
    return (
      <Form.Item label={startCase(displayName)} help={this.describe(property)}>
        {(() => {
          switch(property.type) {
          case 'string': {
            return <Input value={value} readOnly={!editable} disabled={!editable} addonAfter={locked} onChange={(e) => onChange(name, e.target.value)} />
          }
          case 'boolean': {
            return <Checkbox checked={value} readOnly={!editable} disabled={!editable} onChange={(v) => onChange(name, v)} />
          }
          case 'number': {
            return <InputNumber value={value} readOnly={!editable} disabled={!editable} onChange={(v) => onChange(name, v)} />
          }
          case 'array': {
            const values = value ? value : []
            if (property.items.type !== 'array' && property.items.type !== 'object') {
              return <Select mode="tags" tokenSeparators={[',']} value={values} readOnly={!editable} disabled={!editable} onChange={(v) => onChange(name, v)} />
            } else {
              return (
                <>
                  {values.map((val, ind) => 
                    <React.Fragment key={`${name}[${ind}]`}>
                      <PlanOption 
                        name={`${name}[${ind}]`} 
                        property={property.items} 
                        value={val} 
                        editable={editable} 
                        onChange={onChange} 
                        validationErrors={this.props.validationErrors} />
                      <Button disabled={!editable} icon="delete" title={`Remove ${startCase(displayName)} ${ind}`} onClick={() => onChange(name, this.removeFromArray(values, ind))}>
                        {`Remove ${startCase(displayName)} ${ind}`}                    
                      </Button>
                    </React.Fragment>
                  )}
                  {(values.length === 0) ?
                    <Alert type="info" message={`No ${startCase(displayName)} currently defined.`}/>
                    : null}
                  <Button disabled={!editable} icon="plus" title={`Add new ${startCase(displayName)}`} onClick={() => onChange(name, this.addComplexItemToArray(property, values))}>
                    {`Add new ${startCase(displayName)}`}
                  </Button>
                </>
              )
            }
          }
          default: {
            return <Alert type="warning" message={`The property ${displayName} is of an unknown type ${property.type} and cannot be specified through the UI.`}/>
          }
          }
        })()}
        {this.validationErrors(name)}
      </Form.Item>
    )
  }
}

export default PlanOption
