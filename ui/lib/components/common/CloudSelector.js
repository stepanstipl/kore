import * as React from 'react'
import PropTypes from 'prop-types'
import { Typography, Row, Col, Card, Tag } from 'antd'
const { Title, Paragraph } = Typography

class CloudSelector extends React.Component {
  static DEFAULT_ENABLED_CLOUDS = ['GKE', 'EKS']

  static propTypes = {
    selectedCloud: PropTypes.oneOfType([PropTypes.string, PropTypes.bool]).isRequired,
    handleSelectCloud: PropTypes.func.isRequired,
    enabledCloudList: PropTypes.array,
    showCustom: PropTypes.bool
  }

  selectCloud = cloud => () => {
    if ((this.props.enabledCloudList || CloudSelector.DEFAULT_ENABLED_CLOUDS).includes(cloud)) {
      this.props.handleSelectCloud(cloud)
    }
  }

  render() {
    const { selectedCloud, showCustom } = this.props
    let enabledCloudList = this.props.enabledCloudList

    if (!enabledCloudList) {
      enabledCloudList = CloudSelector.DEFAULT_ENABLED_CLOUDS
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
            id="gcp"
            onClick={this.selectCloud('GKE')}
            hoverable={enabledCloudList.includes('GKE')}
            className={ selectedCloud === 'GKE' ? 'cloud-card selected' : 'cloud-card' }
          >
            {enabledCloudList.includes('GKE') ? (
              <>
                <Paragraph className="logo">
                  <img src="/static/images/GCP.png" height="80px" />
                </Paragraph>
                <Paragraph className="name" strong>Google Cloud Platform</Paragraph>
              </>
            ) : (
              <>
                <div className="unavailable">
                  <Paragraph>
                    <img src="/static/images/GCP.png" height="80px" />
                  </Paragraph>
                  <Paragraph strong style={{ textAlign: 'center', marginTop: '20px', marginBottom: '0' }}>Google Cloud Platform</Paragraph>
                </div>
                <ComingSoon />
              </>
            )}
          </Card>
        </Col>
        <Col span={6}>
          <Card
            id="aws"
            onClick={this.selectCloud('EKS')}
            hoverable={enabledCloudList.includes('EKS')}
            className={ selectedCloud === 'EKS' ? 'cloud-card selected' : 'cloud-card' }
          >
            {enabledCloudList.includes('EKS') ? (
              <>
                <Paragraph className="logo">
                  <img src="/static/images/AWS.png" height="80px" />
                </Paragraph>
                <Paragraph className="name" strong>Amazon Web Services</Paragraph>
              </>
            ) : (
              <>
                <div className="unavailable">
                  <Paragraph>
                    <img src="/static/images/AWS.png" height="80px" />
                  </Paragraph>
                  <Paragraph strong style={{ textAlign: 'center', marginTop: '20px', marginBottom: '0' }}>Amazon Web Services</Paragraph>
                </div>
                <ComingSoon />
              </>
            )}
          </Card>
        </Col>
        <Col span={6}>
          <Card
            id="azure"
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
