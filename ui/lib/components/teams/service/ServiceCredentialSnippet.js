import PropTypes from 'prop-types'
import { Alert, Button, Icon, Modal, Typography, Tooltip } from 'antd'
const { Paragraph } = Typography

class ServiceCredentialSnippet extends React.Component {
  static propTypes = {
    serviceCredential: PropTypes.object.isRequired
  }

  state = {
    visible: false
  }
  /* eslint-disable indent */
  render() {
    return (
      <>
      <Tooltip key="snippet" title="See usage snippet"><Icon onClick={() => this.setState({ visible: true })} style={{ marginLeft: '5px', color: '#3d5b58' }} type="eye" /></Tooltip>
      <Modal
        title="Service binding usage"
        visible={this.state.visible}
        onCancel={() => this.setState({ visible: false })}
        footer={[<Button type="primary" key="ok" onClick={() => this.setState({ visible: false })}>Ok</Button>]}
        width={700}
      >
        <div style={{ margin: '0 20px' }}>
          <Alert
            message="You can use the service binding secret in a Kubernetes Pod template using the following snippet."
            type="info"
            showIcon
            style={{ marginBottom: '20px' }}
          />
          <Paragraph copyable className="copy-command" style={{ whiteSpace: 'pre' }}>
{`envFrom:
  - secretRef:
      name: ${this.props.serviceCredential.spec.secretName}`}
          </Paragraph>
          <Paragraph>For more information, see <a target="_blank" rel="noopener noreferrer" href="https://kubernetes.io/docs/concepts/workloads/pods/pod-overview/#pod-templates">https://kubernetes.io/docs/concepts/workloads/pods/pod-overview/#pod-templates</a></Paragraph>
        </div>
      </Modal>
      </>
    )
  }
  /* eslint-enable indent */
}

export default ServiceCredentialSnippet
