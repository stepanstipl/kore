import * as React from 'react'
import { Form, Input, Checkbox, InputNumber, Select, Button, Card, Alert } from 'antd'
import { startCase } from 'lodash'
import CustomPlanOptionRegistry from './custom'
import PlanOptionBase from './PlanOptionBase'

export default class PlanOption extends PlanOptionBase {
  constructor(props) {
    super(props)
  }

  render() {
    const { resourceType, kind, name, property, value, editable, hideNonEditable } = this.props
    if (!editable && hideNonEditable) {
      return null
    }

    // Switch out to a custom option if we have one
    const customControl = CustomPlanOptionRegistry.getCustomPlanOption(resourceType, kind, name, this.props)
    if (customControl) {
      return customControl
    }

    const onChange = this.props.onChange || (() => {})

    const displayName = this.props.displayName || name

    // Special handling for object types - represent as a card with a plan option for each property:
    if (property.type === 'object') {
      const keys = property.properties ? Object.keys(property.properties) : []
      return (
        <Card size="small" title={startCase(displayName)}>
          {keys.map((key) =>
            <PlanOption
              mode={this.props.mode} 
              team={this.props.team} 
              resourceType={resourceType}
              kind={kind}
              key={`${name}.${key}`} 
              name={`${name}.${key}`} 
              displayName={key} 
              property={property.properties[key]} 
              value={value[key]} 
              editable={editable} 
              onChange={onChange} 
              validationErrors={this.props.validationErrors}
            />
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
            return <Input value={value} readOnly={!editable} onChange={(e) => onChange(name, e.target.value)} />
          }
          case 'boolean': {
            return <Checkbox checked={value} disabled={!editable} onChange={(e) => onChange(name, e.target.checked)} />
          }
          case 'number': {
            return <InputNumber value={value} readOnly={!editable} onChange={(v) => onChange(name, v)} />
          }
          case 'integer': {
            return <InputNumber value={value} readOnly={!editable} onChange={(v) => onChange(name, v)} />
          }
          case 'array': {
            const values = value ? value : []
            if (property.items.type !== 'array' && property.items.type !== 'object') {
              return <Select mode="tags" tokenSeparators={[',']} value={values} disabled={!editable} onChange={(v) => onChange(name, v)} />
            } else {
              return (
                <>
                  {values.map((val, ind) => 
                    <React.Fragment key={`${name}[${ind}]`}>
                      <PlanOption 
                        mode={this.props.mode} 
                        team={this.props.team} 
                        resourceType={resourceType}
                        kind={kind}
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
