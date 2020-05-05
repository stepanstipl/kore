import React from 'react'
import PropTypes from 'prop-types'
import { Typography, Collapse, Icon } from 'antd'
const { Paragraph, Text } = Typography

import KoreApi from '../../kore-api'
import ResourceStatusTag from '../resources/ResourceStatusTag'

class ComponentStatusTree extends React.Component {
  static propTypes = {
    team: PropTypes.object.isRequired,
    user: PropTypes.object.isRequired,
    component: PropTypes.object.isRequired,
  }

  constructor(props) {
    super(props)
    this.state = {
      component: props.component,
      children: {},
      loading: {},
      openComponents: []
    }
  }
  
  componentDidUpdate(prevProps) {
    // Refresh open children if the component updated/refreshed.
    if (this.props.component !== prevProps.component) {
      this.loadOpenChildren(this.state.openComponents)
    }
  }

  componentChildren = (component) => {
    let children = []
    // Some things call their children components and some conditions, 
    // join them both together for our purposes:
    if (component.status.components) {
      children.push(...component.status.components)
    }
    if (component.status.conditions) {
      children.push(...component.status.conditions)
    }
    return children
  }

  api = null
  loadChildComponentDetails = async (componentName) => {
    if (!this.api) {
      this.api = await KoreApi.client()
    }
    const res = await this.api.GetTeamResource(this.props.team.metadata.name, componentName)
    this.setState({
      children: {
        ...this.state.children,
        [componentName]: res
      },
      loading: {
        ...this.state.loading,
        [componentName]: false
      }
    })
  }

  loadOpenChildren = (openComponents) => {
    this.setState({ openComponents: openComponents })
    for (let x = 0; x < openComponents.length; x++) {
      // No-op if not something we can load state for:
      if (!openComponents[x].match(/^\w+\/[\w_-]+$/)) {
        continue
      }
      this.setState({ loading: { ...this.state.loading, [openComponents[x]]: true } })
      this.loadChildComponentDetails(openComponents[x]).then(() => {})
    }
  }

  compMsg = (componentStatus) => {
    if (componentStatus.detail) {
      return `${componentStatus.message}: ${componentStatus.detail}`
    }
    return componentStatus.message
  }

  render() {
    const { component } = this.props
    const componentChildren = this.componentChildren(component)
    const statusMsg = this.compMsg(component.status)
    return (
      <>
        {statusMsg ? <Paragraph style={{ marginBottom: '20px', textAlign: 'center' }}>{this.compMsg(component.status)}</Paragraph> : null}
        {componentChildren.length > 0 ? (
          <Collapse onChange={(e) => this.loadOpenChildren(e)}>
            {this.componentChildren(component).map((child => (
              <Collapse.Panel header={child.name} key={child.name} extra={(<ResourceStatusTag resourceStatus={child} />)}>
                <Text>{this.compMsg(child)}</Text>
                {(this.state.children[child.name]) ? (
                  <ComponentStatusTree team={this.props.team} user={this.props.user} component={this.state.children[child.name]} />
                ) : (this.state.loading[child.name] ? <Icon type="loading" /> : null)}
              </Collapse.Panel>
            )))}
          </Collapse>
        ) : null}
      </>
    )
  }
}

export default ComponentStatusTree