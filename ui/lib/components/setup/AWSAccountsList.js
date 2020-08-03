import React from 'react'
import { Alert, Icon, Table, Typography } from 'antd'
const { Paragraph } = Typography

import KoreApi from '../../kore-api'

export default class AWSAccountsList extends React.Component {
  state = {
    dataLoading: true
  }

  async fetchComponentData() {
    // TODO: use API endpoint to get data, awaiting https://github.com/appvia/kore/issues/1174
    return Promise.resolve({
      accounts: [
        {
          metadata: { name: 'kore-project-a-prod' },
          spec: {
            email: 'dave.thompson@appvia.io',
            ouName: 'Custom'
          }
        },
        {
          metadata: { name: 'kore-project-a-notprod' },
          spec: {
            email: 'dave.thompson@appvia.io',
            ouName: 'Custom'
          }
        }
      ]
    })
  }

  componentDidMount() {
    this.fetchComponentData().then(data => {
      this.setState({ ...data, dataLoading: false })
    })
  }

  render() {
    const { dataLoading, accounts } = this.state

    if (dataLoading) {
      return <Icon type="loading" />
    }

    const columns = [
      { title: 'Name', dataIndex: 'metadata.name' },
      { title: 'Email', dataIndex: 'spec.email' },
      { title: 'OU name', dataIndex: 'spec.ouName' }
    ]

    return (
      <>
        <Alert
          message="AWS accounts created by Kore"
          description={<>
            <Paragraph>Listed below are the AWS account created by Kore, these cannot be deleted by Kore once they are no longer used, please following the instructions linked below for information on how to delete them.</Paragraph>
            <Paragraph style={{ marginBottom: 0 }}>
              <a style={{ textDecoration: 'underline' }} href="https://docs.appvia.io/kore/guide/admin/aws_accounting/#deleting-aws-accounts-created-by-kore" target="_blank" rel="noopener noreferrer">https://docs.appvia.io/kore/guide/admin/aws_accounting/#deleting-aws-accounts-created-by-kore</a>
            </Paragraph>
          </>}
          type="info"
          showIcon
          style={{ marginBottom: '30px' }}
        />
        <Table columns={columns} dataSource={accounts} />
      </>
    )

  }
}
