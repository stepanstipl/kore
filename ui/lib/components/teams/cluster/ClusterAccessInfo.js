import React from 'react'
import PropTypes from 'prop-types'
import { Badge, Button, Divider, Icon, Modal, Typography } from 'antd'
const { Paragraph, Text } = Typography
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

class ClusterAccessInfo extends React.Component {
  static propTypes = {
    team: PropTypes.object.isRequired,
    buttonStyle: PropTypes.object
  }

  state = {
    visible: false
  }

  infoItem = ({ num, title }) => (
    <div style={{ marginBottom: '10px' }}>
      <Badge style={{ backgroundColor: '#1890ff', marginRight: '10px' }} count={num} />
      <Text strong>{title}</Text>
    </div>
  )

  render() {
    const apiUrl = new URL(publicRuntimeConfig.koreApiPublicUrl)

    const profileConfigureCommand = `kore profile configure ${apiUrl.hostname}`
    const loginCommand = 'kore login'
    const kubeconfigCommand = `kore kubeconfig -t ${this.props.team.metadata.name}`

    return (
      <>
        <Button style={{ ...this.props.buttonStyle }} type="link" onClick={() => this.setState({ visible: true })}><Icon type="eye" />Access</Button>
        <Modal
          title="Cluster access"
          visible={this.state.visible}
          onCancel={() => this.setState({ visible: false })}
          footer={[<Button type="primary" key="ok" onClick={() => this.setState({ visible: false })}>Ok</Button>]}
          width={700}
        >
          <div style={{ margin: '0 20px' }}>
            {this.infoItem({ num: '1', title: 'Download' })}
            <Paragraph>If you haven&apos;t already, download the CLI from <a href="https://github.com/appvia/kore/releases">https://github.com/appvia/kore/releases</a></Paragraph>

            <Divider />

            {this.infoItem({ num: '2', title: 'Setup profile' })}
            <Paragraph>Create a profile</Paragraph>
            <Paragraph className="copy-command" copyable>{profileConfigureCommand}</Paragraph>
            <Paragraph>Enter the Kore API URL as follows</Paragraph>
            <Paragraph className="copy-command" copyable>{apiUrl.origin}</Paragraph>

            <Divider />

            {this.infoItem({ num: '3', title: 'Login' })}
            <Paragraph>Login to the CLI</Paragraph>
            <Paragraph className="copy-command" copyable>{loginCommand}</Paragraph>

            <Divider />

            {this.infoItem({ num: '4', title: 'Setup access' })}
            <Paragraph>Then, you can use the Kore CLI to setup access to your team&apos;s clusters</Paragraph>
            <Paragraph className="copy-command" copyable>{kubeconfigCommand}</Paragraph>
            <Paragraph>This will add local kubernetes configuration to allow you to use <Text
              style={{ fontFamily: 'monospace' }}>kubectl</Text> to talk to the provisioned cluster(s).</Paragraph>
            <Paragraph>See examples: <a href="https://kubernetes.io/docs/reference/kubectl/overview/" target="_blank" rel="noopener noreferrer">https://kubernetes.io/docs/reference/kubectl/overview/</a></Paragraph>
          </div>
        </Modal>
      </>
    )
  }
}

export default ClusterAccessInfo
