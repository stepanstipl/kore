import React from 'react'
import PropTypes from 'prop-types'
import { Col, Row, Typography } from 'antd'
const { Text } = Typography

import Breadcrumb from '../layout/Breadcrumb'
import TeamSettings from './TeamSettings'

const TeamHeader = ({ team, breadcrumbExt, teamRemoved }) => (
  <>
    <Row gutter={[0, 16]}>
      <Col span={20} style={{ marginTop: '8px' }}>
        <Breadcrumb items={[{ text: team.spec.summary, href: '/teams/[name]', link: `/teams/${team.metadata.name}` }, ...(breadcrumbExt || [])]} />
      </Col>
      <Col span={4} style={{ textAlign: 'right' }}>
        <TeamSettings team={team} teamRemoved={teamRemoved} />
      </Col>
    </Row>
    <Row style={{ marginBottom: '30px' }}>
      <Col span={12}>
        {team.spec.description ? <Text strong>{team.spec.description}</Text> : <Text style={{ fontStyle: 'italic' }} type="secondary">No description</Text>}
      </Col>
      <Col span={12} style={{ textAlign: 'right' }}>
        <Text><Text strong>Team ID: </Text>{team.metadata.name}</Text>
      </Col>
    </Row>
  </>
)

TeamHeader.propTypes = {
  team: PropTypes.object.isRequired,
  breadcrumbExt: PropTypes.array,
  teamRemoved: PropTypes.func.isRequired
}

export default TeamHeader
