import * as React from 'react'
import PropTypes from 'prop-types'
import { Form, Cascader, Icon, Alert } from 'antd'
import { startCase } from 'lodash'

import PlanOptionBase from '../PlanOptionBase'
import KoreApi from '../../../kore-api'
import PlanOption from '../PlanOption'


const FAMILY_INFO = {
  'GKE': {
    families: {
      'c2': 'are compute-optimized offering highest performance per core',
      'e2': 'are cost-optimized, e2-standard are good for general workloads, e2-medium, e2-small and e2-micro offer shared virtual cores for very low usage workloads',
      'f1': 'are shared core with very restricted memory - not recommended',
      'g1': 'are shared core for very low usage workloads',
      'm1': 'are memory-optimized (first generation)',
      'm2': 'are memory-optimized (current generation) - REQUIRE sustained use contract',
      'n1': 'are 1st generation with slightly lower performance per vCPU than N2',
      'n2': 'are current generation general usage',
      'n2d': 'have AMD EPYC Rome processors',
    },
    categories: {
      'General purpose': 'types provide balanced CPU and memory for normal workloads',
      'Compute optimized': 'types provide more CPU and less memory for processor intensive workloads',
      'Memory optimized': 'types provide more memory and less CPU for memory intensive workloads',
    }
  },
  'EKS': {
    families: {
      
    },
    categories: {

    }
  },
  'AKS': {
    families: {

    },
    categories: {

    }
  }
}

export default class PlanOptionClusterMachineType extends PlanOptionBase {
  static propTypes = {
    nodePriceSet: PropTypes.func,
    filterCategories: PropTypes.func
  }

  state = {
    loadingInstances: true,
    noRegion: false,
    types: null,
    typeIndex: null,
    priceIndex: {},
    extInfo: null
  }

  familySplitter = '-'

  componentDidMountComplete = null
  componentDidMount = () => {
    // load regions for provider
    this.componentDidMountComplete = Promise.resolve().then(async() => {
      this.setState({ loadingInstances: true })
      const provider = this.props.kind
      const region = this.props.plan.region
      if (!region) {
        this.setState({ loadingInstances: false, noRegion: true, types: null, typeIndex: null })
        return
      }

      const api = await KoreApi.client()
      try {
        const typeInfo = await api.metadata.GetKubernetesNodeTypes(provider, region)
        const { types, typeIndex, priceIndex } = this.mapTypes(typeInfo, provider)
        this.setState({ loadingInstances: false, noRegion: false, types, typeIndex, priceIndex })
        if (this.props.nodePriceSet && priceIndex[this.props.value]) {
          this.props.nodePriceSet(priceIndex[this.props.value])
        }
      } catch (err) {
        // Couldn't load regions, use free text.
        console.warn('Error loading types - will use plain text entry instead', err)
        this.setState({ loadingInstances: false, noRegion: false, types: null, typeIndex: null })
      }
    })
  }

  mapTypes(typeInfo, provider) {
    if (!typeInfo) {
      return { types: null, typeIndex: null }
    }
    const typeIndex = {}
    const priceIndex = {}
    const typeMap = {}
    const types = []

    let familySplitter
    switch (provider.toUpperCase()) {
    case 'GKE':
      familySplitter = '-'
      break
    case 'EKS':
      familySplitter = '.'
      break
    }

    typeInfo.items.forEach((instType) => {
      if (!typeMap[instType.category]) {
        typeMap[instType.category] = {
          value: instType.category,
          label: instType.category,
          children: []
        }
        types.push(typeMap[instType.category])
      }
      const familyStr = familySplitter ? instType.name.substr(0, instType.name.indexOf(familySplitter)) : 'All'
      let family = typeMap[instType.category].children.find((f) => f.value === familyStr)
      if (!family) {
        family = {
          value: familyStr,
          label: familyStr.toUpperCase(),
          children: []
        }
        typeMap[instType.category].children.push(family)
      }
      family.children.push({
        value: instType.name,
        label: `${instType.name} | ${instType.mCpus/1000}CPU${instType.mCpus > 1000 ? 's' : ''} | ${instType.mem/1000.0}GB | $${(instType.prices.OnDemand / 1000000).toFixed(3)}/hr`,
        mCpus: instType.mCpus,
        mem: instType.mem,
        prices: instType.prices
      })
      typeIndex[instType.name] = [instType.category, familyStr, instType.name]
      priceIndex[instType.name] = instType.prices
    })

    // Sort categories alphabetically
    types.sort((a, b) => (a.value < b.value) ? -1 : (a.value > b.value) ? 1 : 0)
    types.forEach((category) => {
      // Sort families alphabetically within a category
      category.children.sort((a, b) => (a.value < b.value) ? -1 : (a.value > b.value) ? 1 : 0)
      // Sort instances by CPUs then Memory
      category.children.forEach((family) => {
        family.children.sort((a, b) => a.mCpus === b.mCpus ? a.mem - b.mem : a.mCpus - b.mCpus)
      })
    })
    
    return { types, typeIndex, priceIndex }
  }

  getFilteredTypes(types) {
    if (!this.props.filterCategories) {
      return types
    }
    return types.filter((t) => this.props.filterCategories(t.value))
  }  

  setExtInfo = (v) => {
    if (!v || v.length < 3) {
      if (this.state.extInfo !== '') {
        this.setState({ extInfo: '' })
      }
      return
    }
    const { kind } = this.props
    let extInfo = v[0]
    if (FAMILY_INFO[kind] && FAMILY_INFO[kind].categories[v[0]]) {
      extInfo += ' ' + FAMILY_INFO[kind].categories[v[0]]
    }
    extInfo += ', ' + v[1].toUpperCase()
    if (FAMILY_INFO[kind] && FAMILY_INFO[kind].families[v[1]]) {
      extInfo += ' ' + FAMILY_INFO[kind].families[v[1]]
    }
    this.setState({ extInfo })
  }

  render() {
    const { name, editable, property, value, nodePriceSet } = this.props
    const { types, typeIndex, priceIndex, loadingInstances, noRegion, extInfo } = this.state

    if (loadingInstances) {
      return <Icon type="loading" />
    }

    if (noRegion) {
      return (
        <>
          <Alert type="info" message="Set region on plan to see available instance types and prices" />
          <PlanOption {...this.props} disableCustom={true} editable={false} />
        </>
      )
    }

    // If we have no types from the metadata service, just use the default text control
    if (!types) {
      return <PlanOption {...this.props} disableCustom={true} />
    }

    const displayName = this.props.displayName || property.title || startCase(name)
    // const help = this.props.help || this.describe(property)
    const defaultValue = property.const !== undefined && property.const !== null ? property.const : property.default
    const valueOrDefault = value !== undefined && value !== null ? value : defaultValue
    const id = this.props.id || `plan_input_${name}`

    const selectedInstType = valueOrDefault && typeIndex[valueOrDefault] ? [...typeIndex[valueOrDefault]] : null
    if (!extInfo) {
      this.setExtInfo(selectedInstType)
    }

    const onChange = (v) => {
      if (this.props.onChange){
        this.props.onChange(name, v[2]) 
      }
      if (nodePriceSet && priceIndex[v[2]]) {
        nodePriceSet(priceIndex[v[2]])
      }
      this.setExtInfo(v)
    }

    return (
      <>
        <Form.Item label={displayName} help={extInfo}>
          <Cascader id={id} style={{ width: '100%' }} disabled={!editable} options={this.getFilteredTypes(types)} value={selectedInstType} onChange={onChange} />
        </Form.Item>
      </>
    )
  }
}