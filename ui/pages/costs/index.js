import React from 'react'

import Breadcrumb from '../../lib/components/layout/Breadcrumb'
import OverallCostsViewer from '../../lib/components/costs/OverallCostsViewer'

class CostsPage extends React.Component {

  state = {
    summary: null
  }

  static staticProps = {
    title: 'Costs Overview',
    adminOnly: true
  }

  render() {
    return (
      <div>
        <Breadcrumb items={[{ text: 'Costs' }, { text: 'Costs Summary' }]} />
        <OverallCostsViewer />
      </div>
    )
  }
}

export default CostsPage
