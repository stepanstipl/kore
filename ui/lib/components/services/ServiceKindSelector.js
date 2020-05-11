import * as React from 'react'
import PropTypes from 'prop-types'
import { Typography, Row, Col, Card, Icon } from 'antd'
const { Paragraph } = Typography

class ServiceKindSelector extends React.Component {
  static propTypes = {
    serviceKinds: PropTypes.object.isRequired,
    selectedServiceKind: PropTypes.string.isRequired,
    handleSelectKind: PropTypes.func.isRequired,
  }

  selectKind = kind => () => this.props.handleSelectKind(kind)

  render() {
    const { serviceKinds, selectedServiceKind } = this.props

    const cards = serviceKinds.items.filter(s => s.spec.enabled).map(serviceKind =>
      <Col span={6} key={serviceKind.metadata.name}>
        <Card
          id={serviceKind.metadata.name}
          onClick={this.selectKind(serviceKind.metadata.name)}
          hoverable={true}
          className={ selectedServiceKind === serviceKind.metadata.name ? 'service-kind-card selected' : 'service-kind-card' }
        >
          <Paragraph className="logo">
            { serviceKind.spec.imageURL ? (
              <img src={serviceKind.spec.imageURL} height="80px" />
            ) : (
              <Icon type="cloud-server" style={{ fontSize: '80px' }} theme="outlined" />
            ) }
          </Paragraph>
          <Paragraph className="name" strong>{serviceKind.spec.displayName || serviceKind.metadata.name}</Paragraph>
        </Card>
      </Col>
    )

    return (
      <Row gutter={[16,16]} type="flex" justify="center" style={{ marginTop: '40px', marginBottom: '40px' }}>
        { cards }
      </Row>
    )
  }
}

export default ServiceKindSelector
