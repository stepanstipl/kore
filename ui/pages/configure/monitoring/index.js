import React from 'react'
import PropTypes from 'prop-types'
import Link from 'next/link'
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()
import { Alert, Button, Card, Col, Collapse, Divider, List, Popconfirm, Row, Typography } from 'antd'
const { Paragraph, Text } = Typography
import JsYaml from 'js-yaml'
import moment from 'moment'

import Breadcrumb from '../../../lib/components/layout/Breadcrumb'
import UsePlanForm from '../../../lib/components/plans/UsePlanForm'
import KoreApi from '../../../lib/kore-api'
import ResourceStatusTag from '../../../lib/components/resources/ResourceStatusTag'
import ComponentStatusTree from '../../../lib/components/common/ComponentStatusTree'
import { errorMessage, loadingMessage, successMessage } from '../../../lib/utils/message'
import V1Service from '../../../lib/kore-api/model/V1Service'
import V1ServiceSpec from '../../../lib/kore-api/model/V1ServiceSpec'
import copy from '../../../lib/utils/object-copy'
import FormErrorMessage from '../../../lib/components/forms/FormErrorMessage'

export default class ConfigureMonitoringPage extends React.Component {

  static propTypes = {
    user: PropTypes.object.isRequired,
    serviceKind: PropTypes.object.isRequired,
    servicePlan: PropTypes.object.isRequired,
    team: PropTypes.object.isRequired,
    cluster: PropTypes.object.isRequired,
    service: PropTypes.object
  }

  static KORE_CLUSTER_NAME = 'kore'
  static KORE_CLUSTER_NAMESPACE = 'kore'
  static KORE_MONITORING_SERVICE_KIND = 'helm-app'
  static KORE_MONITORING_SERVICE_PLAN_NAME = 'helm-app-kore-monitoring'
  static KORE_MONITORING_SERVICE_NAME = 'kore-monitoring'

  state = {
    service: this.props.service,
    serviceConfig: (this.props.service && this.props.service.spec.configuration.values) ? JsYaml.safeLoad(this.props.service.spec.configuration.values) : {},
    submitting: false,
    validationErrors: [],
    formErrorMessage: false
  }

  static staticProps = {
    title: 'Configure monitoring',
    adminOnly: true
  }

  static getInitialProps = async (ctx) => {
    const api = await KoreApi.client(ctx)
    try {
      const [ serviceKind, servicePlan, team, cluster, service ] = await Promise.all([
        api.GetServiceKind(ConfigureMonitoringPage.KORE_MONITORING_SERVICE_KIND),
        api.GetServicePlan(ConfigureMonitoringPage.KORE_MONITORING_SERVICE_PLAN_NAME),
        api.GetTeam(publicRuntimeConfig.koreAdminTeamName),
        api.GetCluster(publicRuntimeConfig.koreAdminTeamName, ConfigureMonitoringPage.KORE_CLUSTER_NAME),
        api.GetService(publicRuntimeConfig.koreAdminTeamName, ConfigureMonitoringPage.KORE_MONITORING_SERVICE_NAME)
      ])
      return { serviceKind, servicePlan, team, cluster, service }
    } catch (err) {
      console.log('Error getting data for configure monitoring page', err)
    }
  }

  refresh = async () => {
    const service = await (await KoreApi.client()).GetService(publicRuntimeConfig.koreAdminTeamName, ConfigureMonitoringPage.KORE_MONITORING_SERVICE_NAME)
    this.setState(state => ({ service, serviceConfig: !service ? {} : state.serviceConfig }))
    if (!service || service.status.status === 'Success') {
      clearInterval(this.interval)
    }
  }

  startRefreshing = () => {
    if (this.state.service) {
      this.interval = setInterval(this.refresh, 5000)
    }
  }

  componentDidMount() {
    this.startRefreshing()
  }

  componentWillUnmount() {
    if (this.interval) {
      clearInterval(this.interval)
    }
  }

  generateServiceResource = () => {
    //const cluster = this.props.cluster
    const service = this.state.service

    const serviceResource = new V1Service()
    serviceResource.setApiVersion('services.kore.appvia.io/v1')
    serviceResource.setKind('Service')

    const meta = {}
    //const meta = NewV1ObjectMeta(ConfigureMonitoringPage.KORE_MONITORING_SERVICE_NAME, publicRuntimeConfig.koreAdminTeamName)
    if (service) {
      meta.setResourceVersion(service.metadata.resourceVersion)
    }
    serviceResource.setMetadata(meta)

    const serviceSpec = new V1ServiceSpec()
    serviceSpec.setKind(ConfigureMonitoringPage.KORE_MONITORING_SERVICE_KIND)
    serviceSpec.setPlan(ConfigureMonitoringPage.KORE_MONITORING_SERVICE_PLAN_NAME)
    /*
    serviceSpec.setCluster(NewV1Ownership({
      group: cluster.apiVersion.split('/')[0],
      version: cluster.apiVersion.split('/')[1],
      kind: cluster.kind,
      name: cluster.metadata.name,
      namespace: publicRuntimeConfig.koreAdminTeamName
    }))
    */
    serviceSpec.setClusterNamespace(ConfigureMonitoringPage.KORE_CLUSTER_NAMESPACE)
    const config = { ...this.props.servicePlan.spec.configuration }
    if (Object.keys(this.state.serviceConfig).length > 0) {
      config.values = JsYaml.safeDump(this.state.serviceConfig)
    }
    serviceSpec.setConfiguration(config)
    serviceResource.setSpec(serviceSpec)

    return serviceResource
  }

