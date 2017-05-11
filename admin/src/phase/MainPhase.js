// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';
import { BrowserRouter as Router, Route, Redirect, Switch } from 'react-router-dom'

import { Icon } from '../component/Icon'

import './MarkdownEditor.css'
import './StatusBar.css'
import './Tabular.css'
import './TitleBar.css'
import './WorkArea.css'

const moment = require('moment')
const markdown = require('markdown-it')()
  .use(require('markdown-it-footnote'))

class MarkdownEditor extends React.PureComponent {

  constructor(props) {
    super(props)

    this.state = {slugline: "", text: "", showPreview : false}

    this.handleChange = this.handleChange.bind(this)
    this.togglePreview = this.togglePreview.bind(this)
    this.isSubmittable = this.isSubmittable.bind(this)
    this.savePost = this.savePost.bind(this)
  }

  handleChange(event) {
    const name = event.target.name
    const value = event.target.value
    this.setState({[name]: value})
  }

  togglePreview() {
    const show = ! this.state.showPreview
    this.setState({showPreview: show})
  }

  isSubmittable() {
    const t = this.state.text.trim()
    const s = this.state.slugline.trim()
    return (s.length >= 3) && (t.length >= 3)
  }

  savePost() {
    const { onSave } = this.props
    onSave(this.state.slugline, this.state.text)
  }

  render() {
    const { onCancel, onPublish } = this.props
    const { slugline, text, showPreview } = this.state

    const html = showPreview ? markdown.render(text) : ""

    const preview = showPreview ? (
      <div className="Html" dangerouslySetInnerHTML={{__html: html}}/>
    ) : (
      null
    )

    return (
      <section className="MarkdownEditor">
        <div className="Slugline">
          <input name="slugline"
                 autoFocus={true}
                 type="text"
                 value={slugline}
                 placeholder="Summary or slugline...."
                 onChange={this.handleChange}/>
        </div>
        <div className="Editor">
          <div className={"Text" + (showPreview ? "" : " Full")}>
            <textarea name="text"
                      autoFocus={false}
                      placeholder="Thoughts?"
                      value={text}
                      onChange={this.handleChange}/>
          </div>
          { preview }
        </div>
        <div className="Controls">
          <button disabled={!this.isSubmittable()}
                  onClick={this.savePost}>
            Save draft
          </button>
          <button  disabled={!this.isSubmittable()} onClick={onPublish}>
            Publish
          </button>
          <button onClick={this.togglePreview}>
            { showPreview ? "Hide Preview" : "Show Preview" }
          </button>
          <button onClick={onCancel}>
            Cancel
          </button>
        </div>
      </section>
    )
  }
}

class DateShow extends React.PureComponent {
  render () {
    const { date } = this.props
    const show = moment(date).format("D MMM YY - hh:mm A")

    return (
      <span className="DateShow">{ show }</span>
    )
  }
}

class TabularView extends React.PureComponent {

  render() {
    const { columns, render, data } = this.props

    const headers = (
      <thead>
        <tr>
          { columns.map(c => <th key={Math.random()}>{ c }</th>) }
        </tr>
      </thead>
    );

    let table = null
    if (data) {
      table = (
        <tbody>
          { data.map(d => render(d)) }
        </tbody>
      )
    }

    return (
      <div className="Tabular">
        <table>
          { headers }
          { table }
        </table>
      </div>
    )
  }
}

class Posts extends React.PureComponent {

  render() {
    const { posts } = this.props

    const renderPost = p => {
      return (
        <tr key={p.uuid}>
          <td width="1%">
            <Icon type={p.status} color={p.status}/>
          </td>
          <td width="40%"><a>{p.slugline}</a></td>
          <td width="29%"><DateShow date={p.dateCreated}/></td>
          <td width="29%"><DateShow date={p.dateUpdated}/></td>
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
      <TabularView columns={cols} data={posts} render={renderPost}/>
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
    const { viewer } = this.props
    const { showEditor } = this.state

    const editor = showEditor === true ? (
      <MarkdownEditor onPublish={this.publishPost}
                      onCancel={this.hideEditor}
                      onSave={this.savePost}/>
    ) : (
      <button onClick={this.showEditor}>New post</button>
    )

    return (
      <WorkArea>
        <h1>Webl Posts</h1>
        {editor}
        <Posts posts={viewer.posts}/>
      </WorkArea>
    )
  }
}

class StatusBar extends React.PureComponent {
  render() {
    return (
      <section className="StatusBar">
        <div className="Copyright">&copy; 2017, Keith Irwin</div>
      </section>
    )
  }
}

class TitleBar extends React.PureComponent {

  render() {
    const { viewer, logout } = this.props
    const { email, user } = viewer

    return (
      <section className="TitleBar">
        <div className="Title">Administration</div>
        <div className="Name">
          {user + ' <' + email + '>'}
        </div>
        <div className="Options">
          <button onClick={logout}>Sign out</button>
        </div>
      </section>
    )
  }
}

class WorkArea extends React.PureComponent {
  render() {
    return (
      <section className="WorkArea">
        { this.props.children }
      </section>
    )
  }
}

class MainPhase extends React.PureComponent {

  constructor(props) {
    super(props)
    this.state = {viewer: {}}

    this.dispatch = this.dispatch.bind(this)
  }

  componentWillMount() {
    const { client } = this.props
    client.viewerData(response => {
      const data = response.data.viewer
      this.setState({viewer: data ? data : {}})
    })
  }

  dispatch(event, data) {
    console.log("event>", event)
    switch (event) {
      case 'post/add':
        const viewer = this.state.viewer;
        viewer.posts.push(data)
        this.setState({viewer: viewer})
        this.forceUpdate() // yikes!
        break
      default:
        console.log("Unable to handle event.")
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
          <TitleBar viewer={viewer} logout={logout}/>
          <StatusBar/>
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
