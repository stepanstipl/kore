import Link from 'next/link'
import PropTypes from 'prop-types'
import { Layout, Menu, Icon } from 'antd'
const { Sider } = Layout
const { SubMenu } = Menu

import { featureEnabled, KoreFeatures } from '../../utils/features'

class SiderMenu extends React.Component {
  static propTypes = {
    hide: PropTypes.bool.isRequired,
    isAdmin: PropTypes.bool.isRequired,
    userTeams: PropTypes.array,
    otherTeams: PropTypes.array
  }

  state = {
    siderCollapsed: false
  }

  onSiderCollapse = siderCollapsed => {
    this.setState({
      siderCollapsed
    })
  }

  render() {
    const { hide, isAdmin, userTeams, otherTeams } = this.props
    const { siderCollapsed } = this.state

    if (hide) {
      return null
    }

    const menuItem = ({ key, text, href, link, icon }) => (
      <Menu.Item key={key} style={{ margin: '0' }}>
        <Link href={href || link} as={link}>
          <a className="collapsed"><Icon type={icon} />{text}</a>
        </Link>
      </Menu.Item>
    )

    const AdminMenu = () => isAdmin ? (
      <SubMenu key="configure"
        title={
          <span>
            <Icon type="tool" />
            <span>Configure</span>
          </span>
        }
      >
        {menuItem({ key: 'configure_cloud', text: 'Cloud', href: '/configure/cloud/[[...cloud]]', link: '/configure/cloud', icon: 'cloud' })}
        {!featureEnabled(KoreFeatures.APPLICATION_SERVICES) ? null :
          menuItem({ key: 'configure_services', text: 'Services', link: '/configure/services', icon: 'cloud-server' })
        }
        {menuItem({ key: 'configure_users', text: 'Users', link: '/configure/users', icon: 'user' })}
      </SubMenu>
    ) : null

    const AuditMenu = () => isAdmin ? (
      <SubMenu key="audit"
        title={
          <span>
            <Icon type="file-protect" theme="outlined" />
            <span>Audit</span>
          </span>
        }
      >
        {menuItem({ key: 'audit_log', text: 'Events', link: '/audit', icon: 'table' })}
      </SubMenu>
    ) : null

    const MonitoringMenu = () => isAdmin ? (
      <SubMenu key="monitoring"
        title={
          <span>
            <Icon type="file-protect" theme="outlined" />
            <span>Monitoring</span>
          </span>
        }
      >
        {menuItem({ key: 'alerts', text: 'Overview', link: '/monitoring/overview', icon: 'global' })}
        {menuItem({ key: 'rules', text: 'Rules', link: '/monitoring/rules', icon: 'schedule' })}
      </SubMenu>
    ) : null

    return (
      <Sider
        id="sider"
        collapsible
        collapsed={siderCollapsed}
        onCollapse={this.onSiderCollapse}
        width="235"
      >
        <Menu defaultOpenKeys={['configure', 'teams', 'spaces']}  mode="inline">
          <SubMenu key="teams"
            title={
              <span>
                <Icon type="team" />
                <span>Teams</span>
              </span>
            }
          >
            {menuItem({ key: 'new_team', text: 'New team', link: '/teams/new', icon: 'plus-circle' })}
            {(userTeams).concat(otherTeams).map(t => (
              menuItem({ key: t.metadata.name, text: t.spec.summary, href: '/teams/[name]/[tab]', link: `/teams/${t.metadata.name}/clusters`, icon: 'team' })
            ))}
          </SubMenu>
          {AdminMenu()}
          {!featureEnabled(KoreFeatures.MONITORING_SERVICES) ? null :
            MonitoringMenu()
          }
          {AuditMenu()}
          {isAdmin ? (
            <SubMenu key="security"
              title={
                <span>
                  <Icon type="lock" theme="outlined" />
                  <span>Security</span>
                </span>
              }
            >
              {menuItem({ key: 'overview', text: 'Overview', link: '/security', icon: 'global' })}
              {menuItem({ key: 'rules', text: 'Rule Reference', link: '/security/rules', icon: 'schedule' })}
            </SubMenu>
          ) : null}
        </Menu>
      </Sider>
    )
  }
}

export default SiderMenu
