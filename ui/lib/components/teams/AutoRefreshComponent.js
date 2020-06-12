import * as React from 'react'
import PropTypes from 'prop-types'

/*
 Props that must be passed to the parent
   - refreshMs - how often to refresh the state
   - propsResourceDataKey - where is the data located in the props
   - resourceApiRequest - function to call for requesting the resource
   - handleUpdate - function to call when resource is updated from the API
   - handleDelete - function to call when the resource is deleted from the API
 */

class AutoRefreshComponent extends React.Component {

  static propTypes = {
    refreshMs: PropTypes.number.isRequired,
    propsResourceDataKey: PropTypes.string.isRequired,
    resourceApiRequest: PropTypes.func.isRequired,
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

    return Object.keys(AutoRefreshComponent.FINAL_STATES)
      .map(k => AutoRefreshComponent.FINAL_STATES[k] === status ? status : false)
      .find(k => k)
  }

  resourceUpdated(params) {
    params = params || {}
    if (params.deleted) {
      this.finalStateReached && this.finalStateReached({ deleted: true })
      return clearInterval(this.interval)
    }
    const finalState = this.getFinalState()
    if (finalState) {
      this.finalStateReached && this.finalStateReached({ state: finalState })
      clearInterval(this.interval)
    }
  }

  async fetchResource() {
    const resourceData = await this.props.resourceApiRequest()
    return resourceData
  }

  refreshResource = async () => {
    try {
      const resourceName = this.props[this.props.propsResourceDataKey].metadata.name
      const resourceData = await this.fetchResource()
      if (!resourceData) {
        this.resourceUpdated({ deleted: true })
        this.props.handleDelete(resourceName)
      } else {
        this.props.handleUpdate(resourceData, () => this.resourceUpdated())
      }
    } catch (error) {
      // log the error but do nothing else, it will retry the refresh on the next interval
      console.error('Error refreshing resource', error)
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
