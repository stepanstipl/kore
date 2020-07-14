import * as React from 'react'
import { Form, Cascader, Icon } from 'antd'

import PlanOptionBase from '../PlanOptionBase'
import KoreApi from '../../../kore-api'
import PlanOption from '../PlanOption'

export default class PlanOptionClusterRegion extends PlanOptionBase {
  state = {
    loadingRegions: true,
    regions: null,
    regionIndex: null
  }

  componentDidMountComplete = null
  componentDidMount = () => {
    // load regions for provider
    this.componentDidMountComplete = Promise.resolve().then(async() => {
      this.setState({ loadingRegions: true })
      const provider = this.props.kind
      try {
        const regionInfo = await (await KoreApi.client()).metadata.GetKubernetesRegions(provider)
        const { regions, regionIndex } = this.mapRegions(regionInfo)
        this.setState({ loadingRegions: false, regions, regionIndex })
      } catch (err) {
        // Couldn't load regions, use free text.
        console.warn('Error loading regions - will use plain text entry instead', err)
        this.setState({ loadingRegions: false, regions: null, regionIndex: null })
      }
    })
  }

  mapRegions(regionInfo) {
    if (!regionInfo) {
      return { regions: null, regionIndex: null }
    }
    const regionIndex = {}
    const regions = regionInfo.items.map((continent) => {
      return {
        value: continent.name,
        label: continent.name,
        children: continent.regions.map((region) => {
          regionIndex[region.id] = [continent.name, region.id]
          return {
            value: region.id,
            label: `${region.name} - ${region.id}`
          }
        })
      }
    })
    return { regions, regionIndex }
  }

  render() {
    const { name, editable } = this.props
    const { onChange, displayName, help, valueOrDefault, id } = this.prepCommonProps(this.props)
    const { regions, regionIndex, loadingRegions } = this.state

    if (loadingRegions) {
      return <Icon type="loading" />
    }

    // If we have no regions from the metadata service, just use the default text control
    if (!regions) {
      return <PlanOption {...this.props} disableCustom={true} />
    }

    const selectedRegion = valueOrDefault ? regionIndex[valueOrDefault] : null

    return (
      <Form.Item label={displayName} help={help}>
        <Cascader id={id} style={{ width: '100%' }} disabled={!editable} options={regions} value={selectedRegion} onChange={(v) => onChange(name, v[1])} />
      </Form.Item>
    )
  }
}