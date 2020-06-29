import { shallow } from 'enzyme'

import TeamPage from '../../../../../pages/teams/[name]/[tab]'

const props = {
  user: { id: 'jbloggs' },
  team: {
    metadata: { name: 'a-team' },
    spec: { summary: 'A Team' }
  },
  services: { items: [] },
  teamRemoved: jest.fn(),
  config: {
    featureGates: {}
  },
  tabActiveKey: 'clusters'
}

describe('TeamPage', () => {
  let teamPage
  let wrapper

  beforeEach(() => {
    wrapper = shallow(<TeamPage {...props} />)
    teamPage = wrapper.instance()
  })

  describe('#constructor', () => {

    it('sets initial state', () => {
      expect(teamPage.state).toEqual({
        tabActiveKey: 'clusters',
        memberCount: -1,
        clusterCount: -1,
        securityStatus: false
      })
    })

  })

})
