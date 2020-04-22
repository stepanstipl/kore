import { mount } from 'enzyme'

import GKECredentialsForm from '../../../../../lib/components/forms/GKECredentialsForm'
import ApiTestHelpers from '../../../../api-test-helpers'

describe('GKECredentialsForm', () => {
  let props
  let form
  let apiScope
  const gkeCredential = {
    metadata: { name: 'gke' },
    spec: { project: 'project-id', account: 'gke-service-account-cred' }
  }
  const allocation = {
    metadata: { name: 'gke' },
    spec: { resource: { kind: 'GKECredentials' } }
  }

  beforeEach(() => {
    // In case any tests leak to the API, mock out the API at this stage:
    apiScope = (ApiTestHelpers.getScope())

    props = {
      form: {
        isFieldTouched: () => {},
        getFieldDecorator: jest.fn(() => c => c),
        getFieldsError: () => () => {},
        getFieldError: () => {},
        getFieldValue: () => {},
        validateFields: jest.fn()
      },
      team: 'abc',
      allTeams: { items: [] },
      handleSubmit: jest.fn()
    }
    mount(<GKECredentialsForm wrappedComponentRef={component => form = component} {...props} />)
  })

  afterEach(() => {
    // This will check that no calls were made against the API, unless the test registered them:
    apiScope.done()
  })

  describe('#generateGKECredentialsResource', () => {
    it('returns a configured GKECredentials object when given valid values', () => {
      const gkeCredential = form.generateGKECredentialsResource({ name: 'gke', project: 'project-id', account: 'gke-service-account-cred' })
      expect(gkeCredential).toBeDefined()
    })
  })

  describe('#getResource', () => {
    beforeEach(() => {
      apiScope
        .get(`${ApiTestHelpers.basePath}/teams/abc/gkecredentials/gke`).reply(200, gkeCredential)
        .get(`${ApiTestHelpers.basePath}/teams/abc/allocations/gke`).reply(200, allocation)
    })

    it('returns GKE credential and allocation from API', async () => {
      const result = await form.getResource('gke')
      const expected = { ...gkeCredential, allocation }
      expect(result).toEqual(expected)
      apiScope.done()
    })
  })

  describe('#putResource', () => {
    beforeEach(() => {
      apiScope
        .put(`${ApiTestHelpers.basePath}/teams/abc/gkecredentials/gke`).reply(200, gkeCredential)
        .put(`${ApiTestHelpers.basePath}/teams/abc/allocations/gke`).reply(200, allocation)
    })

    it('creates/updates and returns GKE credential and allocation from API', async () => {
      const result = await form.putResource({ name: 'gke', project: 'project-id', account: 'gke-service-account-cred' })
      const expected = { ...gkeCredential, allocation }
      expect(result).toEqual(expected)
      apiScope.done()
    })
  })
})
