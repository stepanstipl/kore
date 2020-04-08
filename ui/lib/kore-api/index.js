const Swagger = require('swagger-client')
const url = require('url')

const config = require('../../config')
const KoreApiClient = require('./kore-api-client')

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
    const spec = await KoreApi._getSpec()

    const api = await Swagger(
      KoreApi.swaggerUrl,
      {
        spec: spec,
        requestInterceptor: (req) => {
          if (process.browser) {
            req.url = req.url.replace(KoreApi.basePath, '/apiproxy')
          } else if (ctx && (ctx.id_token || ctx.req)) {
            req.headers['Authorization'] = `Bearer ${ctx.id_token || ctx.req.session.passport.user.id_token}`
          }
          return req
        }
      }
    )

    return new KoreApiClient(api)
  }

  /**
   * This returns the spec, from the cache if already downloaded, else from the relevant
   * swagger file.
   */
  static _getSpec = async () => {
    // Check if we need to download the swagger, caching it in a static if we do so we 
    // can re-use rather than downloading the swagger for every API call:
    // @TODO: Expire the cache after a while on the server.
    if (KoreApi.spec) { 
      return KoreApi.spec 
    }
    const u = url.parse(config.koreApi.url)
    KoreApi.basePath = u.path
    if (process.browser) {
      KoreApi.swaggerUrl = `${window.location.origin}/swagger.json`
    } else {
      KoreApi.swaggerUrl = `${u.protocol}//${u.host}/swagger.json`
    }
    console.log(`Initialising kore api swagger from ${KoreApi.swaggerUrl}`)
    // Need to disable eslint for this line as it's complaining it's a non-atomic update, doesn't
    // seem to matter how it's changed, it's still considered non-atomic, incorrectly:
    KoreApi.spec = (await Swagger(KoreApi.swaggerUrl)).spec // eslint-disable-line require-atomic-updates
    return KoreApi.spec
  }
}

module.exports = KoreApi