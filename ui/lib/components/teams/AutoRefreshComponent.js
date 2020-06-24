import * as React from 'react'
import PropTypes from 'prop-types'

/*
 Props that must be passed to the parent
   - refreshMs - how often to refresh the state
   - stableRefreshMs - how often to refresh the state once stable, if not defined no stable refresh will occur
   - propsResourceDataKey - where is the data located in the props
   - resourceApiRequest - function to call for requesting the resource
   - handleUpdate - function to call when resource is updated from the API
   - handleDelete - function to call when the resource is deleted from the API
 */

class AutoRefreshComponent extends React.Component {

  static propTypes = {
    refreshMs: PropTypes.number.isRequired,
    stableRefreshMs: PropTypes.number,
    propsResourceDataKey: PropTypes.string.isRequired,
    resourceApiRequest: PropTypes.func.isRequired,
    handleUpdate: PropTypes.func.isRequired,
    handleDelete: PropTypes.func.isRequired
  }

  static STABLE_STATES = {
    SUCCESS: 'Success',
    FAILURE: 'Failure'
  }

  getStableState() {
    const stateResource = this.props[this.props.propsResourceDataKey]
    const status = stateResource.status && stateResource.status.status

    return Object.keys(AutoRefreshComponent.STABLE_STATES)
      .map(k => AutoRefreshComponent.STABLE_STATES[k] === status ? status : false)
      .find(k => k)
  }

  resourceUpdated({ deleted, statusChanged }) {
    if (deleted) {
      this.stableStateReached && this.stableStateReached({ deleted: true })
      return clearInterval(this.interval)
    }
    const stableState = this.getStableState()
    if (stableState) {
      if (statusChanged && this.stableStateReached) {
        this.stableStateReached({ state: stableState })
      }
      clearInterval(this.interval)
      if (this.props.stableRefreshMs) {
        this.interval = setInterval(this.refreshResource,  this.props.stableRefreshMs)
      }
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
        const currentStatus = this.props[this.props.propsResourceDataKey].status.status
        const newStatus = resourceData.status.status
        const statusChanged = currentStatus !== newStatus
        this.props.handleUpdate(resourceData, () => this.resourceUpdated({ statusChanged }))
      }
    } catch (error) {
      // log the error but do nothing else, it will retry the refresh on the next interval
      console.error('Error refreshing resource', error)
    }
  }

  startRefreshing() {
    if (this.getStableState() && !this.props.stableRefreshMs) {
      return
    }
    if (this.interval) {
      clearInterval(this.interval)
    }
    this.interval = setInterval(this.refreshResource, this.getStableState() ? this.props.stableRefreshMs : this.props.refreshMs)
  }

  componentDidMount() {
    this.startRefreshing()
  }

  componentWillUnmount() {
    clearInterval(this.interval)
  }

}

export default AutoRefreshComponent
