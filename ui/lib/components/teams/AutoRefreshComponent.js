import * as React from 'react'
import PropTypes from 'prop-types'
import apiRequest from '../../utils/api-request'

/*
 Props that must be passed to the parent
   - refreshMs - how often to refresh the state
   - propsResourceDataKey - where is the data located in the props
   - resourceApiPath - API path for requesting updated resource data
   - handleUpdate - function to call when resource is updated from the API
   - handleDelete - function to call when the resource is deleted from the API
 */

class AutoRefreshComponent extends React.Component {

  static propTypes = {
    refreshMs: PropTypes.number.isRequired,
    propsResourceDataKey: PropTypes.string.isRequired,
    resourceApiPath: PropTypes.string.isRequired,
    handleUpdate: PropTypes.func.isRequired,
    handleDelete: PropTypes.func.isRequired
  }

  static FINAL_STATES = {
    SUCCESS: 'Success',
    FAILURE: 'Failure'
  }

  getFinalState() {
    const stateResource = this.props[this.props.propsResourceDataKey]
    const status = stateResource.status && stateResource.status.status
    console.log('AutoRefreshComponent | getFinalState', status)

    return Object.keys(AutoRefreshComponent.FINAL_STATES)
      .map(k => AutoRefreshComponent.FINAL_STATES[k] === status ? status : false)
      .find(k => k)
  }

  resourceUpdated(params) {
    params = params || {}
    if (params.deleted) {
      console.log('AutoRefreshComponent | resource deleted')
      this.finalStateReached && this.finalStateReached({ deleted: true })
      return clearInterval(this.interval)
    }
    const finalState = this.getFinalState()
    console.log('AutoRefreshComponent | resource updated', finalState)
    if (finalState) {
      this.finalStateReached && this.finalStateReached({ state: finalState })
      clearInterval(this.interval)
    }
  }

  async fetchResource() {
    const resourceData = await apiRequest(null, 'get', this.props.resourceApiPath)
    return resourceData
  }

  refreshResource = async () => {
    const resourceName = this.props[this.props.propsResourceDataKey].metadata.name
    const resourceData = await this.fetchResource()
    console.log('AutoRefreshComponent | refreshResource', resourceName, resourceData)
    if (Object.keys(resourceData).length === 0) {
      this.resourceUpdated({ deleted: true })
      this.props.handleDelete(resourceName)
    } else {
      this.props.handleUpdate(resourceData, () => this.resourceUpdated())
    }
  }

  startRefreshing() {
    if (!this.getFinalState()) {
      this.interval = setInterval(this.refreshResource, this.props.refreshMs)
    }
  }

  componentDidMount() {
    this.startRefreshing()
  }

  componentWillUnmount() {
    clearInterval(this.interval)
  }

}

export default AutoRefreshComponent
