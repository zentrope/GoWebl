// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react'
import { Map } from 'immutable'

import './Icon.css'

const icon = name => (color) => {
  return (
    <i className={"fa " + name + " Icon-" + color} aria-hidden="true"></i>
  )
}

const icons = Map({
  "new-post":   icon("fa-plus"),
  "edit-site":  icon("fa-cog"),
  "list-posts": icon("fa-list"),
  delete:       icon("fa-trash-o"),
  draft:        icon("fa-toggle-off"),
  edit:         icon("fa-pencil-square-o"),
  published:    icon("fa-toggle-on"),
  question:     icon("fa-question"),
  visit:        icon("fa-external-link"),
  signout:      icon("fa-sign-out")
})

class Icon extends React.PureComponent {

  render() {
    const { type, color } = this.props
    const icon = icons.get(type, icons.get("question"))(color)
    return ( icon )
  }

}

export { Icon }
