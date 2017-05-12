// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';

import { LoadingPhase } from './phase/LoadingPhase'
import { LoginPhase } from './phase/LoginPhase'
import { MainPhase } from './phase/MainPhase'

import { Client } from './Client'

class App extends React.PureComponent {

  constructor(props) {
    super(props)

    const { endpoint } = this.props

    this.client = new Client(endpoint, this.onError)

    this.state = { loggedIn : 0}

    this.onLogout = this.onLogout.bind(this)
    this.onLogin = this.onLogin.bind(this)
  }

  onError(error) {
    console.log(error)
  }

  onLogout() {
    this.setState({loggedIn : -1})
    this.client.invalidateAuthToken()
    localStorage.setItem("auth-token", "no-auth")
  }

  onLogin(token) {
    this.client.setAuthToken(token)
    localStorage.setItem("auth-token", token)
    this.setState({loggedIn : 1})
  }

  componentDidMount() {
    var token = "" + localStorage.getItem("auth-token")

    this.client.validate(token, result => {
      if (! result.data.validate) {
        console.log(result.errors)
        this.onLogout()

      } else {
        this.onLogin(token)
      }
    })
  }

  render() {
    const { loggedIn } = this.state

    switch (loggedIn) {
      case 0:
        return (<LoadingPhase/>)

      case -1:
        return (<LoginPhase client={this.client} login={this.onLogin}/>)

      default:
        return (<MainPhase client={this.client} logout={this.onLogout}/>)
    }
  }
}

export default App
