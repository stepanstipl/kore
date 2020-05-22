import * as React from 'react'
import { Form, Icon, List, Button, Drawer, Input, Descriptions, InputNumber, Checkbox, Collapse, Radio, Modal, Alert } from 'antd'
import { startCase } from 'lodash'
import PlanOptionBase from '../PlanOptionBase'
import ConstrainedDropdown from './ConstrainedDropdown'
import PlanOption from '../PlanOption'

export default class PlanOptionEKSNodeGroups extends PlanOptionBase {
  constructor(props) {
    super(props)
  }

  static AMI_TYPE_GENERAL = 'AL2_x86_64'
  static AMI_TYPE_GPU = 'AL2_x86_64_GPU'

  static defaultNewNodeGroup = {
    minSize: 1,
    maxSize: 10,
    desiredSize: 1,
    amiType: PlanOptionEKSNodeGroups.AMI_TYPE_GENERAL,
    instanceType: 't3.medium',
    diskSize: 10,
    name: '',
  }

  // @TODO: Pull these from AWS
  static supportedInstanceTypes = {
    [PlanOptionEKSNodeGroups.AMI_TYPE_GENERAL]: [
      't3.micro', 't3.small', 't3.medium', 't3.large', 't3.xlarge', 't3.2xlarge',
      't3a.micro', 't3a.small', 't3a.medium', 't3a.large', 't3a.xlarge', 't3a.2xlarge',
      'm5.large', 'm5.xlarge', 'm5.2xlarge', 'm5.4xlarge', 'm5.8xlarge', 'm5.12xlarge',
      'm5a.large', 'm5a.xlarge', 'm5a.2xlarge', 'm5a.4xlarge',
      'c5.large', 'c5.xlarge', 'c5.2xlarge', 'c5.4xlarge', 'c5.9xlarge',
      'r5.large', 'r5.xlarge', 'r5.2xlarge', 'r5.4xlarge',
      'r5a.large', 'r5a.xlarge', 'r5a.2xlarge', 'r5a.4xlarge',
    ],
    [PlanOptionEKSNodeGroups.AMI_TYPE_GPU]: [
      'g4dn.xlarge', 'g4dn.2xlarge', 'g4dn.4xlarge', 'g4dn.8xlarge', 'g4dn.12xlarge',
      'p2.xlarge', 'p2.8xlarge', 'p2.16xlarge',
      'p3.2xlarge', 'p3.8xlarge', 'p3.16xlarge',
      'p3dn.24xlarge',
    ],
  }

  state = {
    selectedNodeGroupIndex: -1,
  }

  addNodeGroup = () => {
    if (!this.props.editable || !this.props.onChange) {
      return
    }

    // Need to handle the value being undefined in the case where this is a new plan or no
    // node groups are defined yet.
    let newValue
    if (this.props.value) {
      newValue = [ ...this.props.value, { ...PlanOptionEKSNodeGroups.defaultNewNodeGroup } ]
    } else {
      newValue = [ { ...PlanOptionEKSNodeGroups.defaultNewNodeGroup } ]
    }

    this.props.onChange(this.props.name, newValue)

    // Open the draw to immediately edit the new node group:
    this.setState({
      selectedNodeGroupIndex: newValue.length - 1
    })
  }

  removeNodeGroup = (idx) => {
    if (!this.props.editable || !this.props.onChange) {
      return
    }

    Modal.confirm({
      title: `Are you sure you want to remove node group ${idx + 1} (${this.props.value[idx].name})?`,
      okText: 'Yes',
      okType: 'danger',
      cancelText: 'No',
      onOk: () => {
        this.setState({
          selectedNodeGroupIndex: -1
        })
    
        this.props.onChange(
          this.props.name, 
          this.props.value.filter((_, i) => i !== idx)
        )
      }
    })
  }

  setAmiType = (idx, value) => {
    if (!this.props.editable || !this.props.onChange) {
      return
    }

    this.props.onChange(
      this.props.name, 
      this.props.value.map((ng, i) => i !== idx ? ng : { ...ng, amiType: value, instanceType: null })
    )
  }

  onReleaseVersionChecked = (idx, checked) => {
    if (!this.props.editable || !this.props.onChange) {
      return
    }

    const releaseVersion = !checked ? `${this.props.plan.version}.` : undefined

    this.props.onChange(
      this.props.name, 
      this.props.value.map((ng, i) => i !== idx ? ng : { ...ng, releaseVersion })
    )
  }

  setNodeGroupProperty = (idx, property, value) => {
    if (!this.props.editable || !this.props.onChange) {
      return
    }

    this.props.onChange(
      this.props.name, 
      this.props.value.map((ng, i) => i !== idx ? ng : { ...ng, [property]: value })
    )
  }

  viewEditNodeGroup = (idx) => {
    this.setState({
      selectedNodeGroupIndex: idx
    })
  }

