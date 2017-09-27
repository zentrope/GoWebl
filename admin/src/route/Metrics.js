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
