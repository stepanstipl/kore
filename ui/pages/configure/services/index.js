import React from 'react'
import ServiceKindList from '../../../lib/components/services/ServiceKindList'
import Breadcrumb from '../../../lib/components/layout/Breadcrumb'
import { Alert } from 'antd'

export default class ServiceIndexPage extends React.Component {
  render() {
    return (
      <>
        <Breadcrumb
          items={[
            { text: 'Cloud Services', href: '/configure/services', link: '/configure/services' }
          ]}
        />
        <Alert 
          type="info" 
          message="Cloud Services"
          description="Enabling services allows teams to provision additional resources, either from cloud providers or directly into their clusters. Each service type can be enabled or disabled, and selecting 'Manage' allows control over the plans for a specific service."
          style={{ marginBottom: '20px' }}
        />
        <ServiceKindList />
      </>
    )
  }
}