// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';

import { Icon } from './Icon'

import './Action.css'

class Action extends React.PureComponent {

  render() {
    const { type, color, onClick } = this.props

    return (
      <span className="Action" onClick={onClick}>
        <Icon type={type} color={color}/>
      </span>
    )
  }
}

export { Action }
