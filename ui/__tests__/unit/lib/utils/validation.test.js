import { patterns } from '../../../../lib/utils/validation'

describe('validation', () => {
  describe('patterns', () => {
    describe('uriCompatible40CharMax', () => {
      const pattern = new RegExp(patterns.uriCompatible40CharMax.pattern)

      test('matches correctly', () => {
        expect('a1').toMatch(pattern)
        expect('a'.repeat(40)).toMatch(pattern)
        expect('sensible-test-string1').toMatch(pattern)
      })

      test('must be 40 chars or less', () => {
        expect('a'.repeat(41)).not.toMatch(pattern)
      })

      test('must be lowercase', () => {
        expect('string-with-UPPER-case').not.toMatch(pattern)
      })

      test('must start with letter', () => {
        expect('1-test-string').not.toMatch(pattern)
      })

      test('must only contain alphanumeric and hyphen', () => {
        expect('not_sensible_test_string1').not.toMatch(pattern)
      })

      test('must end with alphanumeric', () => {
        expect('a1-').not.toMatch(pattern)
      })
    })

    describe('uriCompatible63CharMax', () => {
      const pattern = new RegExp(patterns.uriCompatible63CharMax.pattern)

      test('matches correctly', () => {
        expect('a1').toMatch(pattern)
        expect('a'.repeat(63)).toMatch(pattern)
        expect('sensible-test-string1').toMatch(pattern)
      })

      test('must be 63 chars or less', () => {
        expect('a'.repeat(64)).not.toMatch(pattern)
      })

      test('must be lowercase', () => {
        expect('string-with-UPPER-case').not.toMatch(pattern)
      })

      test('must start with letter', () => {
        expect('1-test-string').not.toMatch(pattern)
      })

      test('must only contain alphanumeric and hyphen', () => {
        expect('not_sensible_test_string1').not.toMatch(pattern)
      })

      test('must end with alphanumeric', () => {
        expect('a1-').not.toMatch(pattern)
      })
    })
  })
})
