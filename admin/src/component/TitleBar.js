// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';

import { Icon } from './Icon.js'

import './TitleBar.css'

class TitleBar extends React.PureComponent {

  render() {
    const { user, logout, title, visit, editSite, newPost } = this.props

    // TODO: Figure out how to hide new post and edit site when either
    //       one is current active. Or at least ignore clicks.

    return (
      <section className="TitleBar">
        <div className="Title">{title}</div>
        <div className="Name">{user}</div>
        <div className="Options">
          <button onClick={newPost}>
            <Icon type="new"/> &nbsp;Post
          </button>
          <button onClick={editSite}>
            <Icon type="settings"/> &nbsp;Site
          </button>
          <button onClick={visit}>
            <Icon type="visit" /> &nbsp;Site
          </button>
          <button onClick={logout}>
            <Icon type="signout"/>
      &nbsp;Sign out
          </button>
        </div>
      </section>
    )
  }
}

export { TitleBar }
