// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';
import { BrowserRouter as Router, Route, Redirect, Switch } from 'react-router-dom'

import { Map, fromJS } from 'immutable'

import { Action } from '../component/Action'
import { DateShow } from '../component/DateShow'
import { Icon } from '../component/Icon'
import { MarkdownEditor } from '../component/MarkdownEditor'
import { StatusBar } from '../component/StatusBar'
import { Tabular } from '../component/Tabular'
import { TitleBar } from '../component/TitleBar'
import { WorkArea } from '../component/WorkArea'

class Posts extends React.PureComponent {

  render() {
    const { posts, dispatch } = this.props

    const toggleStatus = (post) => () => {
      const { status, uuid, slugline } = post.toJS()
      let t = status === "published" ? "draft" : "published"
      if (window.confirm("Set '" + slugline + "' status to '" + t + "'?")) {
        const msg = {uuid: uuid, isPublished: status === "draft"}
        dispatch("post/status", msg)
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
              <Icon type="delete" color="blue"/>
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

class Home extends React.PureComponent {

  constructor(props) {
    super(props)

    this.state = { showEditor: false }

    this.savePost = this.savePost.bind(this)
    this.publishPost = this.publishPost.bind(this)
    this.hideEditor = this.hideEditor.bind(this)
    this.showEditor = this.showEditor.bind(this)
  }

  savePost(slugline, text) {
    const { client, dispatch } = this.props
    client.savePost(slugline, text, "draft", (data) => {
      const newPost = data.data.createPost
      if (newPost) {
        dispatch('post/add', newPost)
      } else {
        console.error(data)
      }
    })
  }

  publishPost() {
    console.log("publish post -> not implemented")
  }

  showEditor() {
    this.setState({showEditor: true})
  }

  hideEditor() {
    this.setState({showEditor: false})
  }

  render() {
    const { viewer, dispatch } = this.props
    const { showEditor } = this.state

    const editor = showEditor === true ? (
      <MarkdownEditor onPublish={this.publishPost}
                      onCancel={this.hideEditor}
                      onSave={this.savePost}/>
    ) : (
      <button onClick={this.showEditor}>New post</button>
    )

    let posts
    if (! viewer.isEmpty()) {
      posts = viewer.get("posts").sortBy(p => p.get("dateCreated")).reverse()
    }

    return (
      <WorkArea>
        <h1>{viewer.get("user")}'s posts</h1>
        {editor}
        <Posts posts={posts} dispatch={dispatch}/>
      </WorkArea>
    )
  }
}

class MainPhase extends React.PureComponent {

  constructor(props) {
    super(props)
    this.state = {viewer: Map({})}

    this.dispatch = this.dispatch.bind(this)
    this.refresh = this.refresh.bind(this)
  }

  componentWillMount() {
    this.refresh()
  }

  refresh() {
    const { client } = this.props
    client.viewerData(response => {
      const data = fromJS(response.data.viewer)
      this.setState({viewer: data})
    })
  }

  dispatch(event, data) {
    const { client } = this.props
    console.log("event>", event)
    switch (event) {
      case 'post/add':
        console.log("post/add is broken at the moment")
        /* const viewer = this.state.viewer;
         * viewer.posts.push(data)
         * this.setState({viewer: viewer})
         * this.forceUpdate() // yikes!*/
        break
      case 'post/status':
        client.setPostStatus(data.uuid, data.isPublished, (response) => {
          const updated = response.data.setPostStatus
          const errors = response.errors
          if (updated) {
            this.refresh() // should fold into existing data
          } else {
            console.log(errors)
          }
        })
        break
      default:
        console.log("Unable to handle event.", data)
    }
  }

  render() {
    const { logout, client } = this.props
    const { viewer } = this.state

    const PropRoute = ({component: Component, path: Path, ...rest}) => (
      <Route exact path={Path} render={(props) => (<Component {...rest} {...props}/> )}/>
    )

    return (
      <Router>
        <section className="App">
          <TitleBar title="Webl Manager" user={viewer.get("email")} logout={logout}/>
          <StatusBar year="2017" copyright="Keith Irwin"/>
          <Switch>
            <PropRoute path="/admin/home" component={Home} viewer={viewer} client={client} dispatch={this.dispatch}/>
            <Redirect to="/admin/home"/>
          </Switch>
        </section>
      </Router>
    )
  }
}

export { MainPhase }
