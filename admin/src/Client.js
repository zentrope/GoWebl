// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// const fl = (s) =>
//   s.replace(/\s+/g, " ")

const validateQL = (token) => {
  return {
      query: "query Validate($input: CredInput) { validate(creds: $input) }",
      operationName: "Validate",
      variables: {input: {token: token}}
  }
}

const loginQL = (user, pass) => {
  return {
    query: "query Authenticate($input: LoginInput) { authenticate(creds: $input) }",
    operationName: "Authenticate",
    variables: {input: {user: user, pass: pass}}
  }
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

}

export { Client }
