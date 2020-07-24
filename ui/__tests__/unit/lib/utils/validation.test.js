import { patterns } from '../../../../lib/utils/validation'

describe('validation', () => {
  describe('patterns', () => {
    describe('uriCompatible40CharMax', () => {
      const pattern = new RegExp(patterns.uriCompatible40CharMax.pattern)

      it('matches correctly', () => {
        expect('a1').toMatch(pattern)
        expect('a'.repeat(40)).toMatch(pattern)
        expect('sensible-test-string1').toMatch(pattern)
      })

      it('must be 40 chars or less', () => {
        expect('a'.repeat(41)).not.toMatch(pattern)
      })

      it('must be lowercase', () => {
        expect('string-with-UPPER-case').not.toMatch(pattern)
      })

      it('must start with letter', () => {
        expect('1-test-string').not.toMatch(pattern)
      })

      it('must only contain alphanumeric and hyphen', () => {
        expect('not_sensible_test_string1').not.toMatch(pattern)
      })

      it('must end with alphanumeric', () => {
        expect('a1-').not.toMatch(pattern)
      })
    })

    describe('uriCompatible63CharMax', () => {
      const pattern = new RegExp(patterns.uriCompatible63CharMax.pattern)

      it('matches correctly', () => {
        expect('a1').toMatch(pattern)
        expect('a'.repeat(63)).toMatch(pattern)
        expect('sensible-test-string1').toMatch(pattern)
      })

      it('must be 63 chars or less', () => {
        expect('a'.repeat(64)).not.toMatch(pattern)
      })

      it('must be lowercase', () => {
        expect('string-with-UPPER-case').not.toMatch(pattern)
      })

      it('must start with letter', () => {
        expect('1-test-string').not.toMatch(pattern)
      })

      it('must only contain alphanumeric and hyphen', () => {
        expect('not_sensible_test_string1').not.toMatch(pattern)
      })

      it('must end with alphanumeric', () => {
        expect('a1-').not.toMatch(pattern)
      })
    })

    describe('amazonIamRoleArn', () => {
      const pattern = new RegExp(patterns.amazonIamRoleArn.pattern)
      it('matches correctly', () => {
        expect('arn:aws:iam::111222333444:role/some-role').toMatch(pattern)
        expect('arn:aws:iam::555555666666:role/anotherrole').toMatch(pattern)
      })

      it('must start with valid prefix', () => {
        expect('4rn:4ws:1am::111222333444:role/some-role').not.toMatch(pattern)
      })

      it('must include 12 digit account number', () => {
        expect('arn:aws:iam::123:role/some-role').not.toMatch(pattern)
      })

      it('must include role suffix', () => {
        expect('arn:aws:iam::111222333444:notrole/some-role').not.toMatch(pattern)
      })

      it('must include the role name', () => {
        expect('arn:aws:iam::111222333444:role/').not.toMatch(pattern)
      })
    })
  })
})
