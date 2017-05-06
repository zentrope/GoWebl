// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// const fl = (s) =>
//   s.replace(/\s+/g, " ")

const headers =
  { "Content-Type": "application/json; charset=utf8" }

const checkCreds = (token) => {
  return {
    method: 'POST',
    headers: headers,
    body: JSON.stringify({
      query: "query CheckCreds($input: CredInput) { checkCreds(creds: $input) }",
      operationName: "CheckCreds",
      variables: {input: {token: token}}
    })
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
    this.url = url + "/query/"
    this.errorDelegate = errorDelegate
  }

  checkCreds(token, callback) {
    fetch(this.url, checkCreds(token))
      .then(res => checkStatus(res))
      .then(res => res.json())
      .then(data => callback(data))
      .catch(err => err.response.json()
                       .then(data => this.errorDelegate(data)))

  }

}

export { Client }
