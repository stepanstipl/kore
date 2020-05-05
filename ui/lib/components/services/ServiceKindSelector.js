import * as React from 'react'
import PropTypes from 'prop-types'
import { Typography, Row, Col, Card, Icon } from 'antd'
const { Paragraph } = Typography

class ServiceKindSelector extends React.Component {
  static propTypes = {
    selectedServiceKind: PropTypes.string.isRequired,
    handleSelectKind: PropTypes.func.isRequired,
  }

  selectKind = kind => () => this.props.handleSelectKind(kind)

  render() {
    const { selectedServiceKind } = this.props

    return (
      <Row gutter={16} type="flex" justify="center" style={{ marginTop: '40px', marginBottom: '40px' }}>
        <Col span={6}>
          <Card
            id="dummy"
            onClick={this.selectKind('dummy')}
            hoverable={true}
            className={ selectedServiceKind === 'dummy' ? 'service-kind-card selected' : 'service-kind-card' }
          >
            <Paragraph className="logo">
              <Icon type="question-circle" style={{ fontSize: '80px' }} theme="outlined" />
            </Paragraph>
            <Paragraph className="name" strong>Dummy service</Paragraph>
          </Card>
        </Col>
      </Row>
    )
  }
}

export default ServiceKindSelector
