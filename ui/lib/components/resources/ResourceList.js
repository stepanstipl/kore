import React from 'react'
import PropTypes from 'prop-types'
import { message } from 'antd'

import copy from '../../utils/object-copy'

class ResourceList extends React.Component {

  static propTypes = {
    style: PropTypes.object
  }

  constructor(props) {
    super(props)
    this.state = {
      dataLoading: true,
      edit: false,
      add: false,
      view: false
    }
  }

  componentDidMount() {
    return this.fetchComponentData()
      .then(data => {
        this.setState({
          ...data,
          dataLoading: false
        })
      })
  }

  refresh = async () => {
    this.setState({ dataLoading: true })
    const data = await this.fetchComponentData()
    this.setState({
      ...data,
      dataLoading: false
    })
  }

  handleStatusUpdated = (updatedResource, done) => {
    const state = copy(this.state)
    const resource = state.resources.items.find(r => r.metadata.name === updatedResource.metadata.name)
    resource.status = updatedResource.status
    this.setState(state, done)
  }

  _setStateKey = (key, data) => {
    const state = copy(this.state)
    state[key] = data ? data : false
    this.setState(state)
  }

  view = resource => async () => this._setStateKey('view', resource)
  edit = resource => async () => this._setStateKey('edit', resource)
  add = enabled => async () => this._setStateKey('add', enabled)

  handleEditSave = updated => {
    const editedName = this.state.edit.metadata.name
    this.setState({
      ...this.state,
      edit: false,
      resources: {
        items: this.state.resources.items.map(resource => {
          if (resource.metadata.name === editedName) {
            return {
              ...resource,
              spec: updated.spec,
              allocation: updated.allocation,
              status: {
                ...resource.status, status: 'Pending'
              }
            }
          }
          return resource
        })
      }
    })
    message.success(this.updatedMessage)
  }

  handleAddSave = async created => {
    this.setState({
      ...this.state,
      add: false,
      resources: {
        items: this.state.resources.items.concat([ created ])
      }
    })
    message.success(this.createdMessage)
  }
}

export default ResourceList
