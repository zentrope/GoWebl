// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';

import './PageSwitcher.css'

class PageSwitcher extends React.PureComponent {

  render() {
    let { title } = this.props

    return (
      <div className="PageSwitcher">
        <div className="Title">
          { title }
        </div>
        <div className="Selector">
          {this.props.children}
        </div>
      </div>
    )
  }
}

export { PageSwitcher }
