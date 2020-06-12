module.exports = {
  launch: {
    headless: process.env.SHOW_BROWSER !== 'true',
    slowMo: process.env.SHOW_BROWSER !== 'true' ? undefined : 15,
    args: ['--no-sandbox', '--start-maximized', '--window-size=1550,950'],
    defaultViewport: {
      width: 1525,
      height: 800
    }
  }
}