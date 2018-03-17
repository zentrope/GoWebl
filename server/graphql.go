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

package server

import (
	"context"
	"fmt"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	graphql "github.com/graph-gophers/graphql-go"
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
	 authenticate(email: String! pass: String!): Viewer!
	 viewer(token: String): Viewer!
	 site(): Site!
 }

 type Mutation {
	 // TODO: Should take objects to allow for easier (and sparse) extension.
	 createPost(slugline: String! status: String! text: String! datePublished: String! token: String): Post!
	 updatePost(uuid: String! slugline: String! text: String! datePublished: String!): Post!
	 deletePost(uuid: String!): String!
	 setPostStatus(uuid: String! isPublished: Boolean!): Post!
	 updateSite(baseUrl: String! description: String! title: String!): Site!

	 // Assumes user id is in auth header
	 updateViewer(name: String! email: String!): Viewer!
	 updateViewerPassword(password: String!): Viewer!
 }

 type Viewer {
	 id: ID!
	 name: String!
	 email: String!
	 type: String!
	 token: String!
	 posts: [Post]!
	 site: Site!
 }

 type Author {
	 id: ID!
	 email: String!
	 name: String!
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
	 datePublished: String!
	 wordCount: Int!
 }

 type Site {
	 baseUrl: String!
	 title: String!
	 description: String!
 }
