// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';

import './MarkdownEditor.css'

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

export { MarkdownEditor }
