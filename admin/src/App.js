import React from 'react';
import { BrowserRouter as Router, Route, Link, Redirect } from 'react-router-dom'

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
    const { endpoint } = this.props
    return (
      <div>
        <h1>Placeholder ({endpoint})</h1>
        <ul>
          <li><Link to="/admin/post">Posts</Link></li>
          <li><Link to="/admin/authors">Authors</Link></li>
        </ul>
      </div>
    )
  }
}

const PropsRoute = ({component: Component, path: Path, ...rest}) => (
  <Route path={Path} exact render={(props) => (<Component {...rest} {...props}/> )}/>
)

class App extends React.PureComponent {
  render() {
    const { endpoint } = this.props
    return (
      <Router>
        <section>
          <PropsRoute path="/" component={Redirect} to="/admin/home" endpoint={endpoint}/>
          <PropsRoute path="/admin" component={Redirect} to="/admin/home" endpoint={endpoint}/>
          <PropsRoute path="/admin/home" component={Home} endpoint={endpoint}/>
          <PropsRoute path="/admin/post" component={Posts}/>
          <PropsRoute path="/admin/authors" component={Authors}/>
        </section>
      </Router>
    );
  }
}

export default App;
