// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';

import { WorkArea } from '../component/WorkArea'
import { DateShow } from '../component/DateShow'
import { Tabular } from '../component/Tabular'

class Activity extends React.PureComponent {

  componentDidMount() {
    const { onRefresh, activity } = this.props
    if (onRefresh && activity.length < 1) {
      onRefresh()
    }
  }

  renderRow(request) {
    const { dateRecorded, path, host, method, userAgent } = request
    return (
      <tr key={Math.random()} title={ userAgent }>
        <td width="20%"><DateShow date={dateRecorded}/></td>
        <td width="40%">{ host }</td>
        <td width="40%">{ method + " " + path }</td>
      </tr>
    )
  }

  render() {
    const { activity } = this.props
    return (
      <WorkArea>
        <Tabular columns={["date", "host", "path"]}
                 data={activity}
                 render={this.renderRow}/>
      </WorkArea>
    )
  }
}

export { Activity }
