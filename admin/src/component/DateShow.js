// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';

const moment = require('moment')

class DateShow extends React.PureComponent {
  render () {
    const { date } = this.props
    const show = moment(date).format("DD MMM YY - hh:mm A")

    return (
      <span className="DateShow">{ show }</span>
    )
  }
}

export { DateShow }
