import * as React from 'react'
import PropTypes from 'prop-types'
import ServiceKindList from './ServiceKindList'
import { getKoreLabel } from '../../utils/crd-helpers';

export default class CloudServiceAdmin extends React.Component {
  static propTypes = {
    cloud: PropTypes.string
  }

  render() {
    const { cloud } = this.props
    return (
      <ServiceKindList filter={ (s) => getKoreLabel(s, 'platform') === cloud } />
    )
  }
}
