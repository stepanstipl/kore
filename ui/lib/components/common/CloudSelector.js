import * as React from 'react'
import PropTypes from 'prop-types'
import { Typography, Row, Col, Card, Tag } from 'antd'
const { Paragraph } = Typography

class CloudSelector extends React.Component {
  static DEFAULT_ENABLED_CLOUDS = ['GCP', 'AWS', 'Azure']

  static propTypes = {
    selectedCloud: PropTypes.oneOfType([PropTypes.string, PropTypes.bool]).isRequired,
    handleSelectCloud: PropTypes.func.isRequired,
    enabledCloudList: PropTypes.array
  }

  selectCloud = cloud => () => {
    if ((this.props.enabledCloudList || CloudSelector.DEFAULT_ENABLED_CLOUDS).includes(cloud)) {
      this.props.handleSelectCloud(cloud)
    }
  }

  render() {
    const { selectedCloud } = this.props
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
      <Row gutter={16} type="flex" justify="center" style={{ marginTop: '10px', marginBottom: '40px' }}>
        <Col span={8}>
          <Card
            id="gcp"
            onClick={this.selectCloud('GCP')}
            hoverable={enabledCloudList.includes('GCP')}
            className={ selectedCloud === 'GCP' ? 'cloud-card selected' : 'cloud-card' }
          >
            {enabledCloudList.includes('GCP') ? (
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
        <Col span={8}>
          <Card
            id="aws"
            onClick={this.selectCloud('AWS')}
            hoverable={enabledCloudList.includes('AWS')}
            className={ selectedCloud === 'AWS' ? 'cloud-card selected' : 'cloud-card' }
          >
            {enabledCloudList.includes('AWS') ? (
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
        <Col span={8}>
          <Card
            id="azure"
            onClick={this.selectCloud('Azure')}
            hoverable={enabledCloudList.includes('Azure')}
            className={ selectedCloud === 'Azure' ? 'cloud-card selected' : 'cloud-card' }
          >
            {enabledCloudList.includes('Azure') ? (
              <>
                <Paragraph className="logo">
                  <img src="/static/images/Azure.svg" height="80px" />
                </Paragraph>
                <Paragraph className="name" strong>Microsoft Azure</Paragraph>
              </>
            ) : (
              <>
                <div className="unavailable">
                  <Paragraph>
                    <img src="/static/images/Azure.svg" height="80px" />
                  </Paragraph>
                  <Paragraph strong style={{ textAlign: 'center', marginTop: '20px', marginBottom: '0' }}>Microsoft Azure</Paragraph>
                </div>
                <ComingSoon />
              </>
            )}
          </Card>
        </Col>
      </Row>
    )
  }
}

export default CloudSelector
