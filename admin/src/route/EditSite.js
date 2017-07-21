// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';
import { Form, FormControls, FormWidgets, FormWidget,
         FormLabel, FormTitle } from '../component/Form'
import { WorkArea } from '../component/WorkArea'

// /admin/site/edit
class EditSite extends React.PureComponent {

  constructor(props) {
    super(props)

    this.mounted = false
    this.state = {title: "", description: "", baseUrl: ""}
    this.handleChange = this.handleChange.bind(this)
    this.saveSite = this.saveSite.bind(this)
  }

  componentDidMount() {
    this.mounted = true
    let { client } = this.props
    client.viewerData(response => {
      let { title, description, baseURL } = response.data.viewer.site
      if (this.mounted) {
        this.setState({title: title, description: description, baseURL: baseURL})
      }
    })
  }

  componentWillUnmount() {
    this.mounted = false
  }

  handleChange(event) {
    const name = event.target.name
    const value = event.target.value
    this.setState({[name]: value})
  }

  saveSite() {
    let { title, description, baseUrl } = this.state
    let { client, onCancel } = this.props

    client.updateSite(title, description, baseUrl, (response) => {
      if (response.errors) {
        console.error(response.errors)
      }

      onCancel()
    })
  }

  render() {

    const { onCancel } = this.props
    const { title, baseUrl, description } = this.state

    return (
      <WorkArea>
        <Form>
          <FormTitle>Edit site</FormTitle>
          <FormWidgets>

            <FormWidget>
              <FormLabel>Title</FormLabel>
              <input name="title"
                     value={title}
                     autoFocus={true}
                     placeholder="Site title"
                     onChange={this.handleChange}/>
            </FormWidget>

            <FormWidget>
              <FormLabel>Description</FormLabel>
              <textarea name="description"
                        placeholder="Site description"
                        value={description}
                        onChange={this.handleChange}/>
            </FormWidget>

            <FormWidget>
              <FormLabel>Base URL</FormLabel>
              <input name="baseUrl"
                     placeholder="Base URL"
                     value={baseUrl}
                     onChange={this.handleChange}/>
            </FormWidget>

          </FormWidgets>

          <FormControls>
            <button onClick={this.saveSite}>Save</button>
            <button onClick={onCancel}>Done</button>
          </FormControls>

        </Form>
      </WorkArea>
    )
  }
}

export { EditSite }
