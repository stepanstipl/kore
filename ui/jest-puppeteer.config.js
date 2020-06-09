module.exports = {
  launch: {
    headless: process.env.SHOW_BROWSER !== 'true',
    slowMo: process.env.SHOW_BROWSER !== 'true' ? undefined : 15,
    args: ['--no-sandbox', '--start-maximized', '--window-size=1900,1000'],
    defaultViewport: {
      width: 1800,
      height: 900
    }
  }
}