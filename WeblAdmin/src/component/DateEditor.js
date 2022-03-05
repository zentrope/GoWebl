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

import './DateEditor.css'

const moment = require('moment')

const formatDate = (time, template) => {
  return (time ? moment(time) : moment()).format(template)
}

class DateEditor extends React.PureComponent {

  constructor(props) {
    super(props)
    this.state = {
      date: formatDate(props.time, props.template),
      isValid: true
    }

    this.update = this.update.bind(this)
  }

  sampleDate() {
    let f = this.props.template
    return moment().format(f)
  }

  update(e) {
    let template = this.props.template
    let n = e.target.name
    let v = e.target.value

    let m = moment(v, template, true)
    this.setState({[n] : v, isValid: m.isValid()})

    let callback = this.props.onChange
    if (m.isValid() && callback) {
      callback(m.toDate())
    }
  }

  render () {
    const { date, isValid } = this.state

    let validClass = "TSEditor " + (isValid ? "Valid" : "Invalid")

    let suggest = (! isValid) ? (
      <span className="Suggest">
        &laquo; { this.sampleDate() } &raquo;
      </span>
    ) : ""

    return (
      <div className={validClass}>
        <input type="text"
               name="date"
               placeholder={this.sampleDate()}
               value={date}
               onChange={this.update}/>
        { suggest }
      </div>
    )
  }
}


export { DateEditor }
