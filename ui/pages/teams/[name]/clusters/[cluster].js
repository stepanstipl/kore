import React from 'react'
import PropTypes from 'prop-types'
import axios from 'axios'
import moment from 'moment'
import { Typography, Collapse, Row, Col, List, Button, Form } from 'antd'
const { Text } = Typography

import KoreApi from '../../../../lib/kore-api'
import Breadcrumb from '../../../../lib/components/layout/Breadcrumb'
import PlanOptionsForm from '../../../../lib/components/plans/PlanOptionsForm'
import ComponentStatusTree from '../../../../lib/components/common/ComponentStatusTree'
import ResourceStatusTag from '../../../../lib/components/resources/ResourceStatusTag'
import { clusterProviderIconSrcMap } from '../../../../lib/utils/ui-helpers'
import copy from '../../../../lib/utils/object-copy'
import FormErrorMessage from '../../../../lib/components/forms/FormErrorMessage'
import { inProgressStatusList } from '../../../../lib/utils/ui-helpers'

class ClusterPage extends React.Component {
  static propTypes = {
    team: PropTypes.object.isRequired,
    user: PropTypes.object.isRequired,
    cluster: PropTypes.object.isRequired
  }

  constructor(props) {
    super(props)
    this.state = {
      cluster: props.cluster,
      components: {},
      editMode: false,
      clusterParams: props.cluster.spec.configuration,
      formErrorMessage: null,
      validationErrors: null
    }
  }

  static getInitialProps = async ctx => {
    const api = await KoreApi.client(ctx)
    const { team, cluster } = await (axios.all([
      api.GetTeam(ctx.query.name), 
      api.GetCluster(ctx.query.name, ctx.query.cluster)
    ]).then(axios.spread((team, cluster) => { 
      return { team, cluster } 
    })))

    if ((!cluster || !team) && ctx.res) {
      /* eslint-disable-next-line require-atomic-updates */
      ctx.res.statusCode = 404
    }
    return { team, cluster }
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

  componentDidMount = () => {
    this.startRefreshing()
  }

  componentWillUnmount = () => {
    if (this.interval) {
      clearInterval(this.interval)
    }
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
      // await this.refreshCluster()
    } catch (err) {
      this.setState({
        saving: false,
        formErrorMessage: (err.fieldErrors && err.message) ? err.message : 'An error occurred saving the cluster, please try again',
        validationErrors: err.fieldErrors // This will be undefined on non-validation errors, which is fine.
      })
    }
  }

  render = () => {
    const { team, user } = this.props
    const { cluster } = this.state
    const created = moment(cluster.metadata.creationTimestamp).fromNow()
    const deleted = cluster.metadata.deletionTimestamp ? moment(cluster.metadata.deletionTimestamp).fromNow() : false
    const clusterNotEditable = !cluster || !cluster.status || inProgressStatusList.includes(cluster.status.status)
    const editClusterFormConfig = {
      layout: 'horizontal', labelAlign: 'left', hideRequiredMark: true,
      labelCol: { xs: 24, xl: 10 }, wrapperCol: { xs: 24, xl: 14 }
    }
  
    return (
      <div>
        <Breadcrumb
          items={[
            { text: team.spec.summary, href: '/teams/[name]', link: `/teams/${team.metadata.name}` },
            { text: `Cluster: ${cluster.metadata.name}` }
          ]}
        />

        <List.Item actions={[<ResourceStatusTag key="status" resourceStatus={cluster.status} />]}>
          <List.Item.Meta
            avatar={<img src={clusterProviderIconSrcMap[cluster.spec.kind]} height="32px" />}
            title={<Text>{cluster.spec.kind} <Text style={{ fontFamily: 'monospace', marginLeft: '15px' }}>{cluster.metadata.name}</Text></Text>}
            description={
              <div>
                <Text type='secondary'>Created {created}</Text>
                {deleted ? <Text type='secondary'><br/>Deleted {deleted}</Text> : null }
              </div>
            }
          />
        </List.Item>

        <Row type="flex" gutter={[16,16]}>
          <Col span={24} xl={12}>
            <Collapse defaultActiveKey={['0']}>
              <Collapse.Panel header="Detailed Cluster Status" extra={(<ResourceStatusTag resourceStatus={cluster.status} />)}>
                <ComponentStatusTree team={team} user={user} component={cluster} />
              </Collapse.Panel>
            </Collapse>
          </Col>
          <Col span={24} xl={12}>
            <Collapse>
              <Collapse.Panel header="Cluster Parameters">
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
                  <PlanOptionsForm
                    team={team}
                    resourceType="cluster"
                    plan={cluster.spec.plan}
                    planValues={this.state.clusterParams}
                    mode={this.state.editMode ? 'edit' : 'view'}
                    validationErrors={this.state.validationErrors}
                    onPlanChange={this.onClusterConfigChanged}
                  />
                </Form>
              </Collapse.Panel>
            </Collapse>
          </Col>
        </Row>
      </div>
    )
  }
}
export default ClusterPage
