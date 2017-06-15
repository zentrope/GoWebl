// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
