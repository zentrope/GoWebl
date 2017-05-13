// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';

import { Action } from '../component/Action'
import { DateShow } from '../component/DateShow'
import { Icon } from '../component/Icon'
import { Tabular } from '../component/Tabular'
import { WorkArea } from '../component/WorkArea'

class Posts extends React.PureComponent {

  render() {
    const { posts, dispatch } = this.props

    const toggleStatus = (post) => () => {
      const { status, uuid, slugline } = post.toJS()
      let t = status === "published" ? "draft" : "published"
      if (window.confirm("Set '" + slugline + "' status to '" + t + "'?")) {
        const msg = {uuid: uuid, isPublished: status === "draft"}
        dispatch('post/status', msg)
      }
    }

    const deletePost = (post) => () => {
      const { uuid, slugline } = post.toJS()
      if (window.confirm("Delete '" + slugline + "' for all time?")) {
        dispatch('post/delete', {uuid: uuid})
      }
    }

    const renderPost = p => {
      const { status, uuid, slugline, dateCreated, dateUpdated } = p.toJS()
      return (
        <tr key={uuid}>
          <td width="1%">
            <Action type={status} color={status} onClick={toggleStatus(p)}/>
          </td>
          <td width="40%"><a>{slugline}</a></td>
          <td width="29%"><DateShow date={dateCreated}/></td>
          <td width="29%"><DateShow date={dateUpdated}/></td>
          <td width="1%">
            <center>
              <Icon type="edit" color="blue"/>
              <span> </span>
              <Action type="delete" color="blue" onClick={deletePost(p)}/>
            </center>
          </td>
        </tr>
      )
    }

    const cols = [null, "slugline", "created", "updated", null]

    return (
      <Tabular columns={cols} data={posts} render={renderPost}/>
    )
  }
}

// /admin/home
class Home extends React.PureComponent {

  render() {
    const { viewer, dispatch } = this.props

    let posts
    if (! viewer.isEmpty()) {
      posts = viewer.get("posts").sortBy(p => p.get("dateCreated")).reverse()
    }

    const newPost = () => {
      const { history } = this.props
      history.push("/admin/post/new")
    }

    return (
      <WorkArea>
        <h1>{viewer.get("user")}'s posts</h1>
        <button onClick={newPost}>New post</button>
        <Posts posts={posts} dispatch={dispatch}/>
      </WorkArea>
    )
  }
}

export { Home }
