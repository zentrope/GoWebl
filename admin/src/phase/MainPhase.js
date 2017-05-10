// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';
import { BrowserRouter as Router, Route, Redirect, Switch } from 'react-router-dom'

import './StatusBar.css'
import './TitleBar.css'
import './WorkArea.css'

const moment = require('moment')

class DateShow extends React.PureComponent {
  render () {
    const { date } = this.props
    const show = moment(date).format("D MMM YY - hh:mm A")

    return (
      <span className="DateShow">{ show }</span>
    )
  }
}

class Posts extends React.PureComponent {

  render() {
    const { posts } = this.props

    const renderPost = p =>
      <tr key={p.uuid}>
        <td className={p.status}>{p.status}</td>
        <td><a>{p.slugline}</a></td>
        <td><DateShow date={p.dateCreated}/></td>
        <td><DateShow date={p.dateUpdated}/></td>
      </tr>

    return (
      <table>
        <thead>
          <tr>
            <th>status</th>
            <th>slugline</th>
            <th>created</th>
            <th>updated</th>
          </tr>
        </thead>
        <tbody>
          { posts ? posts.map(p => renderPost(p)) : null }
        </tbody>
      </table>
    )
  }
}

class Home extends React.PureComponent {


  render() {
    const { viewer } = this.props

    return (
      <WorkArea>
        <h1>Webl Posts</h1>
        <Posts posts={viewer.posts}/>
      </WorkArea>
    )
  }
}

class StatusBar extends React.PureComponent {
  render() {
    return (
      <section className="StatusBar">
        <div className="Copyright">&copy; 2017, Keith Irwin</div>
      </section>
    )
  }
}

class TitleBar extends React.PureComponent {

  render() {
    const { viewer, logout } = this.props
    const { email, user } = viewer

    return (
      <section className="TitleBar">
        <div className="Title">Administration</div>
        <div className="Name">
          {user + ' <' + email + '>'}
        </div>
        <div className="Options">
          <button onClick={logout}>Sign out</button>
        </div>
      </section>
    )
  }
}

class WorkArea extends React.PureComponent {
  render() {
    return (
      <section className="WorkArea">
        { this.props.children }
      </section>
    )
  }
}

class MainPhase extends React.PureComponent {

  constructor(props) {
    super(props)
    this.state = {viewer: {}}
  }

  componentWillMount() {
    const { client } = this.props
    client.viewerData(response => {
      const data = response.data.viewer
      this.setState({viewer: data ? data : {}})
    })
  }

  render() {
    const { logout } = this.props
    const { viewer } = this.state

    const PropRoute = ({component: Component, path: Path, ...rest}) => (
      <Route exact path={Path} render={(props) => (<Component {...rest} {...props}/> )}/>
    )

    return (
      <Router>
        <section>
          <TitleBar viewer={viewer} logout={logout}/>
          <StatusBar/>
          <Switch>
            <PropRoute path="/admin/home" component={Home} viewer={viewer}/>
            <Redirect to="/admin/home"/>
          </Switch>
        </section>
      </Router>
    )
  }
}

export { MainPhase }
