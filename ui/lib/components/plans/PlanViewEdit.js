import * as React from 'react'
import PropTypes from 'prop-types'
import { Form, Checkbox } from 'antd'
import PlanOption from './PlanOption'

/**
 * PlanViewEdit is the underlying component which handles all plan viewing and editing. Most likely, you
 * want to use UsePlanForm or ManagePlanForm instead of using this directly.
 */
export default class PlanViewEdit extends React.Component {
  static propTypes = {
    resourceType: PropTypes.oneOf(['cluster', 'service', 'servicecredential']).isRequired,
    mode: PropTypes.oneOf(['create','edit','view']).isRequired,
    manage: PropTypes.bool,
    team: PropTypes.object,
    kind: PropTypes.string.isRequired,
    plan: PropTypes.object.isRequired,
    schema: PropTypes.object.isRequired,
    editableParams: PropTypes.array.isRequired,
    onPlanValueChange: PropTypes.func,
    validationErrors: PropTypes.array
  }

  state = {
    showReadOnly: false
  }

  constructor(props) {
    super(props)
    this.state = {
      showReadOnly: props.mode === 'view' || props.manage === true
    }
  }

  componentDidUpdate(prevProps) {
    if (this.props.mode !== prevProps.mode) {
      this.setState({ showReadOnly: this.props.mode === 'view' || this.props.manage === true })
    }
  }

  render() {
    const { resourceType, mode, manage, team, kind, plan, schema, editableParams, onPlanValueChange, validationErrors } = this.props
    const { showReadOnly } = this.state
    return (
      <>
        {mode !== 'view' && !editableParams.includes('*') ? (
          <Form.Item label="Show read-only parameters">
            <Checkbox onChange={(v) => this.setState({ showReadOnly: v.target.checked })} checked={showReadOnly} />
          </Form.Item>
        ): null}

        {Object.keys(schema.properties).map((name) => {
          const editable = mode !== 'view' &&
            (editableParams.includes('*') || editableParams.includes(name)) &&
            (schema.properties[name].const === undefined || schema.properties[name].const === null) &&
            (mode === 'create' || manage || !schema.properties[name].immutable) // Disallow editing of params which can only be set at create time when in 'use' mode

          return (
            <PlanOption
              manage={manage}
              mode={mode}
              team={team}
              resourceType={resourceType}
              kind={kind}
              plan={plan}
              key={name}
              name={name}
              property={schema.properties[name]}
              value={plan[name]}
              hideNonEditable={!showReadOnly}
              editable={editable}
              onChange={(n, v) => onPlanValueChange(n, v)}
              validationErrors={validationErrors} />
          )
        })}
      </>
    )
  }
}
