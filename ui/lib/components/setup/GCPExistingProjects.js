import React from 'react'
import Link from 'next/link'
import { Button, Card, Icon, Result, Typography } from 'antd'
const { Paragraph } = Typography

import GKECredentialsList from '../credentials/GKECredentialsList'
import ExistingCloudAccounts from './ExistingCloudAccounts'

class GCPExistingProjects extends ExistingCloudAccounts {

  stepsContentCreds = () => (
    <>
      <Paragraph style={{ fontSize: '16px', fontWeight: '600' }}>Add one or more GCP project credentials</Paragraph>
      <GKECredentialsList getResourceItemList={this.setCredsCount} />
    </>
  )

  setupCompleteContent = () => (
    <Card>
      <Result
        status="success"
        title="Setup complete!"
        subTitle="Kore will use existing GCP projects that it's given access to"
        extra={<Link href="/setup/kore/complete"><Button type="primary" key="continue">Continue</Button></Link>}
      >
        <Paragraph><Icon type="check-circle" theme="twoTone" twoToneColor="#52c41a" /> GCP project credentials</Paragraph>
        <Paragraph><Icon type="check-circle" theme="twoTone" twoToneColor="#52c41a" /> Project access guidance</Paragraph>
      </Result>
    </Card>
  )
}

export default GCPExistingProjects
