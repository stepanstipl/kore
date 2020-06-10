import React from 'react'
import ServiceKindList from '../../../lib/components/services/ServiceKindList'
import Breadcrumb from '../../../lib/components/layout/Breadcrumb'
import { Alert, Card, Checkbox, Form, Input, Typography } from 'antd'
const { Text } = Typography

import { getKoreLabel } from '../../../lib/utils/crd-helpers'

export default class ServiceIndexPage extends React.Component {
  state = {
    platformFilter: [],
    serviceNameFilter: ''
  }

  changeFilters = (platform) => () => {
    if (this.state.platformFilter.indexOf(platform) >= 0) {
      this.setState({ platformFilter: this.state.platformFilter.filter(f => f !== platform) })
    } else {
      this.setState({ platformFilter: [ ...this.state.platformFilter, platform ] })
    }
  }

  checkboxFilter = (name) => (
    <Checkbox checked={this.state.platformFilter.indexOf(name) >= 0} onClick={this.changeFilters(name)}>{name}</Checkbox>
  )

  render() {
    return (
      <>
        <Breadcrumb items={[{ text: 'Configure' }, { text: 'Services' }]}/>
        <Alert
          type="info"
          message="Services"
          description="Enabling services allows teams to provision additional resources, either from cloud providers or directly into their clusters. Each service type can be enabled or disabled, and selecting 'Manage' allows control over the plans for a specific service."
          style={{ marginBottom: '20px' }}
        />
        <Card size="small" style={{ marginBottom: '20px' }}>
          <Form.Item labelAlign="left" labelCol={{ span: 4 }} label={<Text strong>Filter by platform</Text>} style={{ marginBottom: 0 }}>
            {this.checkboxFilter('AWS')}{this.checkboxFilter('Kubernetes')}
          </Form.Item>
          <Form.Item labelAlign="left" labelCol={{ span: 4 }} wrapperCol={{ span: 6 }} label={<Text strong>Filter by service name</Text>} style={{ marginBottom: 0 }}>
            <Input onChange={(e) => this.setState({ serviceNameFilter: e.target.value })} value={this.state.serviceNameFilter} placeholder="Filter by service name"/>
          </Form.Item>
          <a style={{ display: 'block', marginTop: '10px', marginBottom: '5px', textDecoration: 'underline' }} onClick={() => this.setState({ serviceNameFilter: '', platformFilter: [] })}>Clear filters</a>
        </Card>
        <ServiceKindList filter={ (s) => (this.state.platformFilter.length === 0 || this.state.platformFilter.includes(getKoreLabel(s, 'platform'))) && s.spec.displayName.toLowerCase().indexOf(this.state.serviceNameFilter.toLowerCase()) >= 0} />
      </>
    )
  }
}
