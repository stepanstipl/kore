import * as React from 'react'
import PropTypes from 'prop-types'
import { Typography, Row, Col, Card, Tag } from 'antd'
const { Paragraph, Text } = Typography

class CloudSelector extends React.Component {
  static propTypes = {
    selectedCloud: PropTypes.string.isRequired,
    handleSelectCloud: PropTypes.func.isRequired,
    credentials: PropTypes.object
  }

  selectCloud = cloud => () => this.props.handleSelectCloud(cloud)

  render() {
    const { selectedCloud } = this.props

    const ComingSoon = () => (
      <div style={{
        position: 'absolute',
        left: '0',
        width: '100%',
        textAlign: 'center',
        top: '10px'
      }}>
        <Tag color="#2db7f5">Coming soon!</Tag>
      </div>
    )

    return (
      <Row gutter={16} type="flex" justify="center" style={{ marginBottom: '20px' }}>
        <Col span={12}>
          <Card
            id="aws"
            onClick={this.selectCloud('AWS')}
            hoverable={true}
            className={ selectedCloud === 'AWS' ? 'cloud-card selected' : 'cloud-card' }
          >
            <Paragraph className="logo" style={{ marginBottom: '0' }}>
              <img src="/static/images/AWS.png" height="80px" />
              <Text className="name" strong style={{ marginLeft: '15px', fontSize: '16px' }}>Amazon Web Services</Text>
            </Paragraph>
          </Card>
        </Col>
        <Col span={12}>
          <Card
            id="gcp"
            hoverable={false}
            className={ selectedCloud === 'GCP' ? 'cloud-card selected' : 'cloud-card' }
          >
            <Paragraph className="logo" style={{ marginBottom: '0' }}>
              <img src="/static/images/GCP.png" height="80px" />
              <Text className="name" strong style={{ marginLeft: '15px', fontSize: '16px' }}>
                Google Cloud Platform
                <ComingSoon />
              </Text>
            </Paragraph>
          </Card>
        </Col>
      </Row>
    )
  }
}

export default CloudSelector
