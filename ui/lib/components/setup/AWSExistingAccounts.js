import Link from 'next/link'
import { Button, Card, Icon, Result, Typography } from 'antd'
const { Paragraph } = Typography

import EKSCredentialsList from '../credentials/EKSCredentialsList'
import ExistingCloudAccounts from './ExistingCloudAccounts'

class AWSExistingAccounts extends ExistingCloudAccounts {

  stepsContentCreds = () => (
    <>
      <Paragraph style={{ fontSize: '16px', fontWeight: '600' }}>Add one or more AWS account credentials</Paragraph>
      <EKSCredentialsList getResourceItemList={this.setCredsCount} />
    </>
  )

  setupCompleteContent = () => (
    <Card>
      <Result
        status="success"
        title="Setup complete!"
        subTitle="Kore will use existing AWS accounts that it's given access to"
        extra={<Link href="/setup/kore/complete"><Button type="primary" key="continue">Continue</Button></Link>}
      >
        <Paragraph><Icon type="check-circle" theme="twoTone" twoToneColor="#52c41a" /> AWS account credentials</Paragraph>
        <Paragraph><Icon type="check-circle" theme="twoTone" twoToneColor="#52c41a" /> Account access guidance</Paragraph>
      </Result>
    </Card>
  )

}

export default AWSExistingAccounts
