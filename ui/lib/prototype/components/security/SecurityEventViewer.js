import React from 'react'
import PropTypes from 'prop-types'
import { List, Icon } from 'antd'

class SecurityEventViewer extends React.Component {
  static propTypes = {
    items: PropTypes.arrayOf(PropTypes.object).isRequired
  }

  static levelIcons = {
    'info': { 'icon': 'safety-certificate', 'color': 'LimeGreen' },
    'warn': { 'icon': 'warning', 'color': 'Orange' },
    'critical': { 'icon': 'stop', 'color': 'Red' }
  }
  
  static IconText = ({ icon, text }) => (
    <span>
      <Icon type={icon} style={{ marginRight: 8 }} />
      {text}
    </span>
  )

  static LevelIcon = ({ level }) => (
    <Icon type={SecurityEventViewer.levelIcons[level].icon} style={{ fontSize: '30px', paddingLeft: '10px' }} theme="twoTone" twoToneColor={SecurityEventViewer.levelIcons[level].color} />
  )
    
  render() {
    return (
      <List 
        itemLayout="vertical"
        size="large"
        dataSource={this.props.items}
        renderItem={item => (
          <List.Item
            key={item.id}
            actions={[
              <SecurityEventViewer.IconText icon="read" text="View rule" key="list-actions-ruledef" />,
              <SecurityEventViewer.IconText icon="arrow-right" text="Go to resource" key="list-actions-resource" />,
            ]}
          >
            <List.Item.Meta
              title={item.headline}
              description={item.ruleSummary}
              avatar={<SecurityEventViewer.LevelIcon level={item.level} />}
            />
            {item.description}
          </List.Item>
        )}
      />
    )
  }
}

export default SecurityEventViewer