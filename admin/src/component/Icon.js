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

import React from 'react'
import { Map } from 'immutable'

import './Icon.css'

const icon = name => (color) => {
  return (
    <i className={"fa " + name + " Icon-" + color + " fa-fw"} aria-hidden="true"></i>
  )
}

const icons = Map({
  "new-post": icon("fa-plus"),
  "edit-site": icon("fa-cog"),
  "list-posts": icon("fa-file-text-o"),
  "edit-account": icon("fa-user-circle"),
  "change-password": icon("fa-lock"),
  "list-activity": icon("fa-area-chart"),
  "metrics": icon("fa-line-chart"),
  delete: icon("fa-trash-o"),
  draft: icon("fa-toggle-off"),
  edit: icon("fa-pencil-square-o"),
  published: icon("fa-toggle-on"),
  question: icon("fa-question"),
  visit: icon("fa-external-link"),
  signout: icon("fa-sign-out")
})

class Icon extends React.PureComponent {

  render() {
    const { type, color } = this.props
    const icon = icons.get(type, icons.get("question"))(color)
    return ( icon )
  }

}

export { Icon }
