import * as React from 'react'
import PropTypes from 'prop-types'
import { Form, Input, InputNumber, Select, Button, Card, Alert, Icon, Switch } from 'antd'
const { TextArea } = Input
import { startCase } from 'lodash'
import CustomPlanOptionRegistry from './custom'
import PlanOptionBase from './PlanOptionBase'
import KeyMap from './custom/KeyMap'
import ConstrainedDropdown from './custom/ConstrainedDropdown'

/**
 * PlanOption represents a single option on a plan form. Use Manage(Service/Cluster)PlanForm or UsePlanForm to manage/use a plan rather than using this directly.
 */
export default class PlanOption extends PlanOptionBase {
  static propTypes = {
    help: PropTypes.string,
  }

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
    const displayName = this.props.displayName || property.title || startCase(name)
    const help = this.props.help || this.describe(property)
    const valueOrDefault = value !== undefined && value !== null ? value : property.default

    // Special handling for object types - represent as a card with a plan option for each property:
    if (property.type === 'object' && property.properties) {
      const keys = Object.keys(property.properties)
      return (
        <Card size="small" title={displayName}>
          {keys.map((key) =>
            <PlanOption
              {...this.props}
              key={`${name}.${key}`}
              name={`${name}.${key}`}
              displayName={property.properties[key].title || startCase(key)} 
              property={property.properties[key]}
              value={value ? value[key] : null}
              onChange={onChange}
            />
          )}
        </Card>
      )
    }

    // Special handling for 'key map' object types, represented in json schema as having no properties list and additionalProperties of type string
    if (property.type === 'object' && property.additionalProperties && property.additionalProperties.type === 'string') {
      return (
        <Form.Item label={displayName} help={help}>
          <KeyMap value={valueOrDefault} property={property} editable={editable} onChange={(v) => onChange(name, v)} />
        </Form.Item>
      )
    }

    // Handle all other types:
    return (
      <Form.Item label={displayName} help={help}>
        {(() => {
          switch(property.type) {
          case 'string': {
            if (property.format === 'multiline') {
              return <TextArea value={valueOrDefault} readOnly={!editable} onChange={(e) => onChange(name, e.target.value)} rows='20' />
            } else if (property.enum) {
              return <ConstrainedDropdown readOnly={!editable} value={valueOrDefault} allowedValues={property.enum} onChange={(e) => onChange(name, e)} />
            } else {
              return <Input value={valueOrDefault} readOnly={!editable} pattern={property.pattern} onChange={(e) => onChange(name, e.target.value)} />
            }
          }
          case 'boolean': {
            return <Switch checked={valueOrDefault} disabled={!editable} onChange={(v) => onChange(name, v)} checkedChildren={<Icon type="check" />} unCheckedChildren={<Icon type="close" />} />
          }
          case 'number': {
            return <InputNumber value={valueOrDefault} readOnly={!editable} onChange={(v) => onChange(name, v)} />
          }
          case 'integer': {
            return <InputNumber value={valueOrDefault} readOnly={!editable} onChange={(v) => onChange(name, v)} />
          }
          case 'array': {
            const values = valueOrDefault ? valueOrDefault : []
            if (property.items.type !== 'array' && property.items.type !== 'object') {
              return <Select mode="tags" tokenSeparators={[',']} value={values} disabled={!editable} onChange={(v) => onChange(name, v)} />
            } else {
              return (
                <>
                  {values.map((val, ind) =>
                    <React.Fragment key={`${name}[${ind}]`}>
                      <PlanOption
                        {...this.props}
                        name={`${name}[${ind}]`}
                        property={property.items}
                        value={val}
                        onChange={onChange}
                      />
                      <Button disabled={!editable} icon="delete" title={`Remove ${displayName} ${ind}`} onClick={() => onChange(name, this.removeFromArray(values, ind))}>
                        {`Remove ${displayName} ${ind}`}
                      </Button>
                    </React.Fragment>
                  )}
                  {(values.length === 0) ?
                    <Alert type="info" message={`No ${displayName} currently defined.`}/>
                    : null}
                  <Button disabled={!editable} icon="plus" title={`Add new ${displayName}`} onClick={() => onChange(name, this.addComplexItemToArray(property, values))}>
                    {`Add new ${displayName}`}
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
