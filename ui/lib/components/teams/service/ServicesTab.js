import React from 'react'
import Link from 'next/link'
import PropTypes from 'prop-types'
import { Button, Divider, Icon, List, message, Typography } from 'antd'
const { Paragraph } = Typography

import KoreApi from '../../../kore-api'
import copy from '../../../utils/object-copy'
import Service from './Service'

class ServicesTab extends React.Component {

  static propTypes = {
    team: PropTypes.object.isRequired,
    getServiceCount: PropTypes.func
  }

  state = {
    dataLoading: true,
    services: []
  }

  async fetchComponentData() {
    try {
      let services = await (await KoreApi.client()).ListServices(this.props.team.metadata.name)
      services = services.items.filter(s => !s.spec.cluster)
      this.props.getServiceCount && this.props.getServiceCount(services.length)
      return { services }
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

  render() {
    const { team } = this.props
    const { dataLoading, services } = this.state

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
            {services.length > 0 && (
              <List
                dataSource={services}
                renderItem={service => {
                  return (
                    <Service
                      team={team.metadata.name}
                      service={service}
                      deleteService={this.deleteService}
                      handleUpdate={this.handleResourceUpdated('services')}
                      handleDelete={this.handleResourceDeleted('services')}
                      refreshMs={10000}
                      propsResourceDataKey="service"
                      resourceApiPath={`/teams/${team.metadata.name}/services/${service.metadata.name}`}
                    />
                  )
                }}
              />
            )}
          </>
        )}
      </>
    )
  }

}

export default ServicesTab
