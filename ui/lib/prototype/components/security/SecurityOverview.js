import React from 'react'
import PropTypes from 'prop-types'
import { List, Icon, Card, Button } from 'antd'
import Link from 'next/link'

class SecurityOverview extends React.Component {
  static propTypes = {
    overview: PropTypes.object.isRequired
  }

  static IconText = ({ icon, text }) => (
    <span>
      <Icon type={icon} style={{ marginRight: 8 }} />
      {text}
    </span>
  )

  static levelIcons = {
    'info': { 'icon': 'safety-certificate', 'color': 'LimeGreen' },
    'warn': { 'icon': 'warning', 'color': 'Orange' },
    'critical': { 'icon': 'stop', 'color': 'Red' }
  }
  
  static LevelIcon = ({ level, active=true }) => {
    const color = active ? SecurityOverview.levelIcons[level].color : 'Silver'
    return (
      <Icon type={SecurityOverview.levelIcons[level].icon} style={{ fontSize: '30px', paddingLeft: '10px' }} theme="twoTone" twoToneColor={color} />
    )
  }
  
  static LevelIconText = ({ level, text, active=true }) => {
    const color = active ? SecurityOverview.levelIcons[level].color : 'Silver'
    return (
      <span>
        <Icon type={SecurityOverview.levelIcons[level].icon} theme="twoTone" twoToneColor={color} style={{ marginRight: 8 }} />
        {text}
      </span>
    )
  }
  
  render() {
    return (
      <div>
        <Card title="Overview" style={{ marginBottom: '20px' }}>
          <List style={{ marginBottom: '20px' }}>
            <List.Item>
              <SecurityOverview.LevelIcon level="critical" active={this.props.overview.status.critical > 0} /> You have <span>{this.props.overview.status.critical}</span> critical security issues
            </List.Item>
            <List.Item>
              <SecurityOverview.LevelIcon level="warn" active={this.props.overview.status.warn > 0} /> You have <span>{this.props.overview.status.warn}</span> security warnings
            </List.Item>
            <List.Item>
              <SecurityOverview.LevelIcon level="info" active={this.props.overview.status.info > 0} /> You have <span>{this.props.overview.status.info}</span> informational security messages
            </List.Item>
          </List>
          <Link href="/prototype/security/review">
            <Button type="primary">Full security report</Button>
          </Link>
        </Card>
        <Card title="Teams" style={{ marginBottom: '20px' }}>
          <List
            itemLayout="vertical"
            size="large"
            dataSource={this.props.overview.teamSummary}
            renderItem={item => (
              <List.Item
                key={item.name}
                actions={[
                  <SecurityOverview.LevelIconText level="critical" active={item.status.critical > 0} text={`${item.status.critical} critical issues`} key="list-actions-team-crit" />,
                  <SecurityOverview.LevelIconText level="warn" active={item.status.warn > 0} text={`${item.status.warn} warnings`} key="list-actions-team-warn" />,
                  <SecurityOverview.LevelIconText level="info" active={item.status.info > 0} text={`${item.status.info} informational`} key="list-actions-team-info" />,
                  <SecurityOverview.IconText icon="arrow-right" text="Go to team" key="list-actions-team" />,
                ]}
              >
                <List.Item.Meta
                  title={item.name}
                  description={`Security review for ${item.name}`}
                  avatar={<SecurityOverview.LevelIcon level={item.overallStatus} />}
                />
              </List.Item>
            )}
          >
          </List>
        </Card>

        <Card title="Plans" style={{ marginBottom: '20px' }}>
          <List
            itemLayout="vertical"
            size="large"
            dataSource={this.props.overview.planSummary}
            renderItem={item => (
              <List.Item
                key={item.name}
                actions={[
                  <SecurityOverview.LevelIconText level="critical" active={item.status.critical > 0} text={`${item.status.critical} critical issues`} key="list-actions-team-crit" />,
                  <SecurityOverview.LevelIconText level="warn" active={item.status.warn > 0} text={`${item.status.warn} warnings`} key="list-actions-team-warn" />,
                  <SecurityOverview.LevelIconText level="info" active={item.status.info > 0} text={`${item.status.info} informational`} key="list-actions-team-info" />,
                  <SecurityOverview.IconText icon="arrow-right" text="Go to plan" key="list-actions-plan" />,
                ]}
              >
                <List.Item.Meta
                  title={item.name}
                  description={`Security review for ${item.name}`}
                  avatar={<SecurityOverview.LevelIcon level={item.overallStatus} />}
                />
              </List.Item>
            )}
          >
          </List>
        </Card>
      </div>
    )
  }
}

export default SecurityOverview