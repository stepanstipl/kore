import React from 'react'
import Link from 'next/link'
import PropTypes from 'prop-types'
import { Form, Icon, Input, Modal, Table, Tag, Tooltip } from 'antd'

import MonitoringSummary from './MonitoringSummary'
import KoreApi from '../../kore-api'
import FormErrorMessage from '../forms/FormErrorMessage'
import { errorMessage, successMessage } from '../../utils/message'

export default class MonitoringTable extends React.Component {
  static propTypes = {
    alerts: PropTypes.object.isRequired,
    severity: PropTypes.array,
    status: PropTypes.array,
    refreshData: PropTypes.func
  }

  static silenceAlertInitialState = {
    silenceAlert: false,
    silenceAlertSubmitting: false,
    silenceAlertError: false,
    silenceAlertComment: undefined,
    silenceAlertDuration: undefined
  }

  state = {
    ...MonitoringTable.silenceAlertInitialState
  }

  columns = [
    {
      title: 'Severity',
      dataIndex: 'status.rule.spec.severity',
      key: 'severity',
      onFilter: (value, record) => record.status.rule.spec.severity(value) === 0,
      sorter: (a, b) => a.status.rule.spec.severity.length - b.status.rule.spec.severity.length,
      sortDirections: ['descend'],
      render: (text) => (
        <>
          <Tag key="middle" color={text === 'Critical' ? 'red' : 'orange'}>
            {text}
          </Tag>
        </>
      ),
    },
    {
      title: 'Triggered',
      dataIndex: 'metadata.creationTimestamp',
      key: 'triggered',
    },
    {
      title: 'Rule Name',
      dataIndex: 'status.rule.metadata.name',
      key: 'name',
      render: (text) => (
        <Link
          key="view"
          passHref
          href="/docs"
        >
          <a>
            <Tooltip placement="top" title="View the definition of this rule">
              <Icon type="info-circle" />  {text}
            </Tooltip>
          </a>
        </Link>
      )
    },
    {
      title: 'Summary',
      dataIndex: 'spec.summary',
      key: 'summary',
      render: (text, record) => <MonitoringSummary record={record}/>
    },
    {
      title: 'State',
      render: (text, record) => (
        <>
          <Tooltip placement="top" title="The current state of the alert">
            <Tag color="green"><Icon type="info-circle" style={{ marginRight: '5px' }}/>{record.status.status}</Tag>
          </Tooltip>
        </>
      )
    },
    {
      title: 'Actions',
      render: (text, record) => (
        record.status.status === 'Silenced'
          ? <a style={{ textDecoration: 'underline' }} onClick={this.unsilenceAlert(record)}> Unsilence</a>
          : <a style={{ textDecoration: 'underline' }} onClick={() => this.setState({ silenceAlert: record })}> Silence</a>
      )
    }
  ]

  unsilenceAlert = (alert) => () => {
    Modal.confirm({
      title: 'Unsilence this alert?',
      okText: 'Yes',
      cancelText: 'No',
      onOk: async () => {
        try {
          await (await KoreApi.client()).monitoring.UnsilenceAlert(alert.metadata.uid)
          successMessage('Alert unsilenced')
          this.props.refreshData && this.props.refreshData()
        } catch (err) {
          console.error('Error unsilencing alert', err)
          errorMessage('Error unsilencing alert, please try again')
        }
      }
    })
  }

  silenceAlert = async () => {
    const alert = this.state.silenceAlert
    if (!this.state.silenceAlertComment || !this.state.silenceAlertDuration) {
      return this.setState({
        silenceAlertError: 'Please enter a comment and duration'
      })
    }

    this.setState({ silenceAlertSubmitting: true })
    try {
      await (await KoreApi.client()).monitoring.SilenceAlert(alert.metadata.uid, this.state.silenceAlertComment, this.state.silenceAlertDuration)
      this.setState({ ...MonitoringTable.silenceAlertInitialState })
      successMessage('Alert silenced')
      this.props.refreshData && this.props.refreshData()
    } catch (err) {
      console.error('Error silencing alert', err)
      this.setState({ silenceAlertError: 'There was an error trying to silence the alert, please try again.', silenceAlertSubmitting: false })
    }
  }

  renderSilenceAlertModal = () => (
    <Modal
      title="Silence this alert"
      visible={Boolean(this.state.silenceAlert)}
      okType="danger"
      okText="Yes"
      cancelText="No"
      onOk={async () => await this.silenceAlert()}
      confirmLoading={this.state.silenceAlertSubmitting}
      onCancel={() => this.setState({ ...MonitoringTable.silenceAlertInitialState })}
    >
      <FormErrorMessage message={this.state.silenceAlertError} />
      <Form.Item label="Comment" onChange={(e) => this.setState({ silenceAlertComment: e.target.value })}>
        <Input value={this.state.silenceAlertComment} placeholder="Reason for silencing this alert" />
      </Form.Item>
      <Form.Item label="Duration" onChange={(e) => this.setState({ silenceAlertDuration: e.target.value })} help="Enter a duration, such as 5m, 2h or 1d">
        <Input value={this.state.silenceAlertDuration} placeholder="For example, 1h"/>
      </Form.Item>
    </Modal>
  )

  filterOnRules = () => {
    if (!this.props.alerts) {
      return []
    }
    const filtered = this.props.alerts.items.filter(a =>
      (this.props.severity || []).includes(a.status.rule.spec.severity) && (this.props.status || []).includes(a.status.status)
    )
    return filtered
  }


  render() {
    return (
      <>
        <Table
          rowKey={(record) => record.metadata.uid}
          dataSource={this.filterOnRules()}
          columns={this.columns}
          //expandedRowRender={record => <MonitoringAlert rule={record}/> }
        />
        {this.renderSilenceAlertModal()}
      </>
    )
  }
}
