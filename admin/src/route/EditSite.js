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

    this.state = props.site
    this.handleChange = this.handleChange.bind(this)
    this.saveSite = this.saveSite.bind(this)
  }

  componentWillReceiveProps(props) {
    const { title, description, baseUrl } = props.site
    let t = title ? title : ""
    let d = description ? description : ""
    let b = baseUrl ? baseUrl : ""

    this.setState({title: t, description: d, baseUrl: b})
  }

  handleChange(event) {
    const name = event.target.name
    const value = event.target.value
    this.setState({[name]: value})
  }

  saveSite() {
    console.log(this.state)
  }

  render() {

    const { history } = this.props
    const { title, baseUrl, description } = this.state

    const onCancel = () => history.push("/admin/home")

    return (
      <WorkArea>
        <Form>
          <FormTitle>
            Edit site
          </FormTitle>
          <FormWidgets>

            <FormWidget>
              <FormLabel>
                Title
              </FormLabel>
              <input name="title"
                     value={title}
                     autoFocus={true}
                     placeholder="Site title"
                     onChange={this.handleChange}/>
            </FormWidget>

            <FormWidget>
              <FormLabel>
                Description
              </FormLabel>
              <textarea name="description"
                        placeholder="Site description"
                        value={description}
                        onChange={this.handleChange}/>
            </FormWidget>

            <FormWidget>
              <FormLabel>
                Base URL
              </FormLabel>
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
