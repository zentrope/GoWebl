import React from 'react';
import { BrowserRouter as Router, Route, Link, Redirect } from 'react-router-dom'

import { Client } from './Client'

class Posts extends React.PureComponent {

  render() {
    return (
      <section>
        <h1>Posts</h1>
        <Link to="/admin/home">Home</Link>
      </section>
    )
  }
}

class Authors extends React.PureComponent {
  render() {
    return (
      <section>
        <h1>Authors</h1>
        <Link to="/admin/home">Home</Link>
      </section>
    )
  }
}

class Home extends React.PureComponent {

  render() {
    const { endpoint, client } = this.props

    const check = () => {
      var token = "" + localStorage.getItem("auth-token")

      client.checkCreds(token, (result) => {
        console.log("cred check = ", result)
      })
    }
    return (
      <div>
        <h1>Placeholder ({endpoint})</h1>
        <ul>
          <li><Link to="/admin/post">Posts</Link></li>
          <li><Link to="/admin/authors">Authors</Link></li>
        </ul>
        <button onClick={check}>Check Creds</button>
      </div>
    )
  }
}

const PropsRoute = ({component: Component, path: Path, ...rest}) => (
  <Route path={Path} exact render={(props) => (<Component {...rest} {...props}/> )}/>
)

class App extends React.PureComponent {
  constructor(props) {
    super(props)

    const { endpoint } = this.props

    this.client = new Client(endpoint, this.onError)
  }

  componentDidMount() {
  }

  render() {
    const { endpoint } = this.props
    return (
      <Router>
        <section>
          <PropsRoute path="/" component={Redirect} to="/admin/home" endpoint={endpoint}/>
          <PropsRoute path="/admin" component={Redirect} to="/admin/home" endpoint={endpoint}/>
          <PropsRoute path="/admin/home" component={Home} endpoint={endpoint} client={this.client}/>
          <PropsRoute path="/admin/post" component={Posts}/>
          <PropsRoute path="/admin/authors" component={Authors}/>
        </section>
      </Router>
    );
  }
}

export default App;
