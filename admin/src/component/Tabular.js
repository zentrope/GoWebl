// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';

import './Tabular.css'

class Tabular extends React.PureComponent {

  render() {
    const { columns, render, data } = this.props

    const headers = (
      <thead>
        <tr>
          { columns.map(c => <th key={Math.random()}>{ c }</th>) }
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
