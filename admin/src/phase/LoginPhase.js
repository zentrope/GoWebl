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

import "./LoginForm.css"

class LoginForm extends React.PureComponent {

  constructor(props) {
    super(props)

    this.state = {user : "", pass: "", error: ""}

    this.handleChange = this.handleChange.bind(this)
    this.handleSubmit = this.handleSubmit.bind(this)
    this.handleKeyDown = this.handleKeyDown.bind(this)
  }

  handleSubmit() {
    let { user, pass } = this.state
    user = user.trim()

    const { client, login } = this.props

    client.login(user, pass, (result) => {
      let okay = result.data.authenticate !== null
      if (! okay) {
        this.setState({error: "Unable to sign in."})
        document.getElementById("user").focus()
      } else {
        let token = result.data.authenticate.token
        login(token)
      }
    })
  }

  handleChange(e) {
    const name = e.target.name
    const value = e.target.value
    this.setState({[name]: value, error: ""})
  }

  handleKeyDown(e) {
    switch (e.keyCode) {
      case 13:
        if (this.isSubmittable()) {
          this.handleSubmit()
        }
        break;
      case 27:
        this.setState({user: "", pass: ""})
        document.getElementById("user").focus()
        break;
      default:
        ;
    }
  }

  isSubmittable() {
    let { user, pass, error } = this.state
    user = user.trim()
    pass = pass.trim()
    if (error.length > 0) {
      return false
    }
    return (user.length > 0) && (pass.length > 0)
  }

  render() {
    const { site } = this.props
    var { user, pass, error } = this.state

    const submit = this.isSubmittable() ? (
      <button onClick={this.handleSubmit}>Sign in</button>
    ) : (
      null
    )

    return (
      <section className="LoginForm">

        <section className="LoginPanel">
          <h1>Sign in to { site.title }</h1>

          <div className="Error">
            { error }
          </div>

          <div className="Control">
            { submit }
          </div>

          <div className="Widgets">
            <div className="Widget">
              <input id="user"
                     type="text"
                     name="user"
                     value={user}
                     autoComplete="off"
                     autoFocus={true}
                     placeholder="Webl ID"
                     onKeyDown={this.handleKeyDown}
                     onChange={this.handleChange}/>
            </div>
            <div className="Widget Pass">
              <input type="password"
                     name="pass"
                     value={pass}
                     autoComplete="off"
                     autoFocus={false}
                     placeholder="Password"
                     onKeyDown={this.handleKeyDown}
                     onChange={this.handleChange}/>
            </div>
          </div>

          <div className="VisitLink">
            <p><a href={ site.baseUrl }>Visit { site.title }</a></p>
          </div>

        </section>
      </section>
    )
  }
}

class LoginPhase extends React.PureComponent {

  render() {
    const { client, login, site } = this.props

    return (
      <LoginForm login={login} client={client} site={site}/>
    )
  }
}

export { LoginPhase }
