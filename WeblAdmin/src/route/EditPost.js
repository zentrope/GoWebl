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

import { MarkdownEditor } from '../component/MarkdownEditor'
import { WorkArea } from '../component/WorkArea'

const moment = require('moment')

class EditPost extends React.PureComponent {

  constructor(props) {
    super(props)
    this.mounted = false
    this.state = { uuid: "", slugline: "", datePublished: "", text: "" }
    this.update = this.update.bind(this)
    this.load = this.load.bind(this)
  }

  componentDidMount() {
    this.mounted = true
    this.load()
  }

  componentWillUnmount() {
    this.mounted = false
  }

  update(uuid, slugline, text, datePublished) {
    const { client } = this.props

    let pub = moment(datePublished).toISOString()

    client.updatePost(uuid, slugline, text, pub, (response) => {
      if (response.errors) {
        console.error("error", response.errors[0])
        return
      }
    })
  }

  load() {
    let { client, match } = this.props
    let id = match.params.id

    client.viewerData(response => {
      if (response.errors) {
        console.log(response.errors)
        return
      }
      let posts = response.data.viewer.posts
      let post = { uuid: "", slugline: "", datePublished: "", text: "" }
      for (let i = 0; i < posts.length; i++) {
        if (posts[i].uuid === id) {
          post = posts[i]
          break
        }
      }

      let {uuid, slugline, datePublished, text} = post

      if (this.mounted) {
        this.setState({uuid: uuid,
                       slugline: slugline,
                       datePublished: datePublished, text:text})
      }
    })
}

  render() {
    const { onCancel } = this.props
    const { uuid, slugline, text, datePublished } = this.state

    return (
      <WorkArea>
        <MarkdownEditor uuid={uuid}
                        slugline={slugline}
                        datePublished={datePublished}
                        text={text}
                        onCancel={onCancel}
                        onSave={this.update}/>
      </WorkArea>
    )
  }
}

export { EditPost }
