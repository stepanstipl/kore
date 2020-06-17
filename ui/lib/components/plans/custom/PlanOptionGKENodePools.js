import * as React from 'react'
import { Form, Icon, List, Button, Drawer, Input, Descriptions, InputNumber, Collapse, Modal, Alert, Switch, Cascader, Typography, Checkbox } from 'antd'
const { Paragraph } = Typography
import { startCase } from 'lodash'

import PlanOptionBase from '../PlanOptionBase'
import ConstrainedDropdown from './ConstrainedDropdown'
import PlanOption from '../PlanOption'

// @TODO: Pull these from GCP
const supportedMachineTypes = {
  'general2': {
    'name': 'General Purpose (2nd gen)',
    'about': 'General purpose machines offer specified number of vCPUs and 4GB (standard), 1GB (highcpu), or 8GB (highmem) memory per vCPU',
    'flavours': {
      'n2': {
        '_about': 'N2 - General',
        'standard': ['2', '4', '8', '16', '32', '48', '64', '80'],
        'highmem': ['2', '4', '8', '16', '32', '48', '64', '80'],
        'highcpu': ['2', '4', '8', '16', '32', '48', '64', '80']
      },
      'n2d': {
        '_about': 'N2D - AMD EPYC Rome Processors',
        'standard': ['2', '4', '8', '16', '32', '48', '64', '80', '96', '128', '224'],
        'highmem': ['2', '4', '8', '16', '32', '48', '64', '80', '96', '128', '224'],
        'highcpu': ['2', '4', '8', '16', '32', '48', '64', '80', '96', '128', '224']
      },
      'e2': {
        '_about': 'E2 - Cost-optimized',
        'standard': ['2', '4', '8', '16'],
        'highmem': ['2', '4', '8', '16'],
        'highcpu': ['2', '4', '8', '16']
      }
    }
  },
  'general1': {
    'name': 'General Purpose (1st gen)',
    'about': 'First generation machines offer slightly lower performance per vCPU than 2nd gen and 3.75GB (standard), 0.9GB (highcpu), or 6.5GB (highmem) memory per vCPU',
    'flavours': {
      'n1': {
        '_about': 'N1 - General',
        'standard': ['2', '4', '8', '16', '32', '48', '64', '96'],
        'highmem': ['2', '4', '8', '16', '32', '48', '64', '96'],
        'highcpu': ['2', '4', '8', '16', '32', '48', '64', '96']
      }
    }
  },
  'sharedcore': {
    'name': 'Shared Core',
    'about': 'Shared Core machines are cost-optimized for low usage workloads with the specified share of a vCPU and small amounts of memory. Very small instances (such as e2-micro and f1-micro) are unlikely to work well and are not recommended.',
    'flavours': {
      'e2': {
        '_about': 'E2 - Cost-optimized, dual virtual core',
        'micro': ['0.25 (not recommended)'],
        'small': ['0.5'],
        'medium': ['1'],
      },
      'f1': {
        '_about': 'F1 - Older generation, single virtual core',
        'micro': ['0.2 (not recommended)']
      },
      'g1': {
        '_about': 'G1 - Older generation, single virtual core',
        'small': ['0.5']
      }
    }
  },
  'compute': {
    'name': 'Compute Optimized',
    'about': 'Compute-optimized machines offer highest performance per core and 4GB memory per vCPU',
    'flavours': {
      'c2': {
        '_about': 'C2',
        'standard': ['2', '4', '8', '16', '30', '60']
      }
    }
  },
  'memory': {
    'name': 'Memory Optimized',
    'about': 'Memory-optimized machines offer large amounts of memory per vCPU',
    'flavours': {
      'm2': {
        '_about': 'M2 - sustained use contract only, 5TB-11TB total memory',
        'ultramem': ['208', '416']
      },
      'm1': {
        '_about': 'M1 - Older generation, 15-24GB memory per vCPU',
        'ultramem': ['40', '80', '160'],
        'megamem': ['96']
      }
    }
  }
}

// @TODO: Pull these from GCP
const imageTypes = [
  { value: 'COS', display: 'Container-Optimized OS (recommended)' },
  { value: 'COS_CONTAINERD', display: 'Container-Optimized OS with containerd' },
  { value: 'UBUNTU', display: 'Ubuntu' },
  { value: 'UBUNTU_CONTAINERD', display: 'Ubuntu with containerd' },
  { value: 'WINDOWS_LTSC', display: 'Windows (long-term support channel)' },
  { value: 'WINDOWS_SAC', display: 'Windows (semi-annual channel)' }
]

function getSupportedMachineTypes() {
  const index = {}
  const options = Object.keys(supportedMachineTypes).map(
    (k) => ({ value: k, label: supportedMachineTypes[k].name, children: getSupportedMachineTypeFlavours(k, index) })
  )
  return { index, options }
}

