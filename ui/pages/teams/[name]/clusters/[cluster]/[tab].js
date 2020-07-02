import React from 'react'
import PropTypes from 'prop-types'
import Router from 'next/router'
import moment from 'moment'
import {
  Typography,
  Collapse,
  Row,
  Col,
  List,
  Button,
  Form,
  Tabs
} from 'antd'
const { Text } = Typography
const { TabPane } = Tabs

import KoreApi from '../../../../../lib/kore-api'
import TeamHeader from '../../../../../lib/components/teams/TeamHeader'
import UsePlanForm from '../../../../../lib/components/plans/UsePlanForm'
import ComponentStatusTree from '../../../../../lib/components/common/ComponentStatusTree'
import ResourceStatusTag from '../../../../../lib/components/resources/ResourceStatusTag'
import { clusterProviderIconSrcMap } from '../../../../../lib/utils/ui-helpers'
import copy from '../../../../../lib/utils/object-copy'
import { featureEnabled, KoreFeatures } from '../../../../../lib/utils/features'
import FormErrorMessage from '../../../../../lib/components/forms/FormErrorMessage'
import { inProgressStatusList } from '../../../../../lib/utils/ui-helpers'
import ClusterAccessInfo from '../../../../../lib/components/teams/cluster/ClusterAccessInfo'
import { isReadOnlyCRD } from '../../../../../lib/utils/crd-helpers'
import ServicesTab from '../../../../../lib/components/teams/service/ServicesTab'
import NamespacesTab from '../../../../../lib/components/teams/namespace/NamespacesTab'
import TextWithCount from '../../../../../lib/components/utils/TextWithCount'

class ClusterPage extends React.Component {
  static propTypes = {
    team: PropTypes.object.isRequired,
    user: PropTypes.object.isRequired,
    cluster: PropTypes.object.isRequired,
    tabActiveKey: PropTypes.string,
    teamRemoved: PropTypes.func.isRequired
  }

  constructor(props) {
    super(props)
    this.state = {
      tabActiveKey: this.props.tabActiveKey || 'namespaces',
      cluster: props.cluster,
      components: {},
      editMode: false,
      clusterParams: props.cluster.spec.configuration,
      formErrorMessage: null,
      validationErrors: null
    }
  }

  static getInitialProps = async (ctx) => {
    const api = await KoreApi.client(ctx)
    let [ team, cluster ] = await Promise.all([
      api.GetTeam(ctx.query.name),
      api.GetCluster(ctx.query.name, ctx.query.cluster)
    ])
    if ((!cluster || !team) && ctx.res) {
      /* eslint-disable-next-line require-atomic-updates */
      ctx.res.statusCode = 404
    }
    const tabActiveKey = ctx.query.tab || 'namespaces'
    return { team, cluster, tabActiveKey }
  }

  fetchCommonData = async () => {
    if (featureEnabled(KoreFeatures.SERVICES)) {
      const serviceKinds = await (await KoreApi.client()).ListServiceKinds()
      return { serviceKinds: serviceKinds.items }
    }

    return {}
  }

  componentDidMount() {
    this.startRefreshing()
    this.fetchCommonData().then(data => {
      this.setState({ ...data })
    })
  }

  componentDidUpdate() {
    if (this.state.tabActiveKey !== this.props.tabActiveKey) {
      this.setState({ tabActiveKey: this.props.tabActiveKey })
    }
  }

  componentWillUnmount() {
    if (this.interval) {
      clearInterval(this.interval)
    }
  }

  interval = null
  api = null
  startRefreshing = async () => {
    this.api = await KoreApi.client()
    this.interval = setInterval(async () => {
      await this.refreshCluster()
    }, 5000)
  }

  refreshCluster = async () => {
    const cluster = await this.api.GetCluster(this.props.team.metadata.name, this.state.cluster.metadata.name)
    if (cluster) {
      this.setState({
        cluster: cluster,
        // Keep the params up to date with the cluster, unless we're in edit mode.
        clusterParams: this.state.editMode ? this.state.clusterParams : copy(cluster.spec.configuration)
      })
    } else {
      this.setState({ cluster: { ...this.state.cluster, deleted: true } })
    }
  }

  handleTabChange = (key) => {
    const team = this.props.team.metadata.name
    const cluster = this.props.cluster.metadata.name
    Router.push('/teams/[name]/clusters/[cluster]/[tab]', `/teams/${team}/clusters/${cluster}/${key}`)
  }

  onClusterConfigChanged = (updatedClusterParams) => {
    this.setState({
      clusterParams: updatedClusterParams
    })
  }

  onEditClick = (e) => {
    e.stopPropagation()
    this.setState({ editMode: true })
  }

  onCancelClick = (e) => {
    e.stopPropagation()
    this.setState({
      editMode: false,
      clusterParams: copy(this.state.cluster.spec.configuration)
    })
  }

  onSubmit = async (e) => {
    e.preventDefault()
    this.setState({ saving: true, validationErrors: null, formErrorMessage: null })
    const clusterUpdated = copy(this.state.cluster)
    clusterUpdated.spec.configuration = this.state.clusterParams
    try {
      await this.api.UpdateCluster(this.props.team.metadata.name, this.state.cluster.metadata.name, clusterUpdated)
      this.setState({
        cluster: { ...this.state.cluster, status: { ...this.state.cluster.status, status: 'Pending' } },
        saving: false,
        validationErrors: null,
        formErrorMessage: null,
        editMode: false
      })
    } catch (err) {
      this.setState({
        saving: false,
        formErrorMessage: (err.fieldErrors && err.message) ? err.message : 'An error occurred saving the cluster, please try again',
        validationErrors: err.fieldErrors // This will be undefined on non-validation errors, which is fine.
      })
    }
  }

