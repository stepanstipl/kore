import { mount } from 'enzyme'
import NewTeamForm from '../../../../../lib/components/forms/NewTeamForm'
import copy from '../../../../../lib/utils/object-copy'
import ApiTestHelpers from '../../../../api-test-helpers'

describe('NewTeamForm', () => {
  let props
  let newTeamForm
  let wrapper
  let apiScope
  
  beforeEach(() => {
    // In case any tests leak to the API, mock out the API at this stage:
    apiScope = ApiTestHelpers.getScope()
    props = {
      form: {
        isFieldTouched: () => {},
        getFieldDecorator: jest.fn(() => c => c),
        getFieldsError: () => () => {},
        getFieldError: () => {},
        getFieldValue: () => {},
        validateFields: jest.fn()
      },
      handleTeamCreated: jest.fn(),
      user: { id: 'jbloggs' },
      team: false
    }
    wrapper = mount(<NewTeamForm wrappedComponentRef={component => newTeamForm = component} {...props} />)
  })

  afterEach(() => {
    // This will check that no calls were made against the API, unless the test registered them:
    apiScope.done()
  })

  describe('instance methods', () => {

    describe('#constructor', () => {
      it('sets initial state', () => {
        expect(newTeamForm.state).toEqual({
          submitting: false,
          formErrorMessage: false
        })
      })
    })

    describe('#componentDidMount', () => {
      it('validates fields', () => {
        expect(props.form.validateFields).toHaveBeenCalledTimes(1)
      })
    })

    describe('#disableButton', () => {
      it('true when fields have errors', () => {
        const fieldsError = { teamName: 'some error' }
        const disabled = newTeamForm.disableButton(fieldsError)
        expect(disabled).toBe(true)
      })

      it('true when form is submitting', () => {
        const state = copy(newTeamForm.state)
        state.submitting = true
        newTeamForm.setState(state)
        const disabled = newTeamForm.disableButton({})
        expect(disabled).toBe(true)
      })

      it('false when fields valid and not submitting', () => {
        const disabled = newTeamForm.disableButton({})
        expect(disabled).toBe(false)
      })
    })

  })

  describe('#handleSubmit', () => {
    let event

    beforeEach(() => {
      props.form.validateFields.mockClear()
      event = {
        preventDefault: jest.fn()
      }
    })

    it('prevents default action, sets saving state and validates the fields', () => {
      newTeamForm.handleSubmit(event)
      expect(event.preventDefault).toHaveBeenCalledTimes(1)
      expect(newTeamForm.state).toEqual({
        submitting: true,
        formErrorMessage: false
      })
      expect(props.form.validateFields).toHaveBeenCalledTimes(1)
    })

    it('does not submit when the form it is not valid', () => {
      props.form.validateFields = jest.fn(cb => cb(true, {}))
      newTeamForm.handleSubmit(event)
      // have the API scope validate no calls were made:
      apiScope.done()
    })

    it('submits when the form is valid', async () => {
      const team = {
        metadata: {
          name: 'abc'
        },
        spec: {
          summary: 'ABC',
          description: 'ABC team'
        }
      }

      apiScope = (ApiTestHelpers.getScope())
        .get(`${ApiTestHelpers.basePath}/teams/abc`).reply(404)
        .put(`${ApiTestHelpers.basePath}/teams/abc`, team).reply(200, team)

      props.form.validateFields = jest.fn(cb => cb(null, { teamName: 'ABC', teamDescription: 'ABC team' }))
      await newTeamForm.handleSubmit(event)

      // Check expected API calls made:
      apiScope.done()

      expect(newTeamForm.state).toEqual({
        submitting: false,
        formErrorMessage: false
      })
      expect(props.handleTeamCreated).toHaveBeenCalledTimes(1)
      expect(props.handleTeamCreated.mock.calls[0]).toEqual([ team ])
    })

    it('sets error message if api requests fail', async () => {
      apiScope = (ApiTestHelpers.getScope())
        .get(`${ApiTestHelpers.basePath}/teams/abc`).reply(500)

      props.form.validateFields = jest.fn(cb => cb(null, { teamName: 'ABC', teamDescription: 'ABC team' }))
      await newTeamForm.handleSubmit(event)

      // Check expected API calls made:
      apiScope.done()

      expect(props.handleTeamCreated).toHaveBeenCalledTimes(0)
      expect(newTeamForm.state).toEqual({
        submitting: false,
        formErrorMessage: 'An error occurred creating the team, please try again'
      })
    })

  })

  describe('render', () => {
    it('shows form items name, description and save button', () => {
      const formItems = wrapper.find('FormItem')
      expect(formItems).toHaveLength(3)
      expect(formItems.at(0).text()).toEqual('Team name')
      expect(formItems.at(1).text()).toEqual('Team description')
      expect(formItems.at(2).text()).toEqual('Save')
    })
  })

})
