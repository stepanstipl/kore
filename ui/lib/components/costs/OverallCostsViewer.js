import * as React from 'react'
import { DatePicker, Button, Alert, Typography, Card, Icon, Radio, Row, Col } from 'antd'
import moment from 'moment'

import KoreApi from '../../kore-api'
import OverallCostSummary from './OverallCostSummary'
import TeamCostSummary from './TeamCostSummary'
import { getCloudInfo } from '../../utils/cloud'
import { apiDateTime, startOfMonth, endOfMonth } from '../../utils/date-helpers'

export default class OverallCostsViewer extends React.Component {
  state = {
    // Params:
    provider: '',
    from: startOfMonth(),
    to: endOfMonth(),

    // Overall summary:
    loading: true,
    error: null,
    summary: null,

    // Current team summary:
    teamDetailLoading: false,
    teamDetailError: null,
    teamDetail: null
  }

  componentDidMount = () => {
    this.loadSummary()
  }

  loadSummary = async () => {
    this.setState({ loading: true })
    try {
      const summary = await (await KoreApi.client()).costs.GetCostSummary(
        apiDateTime(this.state.from), 
        apiDateTime(this.state.to),
        this.state.provider !== '' ? this.state.provider : null
      )
      this.setState({ loading: false, error: null, summary, teamDetail: null })
    } catch (error) {
      this.setState({ loading: false, error })
    }
  }

  loadTeamDetail = async (teamIdentifier) => {
    this.setState({ teamDetailLoading: true })
    try {
      const teamDetail = await (await KoreApi.client()).costs.GetTeamCostSummary(
        teamIdentifier,
        apiDateTime(this.state.from), 
        apiDateTime(this.state.to),
        this.state.provider !== '' ? this.state.provider : null
      )
      this.setState({ teamDetailLoading: false, teamDetail, teamDetailError: null })
    } catch (teamDetailError) {
      this.setState({ teamDetailLoading: false, teamDetailError })
    }
  }

  prevMonth = () => {
    let { from, to } = this.state
    from = startOfMonth(moment(from).add(-1, 'month'))
    to = endOfMonth(from)
    this.setState({ from, to }, () => this.loadSummary())
  }

  nextMonth = () => {
    let { from, to } = this.state
    from = startOfMonth(moment(from).add(1, 'month'))
    to = endOfMonth(from)
    this.setState({ from, to }, () => this.loadSummary())
  }

  render() {
    const { summary, loading, from, to, error, teamDetail, teamDetailError, teamDetailLoading, provider } = this.state

    const isCurrentMonth = from.format('MMM') === moment().utc(false).format('MMM')
    const toDisplay = isCurrentMonth ? moment().utc(false) : to
    const providerInfo = provider !== '' ? getCloudInfo(provider) : null

    return (
      <>
        <div style={{ marginBottom: '20px' }}>
          <Button disabled={loading} onClick={() => this.prevMonth()}><Icon type="step-backward"/> Previous month</Button>&nbsp;
          <DatePicker.RangePicker
            disabled={loading}
            value={[from, to]}
            format={'D MMM YYYY'}
            allowClear={false}
            onChange={(dates) => this.setState({ from: dates[0], to: dates[1] }, () => this.loadSummary())}
          />&nbsp;
          <Button disabled={loading || isCurrentMonth} onClick={() => this.nextMonth()}>Next month <Icon type="step-forward"/></Button>&nbsp;
          <Radio.Group 
            defaultValue="" 
            onChange={(e) => this.setState({ provider: e.target.value }, () => this.loadSummary())}
            buttonStyle="solid"
          >
            <Radio.Button value="">All Providers</Radio.Button>
            <Radio.Button value="gcp">Google Cloud Platform</Radio.Button>
            <Radio.Button value="aws">Amazon Web Services</Radio.Button>
            <Radio.Button value="azure">Microsoft Azure</Radio.Button>
          </Radio.Group>
        </div>

        <Typography.Title level={2}>Cost Report for {from.format('D MMM YYYY')} to {toDisplay.format('D MMM YYYY')}{providerInfo ? ` on ${providerInfo.cloudLong}` : ''}</Typography.Title>
        <Row gutter={[16, 16]}>
          <Col xs={24} xxl={10}>
            <Card title="Summary">
              {error ? <Alert type="error" message={error.message} /> : null}
              <OverallCostSummary 
                summary={summary} 
                loading={loading}
                onTeamDetail={(teamIdentifier) => this.loadTeamDetail(teamIdentifier)} 
              />
            </Card>
          </Col>
          <Col xs={24} xxl={14}>
            <Card title={`Team costs ${teamDetail ? teamDetail.teamName : ''}`} actions={[
              <Button key="close" disabled={!teamDetail} onClick={() => this.setState({ teamDetail: null })}>Close</Button>
            ]}>
              {teamDetail ? null : <Alert type="info" showIcon={true} message="Select team to inspect their costs in detail" />}
              {teamDetailError ? <Alert type="error" message={teamDetailError.message} /> : null}
              <TeamCostSummary summary={teamDetail} loading={teamDetailLoading} />
            </Card>
          </Col>
        </Row>
      </>
    )
  }
}