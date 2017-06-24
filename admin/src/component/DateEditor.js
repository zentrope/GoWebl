// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
