// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';
import { Map } from 'immutable'
import { MarkdownEditor } from '../component/MarkdownEditor'
import { WorkArea } from '../component/WorkArea'

// /admin/post/<uuid>
class EditPost extends React.PureComponent {

  render() {
    const { posts, onCancel, onSave, match } = this.props

    let postUuid = match.params.id
    let post = (posts) ? (posts.filter(p => p.get("uuid") === postUuid).first()) : Map()

    const { uuid, slugline, text, datePublished } = post.toJS()

    return (
      <WorkArea>
        <MarkdownEditor uuid={uuid} slugline={slugline}
                        datePublished={datePublished}
                        text={text} onCancel={onCancel} onSave={onSave}/>
      </WorkArea>
    )
  }
}

export { EditPost }
