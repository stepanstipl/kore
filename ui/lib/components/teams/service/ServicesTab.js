import React from 'react'
import moment from 'moment'
import PropTypes from 'prop-types'
import { Button, Col, Divider, Drawer, Icon, List, Modal, Row, Tooltip, Tag, Typography } from 'antd'
const { Paragraph, Text } = Typography

import KoreApi from '../../../kore-api'
import copy from '../../../utils/object-copy'
import Service from './Service'
import ServiceBuildForm from './ServiceBuildForm'
import ApplicationServiceForm from './ApplicationServiceForm'
import { inProgressStatusList, statusColorMap, statusIconMap } from '../../../utils/ui-helpers'
import { loadingMessage, errorMessage } from '../../../utils/message'

class ServicesTab extends React.Component {

  static propTypes = {
    team: PropTypes.object.isRequired,
    cluster: PropTypes.object.isRequired,
    serviceType: PropTypes.oneOf(['cloud', 'application']).isRequired,
    getServiceCount: PropTypes.func
  }

  state = {
    dataLoading: true,
    services: [],
    serviceKinds: [],
    serviceCredentials: [],
    createNewService: false
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

      switch (this.props.serviceType) {
      case 'cloud':
        serviceKinds = serviceKinds.items.filter(sk => sk.metadata.labels['kore.appvia.io/platform'] !== 'Kubernetes')
        break
      case 'application':
        serviceKinds = serviceKinds.items.filter(sk => sk.metadata.labels['kore.appvia.io/platform'] === 'Kubernetes')
        break
      }

      services = services.items.filter(s => Boolean(serviceKinds.find(sk => sk.metadata.name === s.spec.kind )) && s.spec.cluster.name === this.props.cluster.metadata.name)
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

  refreshServiceCredentials = async (service) => {
    let serviceCredentials = []
    try {
      const serviceCredentialsResult = await (await KoreApi.client()).ListServiceCredentials(this.props.team.metadata.name, undefined, service.metadata.name)
      serviceCredentials = serviceCredentialsResult.items
    } catch (error) {
      console.error('Failed to get service credentials', error)
    }
    if (serviceCredentials.length > 0) {
      const existingServiceCredentials = copy(this.state.serviceCredentials)
      serviceCredentials.forEach(sc => {
        const found = existingServiceCredentials.find(esc => esc.metadata.name === sc.metadata.name)
        if (found) {
          found.status = sc.status
        } else {
          existingServiceCredentials.push(sc)
        }
      })
      this.setState({ serviceCredentials: existingServiceCredentials })
    }
  }

  handleServiceCreated = async (service) => {
    this.setState((state) => ({
      createNewService: false,
      services: [ ...state.services, service ]
    }), async () => {
      this.props.getServiceCount && this.props.getServiceCount(this.state.services.length)
      await this.refreshServiceCredentials(service)
    })
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
      loadingMessage(`Service deletion requested: ${service.metadata.name}`)
    } catch (err) {
      console.error('Error deleting service', err)
      errorMessage('Error deleting service, please try again.')
    }
  }

  deleteServiceConfirm = async (name, done) => {
    const serviceCredentials = this.state.serviceCredentials.filter(sc => !sc.deleted && sc.spec.service.name === name)
    if (serviceCredentials.length > 0) {
      return Modal.warning({
        title: 'Warning: service cannot be deleted',
        width: 600,
        content: (
          <div>
            <Paragraph strong>The following cluster namespaces currently have access to the service, this access must be removed before the service can be deleted.</Paragraph>
            <List
              size="small"
              dataSource={serviceCredentials}
              renderItem={sc => <List.Item>{sc.spec.cluster.name} / {sc.spec.clusterNamespace}</List.Item>}
            />
          </div>
        ),
        onOk() {}
      })
    }
    await this.deleteService(name, done)
  }

  handleResourceUpdated = (resourceType) => {
    return async (updatedResource, done) => {
      const resourceList = copy(this.state[resourceType])
      const resource = resourceList.find(r => r.metadata.name === updatedResource.metadata.name)
      resource.status = updatedResource.status
      this.setState({ [resourceType]: resourceList }, done)

      if (resourceType === 'services') {
        await this.refreshServiceCredentials(updatedResource)
      }
    }
  }

  handleResourceDeleted = (resourceType) => {
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
          <Text strong style={{ marginRight: '8px' }}>Access: </Text>
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
    const { team, cluster, serviceType } = this.props
    const { dataLoading, services, serviceKinds, serviceCredentials, createNewService } = this.state

    const hasActiveServices =  Boolean(services.filter(c => !c.deleted).length)

    return (
      <>
        <Button type="primary" onClick={() => this.setState({ createNewService: true })}>New {serviceType} service</Button>

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
                    cluster={cluster}
                    service={service}
                    serviceKind={serviceKinds.find(sk => sk.metadata.name === service.spec.kind)}
                    deleteService={this.deleteServiceConfirm}
                    handleUpdate={this.handleResourceUpdated('services')}
                    handleDelete={this.handleResourceDeleted('services')}
                    refreshMs={10000}
                    propsResourceDataKey="service"
                    resourceApiPath={`/teams/${team.metadata.name}/services/${service.metadata.name}`}
                    style={{ paddingTop: 0, paddingBottom: '5px' }}
                  />
                  {!service.deleted && filteredServiceCredentials.length > 0 && this.serviceCredentialList({ serviceCredentials: filteredServiceCredentials })}
                  {!service.deleted && idx < services.length - 1 && <Divider />}
                </React.Fragment>
              )
            })}
          </>
        )}

        <Drawer
          title={`New ${serviceType} service`}
          visible={createNewService}
          onClose={() => this.setState({ createNewService: false })}
          width={900}
        >
          {createNewService && (
            serviceType === 'cloud' ? (
              <ServiceBuildForm
                team={team}
                cluster={cluster}
                teamServices={services}
                handleSubmit={this.handleServiceCreated}
                handleCancel={() => this.setState({ createNewService: false })}
              />
            ) : (
              <ApplicationServiceForm
                team={team}
                cluster={cluster}
                teamServices={this.state.services}
                handleSubmit={this.handleServiceCreated}
                handleCancel={() => this.setState({ createNewService: false })}
              />
            )
          )}
        </Drawer>
      </>
    )
  }

}

export default ServicesTab
