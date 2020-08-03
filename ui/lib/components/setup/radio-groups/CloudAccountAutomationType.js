import PropTypes from 'prop-types'
import { Icon, Modal } from 'antd'
import { pluralize, titleize } from 'inflect'

import RadioGroup from '../../utils/RadioGroup'

const CloudAccountAutomationType = ({ cloud, accountNoun, value, onChange, disabled, inlineHelp, valuesFilter }) => {
  const accountAutomationClusterHelp = () => {
    Modal.info({
      title: `${titleize(accountNoun)} automation: One per cluster`,
      content: `For every cluster a team creates Kore will also create a ${cloud} ${accountNoun} and provision the cluster inside it. The ${cloud} ${accountNoun} will share the name given to the cluster.`,
      onOk() {},
      width: 500
    })
  }

  const accountAutomationCustomHelp = () => {
    Modal.info({
      title: `${titleize(accountNoun)} automation: Custom`,
      content: (
        <div>
          <p>When a team is created in Kore and a cluster is requested, Kore will ensure the associated {cloud} {accountNoun} is also created and the cluster placed inside it.</p>
          <p>You must also specify the plans available for each type of {accountNoun}, this is to ensure the correct cluster specification is being used.</p>
        </div>
      ),
      onOk() {},
      width: 500
    })
  }

  const options = [{
    value: 'CLUSTER',
    title: 'One per cluster',
    description: `Kore will create an ${cloud} ${accountNoun} for each cluster a team provisions`,
    className: 'automated-accounts-cluster',
    extra: inlineHelp ? <Icon style={{ marginTop: '28px' }} type="info-circle" theme="twoTone" onClick={accountAutomationCustomHelp}/> : undefined
  }, {
    value: 'CUSTOM',
    title: 'Custom',
    description: `Configure how Kore will create ${cloud} ${pluralize(accountNoun)} for teams`,
    className: 'automated-accounts-custom',
    extra: inlineHelp ? <Icon style={{ marginTop: '28px' }} type="info-circle" theme="twoTone" onClick={accountAutomationClusterHelp}/> : undefined
  }]
  const filteredOptions = options.filter(o => valuesFilter ? valuesFilter.includes(o.value) : true)

  return (
    <RadioGroup
      heading={`How do you want Kore to automate ${cloud} ${pluralize(accountNoun)} for teams?`}
      onChange={onChange}
      options={filteredOptions}
      value={value || ''}
      disabled={disabled}
      style={{ marginBottom: '15px' }}
    />
  )
}

CloudAccountAutomationType.propTypes = {
  cloud: PropTypes.oneOf(['GCP', 'AWS']),
  accountNoun: PropTypes.string.isRequired,
  value: PropTypes.oneOfType([PropTypes.bool, PropTypes.string]),
  onChange: PropTypes.func,
  disabled: PropTypes.bool,
  inlineHelp: PropTypes.bool,
  valuesFilter: PropTypes.arrayOf(PropTypes.oneOf(['CLUSTER', 'CUSTOM']))
}

export default CloudAccountAutomationType
