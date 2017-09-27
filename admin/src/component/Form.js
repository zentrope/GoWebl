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
import './Form.css'

class Form extends React.PureComponent {
  render() {
    return (
      <section className="Form">
        { this.props.children }
      </section>
    )
  }
}

class FormControls extends React.PureComponent {
  render() {
    return (
      <div className="FormControls">
        { this.props.children }
      </div>
    )
  }
}

class FormWidgets extends React.PureComponent {
  render() {
    return (
      <div className="FormWidgets">
        { this.props.children }
      </div>
    )
  }
}

class FormWidget extends React.PureComponent {
  render() {
    return (
      <div className="FormWidget">
        { this.props.children }
      </div>
    )
  }
}

class FormLabel extends React.PureComponent {
  render() {
    return (
      <div className="FormLabel">
        { this.props.children }
      </div>
    )
  }
}

class FormTitle extends React.PureComponent {
  render() {
    return (
      <div className="FormTitle">
        { this.props.children }
      </div>
    )
  }
}


export { Form, FormControls, FormWidgets, FormWidget, FormLabel, FormTitle }
