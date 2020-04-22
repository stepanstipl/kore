import { mount } from 'enzyme'

import EKSCredentialsForm from '../../../../../lib/components/forms/EKSCredentialsForm'
import ApiTestHelpers from '../../../../api-test-helpers'

describe('EKSCredentialsForm', () => {
  let props
  let form
  let apiScope
  const eksCredential = {
    metadata: { name: 'eks' },
    spec: { accountID: '1234567890', accessKeyID: '123', secretAccessKey: 'aws-account-cred' }
  }
  const allocation = {
    metadata: { name: 'eks' },
    spec: { resource: { kind: 'EKSCredentials' } }
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
    mount(<EKSCredentialsForm wrappedComponentRef={component => form = component} {...props} />)
  })

  afterEach(() => {
    // This will check that no calls were made against the API, unless the test registered them:
    apiScope.done()
  })

  describe('#generateEKSCredentialsResource', () => {
    it('returns a configured EKSCredentials object when given valid values', () => {
      const eksCredential = form.generateEKSCredentialsResource({ name: 'eks', accountID: '1234567890', accessKeyID: '123', secretAccessKey: 'aws-account-cred' })
      expect(eksCredential).toBeDefined()
    })
  })

  describe('#getResource', () => {
    beforeEach(() => {
      apiScope
        .get(`${ApiTestHelpers.basePath}/teams/abc/ekscredentials/eks`).reply(200, eksCredential)
        .get(`${ApiTestHelpers.basePath}/teams/abc/allocations/eks`).reply(200, allocation)
    })

    it('returns EKS credential and allocation from API', async () => {
      const result = await form.getResource('eks')
      const expected = { ...eksCredential, allocation }
      expect(result).toEqual(expected)
      apiScope.done()
    })
  })

  describe('#putResource', () => {
    beforeEach(() => {
      apiScope
        .put(`${ApiTestHelpers.basePath}/teams/abc/ekscredentials/eks`).reply(200, eksCredential)
        .put(`${ApiTestHelpers.basePath}/teams/abc/allocations/eks`).reply(200, allocation)
    })

    it('creates/updates and returns EKS credential and allocation from API', async () => {
      const result = await form.putResource({ name: 'eks', accountID: '1234567890', AccessKeyID: '123', secretAccessKey: 'aws-account-cred' })
      const expected = { ...eksCredential, allocation }
      expect(result).toEqual(expected)
      apiScope.done()
    })
  })
})
