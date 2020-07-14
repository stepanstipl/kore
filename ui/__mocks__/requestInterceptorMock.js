/**
 * No need to for request interceptor logic when running tests
 * just return the request untouched
 */
module.exports = () => (req) => req
