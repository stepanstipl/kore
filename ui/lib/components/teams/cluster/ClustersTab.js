import React from 'react'
import PropTypes from 'prop-types'
import Link from 'next/link'
import { Badge, Button, Card, Collapse, Divider, Drawer, Icon, List, message, Modal, Typography } from 'antd'
const { Paragraph, Text } = Typography
const { Panel } = Collapse
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

import Cluster from './Cluster'
import NamespaceClaim from '../namespace/NamespaceClaim'
import NamespaceClaimForm from '../namespace/NamespaceClaimForm'
import KoreApi from '../../../kore-api'
import copy from '../../../utils/object-copy'

class ClustersTab extends React.Component {

  static propTypes = {
    team: PropTypes.object.isRequired,
    getClusterCount: PropTypes.func
  }

  state = {
    dataLoading: true,
    clusters: [],
    namespaceClaims: [],
    createNamespace: false
  }

  async fetchComponentData () {
    try {
      const team = this.props.team.metadata.name
      const api = await KoreApi.client()
      let [ clusters, namespaceClaims ] = await Promise.all([
        api.ListClusters(team),
        api.ListNamespaces(team)
      ])
      clusters = clusters.items
      namespaceClaims = namespaceClaims.items
      this.props.getClusterCount && this.props.getClusterCount(clusters.length)
      return { clusters, namespaceClaims }
    } catch (err) {
      console.error('Unable to load data for clusters tab', err)
      return {}
    }
  }

  componentDidMount() {
    return this.fetchComponentData().then(data => {
      this.setState({ ...data, dataLoading: false })
    })
  }

  componentDidUpdate(prevProps) {
    if (prevProps.team.metadata.name !== this.props.team.metadata.name) {
      this.setState({ dataLoading: true, clusters: [] })
      return this.fetchComponentData().then(data => this.setState({ ...data, dataLoading: false }))
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
      this.setState({ [resourceType]: this.state[resourceType].filter(r => r.metadata.name !== name) }, () => {
        this.props.getClusterCount && this.props.getClusterCount(this.state.clusters.length)
        done()
      })
    }
  }

  deleteCluster = async (name, done) => {
    const team = this.props.team.metadata.name
    try {
      const clusters = copy(this.state.clusters)
      const cluster = clusters.find(c => c.metadata.name === name)
      await (await KoreApi.client()).RemoveCluster(team, cluster.metadata.name)
      cluster.status.status = 'Deleting'
      cluster.metadata.deletionTimestamp = new Date()
      this.setState({ clusters }, done)
      message.loading(`Cluster deletion requested: ${cluster.metadata.name}`)
    } catch (err) {
      console.error('Error deleting cluster', err)
      message.error('Error deleting cluster, please try again.')
    }
  }

  createNamespace = value => () => this.setState({ createNamespace: value })

  handleNamespaceCreated = namespaceClaim => {
    this.setState({
      namespaceClaims: this.state.namespaceClaims.concat([namespaceClaim]),
      createNamespace: false
    })
    message.loading(`Namespace "${namespaceClaim.spec.name}" requested on cluster "${namespaceClaim.spec.cluster.name}"`)
  }

  deleteNamespace = async (name, done) => {
    const team = this.props.team.metadata.name
    try {
      const namespaceClaims = copy(this.state.namespaceClaims)
      const namespaceClaim = namespaceClaims.find(nc => nc.metadata.name === name)
      await (await KoreApi.client()).RemoveNamespace(team, namespaceClaim.metadata.name)
      namespaceClaim.status.status = 'Deleting'
      namespaceClaim.metadata.deletionTimestamp = new Date()
      this.setState({ namespaceClaims }, done)
      message.loading(`Namespace deletion requested: ${namespaceClaim.spec.name}`)
    } catch (err) {
      console.error('Error deleting namespace', err)
      message.error('Error deleting namespace, please try again.')
    }
  }

