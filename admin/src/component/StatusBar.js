// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';
import './StatusBar.css'

class StatusBar extends React.PureComponent {
  render() {
    const { year, copyright } = this.props

    return (
      <section className="StatusBar">
        <div className="Copyright">&copy; {year}, {copyright}</div>
      </section>
    )
  }
}

export { StatusBar }
