import React from 'react'
import PropTypes from 'prop-types'
import { List, Card, Icon, Tooltip } from 'antd'
import inflect from 'inflect'
import Link from 'next/link'
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

import SecurityStatusIcon from './SecurityStatusIcon'

export default class SecurityOverview extends React.Component {
  static propTypes = {
    overview: PropTypes.object.isRequired
  }

  static statuses = [
    { status: 'Failure', title: 'Security Failures' }, 
    { status: 'Warning', title: 'Security Warnings' }, 
    { status: 'Compliant', title: 'Compliant Checks' }
  ]

  formatStatus = (status, count) => {
    const suffix = count > 0 ? 's' : ''
    if (status === 'Compliant') {
      return `compliant security check${suffix}`
    }
    return `security ${status.toLowerCase()}${suffix}`
  }

  resourceSummary = (overview) => {
    // Sort resources by kind
    const resources = {}
    overview.spec.resources.forEach(resource => {
      if (!resources[resource.resource.kind]) {
        resources[resource.resource.kind] = []
      }
      resources[resource.resource.kind].push({ 
        name: resource.resource.name, 
        namespace: resource.resource.namespace, 
        overallStatus: resource.overallStatus, 
        openIssueCounts: resource.openIssueCounts,
        link: `/security/resources/${resource.resource.group}/${resource.resource.version}/${resource.resource.kind}/${resource.resource.namespace}/${resource.resource.name}`
      })
    })

    // Sort within each kind by severity:
    Object.keys(resources).forEach((kind) => {
      resources[kind] = resources[kind].sort((a, b) => {
        if (a.overallStatus === 'Failure' && b.overallStatus !== 'Failure') {
          return -1
        }
        if (a.overallStatus !== 'Failure' && b.overallStatus === 'Failure') {
          return 1
        }
        if (a.overallStatus === 'Warning' && b.overallStatus === 'Compliant') {
          return -1
        }
        if (a.overallStatus === 'Compliant' && b.overallStatus === 'Warning') {
          return 1
        }
        return 0
      })
    })

    return resources
  }

  safeCount = (possCount) => {
    return !possCount ? 0 : possCount
  }

  getTeamName = (possibleTeamName) => {
    if (!possibleTeamName || publicRuntimeConfig.ignoreTeams.includes(possibleTeamName)) {
      return 'Global'
    }
    return `Team: ${possibleTeamName}`
  }

  render() {
    const { overview } = this.props
    if (!overview) {
      return null
    }
    const resourceSummary = this.resourceSummary(overview)

    return (
      <>
        <Card title="Overview" style={{ marginBottom: '20px' }}>
          <List>
            {SecurityOverview.statuses.map((statusDetails) => {
              const { status, title } = statusDetails
              const count = overview.spec.openIssueCounts[status]
              return (
                <List.Item key={status}>
                  <List.Item.Meta
                    title={title}
                    description={count > 0 ? (
                      <>You have {count} current {this.formatStatus(status, count)}</>
                    ) : (
                      <>You have no current {this.formatStatus(status, 2)}</>
                    )}
                    avatar={<SecurityStatusIcon status={status} inactive={this.safeCount(count) === 0} />}
                  />
                </List.Item>
              )
            })}
          </List>
        </Card>
        {Object.keys(resourceSummary).map((resKind) => (
          <Card title={`${inflect.pluralize(resKind)} Status`} style={{ marginBottom: '20px' }} key={resKind}>
            <List>
              {resourceSummary[resKind].map((resource) => (
                <List.Item 
                  key={resource.name} 
                  actions={[
                    <SecurityStatusIcon status="Failure" inactive={this.safeCount(resource.openIssueCounts['Failure']) === 0} text={`${this.safeCount(resource.openIssueCounts['Failure'])} failures`} key="Failure" />,
                    <SecurityStatusIcon status="Warning" inactive={this.safeCount(resource.openIssueCounts['Warning']) === 0} text={`${this.safeCount(resource.openIssueCounts['Warning'])} warnings`} key="Warning" />,
                    <SecurityStatusIcon status="Compliant" inactive={this.safeCount(resource.openIssueCounts['Compliant']) === 0} text={`${this.safeCount(resource.openIssueCounts['Compliant'])} compliant`} key="Compliant" />,
                    <Link key="view" href="/security/resources/[group]/[version]/[kind]/[namespace]/[name]" as={resource.link}><a><Tooltip placement="left" title="View resource status report"><Icon type="info-circle" /></Tooltip></a></Link>
                  ]}
                >
                  <List.Item.Meta
                    title={(<Link key="view" href="/security/resources/[group]/[version]/[kind]/[namespace]/[name]" as={resource.link}><a><Tooltip title="View resource status report">{resource.name}</Tooltip></a></Link>)}
                    description={this.getTeamName(resource.namespace)}
                    avatar={<SecurityStatusIcon status={resource.overallStatus} />}
                  />
                </List.Item>
              ))}
            </List>
          </Card>
        ))}
      </>
    )
  }
}