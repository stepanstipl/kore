import { Form } from 'antd'
import { mount } from 'enzyme'

import VerifiedAllocatedResourceForm from '../../../../../lib/components/forms/VerifiedAllocatedResourceForm'
import ApiTestHelpers from '../../../../api-test-helpers'

const sampleResource = {
  metadata: { name: 'sample-resource' },
  spec: { property: 'value' }
}

// test wrapper for the VerifiedAllocatedResourceForm component
class TestForm extends VerifiedAllocatedResourceForm {
  allocationFormFieldsInfo = {
    allocationMissing: {
      infoMessage: 'allocationMissing info',
      infoDescription: 'allocationMissing description'
    },
    nameSection: {
      infoMessage: 'nameSection info',
      infoDescription: 'nameSection description',
      nameHelp: 'nameSection name help',
      descriptionHelp: 'nameSection description help'
    },
    allocationSection: {
      infoMessage: 'allocationSection info',
      infoDescription: 'allocationSection description',
      allTeamsWarningMessage: 'allocationSection allTeams warning',
      allTeamsWarningDescription: 'allocationSection allTeams description',
      allocateExtra: 'allocationSection allocate extra'
    }
  }
  getResource = jest.fn().mockResolvedValue(sampleResource)
  putResource = jest.fn().mockResolvedValue(sampleResource)
  resourceFormFields = () => <p>Form</p>
}
const WrappedTestForm = Form.create({ name: 'test_form' })(TestForm)

describe('VerifiedAllocatedResourceForm', () => {
  let props
  let form
  let apiScope

  beforeEach(() => {
    // In case any tests leak to the API, mock out the API at this stage:
    apiScope = (ApiTestHelpers.getScope())

    props = {
      form: {
        isFieldTouched: () => {},
        getFieldDecorator: jest.fn(() => c => c),
        getFieldsError: () => () => {},
        getFieldError: () => {},
        getFieldsValue: () => {},
        getFieldValue: () => {},
        validateFields: jest.fn()
      },
      team: 'abc',
      allTeams: { items: [] },
      handleSubmit: jest.fn()
    }
    mount(<WrappedTestForm wrappedComponentRef={component => form = component} {...props} />)
  })

  afterEach(() => {
    // This will check that no calls were made against the API, unless the test registered them:
    apiScope.done()
  })

  describe('#getMetadataName', () => {
    it('returns the canonical version of the name', () => {
      const name = form.getMetadataName({ name: 'This is the name!' })
      expect(name).toBe('this-is-the-name')
    })

    it('returns name from data if it exists', () => {
      let formWithData
      const propsWithData = { ...props, data: { metadata: { name: 'existing-name' } } }
      mount(<WrappedTestForm wrappedComponentRef={component => formWithData = component} {...propsWithData} />)
      const name = formWithData.getMetadataName({ name: 'This is the name!' })
      expect(name).toBe('existing-name')
    })
  })

  describe('#generateAllocationResource', () => {
    it('returns a configured Allocation object when given valid values', () => {
      const allocation = form.generateAllocationResource(
        { group: 'aws.compute.kore.appvia.io', version: 'v1alpha1', kind: 'EKSCredentials' },
        { name: 'Allocation name', summary: 'Summary of allocation' }
      )
      expect(allocation).toBeDefined()
    })
  })

  describe('#handleSubmit', () => {
    let event
    beforeEach(() => {
      event = { preventDefault: jest.fn() }
      form.setFormSubmitting = jest.fn()
      props.form.validateFields.mockClear()
    })

    it('prevents default', () => {
      form.handleSubmit(event)
      expect(event.preventDefault).toHaveBeenCalledTimes(1)
    })

    it('sets form submitting in state', () => {
      form.handleSubmit(event)
      expect(form.setFormSubmitting).toHaveBeenCalledTimes(1)
      expect(form.setFormSubmitting.mock.calls[0]).toEqual([])
    })

    it('validates fields', () => {
      form.handleSubmit(event)
      expect(props.form.validateFields).toHaveBeenCalledTimes(1)
    })
  })

  describe('#_process', () => {
    beforeEach(() => {
      form.putResource.mockClear()
      form.setFormSubmitting = jest.fn()
      props.handleSubmit.mockClear()
    })

    it('handles form validation errors', async () => {
      await form._process('error', null)
      expect(form.setFormSubmitting).toHaveBeenCalledTimes(1)
      expect(form.setFormSubmitting.mock.calls[0]).toEqual([false, 'Validation failed'])
    })

    it('creates the resource and calls the wrapper component handleSubmit function', async () => {
      await form._process(null, { property: 'value' })
      expect(form.putResource).toHaveBeenCalledTimes(1)
      expect(form.putResource.mock.calls[0]).toEqual([{ property: 'value' }])
      expect(props.handleSubmit).toHaveBeenCalledTimes(1)
      expect(props.handleSubmit.mock.calls[0]).toEqual([sampleResource])
    })

    it('handles errors creating the resource', async () => {
      form.putResource.mockRejectedValue(new Error('PUT resource error'))
      await form._process(null, { property: 'value' })
      expect(form.setFormSubmitting).toHaveBeenCalledTimes(1)
      expect(form.setFormSubmitting.mock.calls[0]).toEqual([false, 'An error occurred saving the form, please try again'])
    })
  })

  describe('#continueWithoutVerification', () => {
    beforeEach(() => {
      form.getResource.mockClear()
      props.form.getFieldsValue = jest.fn().mockReturnValue({ name: 'Resource name' })
      props.handleSubmit.mockClear()
    })

    it('gets resource amd calls the wrapper component handleSubmit function', async () => {
      await form.continueWithoutVerification()
      expect(form.getResource).toHaveBeenCalledTimes(1)
      expect(form.getResource.mock.calls[0]).toEqual(['resource-name'])
      expect(props.handleSubmit).toHaveBeenCalledTimes(1)
      expect(props.handleSubmit.mock.calls[0]).toEqual([sampleResource])
    })
  })
})
