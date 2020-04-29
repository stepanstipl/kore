import config from '../../config'

function pageView(url) {
  window.dataLayer = window.dataLayer || []
  window.dataLayer.push({
    'event': 'PAGE_VIEW',
    'pagePath': url
  })
}

module.exports = {
  GTM_ID: config.kore.gtmId,
  pageView
}
