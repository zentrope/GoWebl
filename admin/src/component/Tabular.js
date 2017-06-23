// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';

import './Tabular.css'

class Tabular extends React.PureComponent {

  render() {
    const { columns, render, data } = this.props


    let headFn = (name) => {
      let className = ""
      if (name && name.endsWith("+")) {
        className = "Right"
      }

      let colName = name ? name.replace("+", "") : ""

      return (
        <th key={Math.random()} className={className}>
          {colName}
        </th>
      )
    }

    const headers = (
      <thead>
        <tr>
          { columns.map(c => headFn(c)) }
        </tr>
      </thead>
    );

    let table = null
    if (data) {
      table = (
        <tbody>
          { data.map(d => render(d)) }
        </tbody>
      )
    }

    return (
      <div className="Tabular">
        <table>
          { headers }
          { table }
        </table>
      </div>
    )
  }
}

export { Tabular }
