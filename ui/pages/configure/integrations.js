import React from 'react'

import Breadcrumb from '../../lib/components/Breadcrumb'
import GKECredentialsList from '../../lib/components/configure/GKECredentialsList'
import GCPOrganizationsList from '../../lib/components/configure/GCPOrganizationsList'

const ConfigureIntegrationsPage = () => (
  <>
    <Breadcrumb items={[{ text: 'Configure' }, { text: 'Integrations' }]} />
    <GKECredentialsList style={{ marginBottom: '20px' }} />
    <GCPOrganizationsList />
  </>
)

export default ConfigureIntegrationsPage
