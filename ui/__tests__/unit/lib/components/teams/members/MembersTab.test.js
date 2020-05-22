import { mount } from 'enzyme'

import MembersTab from '../../../../../../lib/components/teams/members/MembersTab'
import ApiTestHelpers from '../../../../../api-test-helpers'

describe('MembersTab', () => {
  const initialMembers = ['bob', 'jane', 'harold']
  let apiScope
  let props
  let membersTab

  beforeEach(async () => {
    apiScope = (ApiTestHelpers.getScope())
      .get(`${ApiTestHelpers.basePath}/teams/a-team/members`).reply(200, { items: initialMembers })
      .get(`${ApiTestHelpers.basePath}/users`).reply(200, { items: [{ spec: { username: 'admin' } }, { spec: { username: 'julie' } }, { spec: { username: 'ellis' } }] })
      .get(`${ApiTestHelpers.basePath}/teams/a-team/invites/generate?expire=1h`).reply(200, 'https://kore.appvia.io/process/teams/invitation/123')
    props = {
      user: { id: 'jbloggs' },
      team: {
        metadata: { name: 'a-team' },
        spec: { summary: 'A Team' }
      },
      getMemberCount: jest.fn()
    }

    membersTab = mount(<MembersTab {...props} />).instance()
    await membersTab.componentDidMountComplete
  })

  afterEach(() => {
    // This will check that no calls were made against the API, unless the test registered them:
    apiScope.done()
  })

  test('sets up initial state', () => {
    expect(membersTab.state).toEqual({
      dataLoading: false,
      members: initialMembers,
      membersToAdd: [],
      users: ['julie', 'ellis']
    })
  })

  test('updates member count', () => {
    expect(props.getMemberCount).toHaveBeenCalledTimes(1)
    expect(props.getMemberCount.mock.calls[0]).toEqual([ 3 ])
  })

  describe('#addTeamMembers', () => {
    beforeEach(() => {
      apiScope
        .put(`${ApiTestHelpers.basePath}/teams/a-team/members/one`).reply(200)
        .put(`${ApiTestHelpers.basePath}/teams/a-team/members/two`).reply(200)
      props.getMemberCount.mockClear()
    })

    test('makes api request to add each member, sets state and updates member count', async () => {
      membersTab.state.membersToAdd = ['one', 'two']
      await membersTab.addTeamMembers()

      apiScope.done()
      expect(membersTab.state.membersToAdd).toEqual([])
      expect(membersTab.state.members).toEqual(['bob', 'jane', 'harold', 'one', 'two'])
      expect(props.getMemberCount).toHaveBeenCalledTimes(1)
      expect(props.getMemberCount.mock.calls[0]).toEqual([ 5 ])
    })
  })

  describe('#deleteTeamMember', () => {
    beforeEach(() => {
      apiScope.delete(`${ApiTestHelpers.basePath}/teams/a-team/members/bob`).reply(200)
      props.getMemberCount.mockClear()
    })

    test('makes api request to remove the member, sets state and updates member count', async () => {
      await membersTab.deleteTeamMember('bob')()

      apiScope.done()
      expect(membersTab.state.members).toEqual(['jane', 'harold'])
      expect(props.getMemberCount).toHaveBeenCalledTimes(1)
      expect(props.getMemberCount.mock.calls[0]).toEqual([ 2 ])

    })
  })

})
