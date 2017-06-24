// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';

import { Link } from 'react-router-dom'

import { Action } from '../component/Action'
import { DateShow } from '../component/DateShow'
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
      const { status, uuid, slugline, datePublished, wordCount } = p.toJS()
      return (
        <tr key={uuid} className={status}>
          <td width="1%">
            <Action type={status} color={status} onClick={toggleStatus(p)}/>
          </td>
          <td width="10%" className="Right">{wordCount}</td>
          <td width="44%"><Link to={"/admin/post/" + uuid}>{slugline}</Link></td>
          <td width="44%" className="Date"><DateShow date={datePublished}/></td>
          <td width="1%">
            <Action type="delete" color="blue" onClick={deletePost(p)}/>
          </td>
        </tr>
      )
    }

    const cols = [null, "words+", "post", "published", null]

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
      posts = viewer.get("posts").sortBy(p => p.get("datePublished")).reverse()
    }


    return (
      <WorkArea>
        <Posts posts={posts} dispatch={dispatch}/>
      </WorkArea>
    )
  }
}

export { Home }