  clusterAccess = async () => {
    const apiUrl = new URL(publicRuntimeConfig.koreApiPublicUrl)

    const profileConfigureCommand = `kore profile configure ${apiUrl.hostname}`
    const loginCommand = 'kore login'
    const kubeconfigCommand = `kore kubeconfig -t ${this.props.team.metadata.name}`

    const InfoItem = ({ num, title }) => (
      <div style={{ marginBottom: '10px' }}>
        <Badge style={{ backgroundColor: '#1890ff', marginRight: '10px' }} count={num} />
        <Text strong>{title}</Text>
      </div>
    )
    Modal.info({
      title: 'Cluster access',
      content: (
        <div style={{ marginTop: '20px' }}>
          <InfoItem num="1" title="Download" />
          <Paragraph>If you haven&apos;t already, download the CLI from <a href="https://github.com/appvia/kore/releases">https://github.com/appvia/kore/releases</a></Paragraph>

          <InfoItem num="2" title="Setup profile" />
          <Paragraph>Create a profile</Paragraph>
          <Paragraph className="copy-command" style={{ marginRight: '40px' }} copyable>{profileConfigureCommand}</Paragraph>
          <Paragraph>Enter the Kore API URL as follows</Paragraph>
          <Paragraph className="copy-command" style={{ marginRight: '40px' }} copyable>{apiUrl.origin}</Paragraph>

          <InfoItem num="3" title="Login" />
          <Paragraph>Login to the CLI</Paragraph>
          <Paragraph className="copy-command" style={{ marginRight: '40px' }} copyable>{loginCommand}</Paragraph>

          <InfoItem num="4" title="Setup access" />
          <Paragraph>Then, you can use the Kore CLI to setup access to your team&apos;s clusters</Paragraph>
          <Paragraph className="copy-command" style={{ marginRight: '40px' }} copyable>{kubeconfigCommand}</Paragraph>
          <Paragraph>This will add local kubernetes configuration to allow you to use <Text
            style={{ fontFamily: 'monospace' }}>kubectl</Text> to talk to the provisioned cluster(s).</Paragraph>
          <Paragraph>See examples: <a href="https://kubernetes.io/docs/reference/kubectl/overview/" target="_blank" rel="noopener noreferrer">https://kubernetes.io/docs/reference/kubectl/overview/</a></Paragraph>
        </div>
      ),
      width: 700,
      onOk() {}
    })
  }

  render() {
    const { team } = this.props
    const { dataLoading, clusters, namespaceClaims, createNamespace } = this.state

    const hasActiveClusters = Boolean(clusters.filter(c => c.status && c.status.status === 'Success').length)

    return (
      <>
        <div>
          <Button type="primary">
            <Link href="/teams/[name]/clusters/new" as={`/teams/${team.metadata.name}/clusters/new`}>
              <a>New cluster</a>
            </Link>
          </Button>
          {hasActiveClusters && <Button style={{ marginLeft: '10px' }} type="primary" onClick={this.createNamespace(true)}>New namespace</Button>}
          {!dataLoading && hasActiveClusters && <Button style={{ float: 'right' }} type="link" onClick={this.clusterAccess}><Icon type="eye" />Access</Button>}
        </div>

        <Divider />

        {dataLoading ? (
          <Icon type="loading" />
        ) : (
          <>
            {clusters.length === 0 && <Paragraph type="secondary">No clusters found for this team</Paragraph>}
            {clusters.map(cluster => {
              const filteredNamespaceClaims = (namespaceClaims || []).filter(nc => nc.spec.cluster.name === cluster.metadata.name && !nc.deleted)
              return (
                <>
                  <Cluster
                    key={cluster.metadata.name}
                    team={team.metadata.name}
                    cluster={cluster}
                    namespaceClaims={filteredNamespaceClaims}
                    handleCreateNamespace={this.createNamespace(true)}
                    deleteCluster={this.deleteCluster}
                    handleUpdate={this.handleResourceUpdated('clusters')}
                    handleDelete={this.handleResourceDeleted('clusters')}
                    refreshMs={10000}
                    propsResourceDataKey="cluster"
                    resourceApiPath={`/teams/${team.metadata.name}/clusters/${cluster.metadata.name}`}
                  />
                  {filteredNamespaceClaims.length > 0 ? (
                    <>
                      <Collapse style={{ marginLeft: '50px' }} defaultActiveKey={['namespaces']}>
                        <Panel header={<span>Namespaces <Badge showZero={true} style={{ marginLeft: '10px', backgroundColor: '#1890ff' }} count={filteredNamespaceClaims.length} /></span>} key="namespaces">
                          <List
                            size="small"
                            style={{ marginTop: 0, marginBottom: 0 }}
                            dataSource={filteredNamespaceClaims}
                            renderItem={namespaceClaim =>
                              <NamespaceClaim
                                team={team.metadata.name}
                                namespaceClaim={namespaceClaim}
                                deleteNamespace={this.deleteNamespace}
                                handleUpdate={this.handleResourceUpdated('namespaceClaims')}
                                handleDelete={this.handleResourceDeleted('namespaceClaims')}
                                refreshMs={15000}
                                propsResourceDataKey="namespaceClaim"
                                resourceApiPath={`/teams/${team.metadata.name}/namespaceclaims/${namespaceClaim.metadata.name}`}
                              />
                            }
                          >
                          </List>
                        </Panel>
                      </Collapse>
                    </>
                  ) : (
                    <Paragraph style={{ marginLeft: '50px' }}>No namespaces</Paragraph>
                  )}
                  <Divider />
                </>
              )
            })}

            {hasActiveClusters && (
              <Drawer
                title="Create namespace"
                placement="right"
                closable={false}
                onClose={this.createNamespace(false)}
                visible={createNamespace}
                width={700}
              >
                <NamespaceClaimForm team={team.metadata.name} clusters={clusters} handleSubmit={this.handleNamespaceCreated} handleCancel={this.createNamespace(false)}/>
              </Drawer>
            )}

          </>
        )}
      </>
    )
  }
}

export default ClustersTab
