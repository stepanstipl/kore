import * as React from 'react'
import { Form, Typography, Input } from 'antd'
const { Paragraph } = Typography
import { startCase } from 'lodash'

import PlanOptionBase from '../PlanOptionBase'

export default class PlanOptionGKEVersion extends PlanOptionBase {
  render() {
    const { name, editable, property, value, plan } = this.props

    // Drop the version control all together if release channel set.
    if (plan.releaseChannel && plan.releaseChannel !== '') {
      return null
    }

    const onChange = this.props.onChange || (() => {})
    const displayName = this.props.displayName || property.title || startCase(name)
    const valueOrDefault = value !== undefined && value !== null ? value : property.default

    const help = (
      <>
        <Paragraph type="secondary">Following a release channel above is recommended, but if a specific version is required, setting this to a Kubernetes minor version (e.g. 1.15) is recommended over specific patch versions (e.g. 1.15.1 or 1.15.1-gke.2).</Paragraph>
        <Paragraph type="secondary">This sets the version to deploy; note the master <b>will be auto-upgraded by GCP</b> and master auto-upgrades cannot be disabled. Node pools will also be auto-upgraded unless they have auto-upgrade disabled.</Paragraph>
        {this.props.mode === 'edit' ? <Paragraph type="secondary">You may upgrade this cluster by specifying a later version, but you cannot downgrade the cluster. It is only valid to upgrade by a single minor version (e.g. 1.15 to 1.16) and all node pools must be within two minor versions of the master.</Paragraph> : null}
        <Paragraph type="secondary">Examples: - (current default), 1.15 (recommended - latest 1.15.x-gke.y), 1.15.1 (latest 1.15.1-gke.y), 1.15.1-gke.1 (exact GKE patch version, not recommended), latest (bleeding edge, not recommended)</Paragraph>
      </>
    )

    return (
      <Form.Item label={displayName} help={help}>
        <Input value={valueOrDefault} readOnly={!editable} pattern={property.pattern} onChange={(e) => onChange(name, e.target.value)} />
        {this.validationErrors(name)}
      </Form.Item>
    )
  }
}