// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import './index.css';

const isDevMode = window.location.port === "3000"
const endpoint = isDevMode ? (
  window.location.origin.replace("3000", "8080")
) : (
  window.location.origin
)

const docRoot =
  document.getElementById('root')

const render = () =>
  ReactDOM.render(<App endpoint={endpoint}/>, docRoot)

const main = () => {
  console.log("Welcome to Webl Admin App")
  console.log("Using '" + endpoint + "'.")
  render()
}

window.onload = main
