// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';

import { Icon } from './Icon.js'

import './TitleBar.css'

class TitleBar extends React.PureComponent {

  render() {
    const { user, logout, title, visit } = this.props

    return (
      <section className="TitleBar">
        <div className="Meta">
          <div className="Title">{title}</div>
          <div className="Name">{user}</div>
        </div>
        <div className="Options">
          <button onClick={visit}>
            <Icon type="visit" /> Site
          </button>
          <button onClick={logout}>
            <Icon type="signout"/>&nbsp;Bye
          </button>
        </div>
      </section>
    )
  }
}

export { TitleBar }
