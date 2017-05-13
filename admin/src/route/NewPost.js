// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';

import { MarkdownEditor } from '../component/MarkdownEditor'
import { WorkArea } from '../component/WorkArea'

// /admin/post/new
class NewPost extends React.PureComponent {

  render() {
    const { dispatch, history } = this.props

    const onCancel = () => {
      history.push("/admin/home")
    }

    const onSave = (slugline, text) => {
      dispatch('post/save', {slugline: slugline, text: text})
    }

    return (
      <WorkArea>
        <h1>New post</h1>
        <MarkdownEditor onCancel={onCancel} onSave={onSave}/>
      </WorkArea>
    )
  }
}

export { NewPost }
