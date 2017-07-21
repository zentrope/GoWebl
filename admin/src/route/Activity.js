// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';

import { WorkArea } from '../component/WorkArea'
import { DateShow } from '../component/DateShow'
import { Tabular } from '../component/Tabular'

class Activity extends React.PureComponent {

  constructor(props) {
    super(props)
    this.mounted = false
    this.state = { activity: [] }
  }

  componentDidMount() {
    this.mounted = true
    this.refreshActivity()
  }

  componentWillUnmount() {
    this.mounted = false
  }

  refreshActivity() {
    const { client } = this.props
    client.requestData(100, (response) => {
      if (response.errors) {
        console.log(response.errors[0])
        return
      }

      if (this.mounted) {
        this.setState({activity: response.data.viewer.requests})
      }
    })
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
    const { activity } = this.state
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
