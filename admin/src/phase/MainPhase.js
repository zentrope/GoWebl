// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';
import { Router, Route, Redirect, Switch } from 'react-router-dom'

import { Map, fromJS } from 'immutable'

import { StatusBar } from '../component/StatusBar'
import { TitleBar } from '../component/TitleBar'

// Routes
import { EditPost } from '../route/EditPost'
import { EditSite } from '../route/EditSite'
import { Home } from '../route/Home'
import { NewPost } from '../route/NewPost'

import createBrowserHistory from 'history/createBrowserHistory'

class MainPhase extends React.PureComponent {

  constructor(props) {
    super(props)

    this.state = {viewer: Map({}), site: props.site}
    this.history = createBrowserHistory()
    this.dispatch = this.dispatch.bind(this)
    this.refresh = this.refresh.bind(this)

    this.savePost = this.savePost.bind(this)
    this.updatePost = this.updatePost.bind(this)
  }

  componentDidMount() {
    this.refresh()
  }

  updatePost(data) {
    const { uuid, slugline, text } = data
    const { client } = this.props

    client.updatePost(uuid, slugline, text, (response) => {
      const post = response.data.updatePost
      if (post) {
        const posts = this.state.viewer
                          .get("posts")
                          .filter(p => p.get("uuid") !== uuid)
                          .push(fromJS(post))
        this.setState({viewer: this.state.viewer.set("posts", posts)})
        return
      }
      console.error(response.errors)
    })
  }

  savePost(data) {
    const { client } = this.props
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
  }

  refresh() {
    const { client } = this.props
    client.viewerData(response => {
      const data = fromJS(response.data.viewer)
      this.setState({viewer: data, site: response.data.viewer.site})
    })
  }

  dispatch(event, data, callback) {
    const { client } = this.props
    console.log("event>", event)

    switch (event) {

      case 'site/edit':
        this.history.push("/admin/site/edit")
        break

      case 'site/update':
        client.updateSite(data.title, data.description, data.baseUrl, (response) => {
          if (response.errors) {
            console.error(response.errors)
          } else {
            this.setState({ site : response.data.updateSite })
          }
          this.history.push("/admin/home")
        })
        break

      case 'post/get':
        if (callback) {
          callback(this.state.viewer.get("posts")
                       .filter(p => p.get("uuid") === data.uuid)
                       .first())
        }
        break

      case 'post/save':
        this.savePost(data)
        break

      case 'post/update':
        this.updatePost(data)
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
    const { viewer, site } = this.state

    const PropRoute = ({component: Component, path: Path, ...rest}) => (
      <Route exact path={Path} render={(props) => (<Component {...rest} {...props}/> )}/>
    )

    const visit = () => {
      window.location.href = site.baseUrl
    }

    const editSite = () => {
      this.dispatch('site/edit', {})
    }

    const saveSite = (siteData) => {
      this.dispatch('site/update', siteData)
    }

    const newPost = () => {
      this.history.push("/admin/post/new")
    }

    return (
      <Router history={this.history}>
        <section className="App">
          <TitleBar title={ site.title } user={viewer.get("email")}
                    editSite={editSite} visit={visit} logout={logout}
                    newPost={newPost}/>
          <StatusBar year="2017" copyright={ site.title }/>
          <Switch>
            <PropRoute path="/admin/home" component={Home} viewer={viewer} client={client} dispatch={this.dispatch}/>
            <PropRoute path="/admin/post/new" component={NewPost} dispatch={this.dispatch}/>
            <PropRoute path="/admin/post/:id" component={EditPost} dispatch={this.dispatch}/>
            <PropRoute path="/admin/site/edit" component={EditSite} site={site} onSave={saveSite}/>
            <Redirect to="/admin/home"/>
          </Switch>
        </section>
      </Router>
    )
  }
}

export { MainPhase }
