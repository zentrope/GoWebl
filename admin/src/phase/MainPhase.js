// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';
import { Router, Route, Redirect, Switch } from 'react-router-dom'

import { Map, fromJS } from 'immutable'

import { MenuBar } from '../component/MenuBar'
import { StatusBar } from '../component/StatusBar'
import { TitleBar } from '../component/TitleBar'

// Routes
import { EditAccount } from '../route/EditAccount'
import { ChangePassword } from '../route/ChangePassword'
import { EditPost } from '../route/EditPost'
import { EditSite } from '../route/EditSite'
import { Home } from '../route/Home'
import { NewPost } from '../route/NewPost'

import createBrowserHistory from 'history/createBrowserHistory'

class MainPhase extends React.PureComponent {

  constructor(props) {
    super(props)

    this.state = {
      viewer: Map({}),
      site: props.site,
      menu: "list-posts"
    }
    this.history = createBrowserHistory()
    this.dispatch = this.dispatch.bind(this)
    this.refresh = this.refresh.bind(this)

    this.savePost = this.savePost.bind(this)
    this.updatePost = this.updatePost.bind(this)
    this.updateAccount = this.updateAccount.bind(this)
    this.setPassword = this.setPassword.bind(this)
  }

  componentDidMount() {
    this.refresh()
  }

  setPassword(password) {
    const { client } = this.props
    client.updateViewerPassword(password, (response) => {
      if (response.errors) {
        console.error(response.errors)
        return
      }
      this.setState({menu: "list-posts"})
      this.history.push("/admin/home")
    })
  }

  updateAccount(name, email) {
    const { client } = this.props

    client.updateViewer(name, email, (response) => {
      if (response.errors) {
        console.error(response.errors)
        return
      }
      let newName = response.data.updateViewer.name
      let newEmail = response.data.updateViewer.email

      let v2 = this.state.viewer.set("name", newName).set("email", newEmail)
      this.setState({viewer: v2, menu: "list-posts"})
      this.history.push("/admin/home")
    })
  }

  updatePost(uuid, slugline, text) {
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
      console.log("refresh", response)
      // TODO: What happens if there are errors here?
      const data = fromJS(response.data.viewer)
      this.setState({viewer: data, site: response.data.viewer.site})
      this.forceUpdate()
    })
  }

  // This stuff is a mess.
  // TODO: Refactor state management
  dispatch(event, data, callback) {
    const { client } = this.props
    console.log("event>", event)

    switch (event) {

      case 'site/update':
        client.updateSite(data.title, data.description, data.baseUrl, (response) => {
          if (response.errors) {
            console.error(response.errors)
          } else {
            this.setState({ site : response.data.updateSite })
          }
          this.setState({menu: "list-posts"})
          this.history.push("/admin/home")
        })
        break

      case 'post/save':
        this.savePost(data)
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
    const { viewer, site, menu } = this.state

    const PropRoute = ({component: Component, path: Path, ...rest}) => (
      <Route exact path={Path} render={(props) => (<Component {...rest} {...props}/> )}/>
    )

    const visit = () => {
      window.location.href = site.baseUrl
    }

    const saveSite = (siteData) => {
      this.dispatch('site/update', siteData)
    }

    const onCancel = () => {
      this.setState({menu: "list-posts"})
      this.history.push("/admin/home")
    }

    const onMenuClick = (event) => {
      switch (event) {
        case "list-posts":
          this.setState({menu: event})
          this.history.push("/admin/home")
          break;
        case "new-post":
          this.setState({menu: event})
          this.history.push("/admin/post/new")
          break;
        case "edit-site":
          this.setState({menu: event})
          this.history.push("/admin/site/edit")
          break;
        case "edit-account":
          this.setState({menu: event})
          this.history.push("/admin/account/edit")
          break;
        case "change-password":
          this.setState({menu: event})
          this.history.push("/admin/account/password/edit")
          break;
        default:
          console.log("Unknown menu event:", event);
      }
    }

    const userName = viewer.get("name") + " <" + viewer.get("email") + ">"

    return (
      <Router history={this.history}>
        <section className="App">
          <TitleBar title={ site.title } user={userName} visit={visit} logout={logout}/>
          <MenuBar onClick={onMenuClick} selected={menu}/>
          <StatusBar year="2017" copyright={ site.title }/>
          <Switch>
            <PropRoute path="/admin/home" component={Home} viewer={viewer} client={client} onGotoPost={this.gotoPost} dispatch={this.dispatch}/>
            <PropRoute path="/admin/post/new" component={NewPost} dispatch={this.dispatch} onCancel={onCancel}/>
            <PropRoute path="/admin/post/:id" component={EditPost} posts={viewer.get("posts")} onSave={this.updatePost} onCancel={onCancel}/>
            <PropRoute path="/admin/site/edit" component={EditSite} site={site} onSave={saveSite} onCancel={onCancel}/>
            <PropRoute path="/admin/account/edit" component={EditAccount} onCancel={onCancel} email={viewer.get("email")} name={viewer.get("name")} onSave={this.updateAccount}/>
            <PropRoute path="/admin/account/password/edit" component={ChangePassword} onCancel={onCancel} onSave={this.setPassword}/>
            <Redirect to="/admin/home"/>
          </Switch>
        </section>
      </Router>
    )
  }
}

export { MainPhase }
