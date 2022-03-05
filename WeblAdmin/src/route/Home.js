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

import { Link } from 'react-router-dom'

import { Action } from '../component/Action'
import { DateShow } from '../component/DateShow'
import { PageTitle } from '../component/PageTitle'
import { Tabular } from '../component/Tabular'
import { WorkArea } from '../component/WorkArea'

/* /admin/home */

class Home extends React.PureComponent {

  constructor(props) {
    super(props)
    this.mounted = false
    this.state = { posts: [] }

    this.load = this.load.bind(this)
    this.setPostStatus = this.setPostStatus.bind(this)
    this.deletePost = this.deletePost.bind(this)
    this.renderPost = this.renderPost.bind(this)
  }

  componentDidMount() {
    this.mounted = true
    this.load()
  }

  componentWillUnmount() {
    this.mounted = false
  }

  load() {
    let { client } = this.props
    client.viewerData(response => {
      if (response.errors) {
        console.log(response.errors)
        return
      }
      if (this.mounted) {
        let posts = response.data.viewer.posts
        posts.sort(
          (a, b) =>
            a.datePublished === b.datePublished ? (
              0
            ) : (
                a.datePublished < b.datePublished ? 1 : -1
            )
        )
        this.setState({posts: posts})
      }
    })
  }

  setPostStatus(post) {
    let { client } = this.props
    let { status, uuid, slugline } = post
    let t = status === "published" ? "draft" : "published"
    let msg = "Set '" + slugline + "' status to '" + t + "'?"
    if (window.confirm(msg)) {
      client.setPostStatus(uuid, status === "draft", response => {
        if (response.errors) {
          console.log(response.errors)
          return
        }
        this.load()
      })
    }
  }

  deletePost(post) {
    let { client } = this.props
    let { uuid, slugline } = post
    let msg = "Delete '" + slugline + "' for all time?"
    if (window.confirm(msg)) {
      client.deletePost(uuid, response => {
        if (response.errors) {
          console.error(response.errors)
          return
        }
        this.load()
      })
    }
  }

  renderPost(p) {
    const toggle = (p) => () => { this.setPostStatus(p) }
    const remove = (p) => () => { this.deletePost(p) }

    const { status, uuid, slugline, datePublished, wordCount } = p

    return (
      <tr key={uuid} className={status}>
        <td width="1%">
          <Action type={status} color={status} onClick={toggle(p)}/>
        </td>
        <td width="10%" className="Right">{wordCount}</td>
        <td width="64%"><Link to={"/admin/post/" + uuid}>{slugline}</Link></td>
        <td width="24%" className="Date"><DateShow date={datePublished}/></td>
        <td width="1%">
          <Action type="delete" color="blue" onClick={remove(p)}/>
        </td>
      </tr>
    )
  }

  render() {
    let { posts } = this.state
    const { title, user } = this.props

    const cols = [null, "words+", "post", "published", null]
    return (
      <WorkArea>
        <PageTitle title={title} subtitle={user}/>
        <Tabular columns={cols} data={posts} render={this.renderPost}/>
      </WorkArea>
    )
  }
}

export { Home }
