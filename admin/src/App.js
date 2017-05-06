// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';
import { BrowserRouter as Router, Route, Link, Redirect } from 'react-router-dom'

import { Client } from './Client'

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

class Authors extends React.PureComponent {
  render() {
    return (
      <section>
        <h1>Authors</h1>
        <Link to="/admin/home">Home</Link>
      </section>
    )
  }
}

class Home extends React.PureComponent {

  constructor(props) {
    super(props)

    this.state = {loggedIn : 0}
  }

  render() {
    const { client } = this.props
    const { loggedIn } = this.state

    const logout = () => {
      this.setState({loggedIn : -1})
      client.invalidateAuthToken()
      localStorage.setItem("auth-token", "no-auth")

    }

    const login = () => {
      client.login("keith", "test1234", (result) => {
        console.log("login=", result)

        if (!result.data.authenticate) {
          console.log("Bad login.")
          this.setState({loggedIn : -1})
        } else {
          client.setAuthToken(result.data.authenticate)
          localStorage.setItem("auth-token", result.data.authenticate)
          this.setState({loggedIn : 1})
        }
      })
    }

    const check = () => {
      var token = "" + localStorage.getItem("auth-token")

      client.validate(token, (result) => {
        if (! result.data.validate) {
          console.log("User isn't logged in any more.")
          this.setState({loggedIn : -1})
        } else {
          console.log("User is logged in.")
          client.setAuthToken(token)
          localStorage.setItem("auth-token", token)
          this.setState({loggedIn : 1})
        }
      })
    }

    const viewer = () => {
      client.viewerData((result) => {
        console.log(result.data.viewer)
      })
    }

    const status = (loggedIn === 1) ? (
      "Logged In" ) : (
        (loggedIn === 0) ? (
          "Maybe"
        ) : (
          "Logged out"
        )
      )

    return (
      <div>
        <h1>Placeholder ({status})</h1>
        <ul>
          <li><Link to="/admin/post">Posts</Link></li>
          <li><Link to="/admin/authors">Authors</Link></li>
        </ul>
        <button onClick={check}>Check Creds</button>
        <br/>
        <button onClick={login}>Test Login</button>
        <br/>
        <button onClick={viewer}>Get Viewer Data</button>
        <br/>
        <button onClick={logout}>Test Log Out</button>
      </div>
    )
  }
}

const PropsRoute = ({component: Component, path: Path, ...rest}) => (
  <Route path={Path} exact render={(props) => (<Component {...rest} {...props}/> )}/>
)

class App extends React.PureComponent {
  constructor(props) {
    super(props)

    const { endpoint } = this.props

    this.client = new Client(endpoint, this.onError)
  }

  componentDidMount() {
  }

  render() {
    const { endpoint } = this.props
    return (
      <Router>
        <section>
          <PropsRoute path="/" component={Redirect} to="/admin/home" endpoint={endpoint}/>
          <PropsRoute path="/admin" component={Redirect} to="/admin/home" endpoint={endpoint}/>
          <PropsRoute path="/admin/home" component={Home} client={this.client}/>
          <PropsRoute path="/admin/post" component={Posts}/>
          <PropsRoute path="/admin/authors" component={Authors}/>
        </section>
      </Router>
    );
  }
}

export default App;
