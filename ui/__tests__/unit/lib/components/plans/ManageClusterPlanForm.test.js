import { mount } from 'enzyme'

import ManageClusterPlanForm from '../../../../../lib/components/plans/ManageClusterPlanForm'
import ApiTestHelpers from '../../../../api-test-helpers'

describe('ManageClusterPlanForm', () => {
  let props
  let form
  let apiScope

  const schema = {
    properties: {
      bool1: { type: 'boolean' },
      bool2: { type: 'boolean' }
    }
  }

  beforeEach(async () => {
    // In case any tests leak to the API, mock out the API at this stage:
    apiScope = (ApiTestHelpers.getScope())
      .get(`${ApiTestHelpers.basePath}/planschemas/GKE`).reply(200, schema)
      .get(`${ApiTestHelpers.basePath}/accountmanagements`).reply(200, { items: [] })

    props = {
      form: {
        isFieldTouched: () => {},
        getFieldDecorator: jest.fn(() => c => c),
        getFieldsError: () => () => {},
        getFieldError: () => {},
        getFieldValue: () => {},
        validateFields: jest.fn()
      },
      kind: 'GKE',
      handleSubmit: jest.fn(),
      mode: 'edit'
    }
    mount(<ManageClusterPlanForm wrappedComponentRef={component => form = component} {...props} />)
    await form.componentDidMountComplete
  })

  afterEach(() => {
    // This will check that no calls were made against the API, unless the test registered them:
    apiScope.done()
  })

  describe('#generatePlanResource', () => {
    test('returns a configured Plan object when given valid values', () => {
      const plan = form.generatePlanResource({ description: 'Plan description', summary: 'Plan summary', configuration: { planProperty: 'value' } })
      expect(plan).toBeDefined()
    })
  })

  describe('#generatePlanConfiguration', () => {
    test('sets default false for boolean properties', () => {
      const config = form.generatePlanConfiguration()
      expect(config).toEqual({ bool1: false, bool2: false })
    })

    test('includes user configured values', () => {
      form.state.planValues = {
        string1: 'hello',
        number1: 123
      }
      const config = form.generatePlanConfiguration()
      expect(config).toEqual({ bool1: false, bool2: false, string1: 'hello', number1: 123 })
    })

    test('user configured values override default boolean values', () => {
      form.state.planValues = {
        bool1: true,
        string1: 'hello',
        number1: 123
      }
      const config = form.generatePlanConfiguration()
      expect(config).toEqual({ bool1: true, bool2: false, string1: 'hello', number1: 123 })
    })
  })

  describe('#handleSubmit', () => {
    let event
    beforeEach(() => {
      event = { preventDefault: jest.fn() }
      form.setFormSubmitting = jest.fn()
      props.form.validateFields.mockClear()
      form.handleSubmit(event)
    })

    test('prevents default', () => {
      expect(event.preventDefault).toHaveBeenCalledTimes(1)
    })

    test('sets form submitting in state', () => {
      expect(form.setFormSubmitting).toHaveBeenCalledTimes(1)
      expect(form.setFormSubmitting.mock.calls[0]).toEqual([])
    })

    test('validates fields', () => {
      expect(props.form.validateFields).toHaveBeenCalledTimes(1)
    })
  })

  describe('#process', () => {
    const planResource = {
      apiVersion: 'config.kore.appvia.io/v1',
      kind: 'Plan',
      metadata: { name: 'test-plan' },
      spec: {
        description: 'Test plan',
        summary: 'Summary of plan',
        kind: 'GKE',
        configuration: {
          bool1: false,
          bool2: false,
          prop1: 'ABC',
          prop2: 123,
          prop3: true
        }
      }
    }
    beforeEach(() => {
      form.setFormSubmitting = jest.fn()
      props.handleSubmit.mockClear()
      form.state.planValues = planResource.spec.configuration
    })

    test('handles form validation errors', async () => {
      await form.process('error', null)
      expect(form.setFormSubmitting).toHaveBeenCalledTimes(1)
      expect(form.setFormSubmitting.mock.calls[0]).toEqual([false, 'Validation failed'])
    })

    test('creates the resource and calls the wrapper component handleSubmit function', async () => {
      apiScope.put(`${ApiTestHelpers.basePath}/plans/test-plan`, planResource).reply(200, planResource)
      await form.process(null, { description: 'Test plan', summary: 'Summary of plan' })
      expect(props.handleSubmit).toHaveBeenCalledTimes(1)
      expect(props.handleSubmit.mock.calls[0]).toEqual([planResource])
      apiScope.done()
    })

    test('handles validation errors when creating the resource', async () => {
      const fieldErrors = [{ field: 'prop1', type: 'required', message: 'prop1 is required' }]
      apiScope.put(`${ApiTestHelpers.basePath}/plans/test-plan`, planResource).reply(400, { message: 'Validation errors', fieldErrors })
      await form.process(null, { description: 'Test plan', summary: 'Summary of plan' })

      expect(form.setFormSubmitting).toHaveBeenCalledTimes(1)
      expect(form.setFormSubmitting.mock.calls[0]).toEqual([false, 'Validation errors', fieldErrors])
      apiScope.done()
    })
  })
})
