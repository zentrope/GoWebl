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
