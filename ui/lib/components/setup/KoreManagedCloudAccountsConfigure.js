import React from 'react'
import PropTypes from 'prop-types'
import { pluralize, titleize } from 'inflect'
import { Button, Drawer, Typography } from 'antd'
const { Paragraph, Title } = Typography
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

import AutomatedCloudAccountList from './AutomatedCloudAccountList'
import AutomatedCloudAccountForm from './forms/AutomatedCloudAccountForm'

class KoreManagedCloudAccountsConfigure extends React.Component {

  static propTypes = {
    cloudAccountList: PropTypes.array.isRequired,
    plans: PropTypes.array.isRequired,
    handleChange: PropTypes.func.isRequired,
    handleDelete: PropTypes.func.isRequired,
    handleEdit: PropTypes.func.isRequired,
    handleAdd: PropTypes.func.isRequired,
    handleReset: PropTypes.func.isRequired,
    cloud: PropTypes.oneOf(['GCP', 'AWS']).isRequired,
    accountNoun: PropTypes.string.isRequired
  }

  getDefaultCloudAccountList = () => [{
    code: 'not-production',
    name: 'Not production',
    description: 'To be used for all environments except production',
    prefix: 'kore',
    suffix: 'notprod',
    plans: publicRuntimeConfig.cloudAccountAutomation.notprod.defaultPlans[this.props.cloud]
  }, {
    code: 'production',
    name: 'Production',
    description: 'Project just for the production environment',
    prefix: 'kore',
    suffix: 'prod',
    plans: publicRuntimeConfig.cloudAccountAutomation.prod.defaultPlans[this.props.cloud]
  }]

  state = {
    addCloudAccount: false,
    editCloudAccount: false
  }

  addCloudAccount = (enabled) => () => this.setState({ addCloudAccount: enabled })

  editCloudAccount = (projectCode) => () => this.setState({ editCloudAccount: this.props.cloudAccountList.find(p => p.code === projectCode) })

  handleAddCloudAccount = (project) => {
    this.addCloudAccount(false)()
    this.props.handleAdd(project)
  }

  handleEditCloudAccount = (projectCode) => {
    return (project) => {
      this.editCloudAccount(false)()
      this.props.handleEdit(projectCode)(project)
    }
  }

  render() {
    const { cloudAccountList, plans } = this.props

    return (
      <>
        <div style={{ display: 'block', marginBottom: '20px', marginTop: '10px' }}>
          <Button type="primary" onClick={this.addCloudAccount(true)}>+ New</Button>
          <Button className="set-kore-defaults" style={{ marginLeft: '10px' }} onClick={() => this.props.handleReset(this.getDefaultCloudAccountList())}>Set to Kore defaults</Button>
        </div>
        {cloudAccountList.length === 0 ? (
          <Paragraph>No automated {pluralize(this.props.accountNoun)} configured, you can &apos;Set to Kore defaults&apos; and/or add new ones. </Paragraph>
        ) : (
          <AutomatedCloudAccountList
            automatedCloudAccountList={cloudAccountList}
            plans={plans}
            handleChange={this.props.handleChange}
            handleDelete={this.props.handleDelete}
            handleEdit={this.editCloudAccount}
          />
        )}
        {this.state.addCloudAccount && (
          <Drawer
            title={<Title level={4}>New {this.props.cloud} automated {this.props.accountNoun}</Title>}
            visible={this.state.addCloudAccount}
            onClose={this.addCloudAccount(false)}
            width={700}
          >
            <AutomatedCloudAccountForm
              alertTitle={`${titleize(this.props.accountNoun)} naming`}
              alertDescription={`The ${this.props.cloud} ${this.props.accountNoun} will name using the optional prefix and suffix specified below, along with the team ID.`}
              handleSubmit={this.handleAddCloudAccount}
              handleCancel={this.addCloudAccount(false)}
            />
          </Drawer>
        )}
        {this.state.editCloudAccount && (
          <Drawer
            title={<Title level={4}>Edit {this.props.cloud} automated {this.props.accountNoun}</Title>}
            visible={Boolean(this.state.editCloudAccount)}
            onClose={this.editCloudAccount(false)}
            width={700}
          >
            <AutomatedCloudAccountForm
              alertTitle={`${titleize(this.props.accountNoun)} naming`}
              alertDescription={`The ${this.props.cloud} ${this.props.accountNoun} will name using the optional prefix and suffix specified below, along with the team ID.`}
              data={this.state.editCloudAccount}
              handleSubmit={this.handleEditCloudAccount(this.state.editCloudAccount.code)}
              handleCancel={this.editCloudAccount(false)}
            />
          </Drawer>
        )}
      </>
    )
  }
}

export default KoreManagedCloudAccountsConfigure
