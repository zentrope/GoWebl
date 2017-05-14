// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';
import { Map } from 'immutable'
import { MarkdownEditor } from '../component/MarkdownEditor'
import { WorkArea } from '../component/WorkArea'

// /admin/post/<uuid>
class EditPost extends React.PureComponent {

  constructor(props) {
    super(props)
    this.state = {post: Map()}
  }

  componentDidMount() {
    const { dispatch, match, history } = this.props
    dispatch('post/get', {uuid: match.params.id}, (data) => {
      if (! data) {
        history.push('/admin/home')
        return
      }
      this.setState({post: data})
    })
  }

  render() {
    const { dispatch, history } = this.props
    const { post } = this.state
    const { uuid, slugline, text } = post.toJS()

    const onCancel = () => {
      history.push('/admin/home')
    }

    const onSave = (uuid, slugline, text) => {
      dispatch('post/update', {uuid: uuid, slugline: slugline, text: text},
      )
    }

    return (
      <WorkArea>
        <h1>Edit post</h1>
        <MarkdownEditor uuid={uuid} slugline={slugline} text={text} onCancel={onCancel} onSave={onSave}/>
      </WorkArea>
    )
  }
}

export { EditPost }
