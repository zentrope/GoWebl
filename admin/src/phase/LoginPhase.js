// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';

class LoginPhase extends React.PureComponent {

  render() {
    const { client, login } = this.props

    const doLogin = () => {
      client.login("keith", "test1234", (result) => {
        const token = result.data.authenticate
        if (!token) {
          console.log("Bad login.")
        } else {
          login(token)
        }
      })
    }

    return (
      <div>
        <h1>Log in</h1>
        <button onClick={doLogin}>Simulate Login</button>
      </div>
    )
  }
}

export { LoginPhase }
