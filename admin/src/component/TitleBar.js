// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';

import './TitleBar.css'

class TitleBar extends React.PureComponent {

  render() {
    const { user, logout, title } = this.props

    return (
      <section className="TitleBar">
        <div className="Title">{title}</div>
        <div className="Name">{user}</div>
        <div className="Options">
          <button onClick={logout}>Sign out</button>
        </div>
      </section>
    )
  }
}

export { TitleBar }
