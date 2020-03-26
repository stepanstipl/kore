import React from 'react'

import Breadcrumb from '../../lib/components/Breadcrumb'
import GKECredentialsList from '../../lib/components/configure/GKECredentialsList'

const ConfigureIntegrationsPage = () => (
  <>
    <Breadcrumb items={[{ text: 'Configure' }, { text: 'Integrations' }]} />
    <GKECredentialsList />
  </>
)

export default ConfigureIntegrationsPage
