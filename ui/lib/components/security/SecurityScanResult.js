import React from 'react'
import PropTypes from 'prop-types'
import { Card, Descriptions, List, Tooltip, Icon, Button, Collapse, Spin, message, Drawer } from 'antd'
import moment from 'moment'
import Link from 'next/link'
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

import SecurityStatusIcon from './SecurityStatusIcon'
import KoreApi from '../../kore-api'
import SecurityRule from './SecurityRule'

export default class SecurityScanResult extends React.Component {
  static propTypes = {
    result: PropTypes.object.isRequired,
    history: PropTypes.object,
    historyLoading: PropTypes.bool,
    onHistoryRequest: PropTypes.func
  }

  state = {
    loadingRule: false,
    ruleDetails: null
  }

  loadRule = async (ruleCode) => {
    this.setState({
      loadingRule: true,
      ruleDetails: null
    })
    try {
      const rule = await (await KoreApi.client()).security.GetSecurityRule(ruleCode)
      this.setState({
        loadingRule: false,
        ruleDetails: rule
      })
      if (!rule) {
        message.error(`Rule ${ruleCode} not found, cannot show details`)
      }
    } catch (err) {
      this.setState({ loadingRule: false })
      console.error(`Failed to load details of rule ${ruleCode}`, err)
      message.error(`Failed to load details of rule ${ruleCode}`)
    }
  }

  closeRule = () => {
    this.setState({
      loadingRule: false,
      ruleDetails: null
    })
  }

  getTeamName = (possibleTeamName) => {
    if (!possibleTeamName || publicRuntimeConfig.ignoreTeams.includes(possibleTeamName)) {
      return 'Global'
    }
    return possibleTeamName
  }

  render() {
    const { result, history, historyLoading, onHistoryRequest } = this.props
    const { loadingRule, ruleDetails } = this.state
    const checkedMoment = moment(result.spec.checkedAt)
    const checked = <Tooltip title={checkedMoment.format('D MMM YYYY HH:mm:ss')}>{checkedMoment.fromNow()}</Tooltip>
    let archived = 'This is the current result'
    if (result.spec.archivedAt) {
      const archivedMoment = moment(result.spec.archivedAt)
      archived = <Tooltip title={archivedMoment.format('D MMM YYYY HH:mm:ss')}>{archivedMoment.fromNow()}</Tooltip>
    }
    const historySorted = history ? history.items.sort((a, b) => a.checkedAt < b.checkedAt ? 1 : -1) : []

    return (
      <>
        <Drawer
          title={`Rule Details${ruleDetails ? ' - ' + ruleDetails.spec.name : ''}`}
          visible={Boolean(loadingRule || ruleDetails)}
          onClose={() => this.closeRule()}
          width={700}
        >
          {loadingRule ? <div style={{ textAlign: 'center' }}><Spin /></div> : null}
          {ruleDetails !== null ? <SecurityRule rule={ruleDetails} /> : null}
        </Drawer>
        <Card title="Scan Details" style={{ marginBottom: '20px' }}>
          <Descriptions>
            <Descriptions.Item label="Type">
              <Tooltip title={`Group: ${result.spec.resource.group} Version: ${result.spec.resource.version} Kind: ${result.spec.resource.kind}`}>
                {result.spec.resource.kind}
              </Tooltip>
            </Descriptions.Item>
            <Descriptions.Item label="Team">{this.getTeamName(result.spec.resource.namespace)}</Descriptions.Item>
            <Descriptions.Item label="Name">{result.spec.resource.name}</Descriptions.Item>
          </Descriptions>
          <Descriptions>
            <Descriptions.Item label="Overall Status"><SecurityStatusIcon status={result.spec.overallStatus} text={result.spec.overallStatus} /></Descriptions.Item>
            <Descriptions.Item label="Checked">{checked}</Descriptions.Item>
            <Descriptions.Item label="Archived">{archived}</Descriptions.Item>
          </Descriptions>
          {onHistoryRequest ? (
            <Collapse onChange={(open) => open.length > 0 ? onHistoryRequest() : () => {}}>
              <Collapse.Panel header={`Security status history for this ${result.spec.resource.kind.toLowerCase()}`}>
                {historyLoading ? (
                  <div style={{ textAlign: 'center' }}>
                    <Spin />
                  </div>
                ) : null}
                {history && !historyLoading ? (
                  <List 
                    dataSource={historySorted}
                    renderItem={(scan) => {
                      const checked = result.spec.id === scan.spec.id ? `Current (last checked ${moment(scan.spec.checkedAt).format('D MMM YYYY HH:mm:ss')})` : moment(scan.spec.checkedAt).format('D MMM YYYY HH:mm:ss')
                      let archived = scan.spec.archivedAt ? `| Superceded at ${moment(scan.spec.archivedAt).format('D MMM YYYY HH:mm:ss')}` : ''
                      return (
                        <List.Item 
                          key={scan.spec.checkedAt}
                          actions={[<Link key="viewscan" href="/security/scans/[id]" as={`/security/scans/${scan.spec.id}`}><a><Tooltip placement="left" title="View scan details"><Icon type="info-circle" /> View Scan</Tooltip></a></Link>]}>
                          <List.Item.Meta
                            title={(<Link href="/security/scans/[id]" as={`/security/scans/${scan.spec.id}`}><a><Tooltip title="View scan details">{checked}</Tooltip></a></Link>)}
                            description={`Status: ${scan.spec.overallStatus} ${archived}`}
                            avatar={<SecurityStatusIcon status={scan.spec.overallStatus} />}
                          />
                        </List.Item>
                      )
                    }}
                  />
                ) : null}
              </Collapse.Panel>
            </Collapse>
          ) : (
            <Link href="/security/resources/[group]/[version]/[kind]/[namespace]/[name]" as={`/security/resources/${result.spec.resource.group}/${result.spec.resource.version}/${result.spec.resource.kind}/${result.spec.resource.namespace}/${result.spec.resource.name}`}>
              <Button>View current security status for this {result.spec.resource.kind.toLowerCase()}</Button>
            </Link>
          )}
        </Card>

        <Card title="Security Rule Compliance" style={{ marginBottom: '20px' }}>
          <List>
            {result.spec.results.map((ruleResult) => {
              return (
                <List.Item
                  key={ruleResult.ruleCode}
                  actions={[<a key="view" onClick={() => this.loadRule(ruleResult.ruleCode)}><Tooltip placement="left" title="View rule details"><Icon type="info-circle" /> View Rule</Tooltip></a>]}>
                  <List.Item.Meta
                    title={(<a onClick={() => this.loadRule(ruleResult.ruleCode)}><Tooltip title="View rule details">{ruleResult.ruleCode}</Tooltip></a>)}
                    description={`${ruleResult.status} - ${ruleResult.message}`}
                    avatar={<SecurityStatusIcon status={ruleResult.status} />}
                  />
                </List.Item>
              )
            })}
          </List>
        </Card>
      </>
    )
  }
}
