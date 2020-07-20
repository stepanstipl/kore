import React from 'react'
import PropTypes from 'prop-types'
import Link from 'next/link'
import { Typography, Statistic, Icon, Row, Col, Card, Alert, Button, Tag } from 'antd'
const { Title, Paragraph, Text } = Typography
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

import KoreApi from '../lib/kore-api'

class IndexPage extends React.Component {
  static propTypes = {
    user: PropTypes.object.isRequired,
    allUsers: PropTypes.object.isRequired,
    allTeams: PropTypes.object.isRequired,
    adminMembers: PropTypes.object,
    gkeCredentials: PropTypes.object,
    gcpOrganizations: PropTypes.object,
    eksCredentials: PropTypes.object,
    awsOrganizations: PropTypes.object,
    aksCredentials: PropTypes.object,
    version: PropTypes.string
  }

  static staticProps = {
    title: 'Appvia Kore Dashboard'
  }

  static async getPageData(ctx) {
    const { user } = ctx
    const api = await KoreApi.client(ctx)

    let allTeams
    let allUsers
    let adminMembers
    let gkeCredentials
    let gcpOrganizations
    let eksCredentials
    let awsOrganizations
    let aksCredentials

    if (user.isAdmin) {
      [ allTeams, allUsers, adminMembers, gkeCredentials, gcpOrganizations, eksCredentials, awsOrganizations, aksCredentials ] = await Promise.all([
        api.ListTeams(),
        api.ListUsers(),
        api.ListTeamMembers(publicRuntimeConfig.koreAdminTeamName),
        api.ListGKECredentials(publicRuntimeConfig.koreAdminTeamName),
        api.ListGCPOrganizations(publicRuntimeConfig.koreAdminTeamName),
        api.ListEKSCredentials(publicRuntimeConfig.koreAdminTeamName),
        api.ListAWSOrganizations(publicRuntimeConfig.koreAdminTeamName),
        api.ListAKSCredentials(publicRuntimeConfig.koreAdminTeamName),
      ])
    } else {
      [ allTeams, allUsers ] = await Promise.all([
        api.ListTeams(),
        api.ListUsers(),
      ])
    }

    allTeams.items = (allTeams.items || []).filter(t => !publicRuntimeConfig.ignoreTeams.includes(t.metadata.name))
    return { allTeams, allUsers, adminMembers, gkeCredentials, gcpOrganizations, eksCredentials, awsOrganizations, aksCredentials }
  }

  static getInitialProps = async (ctx) => {
    const data = await IndexPage.getPageData(ctx)
    return data
  }