  saveConfiguration = async () => {
    this.setState({ submitting: true })
    const team = publicRuntimeConfig.koreAdminTeamName
    const serviceName = ConfigureMonitoringPage.KORE_MONITORING_SERVICE_NAME
    try {
      const service = await (await KoreApi.client()).UpdateService(team, serviceName, this.generateServiceResource())
      if (!this.state.service) {
        loadingMessage('Monitoring configuration being applied...')
      } else {
        successMessage('Monitoring configuration applied')
      }
      this.setState({ submitting: false, service })
      this.startRefreshing()
    } catch (err) {
      console.error('Error saving monitoring configuration', err)
      this.setState({
        submitting: false,
        formErrorMessage: (err.fieldErrors && err.message) ? err.message : 'An error occurred requesting the service, please try again',
        validationErrors: err.fieldErrors // This will be undefined on non-validation errors, which is fine.
      })
    }
  }

  disableMonitoring = async() => {
    const team = publicRuntimeConfig.koreAdminTeamName
    const serviceName = ConfigureMonitoringPage.KORE_MONITORING_SERVICE_NAME
    try {
      await (await KoreApi.client()).DeleteService(team, serviceName)
      this.setState((state) => {
        const service = copy(state.service)
        service.status.status = 'Deleting'
        return { service }
      }, () => this.startRefreshing())
      loadingMessage('Monitoring being disabled...')
    } catch (err) {
      console.error('Error deleting service', err)
      errorMessage('An error occurred disabling monitoring, please try again')
    }
  }

  render() {
    const { user, serviceKind, team, cluster } = this.props
    const { service, submitting, formErrorMessage, serviceConfig } = this.state

    const saveButtonText = service ? 'Save configuration' : 'Enable & save configuration'

    return (
      <>
        <Breadcrumb items={[{ text: 'Configure' }, { text: 'Monitoring' }]}/>

        {serviceKind.spec.enabled ? (
          <>
            <Row gutter={16}>
              <Col span={12}>
                <List.Item actions={service && service.status.status === 'Success' ? [
                  <Popconfirm
                    key="delete"
                    title="Are you sure you want to disable monitoring?"
                    onConfirm={this.disableMonitoring}
                    okText="Yes"
                    cancelText="No"
                  >
                    <a style={{ textDecoration: 'underline' }}>Disable</a>
                  </Popconfirm>] : []}>
                  <List.Item.Meta
                    className="large-list-item"
                    title={<Text>Monitoring</Text>}
                    description={service ? `Created ${moment(service.metadata.creationTimestamp).fromNow()}` : 'Not configured'}
                  />
                </List.Item>
              </Col>
              <Col span={12}>
                {service ? (
                  <Collapse style={{ marginTop: '12px' }}>
                    <Collapse.Panel header="Status" extra={(<ResourceStatusTag resourceStatus={service.status} />)}>
                      <ComponentStatusTree team={team} user={user} component={service} />
                    </Collapse.Panel>
                  </Collapse>
                ) : null}
              </Col>
            </Row>

            <Divider />

            <Card title="Configuration">

              <FormErrorMessage message={formErrorMessage} />

              <UsePlanForm
                team={team}
                cluster={cluster}
                resourceType="monitoring"
                kind={ConfigureMonitoringPage.KORE_MONITORING_SERVICE_KIND}
                plan={ConfigureMonitoringPage.KORE_MONITORING_SERVICE_PLAN_NAME}
                planValues={serviceConfig}
                validationErrors={this.state.validationErrors}
                onPlanValuesChange={(serviceConfig) => this.setState({ serviceConfig })}
                mode={service ? 'edit' : 'create'}
              />

              <Button type="primary" loading={submitting} disabled={service && service.status.status !== 'Success'} style={{ display: 'block', marginTop: '20px' }} onClick={this.saveConfiguration}>{saveButtonText}</Button>

            </Card>
          </>
        ) : (
          <Alert
            message="Service not enabled"
            description={<>
              <Paragraph>The {serviceKind.spec.displayName} service must be enabled to configured the monitoring, this can be enabled on the Configure Services page.</Paragraph>
              <Paragraph style={{ marginBottom: 0 }}>
                <Link href="/configure/services">
                  <Button>Go to Configure Services</Button>
                </Link>
              </Paragraph>
            </>}
            type="warning"
            showIcon
          />
        )}

      </>
    )
  }
}
