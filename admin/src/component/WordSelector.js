// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';

class WordSelector extends React.PureComponent {
  // Move this to components

  constructor(props) {
    super(props)
    this.handleChange = this.handleChange.bind(this)
  }

  handleChange(event) {
    let delegate = this.props.onChange
    let value = event.target.value
    if (delegate) {
      delegate(value)
    }
  }

  render() {
    const { words, selected } = this.props

    return (
      <select value={selected} onChange={this.handleChange}>
        { words.map(w => <option key={w} value={w}>{w}</option>) }
      </select>
    )
  }
}

export { WordSelector }
