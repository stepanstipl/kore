const axios = require('axios')
const Router = require('express').Router
const url = require('url')

function swagger(koreApiUrl) {
  return async (req, res) => {
    const u = url.parse(koreApiUrl)
    const swaggerUrl = `${u.protocol}//${u.host}/swagger.json`
    try {
      const result = await axios['get'](swaggerUrl)
      // Patch base path to our proxy URL:
      const patched = result.data
      for (const [key, value] of Object.entries(patched.paths)) {
        patched.paths[key.replace(u.path, '/apiproxy')] = value
        delete patched.paths[key]
      }
      return res.type('json').send(JSON.stringify(patched, null, 2))
    } catch (err) {
      const status = (err.response && err.response.status) || 500
      const message = (err.response && err.response.data && err.response.data.message) || err.message
      console.error(`Error making request to API with path ${swaggerUrl}`, status, message)
      return res.status(status).send()
    }
  }
}

function initRouter({ koreApiUrl }) {
  const router = Router()
  router.use('/swagger.json', swagger(koreApiUrl))
  return router
}

module.exports = {
  initRouter
}
