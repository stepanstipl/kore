import { mount } from 'enzyme'

import PolicyForm from '../../../../../lib/components/policies/PolicyForm'
import ApiTestHelpers from '../../../../api-test-helpers'

describe('PolicyForm', () => {
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
      .get(`${ApiTestHelpers.basePath}/teams`).reply(200, { items: [] })
      .get(`${ApiTestHelpers.basePath}/teams/kore-admin/allocations/allow-gke-node-type-changes`).reply(404)

    props = {
      form: {
        isFieldTouched: () => {},
        getFieldDecorator: jest.fn(() => c => c),
        getFieldsError: () => () => {},
        getFieldError: () => {},
        getFieldValue: () => {},
        validateFields: jest.fn()
      },
      policy: {
        'kind':'PlanPolicy',
        'apiVersion':'config.kore.appvia.io/v1',
        'metadata':{
          'name': 'allow-gke-node-type-changes',
          'namespace':'kore-admin',
        },
        'spec':{
          'kind':'GKE',
          'summary':'Allow GKE node type changes',
          'description':'Teams with this policy applied may set the cluster node type',
          'properties':[
            { 'name':'machineType', 'allowUpdate':true, 'disallowUpdate':false }
          ]
        }
      },
      allocatedTeams: [],
      kind: 'GKE',
      handleSubmit: jest.fn()
    }
    mount(<PolicyForm wrappedComponentRef={component => form = component} {...props} />)
    await form.componentDidMountComplete
  })

  afterEach(() => {
    // This will check that no calls were made against the API, unless the test registered them:
    apiScope.done()
  })

  describe('#generatePolicyResource', () => {
    test('returns a configured Policy object when given valid values', () => {
      const plan = form.generatePolicyResource({ description: 'Policy description', summary: 'Policy name', properties: [{}] })
      expect(plan).toBeDefined()
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

  describe('#_process', () => {
    const policyResource = {
      'kind':'PlanPolicy',
      'apiVersion':'config.kore.appvia.io/v1',
      'metadata':{
        'name': 'allow-gke-node-type-changes',
        'namespace':'kore-admin',
      },
      'spec':{
        'kind':'GKE',
        'summary':'Allow GKE node type changes',
        'description':'Teams with this policy applied may set the cluster node type',
        'properties':[
          { 'name':'machineType', 'allowUpdate':true, 'disallowUpdate':false }
        ]
      }
    }

    beforeEach(() => {
      form.setFormSubmitting = jest.fn()
      props.handleSubmit.mockClear()
      form.state.policy = policyResource
      form.state.allocatedTeams = ['*']
    })

    test('handles form validation errors', async () => {
      await form._process('error', null)
      expect(form.setFormSubmitting).toHaveBeenCalledTimes(1)
      expect(form.setFormSubmitting.mock.calls[0]).toEqual([false, 'Validation failed'])
    })

    test('creates the resource and calls the wrapper component handleSubmit function', async () => {
      apiScope.put(`${ApiTestHelpers.basePath}/planpolicies/allow-gke-node-type-changes`).reply(200, policyResource)
      apiScope.put(`${ApiTestHelpers.basePath}/teams/kore-admin/allocations/allow-gke-node-type-changes`).reply(200)
      await form._process(null, { summary: 'Allow GKE node type changes', description: 'Description of policy' })
      expect(props.handleSubmit).toHaveBeenCalledTimes(1)
      expect(props.handleSubmit.mock.calls[0]).toEqual([policyResource])
      apiScope.done()
    })
  })
})
