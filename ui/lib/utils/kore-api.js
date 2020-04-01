import Swagger from 'swagger-client'
import url from 'url'
import config from '../../config'
import redirect from './redirect'

class KoreApi {
  static spec = null
  static basePath = null
  static swaggerUrl = null

  /** 
   * Returns a client for accessing the Kore API. If this can be run server-side, you must
   * pass the context object which includes the request / session so this can be layered in
   * to the request. If this is called ONLY from the client side, you do not need to pass 
   * any ctx.
   */
  static client = async (ctx = null) => {
    // Check if we need to download the swagger, caching it in a static if we do so we 
    // can re-use rather than downloading the swagger for every API call:
    if (!(KoreApi.spec)) {
      const u = url.parse(config.koreApi.url)
      KoreApi.basePath = u.path
      if (process.browser) {
        KoreApi.swaggerUrl = `${window.location.origin}/swagger.json`
      } else {
        KoreApi.swaggerUrl = `${u.protocol}//${u.host}/swagger.json`
      }
      console.log(`Initialising kore api swagger from ${KoreApi.swaggerUrl}`)
      KoreApi.spec = (await Swagger(KoreApi.swaggerUrl)).spec
    }

    let api = await Swagger(
      KoreApi.swaggerUrl,
      {
        spec: KoreApi.spec,
        requestInterceptor: (req) => {
          if (process.browser) {
            req.url = req.url.replace(KoreApi.basePath, "/apiproxy")
          } else {
            req.headers['Authorization'] = `Bearer ${ctx.req.session.passport.user.id_token}`
          }
          return req
        }
      }
    )
    // At the moment, all API operations are untagged so sit in the 'default' space:
    return KoreApi.decorateFunctions(api).default
  }

  /**
   * This decorates every operation returned from the swagger with a function which unwraps the
   * returned object, making the usage of the api much cleaner in the rest of the code.
   * 
   * Also doing some global error handling too.
   */
  static decorateFunctions = (api) => {
    let apis = api.apis
    Object.keys(apis).forEach(tagName =>
      Object.entries(apis[tagName]).forEach(([functionName, fnc]) =>
        apis[tagName][functionName] = (...args) => fnc(...args).then(
          (res) => res.body,
          (err) => {
            // Handle not found as an empty object to reflect the way the existing api-request 
            // worked. 
            // @TODO: Would be better to return undefined or null.
            if (err.response && err.response.status === 404) return {}
            // Handle 401 unauth:
            if (err.response && err.response.status === 401)
              redirect(null, '/login/refresh', true)

            throw err
          })
      )
    )
    return apis
  }
}

export default KoreApi