import * as React from 'react'
import PropTypes from 'prop-types'
import { Form, Icon } from 'antd'
import { uniq } from 'lodash'

import PlanOptionBase from '../PlanOptionBase'
import ConstrainedDropdown from './ConstrainedDropdown'
import PlanOption from '../PlanOption'
import KoreApi from '../../../kore-api'

export default class PlanOptionVersion extends PlanOptionBase {
  static propTypes = {
    expandVersions: PropTypes.bool
  }

  state = {
    loadingVersions: true,
    versions: null
  }

  fetchVersions = async () => {
    this.setState({ loadingVersions: true })
    const provider = this.props.kind
    const region = this.props.plan.region
    if (!region) {
      return this.setState({ loadingVersions: false, versions: null })
    }
    try {
      let versions = await (await KoreApi.client()).metadata.GetKubernetesVersions(provider, region)
      if (this.props.expandVersions) {
        versions = this.expandVersions(versions)
      }
      this.setState({ loadingVersions: false, versions })
    } catch (err) {
      console.warn('Error loading versions - will use plain text entry instead', err)
      this.setState({ loadingVersions: false, versions: null })
    }
  }

  componentDidMountComplete = null
  componentDidMount = () => {
    this.componentDidMountComplete = Promise.resolve().then(async() => await this.fetchVersions())
  }

  componentDidUpdateComplete = null
  componentDidUpdate = (prevProps) => {
    // refresh the list of versions when the region is changed
    if (prevProps.plan.region !== this.props.plan.region) {
      this.componentDidUpdateComplete = Promise.resolve().then(async() => await this.fetchVersions())
    }
  }

  expandVersions = (versions) => {
    const minorVersions = []
    versions.forEach(v => {
      const versionParts = v.split('.')
      // eg 1.16.10, add 1.16 as an available version
      if (versionParts.length === 3) {
        minorVersions.push(v.split('.').filter((e, i) => i <= 1).join('.'))
      }
      // eg 1.16.10-gke.1, add 1.16 and 1.16.10 as available versions
      if (versionParts.length > 3) {
        minorVersions.push(v.split('.').filter((e, i) => i <= 1).join('.'))
        minorVersions.push(v.split('.').filter((e, i) => i <= 2).join('.').replace(/[^0-9.]+/g, ''))
      }
    })
    console.log('expanding versions', minorVersions)
    return uniq([ ...minorVersions, ...versions ]).sort()
  }

  renderControl = (overrides) => {
    overrides = overrides || {}
    const { name, editable, plan } = this.props
    const { onChange, displayName, valueOrDefault, id, help } = this.prepCommonProps(this.props)
    const { versions, loadingVersions } = this.state

    const readOnly = !editable || !plan.region
    const value = plan.region ? valueOrDefault : 'Choose region first'

    return (
      <Form.Item label={displayName} help={overrides.help || help}>
        {loadingVersions ? <Icon type="loading" /> : null}
        {!loadingVersions ? <ConstrainedDropdown id={id} readOnly={readOnly} allowedValues={versions || []} value={value} onChange={(v) => onChange(name, v)} /> : null}
        {this.validationErrors(name)}
      </Form.Item>
    )
  }

  render() {
    const region = this.props.plan.region
    const versions = this.state.versions

    // If we have selected a region and have no versions from the metadata service, just use the default text control
    if (!versions && region) {
      return <PlanOption {...this.props} disableCustom={true} />
    }

    return this.renderControl()
  }
}