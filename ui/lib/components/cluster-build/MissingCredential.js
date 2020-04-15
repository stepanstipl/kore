import Link from 'next/link'
import PropTypes from 'prop-types'
import { Alert, Button } from 'antd'

const MissingCredential = ({ team }) => (
  <div>
    <Alert
      message="No credentials found"
      description="No credentials could be found allocated to this team, therefore you cannot request a cluster build at this time. Please continue to the team dashboard."
      type="info"
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

MissingCredential.propTypes = {
  team: PropTypes.string.isRequired
}

export default MissingCredential
