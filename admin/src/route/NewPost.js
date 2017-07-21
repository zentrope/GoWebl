// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
