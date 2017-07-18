// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';

import { PageSwitcher } from '../component/PageSwitcher'
import { Tabular } from '../component/Tabular'
import { WordSelector } from '../component/WordSelector'
import { WorkArea } from '../component/WorkArea'

class Metrics extends React.PureComponent {

  componentDidMount() {
    // Handles the case where we arrive here via a direct refresh
    // on the route.
    let { metrics, onChange, report} = this.props
    if (onChange && metrics.length < 1) {
      onChange(report)
    }
  }

  renderRow(row) {
    const { key, value } = row
    return (
      <tr key={key}>
        <td>{ key }</td>
        <td>{ value }</td>
      </tr>
    )
  }

  render() {
    let { report, reports, metrics, onChange } = this.props

    return (
      <WorkArea>
        <PageSwitcher title={report}>
          <WordSelector words={reports} selected={report} onChange={onChange}/>
        </PageSwitcher>
        <Tabular columns={["key", "value"]} data={metrics} render={this.renderRow}/>
      </WorkArea>
    )
  }
}

export { Metrics }
