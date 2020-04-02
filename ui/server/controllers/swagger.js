const axios = require('axios')
const Router = require('express').Router
const url = require('url')

function swagger(koreApi) {
  return async (req, res) => {
    const u = url.parse(koreApi.url)
    const swaggerUrl = `${u.protocol}//${u.host}/swagger.json`
    try {
      const result = await axios['get'](swaggerUrl)
      return res.json(result.data)
    } catch (err) {
      const status = (err.response && err.response.status) || 500
      const message = (err.response && err.response.data && err.response.data.message) || err.message
      console.error(`Error making request to API with path ${swaggerUrl}`, status, message)
      return res.status(status).send()
    }
  }
}

function initRouter({ koreApi }) {
  const router = Router()
  router.use('/swagger.json', swagger(koreApi))
  return router
}

module.exports = {
  initRouter
}
