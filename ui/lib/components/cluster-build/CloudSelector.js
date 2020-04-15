import * as React from 'react'
import PropTypes from 'prop-types'
import { Typography, Row, Col, Card, Tag } from 'antd'
const { Title, Paragraph, Text } = Typography

class CloudSelector extends React.Component {
  static propTypes = {
    selectedCloud: PropTypes.string.isRequired,
    handleSelectCloud: PropTypes.func.isRequired,
    credentials: PropTypes.object,
    showCustom: PropTypes.bool
  }

  selectCloud = cloud => () => this.props.handleSelectCloud(cloud)

  render() {
    const { selectedCloud, credentials, showCustom } = this.props

    const Credentials = ({ cloud }) => {
      if (!credentials) {
        return null
      }
      const credType = cloud === 'GKE' ? 'project' : 'account'
      const cloudProviderCount = credentials[cloud].length
      return (
        <Paragraph style={{ textAlign: 'center', marginTop: '20px', marginBottom: '0' }}>
          {cloudProviderCount > 0 ?
            <Tag color="#87d068">{cloudProviderCount} {credType} credential{cloudProviderCount > 1 ? 's' : ''}</Tag> :
            <Text type="warning">No credentials</Text>
          }
        </Paragraph>
      )
    }

    const ComingSoon = () => (
      <div className="coming-soon">
        <Tag color="#2db7f5">Coming soon!</Tag>
      </div>
    )

    return (
      <Row gutter={16} type="flex" justify="center" style={{ marginTop: '40px', marginBottom: '40px' }}>
        <Col span={6}>
          <Card
            onClick={this.selectCloud('GKE')}
            hoverable={true}
            className={ selectedCloud === 'GKE' ? 'cloud-card selected' : 'cloud-card' }
          >
            <Paragraph className="logo">
              <img src="/static/images/GCP.png" height="80px" />
            </Paragraph>
            <Paragraph className="name" strong>Google Cloud Platform</Paragraph>
            <Credentials cloud="GKE" />
          </Card>
        </Col>
        <Col span={6}>
          <Card
            onClick={this.selectCloud('EKS')}
            hoverable={true}
            className={ selectedCloud === 'EKS' ? 'cloud-card selected' : 'cloud-card' }
          >
            <Paragraph className="logo">
              <img src="/static/images/AWS.png" height="80px" />
            </Paragraph>
            <Paragraph className="name" strong>Amazon Web Services</Paragraph>
            <Credentials cloud="EKS" />
          </Card>
        </Col>
        <Col span={6}>
          <Card
            hoverable={false}
            className={ selectedCloud === 'AKS' ? 'cloud-card selected' : 'cloud-card' }
          >
            <div className="unavailable">
              <Paragraph style={{ paddingBottom: '15px', marginTop: '15px' }}>
                <img src="/static/images/Azure.svg" height="50px" />
              </Paragraph>
              <Paragraph strong style={{ textAlign: 'center', marginTop: '20px', marginBottom: '0' }}>Microsoft Azure</Paragraph>
            </div>
            <ComingSoon />
          </Card>
        </Col>
        {showCustom ? (
          <Col span={6}>
            <Card
              hoverable={false}
              className={ selectedCloud === 'CUSTOM' ? 'cloud-card selected' : 'cloud-card' }
            >
              <div className="unavailable">
                <Title level={3} style={{ paddingTop: '30px', height: '80px' }}>Custom</Title>
                <Paragraph strong style={{ marginTop: '20px' }}>Bring your own cluster</Paragraph>
              </div>
              <ComingSoon />
            </Card>
          </Col>
        ) : null}
      </Row>
    )
  }
}

export default CloudSelector
