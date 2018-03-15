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
import { Router, Route, Redirect, Switch } from 'react-router-dom'

import { MenuBar } from '../component/MenuBar'
import { MetaBar } from '../component/MetaBar'

// Routes
import { ChangePassword } from '../route/ChangePassword'
import { EditAccount } from '../route/EditAccount'
import { EditPost } from '../route/EditPost'
import { EditSite } from '../route/EditSite'
import { Home } from '../route/Home'
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
      let { name } = response.data.viewer
      let site = response.data.viewer.site
      this.setState({
        title: site.title,
        user: name, // + " <" + email + ">",
        baseUrl: site.baseUrl
      })
    })
  }

  render() {
    const { logout, client } = this.props
    const { baseUrl, menu } = this.state

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
        default:
          console.log("Unknown menu event:", event);
      }
    }

    return (
      <Router history={this.history}>
        <section className="App">
          {
              // <TitleBar title={title} user={user} />
          }
          <MenuBar onClick={onMenuClick} selected={menu}/>
          <MetaBar visit={visit} logout={logout}/>
          <Switch>
            <PropRoute path="/admin/home" component={Home} client={client}/>
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
