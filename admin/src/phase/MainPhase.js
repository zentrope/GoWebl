// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';
import { BrowserRouter as Router, Route, Link, Redirect, Switch } from 'react-router-dom'

class Posts extends React.PureComponent {

  render() {
    return (
      <section>
        <h1>Posts</h1>
        <Link to="/admin/home">Home</Link>
      </section>
    )
  }
}

class Home extends React.PureComponent {

  render() {
    const { viewer } = this.props

    return (
      <div>
        <h1>[{viewer.email}]</h1>
        <ul>
          <li><Link to="/admin/post">Posts</Link></li>
        </ul>
      </div>
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
          <p><button onClick={logout}>log out</button> -- {viewer.id}</p>
          <Switch>
            <PropRoute path="/admin/home" component={Home} viewer={viewer}/>
            <PropRoute path="/admin/post" component={Posts}/>
            <Redirect to="/admin/home"/>
          </Switch>
        </section>
      </Router>
    )
  }
}

export { MainPhase }
