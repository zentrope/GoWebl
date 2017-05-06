// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"context"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	graphql "github.com/neelance/graphql-go"
)

//=============================================================================
// Schema
//=============================================================================

const Schema = `
 schema {
	 query: Query
	 mutation: Mutation
 }

 type Query {
	 validate(creds: CredInput!): Boolean!
	 authenticate(creds: LoginInput!): String!

	 authors: [Author]!
	 posts: [Post]!
 }

 type Mutation {
	 createAuthor(author: AuthorInput!): Author
 }

 input AuthorInput {
	 handle: String!
	 email: String!
	 password: String!
 }

 input CredInput {
	 token: String!
 }

 input LoginInput {
	 user: String!
	 pass: String!
 }

 type Author {
	 id: ID!
	 handle: String!
	 email: String!
	 type: String!
	 status: String!
	 posts: [Post]!
 }

 type Post {
	 uuid: ID!
	 author: Author!
	 status: String!
	 slugline: String!
	 text: String!
	 dateCreated: String!
	 dateUpdated: String!
 }
`

//=============================================================================
// Root Resolver
//=============================================================================

type AuthKeyContextType string

const AUTH_KEY = AuthKeyContextType("auth-key")

type GraphAPI struct {
	Schema *graphql.Schema
}

type Resolver struct {
	Database *Database
}

func NewApi(database *Database) (*GraphAPI, error) {
	schema, err := graphql.ParseSchema(Schema, &Resolver{database})
	if err != nil {
		return nil, err
	}

	return &GraphAPI{schema}, nil
}

//=============================================================================
// Auth Tokens (JWT)
//=============================================================================

// TODO(keith): Grab from env and config
var SECRET = []byte("thirds-and-fifths")

type ViewerClaims struct {
	User string `json:"user"`
	Type string `json:"type"`
	jwt.StandardClaims
}

func mkAuthToken(author *Author) (string, error) {
	claims := ViewerClaims{
		author.Handle, author.Type,
		jwt.StandardClaims{
			Issuer: "vaclav",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SECRET)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func checkAlgKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	return SECRET, nil
}

func isValidAuthToken(tokenString string) (bool, error) {

	token, err := jwt.ParseWithClaims(tokenString, &ViewerClaims{}, checkAlgKey)

	if err != nil {
		return false, err
	}

	return token.Valid, nil
}

func decodeAuthToken(tokenString string) (*ViewerClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &ViewerClaims{}, checkAlgKey)

	if claims, ok := token.Claims.(*ViewerClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

//=============================================================================
// Auth
//=============================================================================

type CredInput struct {
	Token string
}

type LoginInput struct {
	User string
	Pass string
}

func (r *Resolver) Validate(args *struct{ Creds *CredInput }) (bool, error) {
	tokenString := args.Creds.Token
	return isValidAuthToken(tokenString)
}

func (r *Resolver) Authenticate(args *struct{ Creds *LoginInput }) (string, error) {
	user := args.Creds.User
	pass := args.Creds.Pass

	author, err := r.Database.Authentic(user, pass)
	if err != nil {
		return "", err
	}

	return mkAuthToken(author)
}

//=============================================================================
// author
//=============================================================================

type authorResolver struct {
	database *Database
	author   *Author
}

type AuthorInput struct {
	Handle   string
	Email    string
	Password string
}

func (r *Resolver) CreateAuthor(ctx context.Context, args *struct{ Author *AuthorInput }) (*authorResolver, error) {
	fmt.Printf("AUTH: %v\n", ctx.Value(AUTH_KEY))
	err := r.Database.CreateAuthor(args.Author.Handle, args.Author.Email, args.Author.Password)

	if err != nil {
		return nil, err
	}

	author, err := r.Database.Author(args.Author.Handle)

	if err != nil {
		return nil, err
	}

	return &authorResolver{r.Database, author}, nil
}

func (r *Resolver) Authors(ctx context.Context) ([]*authorResolver, error) {
	fmt.Printf("AUTH: %v\n", ctx.Value(AUTH_KEY))
	authors, err := r.Database.Authors()
	if err != nil {
		return nil, err
	}
	var rs []*authorResolver
	for _, a := range authors {
		rs = append(rs, &authorResolver{r.Database, a})
	}
	return rs, nil
}

func (r *authorResolver) ID() graphql.ID {
	return graphql.ID(r.author.Handle)
}

func (r *authorResolver) Handle() string {
	return r.author.Handle
}

func (r *authorResolver) Email() string {
	return r.author.Email
}

func (r *authorResolver) Type() string {
	return r.author.Type
}

func (r *authorResolver) Status() string {
	return r.author.Status
}

func (r *authorResolver) Posts() ([]*postResolver, error) {
	posts, err := r.database.PostsByAuthor(r.author.Handle)
	if err != nil {
		return nil, err
	}
	var rs []*postResolver
	for _, p := range posts {
		rs = append(rs, &postResolver{r.database, p})
	}
	return rs, nil
}

//=============================================================================
// Posts
//=============================================================================

type postResolver struct {
	database *Database
	post     *Post
}

func (r *Resolver) Posts(ctx context.Context) ([]*postResolver, error) {
	fmt.Printf("AUTH: %v\n", ctx.Value(AUTH_KEY))
	posts, err := r.Database.Posts()
	if err != nil {
		return nil, err
	}

	var rs []*postResolver
	for _, p := range posts {
		rs = append(rs, &postResolver{r.Database, p})
	}
	return rs, nil
}

func (r *postResolver) UUID() graphql.ID {
	return graphql.ID(fmt.Sprintf("%v", r.post.UUID))
}

func (r *postResolver) Author() (*authorResolver, error) {
	author, err := r.database.Author(r.post.Author)
	if err != nil {
		return nil, err
	}
	return &authorResolver{r.database, author}, nil
}

func (r *postResolver) Status() string {
	return r.post.Status
}

func (r *postResolver) Slugline() string {
	return r.post.Slugline
}

func (r *postResolver) Text() string {
	return r.post.Text
}

func (r *postResolver) DateCreated() string {
	return r.post.DateCreated.Format(time.RFC3339)
}

func (r *postResolver) DateUpdated() string {
	return r.post.DateUpdated.Format(time.RFC3339)
}

//=============================================================================
