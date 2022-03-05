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

import { MarkdownEditor } from '../component/MarkdownEditor'
import { WorkArea } from '../component/WorkArea'

const moment = require('moment')

class NewPost extends React.PureComponent {

  constructor(props) {
    super(props)

    this.save = this.save.bind(this)
  }

  save(slugline, text, datePublished) {
    const { client, onCancel } = this.props
    let date = moment(datePublished).toISOString()
    console.log("save post with date:", date)
    client.savePost(slugline, text, date, "draft", (response) => {
      if (response.errors) {
        response.errors.map(e => console.log("err:", e))
        return
      }
      onCancel()
    })
  }

  render() {
    const { onCancel } = this.props

    return (
      <WorkArea>
        <MarkdownEditor onCancel={onCancel} onSave={this.save}/>
      </WorkArea>
    )
  }
}

export { NewPost }
