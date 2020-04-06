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
        let state = copy(this.state)
        state = { ...state, ...data }
        state.dataLoading = false
        this.setState(state)
      })
  }

  handleStatusUpdated = (updatedResource, done) => {
    const state = copy(this.state)
    const resource = state.resources.items.find(r => r.metadata.name === updatedResource.metadata.name)
    resource.status = updatedResource.status
    this.setState(state, done)
  }

  edit = resource => {
    return async () => {
      const state = copy(this.state)
      state.edit = resource ? resource : false
      this.setState(state)
    }
  }

  view = resource => {
    return async () => {
      const state = copy(this.state)
      state.view = resource ? resource : false
      this.setState(state)
    }
  }

  handleEditSave = updated => {
    const state = copy(this.state)

    const edited = state.resources.items.find(c => c.metadata.name === state.edit.metadata.name)
    edited.spec = updated.spec
    edited.allocation = updated.allocation
    edited.status.status = 'Pending'

    state.edit = false
    this.setState(state)
    message.success(this.updatedMessage)
  }

  add = enabled => {
    return () => {
      const state = copy(this.state)
      state.add = enabled
      this.setState(state)
    }
  }

  handleAddSave = async created => {
    const state = copy(this.state)
    state.resources.items.push(created)
    state.add = false
    this.setState(state)
    message.success(this.createdMessage)
  }
}

export default ResourceList
