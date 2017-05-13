// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';
import { Router, Route, Redirect, Switch } from 'react-router-dom'

import { Map, fromJS } from 'immutable'

import { Action } from '../component/Action'
import { DateShow } from '../component/DateShow'
import { Icon } from '../component/Icon'
import { MarkdownEditor } from '../component/MarkdownEditor'
import { StatusBar } from '../component/StatusBar'
import { Tabular } from '../component/Tabular'
import { TitleBar } from '../component/TitleBar'
import { WorkArea } from '../component/WorkArea'

import createBrowserHistory from 'history/createBrowserHistory'

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

class MainPhase extends React.PureComponent {

  constructor(props) {
    super(props)
    this.state = {viewer: Map({})}
    this.history = createBrowserHistory()
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

      case 'post/save':
        client.savePost(data.slugline, data.text, "draft", (response) => {
          const newPost = response.data.createPost
          if (newPost) {
            const v = this.state.viewer.update("posts", ps => ps.push(fromJS(newPost)))
            this.setState({viewer: v})
            this.history.push('/admin/home')
            return
          }
          console.error(response.errors)
        })
        break

      case 'post/delete':
        client.deletePost(data.uuid, (response) => {
          const uuid = response.data.deletePost
          if (uuid) {
            const posts = this.state.viewer
                              .get("posts")
                              .filter(p => p.get("uuid") !== uuid)
            this.setState({viewer: this.state.viewer.set("posts", posts)})
            return
          }
          console.error(response.errors)
        })
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
      <Router history={this.history}>
        <section className="App">
          <TitleBar title="Webl Manager" user={viewer.get("email")} logout={logout}/>
          <StatusBar year="2017" copyright="Keith Irwin"/>
          <Switch>
            <PropRoute path="/admin/home" component={Home} viewer={viewer} client={client} dispatch={this.dispatch}/>
            <PropRoute path="/admin/post/new" component={NewPost} dispatch={this.dispatch}/>
            <Redirect to="/admin/home"/>
          </Switch>
        </section>
      </Router>
    )
  }
}

export { MainPhase }
