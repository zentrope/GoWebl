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
