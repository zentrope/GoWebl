//
// Copyright (c) 2017-2018 Keith Irwin
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
  {name: "Password", event: "change-password"}
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
