import * as React from 'react'
import { Form, Typography } from 'antd'
const { Paragraph } = Typography
import { startCase } from 'lodash'

import PlanOptionBase from '../PlanOptionBase'
import ConstrainedDropdown from './ConstrainedDropdown'

const releaseChannels = [
  { value: '', display: 'None (not recommended)' },
  { value: 'REGULAR', display: 'Regular (recommended)' },
  { value: 'STABLE', display: 'Stable' },
  { value: 'RAPID', display: 'Rapid (not recommended for production workloads)' },
]

const releaseChannelInfo = {
  '': 'No release channel set - specify version below or choose a channel (recommended)',
  'REGULAR': 'Regular updates to stable, tested new Kubernetes releases - steady and predictable release cadence, recommended for most workloads',
  'STABLE': 'Stable prioritizes stability over new functionality - less frequent updates which have undergone the most testing before being rolled out, but last to get new Kubernetes features. Recommended only for workloads with a very low tolerance for upgrade disruption.',
  'RAPID': 'Rapid provides early access to new Kubernetes features but with the risks inherent in frequent new feature updates. Not recommended for production workloads, but ideal to test new Kubernetes features'
}

export default class PlanOptionGKEReleaseChannel extends PlanOptionBase {
  onChannelChange = (channel) => {
    const { onChange, plan } = this.props
    if (channel !== '') {
      // Blank out the versions as well as setting the the channel, and force auto-upgrade
      // to true on all node pools.
      onChange('version', '')
      if (plan.nodePools) {
        plan.nodePools.forEach((np, idx) => {
          onChange(`nodePools[${idx}].version`, '')
          onChange(`nodePools[${idx}].enableAutoupgrade`, true)
        })
      }
    }
    onChange('releaseChannel', channel)
  }

  render() {
    const { name, editable, property, value } = this.props

    const displayName = this.props.displayName || property.title || startCase(name)
    const valueOrDefault = value !== undefined && value !== null ? value : property.default

    const help = (
      <>
        {!releaseChannelInfo[valueOrDefault] ? null : <Paragraph type="secondary">{releaseChannelInfo[valueOrDefault]}</Paragraph>}
        {valueOrDefault === '' ? null : <Paragraph type="secondary">All channels receive critical security patches as required.</Paragraph>}
      </>
    )

    return (
      <Form.Item label={displayName} help={help}>
        <ConstrainedDropdown readOnly={!editable} value={valueOrDefault} allowedValues={releaseChannels} onChange={(e) => this.onChannelChange(e)} />
      </Form.Item>
    )
  }
}