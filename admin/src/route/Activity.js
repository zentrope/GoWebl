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
