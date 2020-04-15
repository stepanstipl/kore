import React from 'react'
import App from 'next/app'
import Head from 'next/head'
import Router from 'next/router'
import Link from 'next/link'
import axios from 'axios'
import { Layout } from 'antd'
const { Header, Content } = Layout

import User from '../lib/components/User'
import SiderMenu from '../lib/components/SiderMenu'
import redirect from '../lib/utils/redirect'
import KoreApi from '../lib/kore-api'
import copy from '../lib/utils/object-copy'
import { kore, server } from '../config'
import OrgService from '../server/services/org'
import userExpired from '../server/lib/user-expired'
import gtag from '../lib/utils/gtag'
import '../assets/styles.less'
import Paragraph from 'antd/lib/typography/Paragraph'

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
    const userTeams = (user.teams.userTeams || []).filter(t => !kore.ignoreTeams.includes(t.metadata.name))
    const otherTeams = (user.teams.otherTeams || []).filter(t => !kore.ignoreTeams.includes(t.metadata.name))
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
      const intervalMs = (server.session.ttlInSeconds + 5) * 1000
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
    axios.get(`${window.location.origin}/version`).then((v) => {
      this.setState({ version: v.data.version })
    })
  }

  componentDidUpdate() {
    this.setSessionTimeout()
  }

  componentWillUnmount() {
    clearInterval(this.interval)
  }

  teamAdded = (team) => {
    const state = copy(this.state)
    state.userTeams.push(team)
    this.setState(state)
  }

  render() {
    const { Component } = this.props
    const props = { ...this.props, ...this.props.pageProps }
    const isAdmin = Boolean(props.user && props.user.isAdmin)
    const hideSider = Boolean(props.hideSider || props.unrestrictedPage)
    const hidePage = Boolean(!props.unrestrictedPage && !props.user)
    const { version } = this.state

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
          <title>{props.title || 'Appvia Kore'}</title>
          <meta charSet="utf-8"/>
          <meta name="viewport" content="initial-scale=1.0, width=device-width"/>
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
                  <a style={{ color: '#FFFFFF' }}>Appvia Kore</a>
                </Link>
              </div>
            </div>
            <User user={props.user}/>
          </Header>
          <Layout hasSider="true">
            <SiderMenu hide={hideSider} isAdmin={isAdmin} userTeams={this.state.userTeams} otherTeams={props.otherTeams}/>
            <Content style={{ background: '#fff', padding: 24, minHeight: 280 }}>
              <Component {...this.props.pageProps} user={this.props.user} teamAdded={this.teamAdded} version={version} />
              <Paragraph style={{ textAlign: 'right', fontSize: '0.8em', padding: 0, margin: 0 }}>Appvia Kore {version}</Paragraph>
            </Content>
          </Layout>
        </Layout>
      </div>
    )
  }
}

export default MyApp