function getSupportedMachineTypeFlavours(family, index) {
  return Object.keys(supportedMachineTypes[family].flavours).map(
    (k) => ({ value: k, label: supportedMachineTypes[family].flavours[k]._about || k.toUpperCase(), children: getSupportedMachineTypeTypes(family, k, index) })
  )
}

function getSupportedMachineTypeTypes(family, flavour, index) {
  return Object.keys(supportedMachineTypes[family].flavours[flavour]).filter((k) => !k.startsWith('_')).map(
    (k) => { 
      return {
        value: k, 
        label: startCase(k), 
        children: supportedMachineTypes[family].flavours[flavour][k].map((cpus) => {
          const mt = getMachineType(family, flavour, k, cpus)
          index[mt] = [family, flavour, k, mt]
          return { value: mt, label: `${cpus} vCPUs` }
        })
      }
    }
  )
}

function getMachineType(family, flavour, type, cpus) {
  return family === 'sharedcore' ? `${flavour}-${type}`.toLowerCase() : `${flavour}-${type}-${cpus}`.toLowerCase()
}

export default class PlanOptionGKENodePools extends PlanOptionBase {
  state = {
    selectedIndex: -1
  }

  addNodePool = () => {
    if (!this.props.editable || !this.props.onChange) {
      return
    }

    // Create the default from the defaults on the plan schema
    const newNodePool = {
      name: ''
    }
    const properties = this.props.property.items.properties
    Object.keys(properties).forEach((k) => {
      if (properties[k].default !== undefined) {
        newNodePool[k] = properties[k].default
      }
    })

    // Need to handle the value being undefined in the case where this is a new plan or no
    // node pools are defined yet.
    let newValue
    if (this.props.value) {
      newValue = [ ...this.props.value, newNodePool ]
    } else {
      newValue = [ newNodePool ]
    }

    this.props.onChange(this.props.name, newValue)

    // Open the drawer to immediately edit the new node pool:
    this.setState({
      selectedIndex: newValue.length - 1
    })
  }

  removeNodePool = (idx) => {
    if (!this.props.editable || !this.props.onChange) {
      return
    }

    Modal.confirm({
      title: `Are you sure you want to remove node pool ${idx + 1} (${this.props.value[idx].name})?`,
      okText: 'Yes',
      okType: 'danger',
      cancelText: 'No',
      onOk: () => {
        this.setState({
          selectedIndex: -1
        })
    
        this.props.onChange(
          this.props.name, 
          this.props.value.filter((_, i) => i !== idx)
        )
      }
    })
  }

  setNodePoolProperty = (idx, property, value) => {
    if (!this.props.editable || !this.props.onChange) {
      return
    }

    this.props.onChange(
      this.props.name, 
      this.props.value.map((ng, i) => i !== idx ? ng : { ...ng, [property]: value })
    )
  }

  viewEditNodePool = (idx) => {
    this.setState({
      selectedIndex: idx
    })
  }

  closeNodePool = () => {
    this.setState({
      selectedIndex: -1
    })
  }

  nodePoolActions = (idx, editable) => {
    const actions = [
      <a id={`plan_nodepool_${idx}_viewedit`} key="viewedit" onClick={() => this.viewEditNodePool(idx)}><Icon type={editable ? 'edit' : 'eye'}></Icon></a>
    ]
    
    // Only show delete if we have more than one node pool
    if (editable && this.props.value && this.props.value.length > 1) {
      actions.push(<a id={`plan_nodepool_${idx}_del`} key="delete" onClick={() => this.removeNodePool(idx)}><Icon type="delete"></Icon></a>)
    }
    return actions
  }

