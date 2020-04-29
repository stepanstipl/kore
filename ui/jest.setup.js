import { configure } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import { setConfig } from 'next/config'

configure({ adapter: new Adapter() })
setConfig({ publicRuntimeConfig: require('./config.public') })
