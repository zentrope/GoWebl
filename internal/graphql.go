package internal

import (
	"fmt"
	"time"

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
	 hello: String!
	 nums: [Int!]!
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
// Random trial balloons
//=============================================================================

func (r *Resolver) Hello() (string, error) {
	return "Hello world!", nil
}

func (r *Resolver) Nums() ([]int32, error) {
	l := make([]int32, 3)
	l[0] = 1
	l[1] = 2
	l[2] = 3
	return l, nil
}

//=============================================================================
// Author
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

func (r *Resolver) CreateAuthor(args *struct{ Author *AuthorInput }) (*authorResolver, error) {
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

func (r *Resolver) Authors() ([]*authorResolver, error) {
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

func (r *Resolver) Posts() ([]*postResolver, error) {
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