  render() {
    const { user, allTeams, allUsers, adminMembers, gkeCredentials, gcpOrganizations, eksCredentials, awsOrganizations, aksCredentials, version } = this.props
    const userTeams = (user.teams.userTeams || []).filter(t => !publicRuntimeConfig.ignoreTeams.includes(t.metadata.name))
    const noUserTeamsExist = userTeams.length === 0
    const gcpCredsMissing = (gkeCredentials && gkeCredentials.items.length === 0) && (gcpOrganizations && gcpOrganizations.items.length === 0)
    const awsCredsMissing = eksCredentials && eksCredentials.items.length === 0 && (awsOrganizations && awsOrganizations.items.length === 0)
    const cloudIntegrationMissing = gcpCredsMissing && awsCredsMissing

    const NoTeamInfoAlert = () => noUserTeamsExist ? (
      <Alert
        message="You are not part of a team"
        description={
          <div>
            <Paragraph style={{ marginTop: '10px' }}>Teams are everything in Kore, we recommend creating a team now to get started.</Paragraph>
            <Button type="secondary">
              <Link href="/teams/new">
                <a>Create a new team</a>
              </Link>
            </Button>
          </div>
        }
        type="info"
        showIcon
        style={{ marginTop: '30px' }}
      />
    ) : null

    const CloudIntegrationWarning = () => cloudIntegrationMissing ? (
      <Alert
        message="No cloud access configured"
        description={
          <div>
            <Paragraph style={{ marginTop: '10px' }}>Without Cloud provider access Kore will be unable to create clusters for teams.</Paragraph>
            <Button type="secondary">
              <Link href="/setup/kore/cloud-access">
                <a>Go to cloud access setup</a>
              </Link>
            </Button>
          </div>
        }
        type="warning"
        showIcon
        style={{ marginTop: '30px' }}
      />
    ) : null

    const TeamStats = () => (
      <Card title="Teams" extra={<Icon type="team" />} bordered={false}>
        <Row gutter={16}>
          <Col span={12}>
            <Statistic style={{ textAlign: 'center' }} title="Yours" value={userTeams.length} valueStyle={{ color: noUserTeamsExist ? 'orange' : '' }} />
          </Col>
          <Col span={12}>
            <Statistic style={{ textAlign: 'center' }} title="Total" value={allTeams.items.length} />
          </Col>
        </Row>
      </Card>
    )

    const UserStats = () => (
      <Card title="Users" extra={<Icon type="user" />} bordered={false}>
        <Row gutter={16}>
          {user.isAdmin ? (
            <div>
              <Col span={12}>
                <Statistic style={{ textAlign: 'center' }} title="Total" value={allUsers.items.length} />
              </Col>
              <Col span={12}>
                <Statistic style={{ textAlign: 'center' }} title="Admins" value={adminMembers.items.length} />
              </Col>
            </div>
          ) : (
            <Col span={24}>
              <Statistic style={{ textAlign: 'center' }} title="Total" value={allUsers.items.length} />
            </Col>
          )}
        </Row>
      </Card>
    )

    const AdminView = () => (
      <div>
        <CloudIntegrationWarning/>
        <NoTeamInfoAlert />
        <Row gutter={16} type="flex" style={{ marginTop: '40px', marginBottom: '40px' }}>
          <Col span={12} xl={4}>
            <TeamStats />
          </Col>
          <Col span={12} xl={4}>
            <UserStats />
          </Col>
          <Col span={24} xl={16}>
            <Card title="Cloud provider integrations" extra={<Icon type="cloud" />} bordered={false}>
              <Row gutter={16}>
                <Col span={8} xs={24} xl={8}>
                  <Statistic
                    title={
                      <>
                        <span>
                          <img src="/static/images/GCP.png" height="25px" style={{ marginRight: '5px' }}/>
                          <Text strong>Google Cloud Platform</Text>
                        </span>
                        <br/>
                        <Text style={{ display: 'inline-block', marginTop: '10px' }}>Organization / Projects</Text>
                      </>
                    }
                    valueRender={() => <>{gcpOrganizations.items.length === 0 ? <Icon type="exclamation-circle" theme="twoTone" twoToneColor="orange"  /> : <Icon type="check-circle" theme="twoTone" twoToneColor="#52c41a" />}<Text> / {gkeCredentials.items.length}</Text></>}
                    style={{ textAlign: 'center', marginBottom: '20px' }}
                    valueStyle={{ color: cloudIntegrationMissing ? 'orange' : '' }}
                  />
                </Col>
                <Col span={8} xs={24} xl={8}>
                  <Statistic
                    title={
                      <>
                        <span>
                          <img src="/static/images/AWS.png" height="25px" style={{ marginRight: '5px' }}/>
                          <Text strong>Amazon Web Services</Text>
                        </span>
                        <br/>
                        <Text style={{ display: 'inline-block', marginTop: '10px' }}>Organization / Accounts</Text>
                      </>
                    }
                    valueRender={() => <>{awsOrganizations.items.length === 0 ? <Icon type="exclamation-circle" theme="twoTone" twoToneColor="orange"  /> : <Icon type="check-circle" theme="twoTone" twoToneColor="#52c41a" />}<Text> / {eksCredentials.items.length}</Text></>}
                    style={{ textAlign: 'center', marginBottom: '20px' }}
                    valueStyle={{ color: cloudIntegrationMissing ? 'orange' : '' }}
                  />
                </Col>
                <Col span={8} xs={24} xl={8}>
                  <Statistic
                    title={
                      <>
                        <span>
                          <img src="/static/images/Azure.svg" height="18px" style={{ marginRight: '5px' }}/>
                          <Text strong>Microsoft Azure</Text>
                        </span>
                        <br/>
                        <Text style={{ display: 'inline-block', marginTop: '10px' }}>Subscriptions</Text>
                      </>
                    }
                    value={aksCredentials.items.length}
                    style={{ textAlign: 'center', marginBottom: '20px' }}
                    valueStyle={{ color: cloudIntegrationMissing ? 'orange' : '' }}
                  />
                </Col>
              </Row>
            </Card>
          </Col>
        </Row>
      </div>
    )

    const UserView = () => (
      <div>
        <NoTeamInfoAlert />
        <Row gutter={16} type="flex" style={{ marginTop: '40px', marginBottom: '40px' }}>
          <Col span={8}>
            <TeamStats />
          </Col>
          <Col span={5}>
            <UserStats />
          </Col>
        </Row>
      </div>
    )

    return (
      <div>
        <Tag style={{ float: 'right' }}>{version}</Tag>
        <Title level={1} style={{ marginBottom: '0' }}>Appvia Kore</Title>
        <Title level={4} type="secondary" style={{ marginTop: '10px' }}>Kubernetes for Teams, Making Cloud Simple for Developers and DevOps</Title>
        {user.isAdmin ? <AdminView /> : <UserView />}
      </div>
    )
  }
}

export default IndexPage
