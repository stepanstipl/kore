import React from 'react'
import App from 'next/app'
import Head from 'next/head'
import Router from 'next/router'
import Link from 'next/link'
import axios from 'axios'
import Paragraph from 'antd/lib/typography/Paragraph'
import { Layout, Tag } from 'antd'
const { Header, Content, Footer } = Layout
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

// style imports
import 'antd/dist/antd.less'
import '../assets/styles.less'

import User from '../lib/components/layout/User'
import SiderMenu from '../lib/components/layout/SiderMenu'
import redirect from '../lib/utils/redirect'
import KoreApi from '../lib/kore-api'
import copy from '../lib/utils/object-copy'
import OrgService from '../server/services/org'
import userExpired from '../server/lib/user-expired'
import gtag from '../lib/utils/gtag'

Router.events.on('routeChangeComplete', url => {
  gtag.pageView(url)
})

class MyApp extends App {
  static async getUserSession(ctx) {
    const req = ctx && ctx.req
    const res = ctx && ctx.res
    const asPath = ctx && ctx.asPath
    if (req) {
      const session = req.session
      const user = session && session.passport && session.passport.user
      if (user) {
        session.requestedPath = asPath

        if (userExpired(session.passport.user)) {
          return res.redirect('/login/refresh')
        }

        const orgService = new OrgService(KoreApi)
        try {
          await orgService.refreshUser(user)
          return user
        } catch (err) {
          console.log('Failed to refresh user in _app.js', err)
          return false
        }
      }
      return false
    }
    try {
      const result = await axios.get(`${window.location.origin}/session/user?requestedPath=${asPath}`)
      return result.data
    } catch (err) {
      return false
    }
  }

  static async getInitialProps({ Component, ctx }) {
    let pageProps = ((Component.staticProps && typeof Component.staticProps === 'function') ? Component.staticProps(ctx) : Component.staticProps) || {}
    if (pageProps.unrestrictedPage) {
      return { pageProps }
    }
    const user = await MyApp.getUserSession(ctx)
    if (!user) {
      return redirect({
        res: ctx.res,
        path: `/login/refresh?requestedPath=${ctx.asPath}`,
        ensureRefreshFromServer: true
      })
    }
    const userTeams = (user.teams.userTeams || []).filter(t => !publicRuntimeConfig.ignoreTeams.includes(t.metadata.name))
    const otherTeams = (user.teams.otherTeams || []).filter(t => !publicRuntimeConfig.ignoreTeams.includes(t.metadata.name))
    if (pageProps.adminOnly && !user.isAdmin) {
      return redirect({
        res: ctx.res,
        router: Router,
        path: '/'
      })
    }
    if (Component.getInitialProps) {
      const initialProps = await Component.getInitialProps({ ...ctx, user })
      pageProps = { ...pageProps, ...initialProps }
    }
    return { pageProps, user, userTeams, otherTeams }
  }

  state = {
    userTeams: this.props.userTeams,
    version: null
  }

  setSessionTimeout() {
    clearInterval(this.interval)
    if (this.props.pageProps && !this.props.pageProps.unrestrictedPage) {
      // using session TTL + 5 seconds
      const intervalMs = (publicRuntimeConfig.sessionTtlInSeconds + 5) * 1000
      this.interval = setInterval(async () => {
        const user = await MyApp.getUserSession()
        if (!user) {
          redirect({
            path: '/login',
            ensureRefreshFromServer: true
          })
        }
      }, intervalMs)
    }
  }

  componentDidMount() {
    this.setSessionTimeout()
    if (publicRuntimeConfig.showPrototypes) {
      this.setState({ prototypePath: window.location.pathname.indexOf('/prototype') === 0 })
    }
    axios.get(`${window.location.origin}/version`).then((v) => {
      this.setState({ version: v.data.version })
    })
  }

  componentDidUpdate(prevProps, prevState) {
    this.setSessionTimeout()
    if (publicRuntimeConfig.showPrototypes) {
      const prototypePath = window.location.pathname.indexOf('/prototype') === 0
      if (prevState.prototypePath !== prototypePath) {
        this.setState({ prototypePath })
      }
    }
  }

  componentWillUnmount() {
    clearInterval(this.interval)
  }

  teamAdded = (team) => {
    const state = copy(this.state)
    state.userTeams.push(team)
    this.setState(state)
  }

  teamRemoved = (team) => {
    this.setState({
      userTeams: this.state.userTeams.filter(t => t.metadata.name !== team),
      otherTeams: (this.state.otherTeams || []).filter(t => t.metadata.name !== team)
    })
  }

  render() {
    const { Component } = this.props
    const props = { ...this.props, ...this.props.pageProps }
    const isAdmin = Boolean(props.user && props.user.isAdmin)
    const hideSider = Boolean(props.hideSider || props.unrestrictedPage)
    const hidePage = Boolean(!props.unrestrictedPage && !props.user)
    const { version, prototypePath } = this.state

    if (hidePage) {
      return null
    }

    return (
      <div>
        <Head>
          <script
            async
            src={`https://www.googletagmanager.com/gtm.js?id=${gtag.GTM_ID}`}
          />
          <script
            dangerouslySetInnerHTML={{
              __html: `
            window.dataLayer = window.dataLayer || [];
            function gtag(){dataLayer.push(arguments);}
            gtag('js', new Date());
            gtag('config', '${gtag.GTM_ID}');
          `,
            }}
          />
          <title>{props.title || 'Kore'}</title>
          <meta charSet="utf-8"/>
          <meta name="viewport" content="initial-scale=1.0, width=device-width"/>
          {!publicRuntimeConfig.disableAnimations ? null : (
            <style type="text/css" dangerouslySetInnerHTML={{
              __html: `
                *, *::after, *::before {
                  transition-delay: 0s !important;
                  transition-duration: 0s !important;
                  animation-delay: -0.0001s !important;
                  animation-duration: 0s !important;
                  animation-play-state: paused !important;
                  caret-color: transparent !important;
                  color-adjust: exact !important;
                }
              `
            }} />
          )}
        </Head>
        <Layout style={{ minHeight:'100vh' }}>
          <Header className='top-header'>
            <div style={{ color: '#FFFFFF', float: 'left', marginLeft: '-25px' }}>
              <div style={{ float: 'left' }}>
                <Link href="/">
                  <a style={{ color: '#FFFFFF' }}>
                    <img src="/static/images/appvia-white.svg" height="28px" />
                  </a>
                </Link>
              </div>
              <div style={{ float: 'left', paddingLeft: '15px', paddingTop: '1px', fontSize: '20px' }}>
                <Link href="/">
                  <a style={{ color: '#FFFFFF' }}>Kore</a>
                </Link>
              </div>
              {prototypePath ? <Link href="/prototype"><Tag style={{ marginLeft: '20px' }}>PROTOTYPE</Tag></Link> : null}
            </div>
            <User user={props.user}/>
          </Header>
          <Layout hasSider="true">
            <SiderMenu hide={hideSider} isAdmin={isAdmin} userTeams={this.state.userTeams} otherTeams={props.otherTeams}/>
            <Layout>
              <Content style={{ background: '#fff', padding: 24 }}>
                <Component {...this.props.pageProps} user={this.props.user} teamAdded={this.teamAdded} teamRemoved={this.teamRemoved} version={version} />
              </Content>
              <Footer className="footer">
                <Paragraph className="version">Kore {version}</Paragraph>
              </Footer>
            </Layout>
          </Layout>
        </Layout>
      </div>
    )
  }
}

export default MyApp
