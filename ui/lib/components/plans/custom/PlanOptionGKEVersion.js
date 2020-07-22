import * as React from 'react'
import { Typography } from 'antd'
const { Paragraph } = Typography

import PlanOptionVersion from './PlanOptionVersion'

export default class PlanOptionGKEVersion extends PlanOptionVersion {

  render() {
    const releaseChannel = this.props.plan.releaseChannel

    // Drop the version control all together if release channel set.
    if (releaseChannel && releaseChannel !== '') {
      return null
    }

    const help = (
      <>
        <Paragraph type="secondary">Following a release channel above is recommended, but if a specific version is required, setting this to a Kubernetes minor version (e.g. 1.15) is recommended over specific patch versions (e.g. 1.15.1 or 1.15.1-gke.2).</Paragraph>
        <Paragraph type="secondary">This sets the version to deploy; note the master <b>will be auto-upgraded by GCP</b> and master auto-upgrades cannot be disabled. Node pools will also be auto-upgraded unless they have auto-upgrade disabled.</Paragraph>
        {this.props.mode === 'edit' ? <Paragraph type="secondary">You may upgrade this cluster by specifying a later version, but you cannot downgrade the cluster. It is only valid to upgrade by a single minor version (e.g. 1.15 to 1.16) and all node pools must be within two minor versions of the master.</Paragraph> : null}
        <Paragraph type="secondary">Examples: - (current default), 1.15 (recommended - latest 1.15.x-gke.y), 1.15.1 (latest 1.15.1-gke.y), 1.15.1-gke.1 (exact GKE patch version, not recommended), latest (bleeding edge, not recommended)</Paragraph>
      </>
    )

    return this.renderControl({ help })
  }
}