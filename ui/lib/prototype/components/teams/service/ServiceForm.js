import React from 'react'
import PropTypes from 'prop-types'
import { Alert, Avatar, Button, Card, Form, List, Select, Tag } from 'antd'

import copy from '../../../../../lib/utils/object-copy'

// prototype imports
import CloudSelector from '../../../../../lib/prototype/components/common/CloudSelector'

class ServiceForm extends React.Component {
  static propTypes = {
    form: PropTypes.object.isRequired,
    clusters: PropTypes.object.isRequired,
    namespaceClaims: PropTypes.object.isRequired,
    handleSubmit: PropTypes.func.isRequired,
    handleCancel: PropTypes.func.isRequired
  }

  availableServiceList = [{
    name: 'Amazon S3',
    description: `Amazon Simple Storage Service (Amazon S3) is storage for the
        Internet. You can use Amazon S3 to store and retrieve any amount of data at
        any time, from anywhere on the web. You can accomplish these tasks using the
        simple and intuitive web interface of the AWS Management Console.`,
    plan: 'Production',
    imageURL: 'https://s3.amazonaws.com/awsservicebroker/icons/Storage_AmazonS3_LARGE.png'
  }, {
    name: 'Amazon RDS for PostgreSQL',
    description: `PostgreSQL has become the preferred open source relational database
        for many enterprise developers and start-ups, powering leading geospatial and
        mobile applications. Amazon RDS makes it easy to set up, operate, and scale
        PostgreSQL deployments in the cloud. With Amazon RDS, you can deploy scalable
        PostgreSQL deployments in minutes with cost-efficient and resizable hardware
        capacity. Amazon RDS manages complex and time-consuming administrative tasks
        such as PostgreSQL software installation and upgrades; storage management; replication
        for high availability and read throughput; and backups for disaster recovery.`,
    plan: 'Production',
    imageURL: 'https://s3.amazonaws.com/awsservicebroker/icons/AmazonRDS_LARGE.png'
  }, {
    name: 'Amazon RDS for PostgreSQL',
    description: `PostgreSQL has become the preferred open source relational database
        for many enterprise developers and start-ups, powering leading geospatial and
        mobile applications. Amazon RDS makes it easy to set up, operate, and scale
        PostgreSQL deployments in the cloud. With Amazon RDS, you can deploy scalable
        PostgreSQL deployments in minutes with cost-efficient and resizable hardware
        capacity. Amazon RDS manages complex and time-consuming administrative tasks
        such as PostgreSQL software installation and upgrades; storage management; replication
        for high availability and read throughput; and backups for disaster recovery.`,
    plan: 'Development',
    imageURL: 'https://s3.amazonaws.com/awsservicebroker/icons/AmazonRDS_LARGE.png'
  }]

  state = {
    selectedCloud: 'AWS',
    selectedService: '',
    formErrorMessage: false
  }

  handleSelectCloud = cloud => {
    if (this.state.selectedCloud !== cloud) {
      const state = copy(this.state)
      state.selectedCloud = cloud
      this.setState(state)
    }
  }

  selectService = (name, plan) => () => this.setState({ selectedService: `${name}_${plan}` })

  handleSubmit = (e) => {
    e.preventDefault()

    this.props.form.validateFields((err, values) => {
      if (err) {
        this.setState({ formErrorMessage: 'Validation failed' })
        return
      }
      if (!this.state.selectedService) {
        this.setState({ formErrorMessage: 'Please select a cloud service' })
        return
      }

      const selectedService = this.availableServiceList.find(s => `${s.name}_${s.plan}` === this.state.selectedService)
      this.props.handleSubmit({
        ...selectedService,
        ...values
      })
      this.resetForm()
    })
  }

  resetForm = () => {
    this.setState({ selectedService: '' })
    this.props.form.resetFields()
  }

  cancel = () => {
    this.resetForm()
    this.props.handleCancel()
  }

  render() {
    const { getFieldDecorator } = this.props.form
    const clusters = this.props.clusters.items
    const namespaces = this.props.namespaceClaims.items
    const { selectedCloud, selectedService } = this.state

    const formConfig = {
      layout: 'horizontal',
      labelAlign: 'left',
      hideRequiredMark: true,
      labelCol: {
        sm: { span: 24 },
        md: { span: 6 },
        lg: { span: 4 }
      },
      wrapperCol: {
        span: 12
      }
    }

    return (
      <div>
        <CloudSelector selectedCloud={selectedCloud} handleSelectCloud={this.handleSelectCloud} />
        <Form {...formConfig} onSubmit={this.handleSubmit}>
          <Card style={{ marginBottom: '20px' }}>
            <Alert
              message="Environment"
              description="Select the cluster and namespace combination you would like for the cloud service."
              type="info"
              showIcon
              style={{ marginBottom: '20px' }}
            />
            <Form.Item label="Cluster" help="Choose your cluster">
              {getFieldDecorator('cluster', {
                rules: [{ required: true, message: 'Please select the cluster!' }],
                initialValue: clusters.length === 1 ? clusters[0].metadata.name : undefined
              })(
                <Select placeholder="Cluster">
                  {clusters.map(c => <Select.Option key={c.metadata.name} value={c.metadata.name}>{c.metadata.name}</Select.Option>)}
                </Select>
              )}
            </Form.Item>
            <Form.Item label="Namespace" help="Choose your namespace">
              {getFieldDecorator('namespaceClaim', {
                rules: [{ required: true, message: 'Please select the namespace!' }],
                initialValue: namespaces.length === 1 ? namespaces[0].metadata.name : undefined
              })(
                <Select placeholder="Namespace">
                  {namespaces.map(n => <Select.Option key={n.spec.name} value={n.spec.name}>{n.spec.name}</Select.Option>)}
                </Select>
              )}
            </Form.Item>
          </Card>
          <Card style={{ marginBottom: '20px' }}>
            <Alert
              message="Cloud service"
              description="Select the cloud service you would like to use."
              type="info"
              showIcon
              style={{ marginBottom: '20px' }}
            />
            <List
              dataSource={this.availableServiceList}
              renderItem={service => (
                <List.Item actions={`${service.name}_${service.plan}` !== selectedService ? [<Button key="select" onClick={this.selectService(service.name, service.plan)}>Select</Button>] : [<Tag key="selected">Selected</Tag>]}>
                  <List.Item.Meta
                    avatar={<Avatar src={service.imageURL} />}
                    title={<p>{service.name} - {service.plan}</p>}
                  />
                </List.Item>
              )}
            />
          </Card>

          <Form.Item>
            <Button type="primary" htmlType="submit" disabled={!selectedService}>Save</Button>
            <Button type="link" onClick={this.cancel}>Cancel</Button>
          </Form.Item>
        </Form>
      </div>
    )
  }
}

const WrappedServiceForm = Form.create({ name: 'service_claim' })(ServiceForm)

export default WrappedServiceForm