  closeNodeGroup = () => {
    this.setState({
      selectedNodeGroupIndex: -1
    })
  }

  nodeGroupActions = (idx, editable) => {
    const actions = [
      <a key="viewedit" onClick={() => this.viewEditNodeGroup(idx)}><Icon type={editable ? 'edit' : 'eye'}></Icon></a>
    ]
    
    // Only show delete if we have more than one node group
    if (editable && this.props.value && this.props.value.length > 1) {
      actions.push(<a key="delete" onClick={() => this.removeNodeGroup(idx)}><Icon type="delete"></Icon></a>)
    }
    return actions
  }

  render() {
    const { name, editable, property, plan } = this.props
    const { selectedNodeGroupIndex } = this.state

    const value = this.props.value || []
    const selectedNodeGroup = selectedNodeGroupIndex >= 0 ? value[selectedNodeGroupIndex] : null
    const displayName = this.props.displayName || startCase(name)
    const description = this.props.manage ? 'Set default node groups for clusters created from this plan' : 'Manage node groups for this cluster'

    let instanceTypes = []
    let amiType = null
    let releaseVersionSet = false
    let ngNameClash = false
    if (selectedNodeGroup) {
      amiType = selectedNodeGroup.amiType || 'AL2_x86_64'
      instanceTypes = PlanOptionEKSNodeGroups.supportedInstanceTypes[amiType]
      releaseVersionSet = selectedNodeGroup.releaseVersion && selectedNodeGroup.releaseVersion.length > 0
      // we have duplicate names if another node group with a different index has the same name as this one
      ngNameClash = selectedNodeGroup.name && selectedNodeGroup.name.length > 0 && value.some((v, i) => i !== selectedNodeGroupIndex && v.name === selectedNodeGroup.name)
    }

    return (
      <>
        <Form.Item label={displayName} help={description}>
          <List dataSource={value} renderItem={(ng, idx) => {
            return (
              <List.Item actions={this.nodeGroupActions(idx, editable)}>
                <List.Item.Meta 
                  title={<a onClick={() => this.viewEditNodeGroup(idx)}>{`Node Group ${idx + 1} (${ng.name})`}</a>} 
                  description={`Size: min=${ng.minSize} max=${ng.maxSize} desired=${ng.desiredSize} | Node type: ${ng.instanceType}`} 
                />
                {!this.hasValidationErrors(`${name}[${idx}]`) ? null : <Alert type="error" message="Validation errors - please edit and resolve" />}
              </List.Item>
            )
          }} />
          {!editable ? null : <Button onClick={this.addNodeGroup}>Add node group</Button>}
        </Form.Item>
        <Drawer
          title={`Node Group ${selectedNodeGroup ? selectedNodeGroupIndex + 1 : ''}`}
          visible={Boolean(selectedNodeGroup)}
          closable={!ngNameClash}
          maskClosable={!ngNameClash}
          onClose={() => this.closeNodeGroup()}
          width={700}
        >
          {!selectedNodeGroup ? null : (
            <>
              <Collapse defaultActiveKey={['basics','compute','metadata']}>
                <Collapse.Panel key="basics" header="Basic Configuration (name, sizing)">
                  <Form.Item label="Name" help="Unique name for this group within the cluster">
                    <Input value={selectedNodeGroup.name} onChange={(e) => this.setNodeGroupProperty(selectedNodeGroupIndex, 'name', e.target.value)} readOnly={!editable} />
                    {this.validationErrors(`${name}[${selectedNodeGroupIndex}].name`)}
                    {!ngNameClash ? null : <Alert type="error" message="This name is already used by another node group, it must be changed." />}
                  </Form.Item>
                  <Form.Item label="Group Size">
                    <Descriptions layout="horizontal" size="small">
                      <Descriptions.Item label="Minimum">
                        <InputNumber value={selectedNodeGroup.minSize} size="small" min={property.items.properties.minSize.minimum} max={selectedNodeGroup.maxSize} readOnly={!editable} onChange={(v) => this.setNodeGroupProperty(selectedNodeGroupIndex, 'minSize', v)} />
                        {this.validationErrors(`${name}[${selectedNodeGroupIndex}].minSize`)}
                      </Descriptions.Item>
                      <Descriptions.Item label="Desired">
                        <InputNumber value={selectedNodeGroup.desiredSize} size="small" min={selectedNodeGroup.minSize} max={selectedNodeGroup.maxSize} readOnly={!editable} onChange={(v) => this.setNodeGroupProperty(selectedNodeGroupIndex, 'desiredSize', v)} />
                        {this.validationErrors(`${name}[${selectedNodeGroupIndex}].desiredSize`)}
                      </Descriptions.Item>
                      <Descriptions.Item label="Maximum">
                        <InputNumber value={selectedNodeGroup.maxSize} size="small" min={selectedNodeGroup.minSize} readOnly={!editable} onChange={(v) => this.setNodeGroupProperty(selectedNodeGroupIndex, 'maxSize', v)} />
                        {this.validationErrors(`${name}[${selectedNodeGroupIndex}].maxSize`)}
                      </Descriptions.Item>
                    </Descriptions>
                  </Form.Item>
                </Collapse.Panel>
                <Collapse.Panel key="compute" header="Compute Configuration (instance type, GPU or regular workload)">
                  <Form.Item label="Compute Type" help="Whether this node group is for general purpose or GPU workloads">
                    <Radio.Group value={amiType} onChange={(v) => this.setAmiType(selectedNodeGroupIndex, v.target.value)}>
                      <Radio value={PlanOptionEKSNodeGroups.AMI_TYPE_GENERAL}>General Purpose</Radio>
                      <Radio value={PlanOptionEKSNodeGroups.AMI_TYPE_GPU}>GPU</Radio>
                    </Radio.Group>
                    {this.validationErrors(`${name}[${selectedNodeGroupIndex}].amiType`)}
                  </Form.Item>
                  <Form.Item label="AWS AMI Version" help={!releaseVersionSet ? undefined : <><b>Must</b> be for Kubernetes <b>{plan.version}</b> and <b>{amiType === PlanOptionEKSNodeGroups.AMI_TYPE_GPU ? 'GPU' : 'general'}</b> workloads. Find <a target="_blank" rel="noopener noreferrer" href="https://docs.aws.amazon.com/eks/latest/userguide/eks-linux-ami-versions.html">supported versions</a> in AWS documentation.</>}>
                    <Checkbox disabled={!editable} checked={!releaseVersionSet} onChange={(v) => this.onReleaseVersionChecked(selectedNodeGroupIndex, v.target.checked)}/> Use latest (<b>recommended</b>)
                    {!releaseVersionSet ? null : <Input value={selectedNodeGroup.releaseVersion} readOnly={!editable} onChange={(e) => this.setNodeGroupProperty(selectedNodeGroupIndex, 'releaseVersion', e.target.value)} />}
                    {this.validationErrors(`${name}[${selectedNodeGroupIndex}].releaseVersion`)}
                  </Form.Item>
                  <Form.Item label="AWS Instance Type">
                    <ConstrainedDropdown allowedValues={instanceTypes} value={selectedNodeGroup.instanceType} onChange={(v) => this.setNodeGroupProperty(selectedNodeGroupIndex, 'instanceType', v)} />
                    {this.validationErrors(`${name}[${selectedNodeGroupIndex}].instanceType`)}
                  </Form.Item>
                  <PlanOption {...this.props} displayName="Instance Root Disk Size (GiB)" name={`${name}[${selectedNodeGroupIndex}].diskSize`} property={property.items.properties.diskSize} value={selectedNodeGroup.diskSize} onChange={(_, v) => this.setNodeGroupProperty(selectedNodeGroupIndex, 'diskSize', v)} />
                </Collapse.Panel>
                <Collapse.Panel key="metadata" header="Metadata (labels, tags, etc)">
                  <PlanOption {...this.props} displayName="Labels" help="Labels help kubernetes workloads find this group" name={`${name}[${selectedNodeGroupIndex}].labels`} property={property.items.properties.labels} value={selectedNodeGroup.labels} onChange={(_, v) => this.setNodeGroupProperty(selectedNodeGroupIndex, 'labels', v)} />
                  <PlanOption {...this.props} displayName="Tags" help="AWS tags to apply to the node group" name={`${name}[${selectedNodeGroupIndex}].tags`} property={property.items.properties.tags} value={selectedNodeGroup.tags} onChange={(_, v) => this.setNodeGroupProperty(selectedNodeGroupIndex, 'tags', v)} />
                </Collapse.Panel>
                <Collapse.Panel key="ssh" header="SSH Connectivity (keys, security groups)">
                  <PlanOption {...this.props} displayName="EC2 SSH Key" name={`${name}[${selectedNodeGroupIndex}].eC2SSHKey`} property={property.items.properties.eC2SSHKey} value={selectedNodeGroup.eC2SSHKey} onChange={(_, v) => this.setNodeGroupProperty(selectedNodeGroupIndex, 'eC2SSHKey', v)} />
                  <PlanOption {...this.props} displayName="SSH Security Groups" name={`${name}[${selectedNodeGroupIndex}].sshSourceSecurityGroups`} property={property.items.properties.sshSourceSecurityGroups} value={selectedNodeGroup.sshSourceSecurityGroups} onChange={(_, v) => this.setNodeGroupProperty(selectedNodeGroupIndex, 'sshSourceSecurityGroups', v)} />
                </Collapse.Panel>
              </Collapse>
            </>
          )}
        </Drawer>
      </>
    )
  }
}
