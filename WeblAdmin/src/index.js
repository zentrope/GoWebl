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
