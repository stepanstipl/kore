import * as React from 'react'
import PropTypes from 'prop-types'
import Router from 'next/router'
import { Button, Form } from 'antd'

import redirect from '../../../utils/redirect'
import CloudSelector from '../../common/CloudSelector'
import MissingCredential from './MissingCredential'
import ClusterOptionsForm from './ClusterOptionsForm'
import FormErrorMessage from '../../forms/FormErrorMessage'
import KoreApi from '../../../kore-api'
import getConfig from 'next/config'
import { loadingMessage } from '../../../utils/message'
const { publicRuntimeConfig } = getConfig()

class ClusterBuildForm extends React.Component {
  static propTypes = {
    form: PropTypes.any.isRequired,
    skipButtonText: PropTypes.string,
    team: PropTypes.object.isRequired,
    teamClusters: PropTypes.array.isRequired,
    user: PropTypes.object.isRequired
  }

  constructor(props) {
    super(props)
    this.state = {
      submitButtonText: 'Save',
      skipButtonText: this.props.skipButtonText || 'Skip',
      submitting: false,
      formErrorMessage: false,
      selectedCloud: '',
      selectedProvider: '',
      dataLoading: true,
      credentials: {},
      planValues: null,
      validationErrors: null
    }
  }

  async fetchComponentData() {
    const team = this.props.team.metadata.name
    const api = await KoreApi.client()
    const [ allocations, plans ] = await Promise.all([
      api.ListAllocations(team, true),
      api.ListPlans()
    ])
    return { allocations, plans }
  }

  componentDidMountComplete = null
  componentDidMount() {
    // Assign the promise chain to a variable so tests can wait for it to complete.
    this.componentDidMountComplete = Promise.resolve().then(async () => {
      const { allocations, plans } = await this.fetchComponentData()
      const gkeCredentials = (allocations.items || []).filter(a => a.spec.resource.kind === 'GKECredentials')
      const eksCredentials = (allocations.items || []).filter(a => a.spec.resource.kind === 'EKSCredentials')
      const aksCredentials = (allocations.items || []).filter(a => a.spec.resource.kind === 'AKSCredentials')

      // use the predictable name of the AccountManagement resource that is allocated to filter the allocations
      const gcpAccountManagement = (allocations.items || []).find(a => a.spec.resource.kind === 'AccountManagement' && a.spec.resource.name === 'am-gcp')
      const awsAccountManagement = (allocations.items || []).find(a => a.spec.resource.kind === 'AccountManagement' && a.spec.resource.name === 'am-aws')

      this.setState({
        credentials: {
          GCP: {
            credentials: gkeCredentials,
            accountManagement: gcpAccountManagement
          },
          AWS: {
            credentials: eksCredentials,
            accountManagement: awsAccountManagement
          },
          Azure: {
            credentials: aksCredentials,
            accountManagement: undefined
          }
        },
        plans: plans,
        dataLoading: false
      })
    })
  }

  getClusterResource = (values) => {
    const team = this.props.team.metadata.name
    const selectedPlan = this.state.plans.items.find(p => p.metadata.name === values.plan)

    let selectedCredential
    if (values.credential) {
      selectedCredential = this.state.credentials[this.state.selectedCloud].credentials.find(p => p.metadata.name === values.credential)
    } else {
      selectedCredential = this.state.credentials[this.state.selectedCloud].accountManagement
    }
    const credentials = selectedCredential.spec.resource

    return KoreApi.resources()
      .team(team)
      .generateClusterResource(this.props.user.id, values, selectedPlan, this.state.planValues, credentials)
  }

  handleSubmit = e => {
    e.preventDefault()
    this.setState({ submitting: true, formErrorMessage: false })

    this.clusterOptionsForm.props.form.validateFields(async (err, values) => {
      if (err) {
        this.setState({ formErrorMessage: 'Validation failed', submitting: false })
        return
      }
      try {
        await (await KoreApi.client()).UpdateCluster(
          this.props.team.metadata.name,
          values.clusterName,
          this.getClusterResource(values)
        )
        loadingMessage('Cluster build requested...')
        return redirect({
          router: Router,
          path: `/teams/${this.props.team.metadata.name}/clusters/${values.clusterName}`
        })
      } catch (err) {
        this.setState({
          submitting: false,
          formErrorMessage: (err.fieldErrors && err.message) ? err.message : 'An error occurred requesting the cluster, please try again',
          validationErrors: err.fieldErrors // This will be undefined on non-validation errors, which is fine.
        })
      }
    })
  }

  handleSelectCloud = cloud => {
    this.setState({
      selectedCloud: cloud,
      selectedProvider: publicRuntimeConfig.clusterProviderMap[cloud],
      planValues: null,
      validationErrors: null
    })
  }

  setPlanValues = planValues => {
    this.setState({ planValues })
  }

  clusterBuildForm = () => {
    const { submitting, selectedCloud, selectedProvider, formErrorMessage } = this.state
    const filteredPlans = this.state.plans.items.filter(p => p.spec.kind === selectedProvider)
    const filteredCredentials = this.state.credentials[selectedCloud].credentials
    const accountManagement = this.state.credentials[selectedCloud].accountManagement
    const formConfig = {
      layout: 'horizontal',
      labelAlign: 'left',
      hideRequiredMark: true,
      labelCol: {
        sm: { span: 24 },
        md: { span: 24 },
        lg: { span: 6 }
      },
      wrapperCol: {
        sm: { span: 24 },
        md: { span: 24 },
        lg: { span: 18 }
      }
    }

    return (
      <Form {...formConfig} onSubmit={this.handleSubmit}>
        <FormErrorMessage message={formErrorMessage} />
        <ClusterOptionsForm
          team={this.props.team}
          selectedCloud={selectedCloud}
          selectedProvider={selectedProvider}
          credentials={filteredCredentials}
          accountManagement={accountManagement}
          plans={filteredPlans}
          teamClusters={this.props.teamClusters}
          onPlanValuesChange={this.setPlanValues}
          validationErrors={this.state.validationErrors}
          wrappedComponentRef={inst => this.clusterOptionsForm = inst}
        />
        <Form.Item style={{ marginTop: '20px' }}>
          <Button id="save" type="primary" htmlType="submit" loading={submitting}>
            {this.state.submitButtonText}
          </Button>
        </Form.Item>
      </Form>
    )
  }

  render() {
    if (this.state.dataLoading || !this.props.team) {
      return null
    }

    const { selectedCloud } = this.state
    const filteredCredentials = selectedCloud && this.state.credentials[selectedCloud].credentials
    const accountManagement = selectedCloud && this.state.credentials[selectedCloud].accountManagement

    return (
      <div>
        <CloudSelector selectedCloud={selectedCloud} handleSelectCloud={this.handleSelectCloud} />
        {selectedCloud ? (
          filteredCredentials.length > 0 || accountManagement ?
            this.clusterBuildForm() :
            <MissingCredential team={this.props.team.metadata.name} cloud={selectedCloud} />
        ) : null}
      </div>
    )
  }
}

const WrappedClusterBuildForm = Form.create({ name: 'new_team_cluster_build' })(ClusterBuildForm)

export default WrappedClusterBuildForm
