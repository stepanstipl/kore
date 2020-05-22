import React from 'react'
import Link from 'next/link'
import moment from 'moment'
import PropTypes from 'prop-types'
import { Button, Col, Divider, Icon, message, Row, Tooltip, Tag, Typography } from 'antd'
const { Paragraph, Text } = Typography

import KoreApi from '../../../kore-api'
import copy from '../../../utils/object-copy'
import Service from './Service'
import { inProgressStatusList, statusColorMap, statusIconMap } from '../../../utils/ui-helpers'

class ServicesTab extends React.Component {

  static propTypes = {
    team: PropTypes.object.isRequired,
    getServiceCount: PropTypes.func
  }

  state = {
    dataLoading: true,
    services: [],
    serviceKinds: [],
    serviceCredentials: []
  }

  async fetchComponentData() {
    try {
      const team = this.props.team.metadata.name
      const api = await KoreApi.client()
      let [ services, serviceKinds, serviceCredentials ] = await Promise.all([
        api.ListServices(team),
        api.ListServiceKinds(team),
        api.ListServiceCredentials(team)
      ])
      services = services.items.filter(s => !s.spec.cluster || !s.spec.cluster.name)
      serviceKinds = serviceKinds.items
      serviceCredentials = serviceCredentials.items
      this.props.getServiceCount && this.props.getServiceCount(services.length)
      return { services, serviceKinds, serviceCredentials }
    } catch (err) {
      console.error('Unable to load data for services tab', err)
      return {}
    }
  }

  componentDidMountComplete = null
  componentDidMount() {
    this.componentDidMountComplete = this.fetchComponentData().then(data => {
      this.setState({ ...data, dataLoading: false })
    })
  }

  componentDidUpdate(prevProps) {
    if (prevProps.team.metadata.name !== this.props.team.metadata.name) {
      this.setState({ dataLoading: true })
      return this.fetchComponentData().then(data => this.setState({ ...data, dataLoading: false }))
    }
  }

  deleteService = async (name, done) => {
    const team = this.props.team.metadata.name
    try {
      const services = copy(this.state.services)
      const service = services.find(s => s.metadata.name === name)
      await (await KoreApi.client()).DeleteService(team, service.metadata.name)
      service.status.status = 'Deleting'
      service.metadata.deletionTimestamp = new Date()
      this.setState({ services }, done)
      message.loading(`Service deletion requested: ${service.metadata.name}`)
    } catch (err) {
      console.error('Error deleting service', err)
      message.error('Error deleting service, please try again.')
    }
  }

  handleResourceUpdated = resourceType => {
    return (updatedResource, done) => {
      const resourceList = copy(this.state[resourceType])
      const resource = resourceList.find(r => r.metadata.name === updatedResource.metadata.name)
      resource.status = updatedResource.status
      this.setState({ [resourceType]: resourceList }, done)
    }
  }

  handleResourceDeleted = resourceType => {
    return (name, done) => {
      const resourceList = copy(this.state[resourceType])
      const resource = resourceList.find(r => r.metadata.name === name)
      resource.deleted = true

      this.setState({ [resourceType]: resourceList }, () => {
        this.props.getServiceCount && this.props.getServiceCount(this.state.services.filter(s => !s.deleted).length)
        done()
      })
    }
  }

  serviceCredentialList = ({ serviceCredentials }) => {
    return (
      <Row style={{ marginLeft: '50px' }}>
        <Col>
          <Text strong style={{ marginRight: '8px' }}>Bindings: </Text>
          {serviceCredentials.map(serviceCredential => {
            const status = serviceCredential.status.status || 'Pending'
            const created = moment(serviceCredential.metadata.creationTimestamp).fromNow()
            return (
              <span key={serviceCredential.metadata.name} style={{ marginRight: '5px' }}>
                <Tooltip title={`Created ${created}`}>
                  <Tag color={statusColorMap[status] || 'red'}>{serviceCredential.spec.cluster.name}/{serviceCredential.spec.clusterNamespace} {inProgressStatusList.includes(status) ? <Icon type="loading" /> : <Icon type={statusIconMap[status]} />}</Tag>
                </Tooltip>
              </span>
            )
          })}
        </Col>
      </Row>
    )
  }

  render() {
    const { team } = this.props
    const { dataLoading, services, serviceKinds, serviceCredentials } = this.state

    const hasActiveServices =  Boolean(services.filter(c => !c.deleted).length)

    return (
      <>
        <Button type="primary">
          <Link href="/teams/[name]/services/new" as={`/teams/${team.metadata.name}/services/new`}>
            <a>New service</a>
          </Link>
        </Button>

        <Divider />

        {dataLoading ? (
          <Icon type="loading" />
        ) : (
          <>
            {!hasActiveServices && <Paragraph type="secondary">No services found for this team</Paragraph>}

            {services.map((service, idx) => {
              const filteredServiceCredentials = (serviceCredentials || []).filter(sc => sc.spec.service.name === service.metadata.name)
              return (
                <React.Fragment key={service.metadata.name}>
                  <Service
                    team={team.metadata.name}
                    service={service}
                    serviceKind={serviceKinds.find(sk => sk.metadata.name === service.spec.kind)}
                    deleteService={this.deleteService}
                    handleUpdate={this.handleResourceUpdated('services')}
                    handleDelete={this.handleResourceDeleted('services')}
                    refreshMs={10000}
                    propsResourceDataKey="service"
                    resourceApiPath={`/teams/${team.metadata.name}/services/${service.metadata.name}`}
                  />
                  {!service.deleted && filteredServiceCredentials.length > 0 && this.serviceCredentialList({ serviceCredentials: filteredServiceCredentials })}
                  {idx < services.length - 1 && <Divider />}
                </React.Fragment>
              )
            })}
          </>
        )}
      </>
    )
  }

}

export default ServicesTab
