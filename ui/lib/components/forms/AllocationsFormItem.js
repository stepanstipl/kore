import React from 'react'
import PropTypes from 'prop-types'
import { Form, Select, Radio, Alert } from 'antd'

import KoreApi from '../../kore-api'
import { kore } from '../../../config'

export default class AllocationsFormItem extends React.Component {
  static propTypes = {
    allTeams: PropTypes.array,
    allocatedTeams: PropTypes.array.isRequired,
    onAllocationChange: PropTypes.func.isRequired
  }

  state = {
    allTeams: [],
  }

  componentDidMountComplete = null
  componentDidMount() {
    this.componentDidMountComplete = Promise.resolve().then(async () => {
      // If we've not been handed teams, load them ourselves.
      let teams = this.props.allTeams
      if (!teams) {
        teams = await (await KoreApi.client()).ListTeams()
      }
      this.setState({
        allTeams: teams.items.filter(t => !kore.ignoreTeams.includes(t.metadata.name))
      })
    })
  }

  getAllocateMode = () => {
    return this.props.allocatedTeams.find((t) => t === '*') ? 'all' : 'specified'
  }

  onAllocateModeChange = (v) => {
    this.props.onAllocationChange(v === 'all' ? ['*'] : [])
  }

  onTeamsChange = (v) => {
    this.props.onAllocationChange([...v])
  }

  getHelp = () => {
    const mode = this.getAllocateMode()
    if (mode === 'all') {
      return <>This will be allocated to <span style={{ fontWeight: 'bold' }}>all teams</span></>
    } 
    const hasTeams = this.props.allocatedTeams.length > 0
    if (hasTeams) {
      return <>This will be allocated to the <span style={{ fontWeight: 'bold' }}>specified teams</span> only</>
    }
    return <Alert type='error' message={<>You must specify <span style={{ fontWeight: 'bold' }}>one or more teams</span> this will be allocated to</>} />
  }

  render() {
    if (this.props.allocatedTeams === null) {
      // @TODO: Loading...
      return null
    }
    const { allTeams } = this.state
    const mode = this.getAllocateMode()
    return (
      <Form.Item label="Allocate to teams" help={this.getHelp()} required={true}>
        <Radio.Group onChange={(e) => this.onAllocateModeChange(e.target.value)} value={mode}>
          <Radio value="all">All</Radio>
          <Radio value="specified">Specified teams</Radio>
        </Radio.Group>
        {mode !== 'all' ? (
          <Select
            mode="multiple"
            style={{ width: '100%' }}
            value={this.props.allocatedTeams.filter((t) => t !== '*')}
            onChange={(v) => this.onTeamsChange(v)}
          >
            {allTeams.map(t => (
              <Select.Option key={t.metadata.name} value={t.metadata.name}>{t.spec.summary}</Select.Option>
            ))}
          </Select>
        ) : null}
      </Form.Item>
    )
  }
}