import React from 'react'
import PropTypes from 'prop-types'
import { successMessage } from '../../utils/message'

class ResourceList extends React.Component {

  static propTypes = {
    style: PropTypes.object,
    getResourceItemList: PropTypes.func,
    autoAllocateToAllTeams: PropTypes.bool
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
        this.props.getResourceItemList && this.props.getResourceItemList(data.resources.items)
        this.setState({
          ...data,
          dataLoading: false
        })
      })
  }

  refresh = async () => {
    this.setState({ dataLoading: true })
    const data = await this.fetchComponentData()
    this.props.getResourceItemList && this.props.getResourceItemList(data.resources.items)
    this.setState({
      ...data,
      dataLoading: false
    })
  }

  handleStatusUpdated = (updatedResource, done) => {
    this.setState((state) => {
      return {
        resources: {
          ...state.resources,
          items: state.resources.items.map((r) => r.metadata.name !== updatedResource.metadata.name ? r : { ...r, status: updatedResource.status })
        }
      }
    }, done)
  }

  view = (resource) => () => this.setState({ 'view': resource || false })
  edit = (resource) => () => this.setState({ 'edit': resource || false })
  add = (enabled) => () => this.setState({ 'add': enabled || false })

  handleEditSave = (updated) => {
    const editedName = this.state.edit.metadata.name
    this.setState((state) => {
      return {
        edit: false,
        resources: {
          items: state.resources.items.map(resource => {
            if (resource.metadata.name === editedName) {
              return {
                ...resource,
                spec: updated.spec,
                allocation: updated.allocation,
                status: {
                  ...resource.status, status: 'Pending'
                },
                ...updated.append
              }
            }
            return resource
          })
        }
      }
    })
    successMessage(this.updatedMessage)
  }

  handleAddSave = async (created) => {
    const newItems = this.state.resources.items.concat([ created ])
    this.setState({
      add: false,
      resources: {
        items: newItems
      }
    })
    this.props.getResourceItemList && this.props.getResourceItemList(newItems)
    successMessage(this.createdMessage)
  }
}

export default ResourceList
