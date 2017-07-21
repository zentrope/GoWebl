// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';
import { Router, Route, Redirect, Switch } from 'react-router-dom'

import { MenuBar } from '../component/MenuBar'
import { StatusBar } from '../component/StatusBar'
import { TitleBar } from '../component/TitleBar'

// Routes
import { Activity } from '../route/Activity'
import { ChangePassword } from '../route/ChangePassword'
import { EditAccount } from '../route/EditAccount'
import { EditPost } from '../route/EditPost'
import { EditSite } from '../route/EditSite'
import { Home } from '../route/Home'
import { Metrics } from '../route/Metrics'
import { NewPost } from '../route/NewPost'

import createBrowserHistory from 'history/createBrowserHistory'

class MainPhase extends React.PureComponent {

  constructor(props) {
    super(props)

    this.state = {
      user: "",
      title: "",
      baseUrl: "",
      menu: "list-posts"
    }
    this.history = createBrowserHistory()
    this.refresh = this.refresh.bind(this)
  }

  componentDidMount() {
    this.refresh()
  }

  refresh() {
    const { client } = this.props
    client.viewerData(response => {
      if (response.errors) {
        console.error(response.errors)
        return
      }
      let { name, email } = response.data.viewer
      let site = response.data.viewer.site
      this.setState({
        title: site.title,
        user: name + " <" + email + ">",
        baseUrl: site.baseUrl
      })
    })
  }

  render() {
    const { logout, client } = this.props
    const { title, user, baseUrl, menu } = this.state

    const PropRoute = ({component: Component, path: Path, ...rest}) => (
      <Route exact path={Path} render={(props) => (<Component {...rest} {...props}/> )}/>
    )

    const visit = () => {
      window.location.href = baseUrl
    }

    const onCancel = () => {
      this.setState({menu: "list-posts"})
      this.refresh()
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
        case "list-activity":
          this.setState({menu: event})
          this.history.push("/admin/activity")
          break;
        case "metrics":
          this.setState({menu: event})
          this.history.push("/admin/metrics")
          break;
        default:
          console.log("Unknown menu event:", event);
      }
    }

    return (
      <Router history={this.history}>
        <section className="App">
          <TitleBar title={title} user={user} visit={visit} logout={logout}/>
          <MenuBar onClick={onMenuClick} selected={menu}/>
          <StatusBar year="2017" copyright={ title }/>
          <Switch>
            <PropRoute path="/admin/home" component={Home} client={client}/>
            <PropRoute path="/admin/activity" component={Activity}  client={client}/>
            <PropRoute path="/admin/metrics" component={Metrics} client={client}/>
            <PropRoute path="/admin/post/new" component={NewPost} client={client} onCancel={onCancel}/>
            <PropRoute path="/admin/post/:id" component={EditPost} client={client} onCancel={onCancel}/>
            <PropRoute path="/admin/site/edit" component={EditSite} client={client} onCancel={onCancel}/>
            <PropRoute path="/admin/account/edit" component={EditAccount} client={client} onCancel={onCancel}/>
            <PropRoute path="/admin/account/password/edit" component={ChangePassword} client={client} onCancel={onCancel}/>
            <Redirect to="/admin/home"/>
          </Switch>
        </section>
      </Router>
    )
  }
}

export { MainPhase }
