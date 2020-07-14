module.exports = {
  setupFiles: ['<rootDir>/jest.setup.js'],
  testPathIgnorePatterns: [
    '<rootDir>/.next/',
    '<rootDir>/node_modules/',
    '<rootDir>/__tests__/api-test-helpers.js',
    '<rootDir>/__tests__/e2e/config.js',
    '<rootDir>/__tests__/e2e/page-objects'
  ],
  moduleNameMapper: {
    '\\.(css|less)$': '<rootDir>/__mocks__/styleMock.js',
    './request-interceptor': '<rootDir>/__mocks__/requestInterceptorMock.js'
  }
}
