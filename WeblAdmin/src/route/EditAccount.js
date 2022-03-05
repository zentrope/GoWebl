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
import { Form, FormControls, FormWidgets, FormWidget,
         FormLabel, FormTitle } from '../component/Form'
import { WorkArea } from '../component/WorkArea'

/* /admin/account/edit */

const isBlank = (s) => ((! s) || (s.trim().length === 0))

class EditAccount extends React.PureComponent {

  constructor(props) {
    super(props)

    this.state = {name: "", email: ""}
    this.handleChange = this.handleChange.bind(this)
    this.save = this.save.bind(this)
    this.load = this.load.bind(this)
    this.disabled = this.disabled.bind(this)
  }

  componentDidMount() {
    this.mounted = true
    this.load()
  }

  componentWillUnmount() {
    this.mounted = false
  }

  handleChange(event) {
    let name = event.target.name
    let value = event.target.value
    this.setState({[name]: value})
  }

  load() {
    let { client } = this.props
    client.viewerData(response => {
      let { email, name } = response.data.viewer
      this.oldEmail = email
      this.oldName = name
      if (this.mounted) {
        this.setState({email: email, name: name})
      }
    })
  }

  save() {
    let { name, email } = this.state
    let { client, onCancel } = this.props

    client.updateViewer(name, email, (response) => {
      if (response.errors) {
        console.error(response.errors)
        return
      }
      onCancel()
    })
  }

  disabled() {
    return (
      isBlank(this.state.name) ||
      isBlank(this.state.email) ||
      ((this.state.name === this.oldName) && (this.state.email === this.oldEmail))
    )
  }

  render() {
    const { onCancel } = this.props
    const { name, email } = this.state
    return (
      <WorkArea>
        <Form>
          <FormTitle>Edit account</FormTitle>
          <FormWidgets>
            <FormWidget>
              <FormLabel>Name</FormLabel>
              <input autoFocus={true} name="name" value={name} onChange={this.handleChange}/>
            </FormWidget>
            <FormWidget>
              <FormLabel>Email</FormLabel>
              <input autoFocus={false} name="email" value={email} onChange={this.handleChange}/>
            </FormWidget>
          </FormWidgets>
          <FormControls>
            <button disabled={this.disabled()} onClick={this.save}>Save changes</button>
            <button onClick={onCancel}>Cancel</button>
          </FormControls>
        </Form>
      </WorkArea>
    )
  }
}

export { EditAccount }