  render() {
    const { name, editable, property, plan } = this.props
    const { selectedIndex } = this.state
    const id_prefix = 'plan_nodepool'

    const value = this.props.value || property.default || []
    const selected = selectedIndex >= 0 ? value[selectedIndex] : null
    const displayName = this.props.displayName || startCase(name)
    const description = this.props.manage ? 'Default node pools for clusters created from this plan' : null

    let ngNameClash = false, versionFollowMaster = false, nodePoolCloseable = true
    let allMachineTypes = null, selectedMachineTypeFamilyInfo = null
    if (selected) {
      // we have duplicate names if another node pool with a different index has the same name as this one
      ngNameClash = selected.name && selected.name.length > 0 && value.some((v, i) => i !== selectedIndex && v.name === selected.name)
      versionFollowMaster = !selected.version || selected.version === ''
      nodePoolCloseable = !ngNameClash && selected.name && selected.name.match(property.items.properties.name.pattern)
      allMachineTypes = getSupportedMachineTypes()
      selectedMachineTypeFamilyInfo = supportedMachineTypes[allMachineTypes.index[selected.machineType][0]].about
    }

    let followingReleaseChannel = false
    if (plan.releaseChannel && plan.releaseChannel !== '') {
      followingReleaseChannel = true
    }


    return (
      <>
        <Form.Item label={displayName} help={description}>
          <List id={`${id_prefix}s`} dataSource={value} renderItem={(ng, idx) => {
            return (
              <List.Item actions={this.nodePoolActions(idx, editable)}>
                <List.Item.Meta 
                  title={<a id={`${id_prefix}_${idx}_viewedittitle`} onClick={() => this.viewEditNodePool(idx)}>{`Node Pool ${idx + 1} (${ng.name})`}</a>} 
                  description={ng.enableAutoscaler ? 
                    `Size per zone: min=${ng.minSize} initial=${ng.size} max=${ng.maxSize} | Node type: ${ng.machineType}`
                    :
                    `Size per zone: ${ng.size} | Node type: ${ng.machineType}`
                  }
                />
                {!this.hasValidationErrors(`${name}[${idx}]`) ? null : <Alert type="error" message="Validation errors - please edit and resolve" />}
              </List.Item>
            )
          }} />
          {!editable ? null : <Button id={`${id_prefix}_add`} onClick={this.addNodePool}>Add node pool</Button>}
          {this.validationErrors(name)}
        </Form.Item>
        <Drawer
          title={`Node Pool ${selected ? selectedIndex + 1 : ''}`}
          visible={Boolean(selected)}
          closable={nodePoolCloseable}
          maskClosable={nodePoolCloseable}
          onClose={() => this.closeNodePool()}
          width={800}
        >
          {!selected ? null : (
            <>
              <Collapse defaultActiveKey={['basics','compute','metadata']}>
                <Collapse.Panel key="basics" header="Basic Configuration (name, versions, sizing)">
                  <Form.Item label="Name" help="Unique name for this group within the cluster">
                    <Input id={`${id_prefix}_name`} pattern={property.items.properties.name.pattern} value={selected.name} onChange={(e) => this.setNodePoolProperty(selectedIndex, 'name', e.target.value)} readOnly={!editable} />
                    {this.validationErrors(`${name}[${selectedIndex}].name`)}
                    {!ngNameClash ? null : <Alert type="error" message="This name is already used by another node pool, it must be changed." />}
                    {selected.name && selected.name.match(property.items.properties.name.pattern) ? null : <Alert type="error" message="Name must be minimum 2, maximum 40 alpha-numeric characters and hyphens" />}
                  </Form.Item>
                  {followingReleaseChannel ? (
                    <Form.Item label="Version" help="Set the Kubernetes version for this node pool">
                      <Alert type="info" message="Release channel is set so this node pool will automatically be upgraded in sync with the cluster. You cannot change auto-upgrade and node pool version with a release channel selected." />
                    </Form.Item>
                  ) : (
                    <>
                      <Form.Item label="Auto-upgrade" help="Allow GCP to automatically upgrade nodes in this pool (recommended)">
                        <Switch id={`${id_prefix}_enableAutoupgrade`} checked={selected.enableAutoupgrade} disabled={!editable || followingReleaseChannel} onChange={(v) => this.setNodePoolProperty(selectedIndex, 'enableAutoupgrade', v)} checkedChildren={<Icon type="check" />} unCheckedChildren={<Icon type="close" />} />
                      </Form.Item>
                      <Form.Item label="Version" help="Set the Kubernetes version for this node pool">
                        <Checkbox id={`${id_prefix}_versionFollowMaster`} checked={versionFollowMaster} onChange={(e) => this.setNodePoolProperty(selectedIndex, 'version', e.target.checked ? '' : plan.version)} disabled={!editable} /> Same as master (recommended)
                        {versionFollowMaster ? null : (
                          <>
                            <Input id={`${id_prefix}_version`} pattern={property.items.properties.version.pattern} value={selected.version} readOnly={!editable} onChange={(e) => this.setNodePoolProperty(selectedIndex, 'version', e.target.value)} />
                          </>
                        )}
                      </Form.Item>
                    </>
                  )}
                  <PlanOption id={`${id_prefix}_enableAutoscaler`} {...this.props} displayName="Auto-scale" name={`${name}[${selectedIndex}].enableAutoscaler`} property={property.items.properties.enableAutoscaler} value={selected.enableAutoscaler} onChange={(_, v) => this.setNodePoolProperty(selectedIndex, 'enableAutoscaler', v)} />
                  <Form.Item label="Pool size per zone">
                    <Descriptions layout="horizontal" size="small">
                      {!selected.enableAutoscaler ? null : <Descriptions.Item label="Minimum">
                        <InputNumber id={`${id_prefix}_minSize`} value={selected.minSize} size="small" min={property.items.properties.minSize.minimum} max={selected.maxSize} readOnly={!editable} onChange={(v) => this.setNodePoolProperty(selectedIndex, 'minSize', v)} />
                        {this.validationErrors(`${name}[${selectedIndex}].minSize`)}
                      </Descriptions.Item>}
                      <Descriptions.Item label={selected.enableAutoscaler ? 'Initial' : null}>
                        <InputNumber id={`${id_prefix}_size`} value={selected.size} size="small" min={selected.enableAutoscaler ? selected.minSize : 1} max={selected.enableAutoscaler ? selected.maxSize : 99999} readOnly={!editable} onChange={(v) => this.setNodePoolProperty(selectedIndex, 'size', v)} />
                        {this.validationErrors(`${name}[${selectedIndex}].size`)}
                      </Descriptions.Item>
                      {!selected.enableAutoscaler ? null : <Descriptions.Item label="Maximum">
                        <InputNumber id={`${id_prefix}_maxSize`} value={selected.maxSize} size="small" min={selected.minSize} readOnly={!editable} onChange={(v) => this.setNodePoolProperty(selectedIndex, 'maxSize', v)} />
                        {this.validationErrors(`${name}[${selectedIndex}].maxSize`)}
                      </Descriptions.Item>}
                    </Descriptions>
                  </Form.Item>
                  <PlanOption id={`${id_prefix}_maxPodsPerNode`} {...this.props} displayName="Max pods per node" name={`${name}[${selectedIndex}].maxPodsPerNode`} property={property.items.properties.maxPodsPerNode} value={selected.maxPodsPerNode} onChange={(_, v) => this.setNodePoolProperty(selectedIndex, 'maxPodsPerNode', v)} />
                </Collapse.Panel>
                <Collapse.Panel key="compute" header="Compute Configuration (machine type, disk size, image type, auto-repair)">
                  <Form.Item label="Image Type" help={<>For help choosing an image type, see <a target="_blank" rel="noopener noreferrer" href="https://cloud.google.com/kubernetes-engine/docs/concepts/node-images">the GCP documentation</a></>}>
                    <ConstrainedDropdown id={`${id_prefix}_imageType`} allowedValues={imageTypes} value={selected.imageType} onChange={(v) => this.setNodePoolProperty(selectedIndex, 'imageType', v)} />
                  </Form.Item>
                  <Form.Item label="GCP Machine Type" help={`${selected.machineType}`}>
                    <Cascader id={`${id_prefix}_machineType`} style={{ width: '100%' }} disabled={!editable} options={allMachineTypes.options} value={allMachineTypes.index[selected.machineType]} onChange={(v) => this.setNodePoolProperty(selectedIndex, 'machineType', v[3] )} />
                    {!selectedMachineTypeFamilyInfo ? null : <Paragraph type="secondary" style={{ lineHeight: '20px' }}>{selectedMachineTypeFamilyInfo}</Paragraph>}
                  </Form.Item>
                  <PlanOption id={`${id_prefix}_diskSize`} {...this.props} displayName="Instance Root Disk Size (GiB)" name={`${name}[${selectedIndex}].diskSize`} property={property.items.properties.diskSize} value={selected.diskSize} onChange={(_, v) => this.setNodePoolProperty(selectedIndex, 'diskSize', v)} />
                  <PlanOption id={`${id_prefix}_enableAutorepair`} {...this.props} displayName="Auto-repair" name={`${name}[${selectedIndex}].enableAutorepair`} property={property.items.properties.enableAutorepair} value={selected.enableAutorepair} onChange={(_, v) => this.setNodePoolProperty(selectedIndex, 'enableAutorepair', v)} />
                  <PlanOption id={`${id_prefix}_preemptible`} {...this.props} displayName="Pre-emptible" name={`${name}[${selectedIndex}].preemptible`} property={property.items.properties.preemptible} value={selected.preemptible} onChange={(_, v) => this.setNodePoolProperty(selectedIndex, 'preemptible', v)} />
                </Collapse.Panel>
                <Collapse.Panel key="metadata" header="Labels">
                  <PlanOption id={`${id_prefix}_labels`} {...this.props} displayName="Labels" help="Labels help kubernetes workloads find this group" name={`${name}[${selectedIndex}].labels`} property={property.items.properties.labels} value={selected.labels} onChange={(_, v) => this.setNodePoolProperty(selectedIndex, 'labels', v)} />
                </Collapse.Panel>
              </Collapse>
              <Form.Item>
                <Button id={`${id_prefix}_close`} type="primary" disabled={!nodePoolCloseable} onClick={() => this.closeNodePool()}>{nodePoolCloseable ? 'Close' : 'Node pool not valid - please correct errors'}</Button>
              </Form.Item>
            </>
          )}
        </Drawer>
      </>
    )
  }
}
