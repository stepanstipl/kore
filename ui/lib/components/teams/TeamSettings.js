import React from 'react'
import PropTypes from 'prop-types'
import Link from 'next/link'
import Router from 'next/router'
import { Button, Dropdown, Icon, List, Menu, Modal, Typography } from 'antd'
const { Paragraph } = Typography

import KoreApi from '../../kore-api'
import redirect from '../../utils/redirect'
import { errorMessage, successMessage } from '../../utils/message'

class TeamSettings extends React.Component {

  static propTypes = {
    team: PropTypes.object.isRequired,
    teamRemoved: PropTypes.func.isRequired
  }

  deleteTeam = async () => {
    try {
      const team = this.props.team.metadata.name
      await (await KoreApi.client()).RemoveTeam(team)
      this.props.teamRemoved(team)
      successMessage(`Team "${team}" deleted`)
      return redirect({ router: Router, path: '/' })
    } catch (err) {
      if (err.statusCode === 409 && err.dependents) {
        return Modal.warning({
          title: 'The team cannot be deleted',
          content: (
            <>
              <Paragraph strong>Error: {err.message}</Paragraph>
              <List
                size="small"
                dataSource={err.dependents}
                renderItem={d => <List.Item>{d.kind}: {d.name}</List.Item>}
              />
            </>
          ),
          onOk() {}
        })
      }
      console.log('Error deleting team', err)
      errorMessage('Team could not be deleted, please try again later')
    }
  }

  deleteTeamConfirm = () => {
    Modal.confirm({
      title: 'Are you sure you want to delete this team?',
      content: 'This cannot be undone',
      okText: 'Yes',
      okType: 'danger',
      cancelText: 'No',
      onOk: this.deleteTeam
    })
  }

  render() {
    const team = this.props.team
    const menu = (
      <Menu>
        <Menu.Item key="audit">
          <Link href="/teams/[name]/audit" as={`/teams/${team.metadata.name}/audit`}>
            <a>
              <Icon type="table" style={{ marginRight: '5px' }} />
              Team audit viewer
            </a>
          </Link>
        </Menu.Item>
        <Menu.Item key="delete" id="delete_team" className="ant-btn-danger" onClick={this.deleteTeamConfirm}>
          <Icon type="delete" style={{ marginRight: '5px' }} />
          Delete team
        </Menu.Item>
      </Menu>
    )
    return (
      <Dropdown trigger={['click']} overlay={menu}>
        <Button id="team_settings">
          <Icon type="setting" style={{ marginRight: '10px' }} />
          <Icon type="down" />
        </Button>
      </Dropdown>
    )
  }

}

export default TeamSettings