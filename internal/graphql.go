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
	 validate(token: String!): Boolean!
	 authenticate(user: String! pass: String!): String!
	 viewer(token: String): Viewer!
 }

 type Mutation {
	 createPost(slugline: String! status: String! text: String! token: String): Post!
	 updatePost(uuid: String! slugline: String! text: String!): Post!
	 deletePost(uuid: String!): String!
	 setPostStatus(uuid: String! isPublished: Boolean!): Post!
 }

 type Viewer {
	 id: ID!
	 user: String!
	 email: String!
	 type: String!
	 posts: [Post]!
 }

 type Author {
	 id: ID!
	 handle: String!
	 email: String!
	 type: String!
	 status: String!
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

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*ViewerClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

func findAuthClaims(ctx context.Context, token *string) (*ViewerClaims, error) {
	auth := ctx.Value(AUTH_KEY).(string)
	if token != nil {
		auth = *token
	}

	claims, err := decodeAuthToken(auth)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

//=============================================================================
// Auth
//=============================================================================

func (r *Resolver) Validate(args *struct{ Token string }) (bool, error) {
	tokenString := args.Token
	// TODO: Make sure the user in the token still exists
	return isValidAuthToken(tokenString)
}

func (r *Resolver) Authenticate(args *struct{ User, Pass string }) (string, error) {
	user := args.User
	pass := args.Pass

	author, err := r.Database.Authentic(user, pass)
	if err != nil {
		return "", err
	}

	return mkAuthToken(author)
}

//=============================================================================
// Viewer
//=============================================================================

type Viewer struct {
	ID    graphql.ID
	User  string
	Type  string
	Posts []*Post
}

type viewerResolver struct {
	database *Database
	author   *Author
}

func (r *Resolver) Viewer(ctx context.Context, args *struct{ Token *string }) (*viewerResolver, error) {

	claims, err := findAuthClaims(ctx, args.Token)
	if err != nil {
		return nil, err
	}

	author, err := r.Database.Author(claims.User)

	if err != nil {
		return nil, err
	}

	return &viewerResolver{r.Database, author}, nil
}

func (v *viewerResolver) ID() graphql.ID {
	return graphql.ID(v.author.Handle)
}

func (v *viewerResolver) User() string {
	return v.author.Handle
}

func (v *viewerResolver) Email() string {
	return v.author.Email
}

func (v *viewerResolver) Type() string {
	return v.author.Type
}

func (v *viewerResolver) Posts() ([]*postResolver, error) {
	posts, err := v.database.PostsByAuthor(v.author.Handle)
	if err != nil {
		return nil, err
	}
	var rs []*postResolver
	for _, p := range posts {
		rs = append(rs, &postResolver{v.database, p})
	}
	return rs, nil
}

//=============================================================================
// Author
//=============================================================================

type authorResolver struct {
	database *Database
	author   *Author
}

func (r *Resolver) Authors(ctx context.Context) ([]*authorResolver, error) {
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

//=============================================================================
// Posts
//=============================================================================

type postResolver struct {
	database *Database
	post     *Post
}

func (r *Resolver) CreatePost(ctx context.Context, args *struct {
	Slugline string
	Status   string
	Text     string
	Token    *string
}) (*postResolver, error) {

	claims, err := findAuthClaims(ctx, args.Token)
	if err != nil {
		return nil, err
	}

	uuid, err := r.Database.CreatePost(
		claims.User,
		args.Slugline,
		args.Status,
		args.Text,
	)

	if err != nil {
		return nil, err
	}

	post, err := r.Database.Post(uuid)
	if err != nil {
		return nil, err
	}

	return &postResolver{r.Database, post}, nil
}

func (r *Resolver) UpdatePost(ctx context.Context, args *struct {
	Uuid     string
	Slugline string
	Text     string
}) (*postResolver, error) {

	claims, err := findAuthClaims(ctx, nil)
	if err != nil {
		return nil, err
	}

	post, err := r.Database.UpdatePost(args.Uuid, args.Slugline, args.Text, claims.User)
	if err != nil {
		return nil, err
	}

	return &postResolver{r.Database, post}, nil
}

func (r *Resolver) DeletePost(ctx context.Context, args *struct{ Uuid string }) (string, error) {

	claims, err := findAuthClaims(ctx, nil)
	if err != nil {
		return "", err
	}

	if err := r.Database.DeletePost(args.Uuid, claims.User); err != nil {
		return "", err
	}
	return args.Uuid, nil
}

func (r *Resolver) SetPostStatus(ctx context.Context, args *struct {
	Uuid        string
	IsPublished bool
}) (*postResolver, error) {

	claims, err := findAuthClaims(ctx, nil)
	if err != nil {
		return nil, err
	}

	status := PS_Draft
	if args.IsPublished {
		status = PS_Published
	}

	post, err := r.Database.SetPostStatus(args.Uuid, claims.User, status)
	if err != nil {
		return nil, err
	}

	return &postResolver{r.Database, post}, nil
}

func (r *Resolver) Posts(ctx context.Context) ([]*postResolver, error) {
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
