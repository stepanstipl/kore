import PropTypes from 'prop-types'
import { Typography } from 'antd'
import { pluralize } from 'inflect'

import RadioGroup from '../../utils/RadioGroup'

const KoreTeamCloudIntegration = ({ cloud, accountNoun, value, onChange, disabled }) => {
  const options = [{
    value: 'KORE',
    title: <>Kore managed {pluralize(accountNoun)} <Typography.Text type="secondary"> (recommended)</Typography.Text></>,
    description: `Kore will manage the ${cloud} ${pluralize(accountNoun)} required for teams`,
    className: 'use-kore-managed-projects'
  }, {
    value: 'EXISTING',
    title: `Use existing ${pluralize(accountNoun)}`,
    description: `Kore teams will use existing ${cloud} ${pluralize(accountNoun)} that it's given access to`,
    className: 'use-existing-projects'
  }]

  return (
    <RadioGroup
      heading={`How do you want Kore teams to integrate with ${cloud} ${pluralize(accountNoun)}?`}
      onChange={onChange}
      options={options}
      value={value}
      disabled={disabled}
      style={{ marginBottom: '15px' }}
    />
  )
}

KoreTeamCloudIntegration.propTypes = {
  cloud: PropTypes.oneOf(['GCP', 'AWS']),
  accountNoun: PropTypes.string.isRequired,
  value: PropTypes.string,
  onChange: PropTypes.func,
  disabled: PropTypes.bool
}

export default KoreTeamCloudIntegration
