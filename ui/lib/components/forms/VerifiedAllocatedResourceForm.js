import * as React from 'react'
import PropTypes from 'prop-types'
import canonical from '../../utils/canonical'
import V1Allocation from '../../kore-api/model/V1Allocation'
import V1AllocationSpec from '../../kore-api/model/V1AllocationSpec'
import V1ObjectMeta from '../../kore-api/model/V1ObjectMeta'
import V1Ownership from '../../kore-api/model/V1Ownership'
import { message, Typography } from 'antd'
const { Paragraph, Text } = Typography

class VerifiedAllocatedResourceForm extends React.Component {
  static propTypes = {
    form: PropTypes.any.isRequired,
    team: PropTypes.string.isRequired,
    allTeams: PropTypes.object,
    data: PropTypes.object,
    handleSubmit: PropTypes.func.isRequired,
    saveButtonText: PropTypes.string,
    inlineVerification: PropTypes.bool
  }

  constructor(props) {
    super(props)
    let allocations = []
    if (props.data && props.data.allocation) {
      allocations = props.data.allocation.spec.teams.filter(a => a !== '*')
    }
    this.state = {
      submitting: false,
      formErrorMessage: false,
      allocations,
      inlineVerificationFailed: false
    }
  }

  componentDidMount() {
    // To disabled submit button at the beginning.
    this.props.form.validateFields()
  }

  disableButton = fieldsError => {
    if (this.state.submitting) {
      return true
    }
    return Object.keys(fieldsError).some(field => fieldsError[field])
  }

  onAllocationsChange = value => {
    this.setState({
      ...this.state,
      allocations: value
    })
  }

  getResource = () => {
    throw new Error('getResource must be implemented')
  }

  putResource = () => {
    throw new Error('putResource must be implemented')
  }

  async verify(resource, tryCount) {
    const messageKey = 'verify'
    tryCount = tryCount || 0
    if (tryCount === 0) {
      message.loading({ content: 'Verifying credentials', key: messageKey, duration: 0 })
    }
    if (tryCount === 3) {
      message.error({ content: 'Credentials verification failed', key: messageKey })
      this.setState({
        ...this.state,
        inlineVerificationFailed: true,
        submitting: false,
        formErrorMessage: (
          <>
            <Paragraph>The credentials have been saved but could not be verified, see the error below. Please try again or click &quot;Continue without verification&quot;.</Paragraph>
            {(resource.status.conditions || []).map((c, idx) =>
              <Paragraph key={idx} style={{ marginBottom: '0' }}>
                <Text strong>{c.message}</Text>
                <br/>
                <Text>{c.detail}</Text>
              </Paragraph>
            )}
          </>
        )
      })
    } else {
      setTimeout(async () => {
        const resourceResult = await this.getResource(resource.metadata.name)
        if (resourceResult.status.status === 'Success') {
          message.success({ content: 'Credentials verification successful', key: messageKey })
          return await this.props.handleSubmit(resourceResult)
        }
        return await this.verify(resourceResult, tryCount + 1)
      }, 2000)
    }
  }

  setFormSubmitting = (submitting = true, formErrorMessage = false) => {
    this.setState({
      ...this.state,
      submitting,
      formErrorMessage
    })
  }

  getMetadataName = values => {
    const data = this.props.data
    return (data && data.metadata && data.metadata.name) || canonical(values.name)
  }

  generateAllocationResource = (ownership, values) => {
    const metadataName = this.getMetadataName(values)

    const resource = new V1Allocation()
    resource.setApiVersion('config.kore.appvia.io/v1')
    resource.setKind('Allocation')

    const meta = new V1ObjectMeta()
    meta.setName(metadataName)
    meta.setNamespace(this.props.team)
    resource.setMetadata(meta)

    const spec = new V1AllocationSpec()
    spec.setName(values.name)
    spec.setSummary(values.summary)
    spec.setTeams(this.state.allocations.length > 0 ? this.state.allocations : ['*'])
    const owner = new V1Ownership()
    owner.setGroup(ownership.group)
    owner.setVersion(ownership.version)
    owner.setKind(ownership.kind)
    owner.setName(metadataName)
    owner.setNamespace(this.props.team)
    spec.setResource(owner)

    resource.setSpec(spec)

    return resource
  }

  handleSubmit = e => {
    e.preventDefault()

    this.setFormSubmitting()

    return this.props.form.validateFields(async (err, values) => {
      if (err) {
        this.setFormSubmitting(false, 'Validation failed')
        return
      }

      try {
        const resourceResult = await this.putResource(values)
        if (this.props.inlineVerification) {
          return await this.verify(resourceResult)
        }
        return await this.props.handleSubmit(resourceResult)
      } catch (err) {
        console.error('Error submitting form', err)
        this.setFormSubmitting(false, 'An error occurred saving the form, please try again')
      }
    })
  }

  continueWithoutVerification = async () => {
    try {
      const metadataName = this.getMetadataName(this.props.form.getFieldsValue())
      const resourceResult = await this.getResource(metadataName)
      await this.props.handleSubmit(resourceResult)
    } catch (err) {
      console.error('Error getting data', err)
      this.setFormSubmitting(false, 'An error occurred saving the form, please try again')
    }
  }

  render() {
    throw new Error('render must be implemented')
  }
}

export default VerifiedAllocatedResourceForm
