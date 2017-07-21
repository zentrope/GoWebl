// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';

import { PageSwitcher } from '../component/PageSwitcher'
import { Tabular } from '../component/Tabular'
import { WordSelector } from '../component/WordSelector'
import { WorkArea } from '../component/WorkArea'

class Metrics extends React.PureComponent {

  constructor(props) {
    super(props)

    this.mounted = false

    this.state = {
      report: "Top Routes",
      metrics: []
    }

    this.handleChange = this.handleChange.bind(this)
  }

  componentDidMount() {
    this.mounted = true
    let { report } = this.state
    this.handleChange(report)
  }

  componentWillUnmount() {
    this.mounted = false
  }

  handleChange(reportName) {
    let { client } = this.props
    let key = reportName.toLowerCase()
    let resolver = {
      "top hits": "topHits",
      "top routes": "topRoutes",
      "top refers": "topRefers",
      "hits per day": "hitsPerDay"
    }

    let report = resolver[key] ? resolver[key] : "topRoutes";

    client.metricsReport(report, (response) => {
      if (response.errors) {
        console.log(response.errors)
        return
      }

      let values = response.data.viewer.metrics[report]
      if (this.mounted) {
        this.setState({metrics: values, report: reportName})
      }
    })
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
    let { report, metrics } = this.state
    const reports = ["Top Hits", "Top Routes", "Top Refers", "Hits per Day"]

    return (
      <WorkArea>
        <PageSwitcher title={report}>
          <WordSelector words={reports} selected={report} onChange={this.handleChange}/>
        </PageSwitcher>
        <Tabular columns={["key", "value"]} data={metrics} render={this.renderRow}/>
      </WorkArea>
    )
  }
}

export { Metrics }
