module.exports = {
  testUrl: process.env.KORE_UI_TEST_URL || 'http://localhost:3000',
  adminPass: process.env.KORE_ADMIN_PASS || 'password',
  // Use longer timeout if we're showing the browser:
  timeout: process.env.SHOW_BROWSER !== 'true' ? 15000 : 60000,
  expectTimeout: 1000,
  drawerOpenClosePause: 500
}
