import Link from 'next/link'
import PropTypes from 'prop-types'
import { Alert, Button, Typography } from 'antd'
const { Paragraph, Text } = Typography

const MissingCredential = ({ team, cloud }) => {
  const message = {
    'GCP': 'GCP project access not found',
    'AWS': 'AWS account access not found'
  }
  const description = {
    'GCP': (
      <>
        <Paragraph>This team does not have access to create clusters in any GCP projects. Please use the contact below to grant this team access to a GCP project.</Paragraph>
        <Text strong>Kore administrator</Text>
      </>
    ),
    'AWS': (
      <>
        <Paragraph>This team does not have access to create clusters in any AWS accounts. Please use the contact below to grant this team access to an AWS account.</Paragraph>
        <Text strong>Kore administrator</Text>
      </>
    )
  }
  return (
    <div>
      <Alert
        message={message[cloud]}
        description={description[cloud]}
        type="warning"
        showIcon
        style={{ marginBottom: '20px' }}
      />
      <Button type="primary">
        <Link href="/teams/[name]" as={`/teams/${team}`}>
          <a>Team dashboard</a>
        </Link>
      </Button>
    </div>
  )
}

MissingCredential.propTypes = {
  team: PropTypes.string.isRequired,
  cloud: PropTypes.string.isRequired,
}

export default MissingCredential
