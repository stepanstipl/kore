import { mount } from 'enzyme'

import EKSCredentialsForm from '../../../../../lib/components/credentials/EKSCredentialsForm'
import ApiTestHelpers from '../../../../api-test-helpers'

describe('EKSCredentialsForm', () => {
  let props
  let form
  let apiScope
  const secret = {
    metadata: { name: 'eks' },
    spec: { type: 'eks-credential' }
  }
  const eksCredential = {
    kind: 'EKSCredentials',
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
      team: 'kore-admin',
      allTeams: { items: [] },
      handleSubmit: jest.fn()
    }
    mount(<EKSCredentialsForm wrappedComponentRef={component => form = component} {...props} />)
  })

  afterEach(() => {
    // This will check that no calls were made against the API, unless the test registered them:
    apiScope.done()
  })

  describe('#generateSecretResource', () => {
    it('returns a configured Secret object', () => {
      const secret = form.generateSecretResource({ name: 'eks', accessKeyID: 'access-key', secretAccessKey: 'secret-key' })
      expect(secret).toBeDefined()
      expect(secret.spec.data.access_key_id).toBe(btoa('access-key'))
      expect(secret.spec.data.access_secret_key).toBe(btoa('secret-key'))
    })
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
        .get(`${ApiTestHelpers.basePath}/teams/kore-admin/ekscredentials/eks`).reply(200, eksCredential)
        .get(`${ApiTestHelpers.basePath}/teams/kore-admin/allocations/eks`).reply(200, allocation)
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
        .put(`${ApiTestHelpers.basePath}/teams/kore-admin/secrets/eks`).reply(200, secret)
        .put(`${ApiTestHelpers.basePath}/teams/kore-admin/ekscredentials/eks`).reply(200, eksCredential)
        .put(`${ApiTestHelpers.basePath}/teams/kore-admin/allocations/eks`).reply(200, allocation)
    })

    it('creates/updates and returns EKS credential and allocation from API', async () => {
      const result = await form.putResource({ name: 'eks', accountID: '1234567890', AccessKeyID: '123', secretAccessKey: 'aws-account-cred' })
      const expected = { ...eksCredential, allocation }
      expect(result).toEqual(expected)
      apiScope.done()
    })
  })
})
