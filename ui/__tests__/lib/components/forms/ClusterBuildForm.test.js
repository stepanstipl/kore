import { mount } from 'enzyme'

import ClusterBuildForm from '../../../../lib/components/forms/ClusterBuildForm'
import ApiTestHelpers from '../../../api-test-helpers'

describe('ClusterBuildForm', () => {
  let props
  let form
  let apiScope

  let plans = { 
    items: [
      { spec: { kind: 'GKE' }, metadata: { name: 'GKE Development' } },
      { spec: { kind: 'GKE' }, metadata: { name: 'GKE Production' } }
    ] 
  }
  let allocations = [
    { spec: { resource: { kind: 'GKECredentials' } }, metadata: { name: 'GKE' } },
    { spec: { resource: { kind: 'EKSCredentials' } }, metadata: { name: 'EKS' } }
  ]

  beforeEach(async () => {
    // In case any tests leak to the API, mock out the API at this stage:
    apiScope = (ApiTestHelpers.getScope())
      .get(`${ApiTestHelpers.basePath}/teams/abc/allocations?assigned=true`).reply(200, { items: allocations })
      .get(`${ApiTestHelpers.basePath}/plans`).reply(200, plans)

    props = {
      form: {
        isFieldTouched: () => {},
        getFieldDecorator: jest.fn(() => c => c),
        getFieldsError: () => () => {},
        getFieldError: () => {},
        getFieldValue: () => {},
        validateFields: jest.fn()
      },
      team: { metadata: { name: 'abc' } },
      teamClusters: [],
      user: { id: 'jbloggs' }
    }
    mount(<ClusterBuildForm wrappedComponentRef={component => form = component} {...props} />)
    // Wait for the mount to complete asynchronously
    await form.componentDidMountComplete
  })

  afterEach(() => {
    // This will check that no calls were made against the API, unless the test registered them:
    apiScope.done()
  })

  describe('#componentDidMount', () => {
    it('should load allocations and plans', async () => {
      // Check API has been accessed as expected.
      apiScope.done()
      expect(form.state.plans).toEqual(plans)
      expect(form.state.credentials.GKE).toEqual([{ ...allocations[0] }])
      expect(form.state.credentials.EKS).toEqual([{ ...allocations[1] }])
    })
  })

  describe('#getClusterResource', () => {
    it('should return a configured cluster object when given valid values', () => {
      form.handleSelectCloud('GKE')
      const cluster = form.getClusterResource({ credential: 'GKE', plan: plans.items[0].metadata.name, clusterName: 'abc-test-cluster' })
      expect(cluster).toBeDefined()
    })
  })

  describe('#handleSubmit', () => {
    // @TODO: Test handleSubmit and validation.

    // let event
    // beforeEach(() => {
    //   props.form.validateFields.mockClear()
    //   event = {
    //     preventDefault: jest.fn()
    //   }
    // })

    // it('prevents default action, sets saving state and validates the fields', () => {
    //   form.handleSubmit(event)
    //   expect(event.preventDefault).toHaveBeenCalledTimes(1)
    //   expect(form.state.submitting).toEqual(true)
    //   expect(form.state.formErrorMessage).toBeFalsy()
    //   expect(props.form.validateFields).toHaveBeenCalledTimes(1)
    // })
  })
})