`

//=============================================================================
// Root Resolver
//=============================================================================

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

const BAD_AUTH_MSG = "Invalid authorization request."

type ViewerClaims struct {
	Uuid string `json:"uuid"`
	Type string `json:"type"`
	jwt.StandardClaims
}

func getSecret(ctx context.Context) []byte {
	site := ctx.Value(SITE_KEY).(*SiteConfig)
	return []byte(site.JwtSecret)
}

func mkAuthToken(ctx context.Context, author *Author) (string, error) {
	secret := getSecret(ctx)

	claims := ViewerClaims{
		author.Uuid,
		author.Type,
		jwt.StandardClaims{
			Issuer: "vaclav",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func checkAlgKey(ctx context.Context) jwt.Keyfunc {
	secret := getSecret(ctx)
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	}
}

func isValidAuthToken(ctx context.Context, tokenString string) (bool, error) {

	token, err := jwt.ParseWithClaims(tokenString, &ViewerClaims{}, checkAlgKey(ctx))

	if err != nil {
		log.Printf("auth.error: %v", err)
		return false, fmt.Errorf(BAD_AUTH_MSG)
	}

	return token.Valid, nil
}

func decodeAuthToken(ctx context.Context, tokenString string) (*ViewerClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &ViewerClaims{}, checkAlgKey(ctx))

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*ViewerClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

func getAuthToken(ctx context.Context) string {
	return ctx.Value(AUTH_KEY).(string)
}

func optionalAuthToken(ctx context.Context, token *string) string {
	auth := getAuthToken(ctx)
	if token != nil {
		auth = *token
	}
	return auth
}

func findAuthClaims(ctx context.Context, token *string) (*ViewerClaims, error) {
	auth := optionalAuthToken(ctx, token)

	claims, err := decodeAuthToken(ctx, auth)
	if err != nil {
		log.Printf("auth.error: %v", err)
		return nil, fmt.Errorf(BAD_AUTH_MSG)
	}

	return claims, nil
}

//=============================================================================
// Auth
//=============================================================================

func (r *Resolver) Validate(ctx context.Context, args *struct{ Token string }) (bool, error) {
	tokenString := args.Token
	// TODO: Make sure the user in the token still exists
	return isValidAuthToken(ctx, tokenString)
}

func (r *Resolver) Authenticate(ctx context.Context, args *struct{ Email, Pass string }) (*viewerResolver, error) {
	email := args.Email
	pass := args.Pass

	author, err := r.Database.Authentic(email, pass)
	if err != nil {
		log.Printf("auth.error: %v -> %v", email, err)
		return nil, fmt.Errorf(BAD_AUTH_MSG)
	}

	token, err := mkAuthToken(ctx, author)
	if err != nil {
		return nil, err
	}

	res := viewerResolver{
		database: r.Database,
		author:   author,
		token:    token,
		site:     ctx.Value(SITE_KEY).(*SiteConfig),
	}

	return &res, nil
}

//=============================================================================
// Viewer
//=============================================================================

type viewerResolver struct {
	database *Database
	author   *Author
	token    string
	site     *SiteConfig
}

func (r *Resolver) UpdateViewerPassword(ctx context.Context, args *struct {
	Password string
}) (*viewerResolver, error) {

	claims, err := findAuthClaims(ctx, nil)
	if err != nil {
		return nil, err
	}

	author, err := r.Database.UpdateAuthorPassword(claims.Uuid, args.Password)
	if err != nil {
		return nil, err
	}

	return &viewerResolver{
		database: r.Database,
		author:   author,
		site:     ctx.Value(SITE_KEY).(*SiteConfig),
		token:    optionalAuthToken(ctx, nil),
	}, nil
}

func (r *Resolver) UpdateViewer(ctx context.Context, args *struct {
	Name  string
	Email string
}) (*viewerResolver, error) {

	claims, err := findAuthClaims(ctx, nil)
	if err != nil {
		return nil, err
	}

	author, err := r.Database.UpdateAuthor(claims.Uuid, args.Name, args.Email)
	if err != nil {
		return nil, err
	}

	return &viewerResolver{
		database: r.Database,
		author:   author,
		site:     ctx.Value(SITE_KEY).(*SiteConfig),
		token:    optionalAuthToken(ctx, nil),
	}, nil
}

func (r *Resolver) Viewer(ctx context.Context, args *struct{ Token *string }) (*viewerResolver, error) {

	claims, err := findAuthClaims(ctx, args.Token)
	if err != nil {
		return nil, err
	}

	author, err := r.Database.Author(claims.Uuid)
	if err != nil {
		return nil, err
	}

	site := ctx.Value(SITE_KEY).(*SiteConfig)

	return &viewerResolver{
		database: r.Database,
		author:   author,
		site:     site,
		token:    optionalAuthToken(ctx, args.Token),
	}, nil
}

func (v *viewerResolver) ID() graphql.ID {
	return graphql.ID(v.author.Uuid)
}

func (v *viewerResolver) Name() string {
	return v.author.Name
}

func (v *viewerResolver) Email() string {
	return v.author.Email
}

func (v *viewerResolver) Type() string {
	return v.author.Type
}

func (r *viewerResolver) Token() string {
	return r.token
}

func (v *viewerResolver) Posts() ([]*postResolver, error) {
	posts, err := v.database.PostsByAuthor(v.author.Uuid)
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
// Site
//=============================================================================

type siteResolver struct {
	site *SiteConfig
}

func (r *Resolver) UpdateSite(ctx context.Context, args *struct {
	Title       string
	Description string
	BaseURL     string
}) (*siteResolver, error) {

	_, err := findAuthClaims(ctx, nil)
	if err != nil {
		return nil, err
	}

	site, err := r.Database.UpdateSite(args.Title, args.Description, args.BaseURL)
	if err != nil {
		return nil, err
	}

	return &siteResolver{site}, nil
}

func (r *Resolver) Site(ctx context.Context) siteResolver {
	return siteResolver{
		site: ctx.Value(SITE_KEY).(*SiteConfig),
	}
}

func (r *viewerResolver) Site() siteResolver {
	return siteResolver{site: r.site}
}

func (r siteResolver) BaseURL() string {
	return r.site.BaseURL
}

func (r siteResolver) Title() string {
	return r.site.Title
}

func (r siteResolver) Description() string {
	return r.site.Description
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
	return graphql.ID(r.author.Uuid)
}

func (r *authorResolver) Name() string {
	return r.author.Name
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
	Slugline      string
	Status        string
	Text          string
	DatePublished string
	Token         *string
}) (*postResolver, error) {

	claims, err := findAuthClaims(ctx, args.Token)
	if err != nil {
		return nil, err
	}

	uuid, err := r.Database.CreatePost(
		claims.Uuid,
		args.Slugline,
		args.Status,
		args.Text,
		args.DatePublished,
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
	Uuid          string
	Slugline      string
	Text          string
	DatePublished string
}) (*postResolver, error) {

	claims, err := findAuthClaims(ctx, nil)
	if err != nil {
		return nil, err
	}

	post, err := r.Database.UpdatePost(args.Uuid, args.Slugline, args.Text, args.DatePublished, claims.Uuid)
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

	if err := r.Database.DeletePost(args.Uuid, claims.Uuid); err != nil {
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

	post, err := r.Database.SetPostStatus(args.Uuid, claims.Uuid, status)
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
	author, err := r.database.Author(r.post.AuthorUuid)
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

func (r *postResolver) DatePublished() string {
	return r.post.DatePublished.Format(time.RFC3339)
}

func (r *postResolver) WordCount() int32 {
	return int32(r.post.WordCount)
}
