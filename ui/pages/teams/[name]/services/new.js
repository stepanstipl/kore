import React from 'react'
import PropTypes from 'prop-types'
import axios from 'axios'
import { Typography } from 'antd'
const { Title } = Typography

import Breadcrumb from '../../../../lib/components/layout/Breadcrumb'
import ServiceBuildForm from '../../../../lib/components/teams/service/ServiceBuildForm'
import apiRequest from '../../../../lib/utils/api-request'
import apiPaths from '../../../../lib/utils/api-paths'

class NewTeamServicePage extends React.Component {
  static propTypes = {
    user: PropTypes.object.isRequired,
    team: PropTypes.object.isRequired,
    services: PropTypes.object.isRequired
  }

  static staticProps = {
    title: 'New team service'
  }

  static async getPageData({ req, res, query }) {
    const name = query.name
    const getTeam = () => apiRequest({ req, res }, 'get', apiPaths.team(name).self)
    const getTeamServices = () => apiRequest({ req, res }, 'get', apiPaths.team(name).services)

    return axios.all([getTeam(), getTeamServices()])
      .then(axios.spread(function (team, services) {
        return { team, services }
      }))
      .catch(err => {
        throw new Error(err.message)
      })
  }

  static getInitialProps = async (ctx) => {
    const data = await NewTeamServicePage.getPageData(ctx)
    return data
  }

  render() {
    const teamName = this.props.team.metadata.name
    const teamServices = this.props.services.items

    return (
      <div>
        <Breadcrumb
          items={[
            { text: this.props.team.spec.summary, href: '/teams/[name]', link: `/teams/${teamName}` },
            { text: 'New service' }
          ]}
        />
        <Title>New Service for {this.props.team.spec.summary}</Title>
        <ServiceBuildForm
          user={this.props.user}
          team={this.props.team}
          teamServices={teamServices}
          skipButtonText="Cancel"
        />
      </div>
    )
  }
}

export default NewTeamServicePage
