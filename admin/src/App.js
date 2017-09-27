//
// Copyright (c) 2017 Keith Irwin
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published
// by the Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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

    this.state = { loggedIn: 0, site: {baseURL: "/", title: "Webl", description: "Webl"}}

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

    this.client.siteData(result => {
      if (! result.data.site) {
        console.log(result.errors)
        return
      }
      this.setState({site: result.data.site})
      this.forceUpdate()
    })

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
    const { loggedIn, site } = this.state

    switch (loggedIn) {
      case 0:
        return (<LoadingPhase/>)

      case -1:
        return (<LoginPhase site={site} client={this.client} login={this.onLogin}/>)

      default:
        return (<MainPhase site={site} client={this.client} logout={this.onLogout}/>)
    }
  }
}

export default App