  render = () => {
    const { team, user, teamRemoved } = this.props
    const { cluster } = this.state

    const created = moment(cluster.metadata.creationTimestamp).fromNow()
    const deleted = cluster.metadata.deletionTimestamp ? moment(cluster.metadata.deletionTimestamp).fromNow() : false
    const clusterNotEditable = !cluster || isReadOnlyCRD(cluster) || !cluster.status || inProgressStatusList.includes(cluster.status.status)
    const editClusterFormConfig = {
      layout: 'horizontal', labelAlign: 'left', hideRequiredMark: true,
      labelCol: { xs: 24, xl: 10 }, wrapperCol: { xs: 24, xl: 14 }
    }

    return (
      <>
        <TeamHeader team={team} breadcrumbExt={[
          { text: 'Clusters', href: '/teams/[name]/[tab]', link: `/teams/${team.metadata.name}/clusters` },
          { text: cluster.metadata.name }
        ]} teamRemoved={teamRemoved} />

        <Row type="flex" gutter={[16,16]} style={{ marginTop: '-20px' }}>
          <Col span={24} xl={12}>
            <List.Item>
              <List.Item.Meta
                className="large-list-item"
                avatar={<img src={clusterProviderIconSrcMap[cluster.spec.kind]} />}
                title={<Text id="cluster_name" style={{ marginTop: '15px', display: 'block' }}>{cluster.metadata.name}</Text>}
                description={
                  <div>
                    <Text type='secondary'>Created {created}</Text>
                    {deleted ? <Text type='secondary'><br/>Deleted {deleted}</Text> : null }
                  </div>
                }
              />
              <div>
                <ClusterAccessInfo team={this.props.team} />
              </div>
            </List.Item>
          </Col>
          <Col span={24} xl={12} style={{ marginTop: '14px' }}>
            <Collapse style={{ marginTop: '12px', marginBottom: '20px' }}>
              <Collapse.Panel header="Detailed Cluster Status" extra={(<ResourceStatusTag id="cluster_status" resourceStatus={cluster.status} />)}>
                <ComponentStatusTree team={team} user={user} component={cluster} />
              </Collapse.Panel>
            </Collapse>
          </Col>
        </Row>

        <Tabs activeKey={this.state.tabActiveKey} onChange={(key) => this.handleTabChange(key)} tabBarStyle={{ marginBottom: '20px' }}>
          <TabPane key="namespaces" tab={<TextWithCount title="Namespaces" count={this.state.namespaceCount} />} forceRender={true}>
            <NamespacesTab user={this.props.user} team={this.props.team} cluster={this.props.cluster} onNamespaceCountChange={(count) => this.setState({ namespaceCount: count })} />
          </TabPane>

          {!featureEnabled(KoreFeatures.SERVICES) ? null : (
            <TabPane key="services" tab={<TextWithCount title="Cloud services" count={this.state.cloudServiceCount} />} forceRender={true}>
              <ServicesTab user={this.props.user} team={this.props.team} cluster={this.props.cluster} serviceType="cloud" getServiceCount={(count) => this.setState({ cloudServiceCount: count })} />
            </TabPane>
          )}

          {!featureEnabled(KoreFeatures.APPLICATION_SERVICES) ? null : (
            <TabPane key="application-services" tab={<TextWithCount title="Application services" count={this.state.applicationServiceCount} />} forceRender={true}>
              <ServicesTab user={this.props.user} team={this.props.team} cluster={this.props.cluster} serviceType="application" getServiceCount={(count) => this.setState({ applicationServiceCount: count })} />
            </TabPane>
          )}

          <TabPane key="settings" tab="Settings" forceRender={true}>
            <Form {...editClusterFormConfig} onSubmit={(e) => this.onSubmit(e)}>
              <FormErrorMessage message={this.state.formErrorMessage} />
              <Form.Item label="" colon={false}>
                {!this.state.editMode ? (
                  <Button icon="edit" htmlType="button" disabled={clusterNotEditable} onClick={(e) => this.onEditClick(e)}>Edit</Button>
                ) : (
                  <>
                    <Button type="primary" icon="save" htmlType="submit" loading={this.state.saving} disabled={this.state.saving || clusterNotEditable}>Save</Button>
                    &nbsp;
                    <Button icon="stop" htmlType="button" onClick={(e) => this.onCancelClick(e)}>Cancel</Button>
                  </>
                )}
              </Form.Item>
              <UsePlanForm
                team={team}
                cluster={cluster}
                resourceType="cluster"
                kind={cluster.spec.kind}
                plan={cluster.spec.plan}
                planValues={this.state.clusterParams}
                mode={this.state.editMode ? 'edit' : 'view'}
                validationErrors={this.state.validationErrors}
                onPlanValuesChange={this.onClusterConfigChanged}
              />
            </Form>
          </TabPane>
        </Tabs>
      </>
    )
  }
}
export default ClusterPage
