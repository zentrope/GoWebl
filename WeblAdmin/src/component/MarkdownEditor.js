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

import { DateEditor } from './DateEditor'

import './MarkdownEditor.css'

const markdown = require('markdown-it')({html: true})
  .use(require('markdown-it-footnote'))

class MarkdownEditor extends React.PureComponent {

  constructor(props) {
    super(props)

    const { uuid, slugline, text, datePublished } = props
    let u = uuid ? uuid : ""
    let s = slugline ? slugline : ""
    let t = text ? text : ""
    let d = datePublished ? datePublished : new Date()

    this.state = {
      uuid: u,
      slugline: s,
      text: t,
      datePublished: d,
      showPreview : false,
      dirty: false
    }

    this.handleChange = this.handleChange.bind(this)
    this.togglePreview = this.togglePreview.bind(this)
    this.isSubmittable = this.isSubmittable.bind(this)
    this.savePost = this.savePost.bind(this)
    this.setDate = this.setDate.bind(this)
  }

  setDate(date) {
    this.setState({datePublished: date})
  }

  UNSAFE_componentWillReceiveProps(nextProps) {
    const { uuid, slugline, text, datePublished } = nextProps
    let u = uuid ? uuid : ""
    let s = slugline ? slugline : ""
    let t = text ? text : ""
    let d = datePublished ? datePublished : new Date()
    this.setState({
      uuid: u,
      slugline: s,
      text: t,
      datePublished: d
    })
  }

  handleChange(event) {
    let name = event.target.name
    let value = event.target.value
    this.setState({[name]: value, dirty: true})
  }

  togglePreview() {
    let show = ! this.state.showPreview
    this.setState({showPreview: show})
  }

  isSubmittable() {
    let t = this.state.text.trim()
    let s = this.state.slugline.trim()
    return (s.length >= 3) && (t.length >= 3)
  }

  savePost() {
    const { onSave } = this.props
    const { uuid, slugline, text, datePublished } = this.state
    if (uuid) {
      onSave(uuid, slugline, text, datePublished)
    } else {
      onSave(slugline, text, datePublished)
    }
    this.setState({dirty: false})
  }

  render() {
    const { onCancel } = this.props
    const { slugline, text, datePublished, showPreview, dirty } = this.state

    const html = showPreview ? markdown.render(text) : ""

    const preview = showPreview ? (
      <div className="Html" dangerouslySetInnerHTML={{__html: html}}/>
    ) : (
      null
    )

    let wordCount = (text === "") ? 0 : text.trim().split(/\s+/).length
    let status = dirty ? "UNSAVED" : "Saved"

    return (
      <section className="MarkdownEditor">

        <div className="TopBar">
          <div className="Slugline">
            <input name="slugline"
                   autoFocus={true}
                   type="text"
                   value={slugline}
                   placeholder="Summary or slugline...."
                   onChange={this.handleChange}/>
          </div>
          <div className="Metadata">
            <DateEditor time={datePublished}
                        template="YYYY-MM-DD hh:mm A"
                        onChange={this.setDate}/>
          </div>
        </div>


        <div className="Viewers">
          <div className="Editor">
            <textarea name="text"
                      autoFocus={false}
                      placeholder="Thoughts?"
                      value={text}
                      onChange={this.handleChange}/>
          </div>
          { preview }
        </div>

        <div className="Controls">
          <div className="Left">
            <button disabled={!this.isSubmittable()}
                    onClick={this.savePost}>
              Save
            </button>
            <button onClick={onCancel}>
              { dirty ? "Cancel" : "Done" }
            </button>
          </div>
          <div className="Badges">
            <div className="Badge">
              { wordCount } words
            </div>
            <div className={"Badge " + status}>
              { status }
            </div>
          </div>
          <div className="Right">
            <button onClick={this.togglePreview}>
              { showPreview ? "Hide Preview" : "Show Preview" }
            </button>
          </div>

        </div>
      </section>
    )
  }
}

export { MarkdownEditor }
