import React from 'react'
import { Alert, Tabs, Card, Typography, Tooltip, Icon, Tag, Row, Col, Button, Checkbox } from 'antd'
const { Title, Text } = Typography
import Breadcrumb from '../../../lib/components/Breadcrumb'

class SettingsPage extends React.Component {

  render() {
    const targets = [
      {
        id: 'notprod',
        title: 'Not production',
        description: 'Target covering all but the production environment',
        plans: {
          GKE: [{
            spec: {
              description: 'GKE Development Cluster',
              summary: 'Provides a development cluster within GKE',
            }
          }, {
            spec: {
              description: 'GKE Development Cluster High-CPU',
              summary: 'Provides a development cluster within GKE, with high-CPU nodes for intensive operations',
            }
          }],
          EKS: [{
            spec: {
              description: 'EKS Development Cluster',
              summary: 'Provides a development cluster within EKS',
            }
          }],
          AKS: []
        }
      },
      {
        id: 'prod',
        title: 'Production',
        description: 'Target covering the production environment',
        plans: {
          GKE: [{
            spec: {
              description: 'GKE Production Cluster',
              summary: 'Provides a production cluster within GKE',
            }
          }],
          EKS: [{
            spec: {
              description: 'EKS Production Cluster',
              summary: 'Provides a production cluster within EKS',
            }
          }],
          AKS: []
        }
      }
    ]

    const IconTooltip = ({ icon, text }) => (
      <Tooltip title={text}>
        <Icon type={icon} theme="twoTone" />
      </Tooltip>
    )

    const IconTooltipButton = ({ icon, text, onClick }) => (
      <Tooltip title={text}>
        <a style={{ marginLeft: '5px' }} onClick={onClick}><Icon type={icon} theme="twoTone" /></a>
      </Tooltip>
    )

    const CloudPlans = ({ imageFilename, cloudName, plans }) => (
      <Col span={8}>
        <Card
          title={
            <span>
              <img src={`/static/images/${imageFilename}`} height="20px" style={{ marginRight: '10px' }}/>
              {cloudName}
            </span>
          }
          size="small"
          bordered={false}
        >
          {plans.length === 0 ? <div style={{ padding: '5px 0' }}>No plans</div> : null}
          {plans.map((plan, i) => (
            <div key={i} style={{ padding: '5px 0' }}>
              <Text style={{ marginRight: '10px' }}>{plan.spec.description}</Text>
              <IconTooltip icon="info-circle" text={plan.spec.summary} />
              <IconTooltipButton icon="eye" text="View plan" onClick={() => {}} />
            </div>
          ))}
          <div style={{ padding: '5px 0' }}>
            <a>+ Associate plan</a>
          </div>
        </Card>
      </Col>
    )

    return (
      <>
        <Breadcrumb items={[{ text: 'Configure' }, { text: 'Settings' }]} />
        <Alert
          message="View and configure the settings for Kore"
          type="info"
          style={{ marginBottom: '20px' }}
        />

        <Tabs defaultActiveKey={'teams'} tabPosition="left">
          <Tabs.TabPane tab="Teams" key="teams">
            <Card
              title="Targets"
              extra={<Button type="primary">+ New</Button>}
            >
              <Alert
                message="This setting controls which targets teams are able to create. When a team sets up a target it will create a project within the chosen cloud provider."
                type="info"
                style={{ marginBottom: '20px' }}
              />

              <Checkbox checked={true} disabled={true}>
                <Text strong style={{ marginRight: '10px' }}>Allow teams to use project level credentials</Text>
                <IconTooltip icon="info-circle" text="Upon creation of a target, this will enable the team to choose a project that is already created within the cloud provider, rather than creating a new one. Note, it must be configured by an admin and allocated to teams to be available." />
              </Checkbox>

              {targets.map(target => (
                <Card style={{ marginTop: '20px' }} key={target.id}>
                  <Title level={4} style={{ display: 'inline', marginRight: '10px' }}>{target.title}</Title>
                  <IconTooltip icon="info-circle" text={target.description} />
                  <Tag style={{ marginLeft: '15px' }}>{target.id}</Tag>
                  <Text strong style={{ fontSize: '16px', display: 'block', marginBottom: '5px', marginTop: '10px' }}>Cluster plans</Text>
                  <Row gutter={16}>
                    <CloudPlans imageFilename="GCP.png" cloudName="Google Cloud Platform" plans={target.plans.GKE} />
                    <CloudPlans imageFilename="AWS.png" cloudName="Amazon Web Services" plans={target.plans.EKS} />
                    <CloudPlans imageFilename="Azure.svg" cloudName="Microsoft Azure" plans={target.plans.AKS} />
                  </Row>
                </Card>
              ))}

            </Card>
          </Tabs.TabPane>
          <Tabs.TabPane tab="Setting 2" key="2" />
          <Tabs.TabPane tab="Setting 3" key="3" />
        </Tabs>
      </>
    )
  }

}

export default SettingsPage
