// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';
import { Form, FormControls, FormWidgets, FormWidget, FormLabel, FormTitle } from '../component/Form'
import { WorkArea } from '../component/WorkArea'

const isBlank = (s) => ((! s) || (s.trim().length === 0))

class ChangePassword extends React.PureComponent {

  constructor(props) {
    super(props)

    this.state = { password: "", confirm: "" }
    this.handleChange = this.handleChange.bind(this)
    this.update = this.update.bind(this)
    this.disabled = this.disabled.bind(this)
  }

  handleChange(event) {
    let { name, value } = event.target
    this.setState({[name]: value})
  }

  update() {
    const { client, onCancel } = this.props
    let { password } = this.state

    client.updateViewerPassword(password, (response) => {
      if (response.errors) {
        console.error(response.errors)
        return
      }
      onCancel()
    })
  }

  disabled() {
    let { password, confirm } = this.state
    return isBlank(password) ||
           isBlank(confirm) ||
           password.length < 8 ||
           confirm !== password
  }

  render() {
    const { onCancel } = this.props
    const { password, confirm } = this.state

    return (
      <WorkArea>
        <Form>
          <FormTitle>Change password</FormTitle>
          <FormWidgets>
            <FormWidget>
              <FormLabel>Password</FormLabel>
              <input value={password} autoFocus={true} name="password" type="password" onChange={this.handleChange}/>
            </FormWidget>
            <FormWidget>
              <FormLabel>Confirm</FormLabel>
              <input value={confirm} name="confirm" type="password" onChange={this.handleChange}/>
            </FormWidget>
          </FormWidgets>
          <FormControls>
            <button disabled={this.disabled()} onClick={this.update}>Set new password</button>
            <button onClick={onCancel}>Cancel</button>
          </FormControls>
        </Form>
      </WorkArea>
    )
  }
}

export { ChangePassword }
