// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';

import { MarkdownEditor } from '../component/MarkdownEditor'
import { WorkArea } from '../component/WorkArea'

const moment = require('moment')

class EditPost extends React.PureComponent {

  constructor(props) {
    super(props)
    this.mounted = false
    this.state = { uuid: "", slugline: "", datePublished: "", text: "" }
    this.update = this.update.bind(this)
    this.load = this.load.bind(this)
  }

  componentDidMount() {
    this.mounted = true
    this.load()
  }

  componentWillUnmount() {
    this.mounted = false
  }

  update(uuid, slugline, text, datePublished) {
    const { client } = this.props

    let pub = moment(datePublished).toISOString()

    client.updatePost(uuid, slugline, text, pub, (response) => {
      if (response.errors) {
        console.error("error", response.errors[0])
        return
      }
    })
  }

  load() {
    let { client, match } = this.props
    let id = match.params.id

    client.viewerData(response => {
      if (response.errors) {
        console.log(response.errors)
        return
      }
      let posts = response.data.viewer.posts
      let post = { uuid: "", slugline: "", datePublished: "", text: "" }
      for (let i = 0; i < posts.length; i++) {
        if (posts[i].uuid === id) {
          post = posts[i]
          break
        }
      }

      let {uuid, slugline, datePublished, text} = post

      if (this.mounted) {
        this.setState({uuid: uuid,
                       slugline: slugline,
                       datePublished: datePublished, text:text})
      }
    })
}

  render() {
    const { onCancel } = this.props
    const { uuid, slugline, text, datePublished } = this.state

    return (
      <WorkArea>
        <MarkdownEditor uuid={uuid}
                        slugline={slugline}
                        datePublished={datePublished}
                        text={text}
                        onCancel={onCancel}
                        onSave={this.update}/>
      </WorkArea>
    )
  }
}

export { EditPost }
