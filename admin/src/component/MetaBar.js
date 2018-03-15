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
import { Icon } from './Icon.js'
import './MetaBar.css'

class MetaBar extends React.PureComponent {
  render() {
    const { logout, visit } = this.props
    return (
      <section className="MetaBar">
        <div className="Option" onClick={visit} title="Visit site">
          <Icon type="visit" />
        </div>
        <div className="Option" onClick={logout} title="Sign out">
          <Icon type="signout"/>
        </div>
      </section>
    )
  }
}

export { MetaBar }
