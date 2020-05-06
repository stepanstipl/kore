import React from 'react'
import PropTypes from 'prop-types'
import { Alert, Button, Drawer, Typography } from 'antd'
const { Paragraph, Title } = Typography

import GCPAutomatedProjectForm from './forms/GCPAutomatedProjectForm'
import GCPAutomatedProjectList from './GCPAutomatedProjectList'

class GCPKoreManagedProjectsCustom extends React.Component {

  static defaultGcpProjectList = [{
    code: 'not-production',
    name: 'Not production',
    description: 'To be used for all environments except production',
    prefix: 'kore',
    suffix: 'notprod',
    plans: ['gke-development']
  }, {
    code: 'production',
    name: 'Production',
    description: 'Project just for the production environment',
    prefix: 'kore',
    suffix: 'prod',
    plans: ['gke-production']
  }]

  static propTypes = {
    gcpProjectList: PropTypes.array.isRequired,
    plans: PropTypes.array.isRequired,
    handleChange: PropTypes.func.isRequired,
    handleDelete: PropTypes.func.isRequired,
    handleEdit: PropTypes.func.isRequired,
    handleAdd: PropTypes.func.isRequired,
    handleReset: PropTypes.func.isRequired,
    hideGuidance: PropTypes.bool
  }

  state = {
    addGcpProject: false,
    editGcpProject: false
  }

  addGcpProject = (enabled) => () => this.setState({ addGcpProject: enabled })

  editGcpProject = (projectCode) => () => this.setState({ editGcpProject: this.props.gcpProjectList.find(p => p.code === projectCode) })

  handleAddGcpProject = (project) => {
    this.addGcpProject(false)()
    this.props.handleAdd(project)
  }

  handleEditGcpProject = (projectCode) => {
    return (project) => {
      this.editGcpProject(false)()
      this.props.handleEdit(projectCode)(project)
    }
  }

  render() {
    const { gcpProjectList, plans, hideGuidance } = this.props

    return (
      <>
        {!hideGuidance && (
          <Alert
            message="When a team is created in Kore and a cluster is requested, Kore will ensure the associated GCP project is also created and the cluster placed inside it. You must also specify the plans available for each type of project, this is to ensure the correct cluster specification is being used."
            type="info"
            showIcon
            style={{ marginBottom: '20px', marginTop: '10px' }}
          />
        )}
        <div style={{ display: 'block', marginBottom: '20px', marginTop: '10px' }}>
          <Button type="primary" onClick={this.addGcpProject(true)}>+ New</Button>
          <Button className="set-kore-defaults" style={{ marginLeft: '10px' }} onClick={() => this.props.handleReset(GCPKoreManagedProjectsCustom.defaultGcpProjectList)}>Set to Kore defaults</Button>
        </div>
        {gcpProjectList.length === 0 ? (
          <Paragraph>No automated projects configured, you can &apos;Set to Kore defaults&apos; and/or add new ones. </Paragraph>
        ) : (
          <GCPAutomatedProjectList
            automatedProjectList={gcpProjectList}
            plans={plans}
            handleChange={this.props.handleChange}
            handleDelete={this.props.handleDelete}
            handleEdit={this.editGcpProject}
          />
        )}
        {this.state.addGcpProject && (
          <Drawer
            title={<Title level={4}>New GCP automated project</Title>}
            visible={this.state.addGcpProject}
            onClose={this.addGcpProject(false)}
            width={700}
          >
            <GCPAutomatedProjectForm handleSubmit={this.handleAddGcpProject} handleCancel={this.addGcpProject(false)} />
          </Drawer>
        )}
        {this.state.editGcpProject && (
          <Drawer
            title={<Title level={4}>Edit GCP automated project</Title>}
            visible={Boolean(this.state.editGcpProject)}
            onClose={this.editGcpProject(false)}
            width={700}
          >
            <GCPAutomatedProjectForm data={this.state.editGcpProject} handleSubmit={this.handleEditGcpProject(this.state.editGcpProject.code)} handleCancel={this.editGcpProject(false)} />
          </Drawer>
        )}
      </>
    )

  }
}

export default GCPKoreManagedProjectsCustom
