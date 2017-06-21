// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';
import { Icon } from './Icon';

import './MenuBar.css'

class MenuItem extends React.PureComponent {

  render() {
    const { name, event, onClick, selected } = this.props

    const doit = (e) => {
      let tag = e.target.getAttribute("name")
      onClick(tag)
    }

    const className = selected === event ? "MenuItem Focus" : "MenuItem"

    return (
      <div className={className} name={event} onClick={doit}>
        <div className="Icon" >
          <Icon type={event}/>
        </div>
        <div className="Name" name={event}>
          {name}
        </div>
      </div>
    )
  }
}

const menus = [
  {name: "Posts", event: "list-posts"},
  {name: "New post", event: "new-post"},
  {name: "Site", event: "edit-site"},
  {name: "Account", event: "edit-account"},
  {name: "Password", event: "change-password"},
]

class MenuBar extends React.PureComponent {

  render() {
    const { onClick, selected } = this.props

    return (
      <section className="MenuBar">
        {
          menus.map(m =>
            <MenuItem key={m.name} name={m.name} event={m.event}
                      onClick={onClick} selected={selected}/>)
        }
      </section>
    )
  }
}


export { MenuBar }
