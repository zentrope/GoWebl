// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

const fl = (s) =>
  s.replace(/\s+/g, " ")

const createPostQL = (slugline, status, text) => {
  const q = fl(`mutation
    CreatePost($slugline: String! $status: String! $text: String! $token: String) {
      createPost(slugline: $slugline, status: $status, text: $text, token: $token) {
        uuid slugline status dateCreated dateUpdated text }}`)

  return {
    query: q,
    operationName: "CreatePost",
    variables: {
      slugline: slugline,
      status: status,
      text: text,
    }
  }
}

const validateQL = (token) => {
  return {
    query: "query Validate($token: String!) { validate(token: $token) }",
    operationName: "Validate",
    variables: {token: token}
  }
}

const loginQL = (user, pass) => {
  const q = fl(`query
    Authenticate($user: String! $pass: String!) {
      authenticate(user: $user pass: $pass) }`)
  return {
    query: q,
    operationName: "Authenticate",
    variables: {user: user, pass: pass}
  }
}

const viewerQL = () => {
  const q = fl(`query {
    viewer { id user type email
      posts { uuid status slugline dateCreated dateUpdated text } } }
  `)
  return { query: q }
}

const checkStatus = (response) => {
  if (response.status >= 200 && response.status < 300) {
    return response
  } else {
    let error = new Error(response.statusText)
    error.response = response
    throw error
  }
}

class Client {

  constructor(url, errorDelegate) {
    this.url = url + "/query"
    this.errorDelegate = errorDelegate
    this.authToken = "no-auth"
  }

  __mkHeaders() {
    return {
      "Content-Type": "application/json; charset=utf8",
      "Authorization": "Bearer " + this.authToken
    }
  }

  __mkPost(body) {
    return {
      method: 'POST',
      headers: this.__mkHeaders(),
      body: JSON.stringify(body)
    }
  }

  __doQuery(ql, callback) {
    let query = this.__mkPost(ql)
    fetch(this.url, query)
      .then(res => checkStatus(res))
      .then(res => res.json())
      .then(data => callback(data))
      .catch(err => err.response.json()
                       .then(data => this.errorDelegate(data)))
  }

  setAuthToken(token) {
    this.authToken = token
  }

  invalidateAuthToken() {
    this.authToken = "no-auth"
  }

  login(user, pass, callback) {
    this.__doQuery(loginQL(user, pass), callback)
  }

  validate(token, callback) {
    this.__doQuery(validateQL(token), callback)
  }

  viewerData(callback) {
    this.__doQuery(viewerQL(), callback)
  }

  savePost(slugline, text, status, callback) {
    this.__doQuery(createPostQL(slugline, status, text), callback)
  }

}

export { Client }
