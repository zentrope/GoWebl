package api

import (
	"fmt"
	"time"

	graphql "github.com/neelance/graphql-go"
	db "github.com/zentrope/webl/database"
)

type Resolver struct {
	Database *db.Database
}

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
	database *db.Database
	author   *db.Author
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

	author := r.Database.Author(args.Author.Handle)

	return &authorResolver{r.Database, author}, nil
}

func (r *Resolver) Authors() []*authorResolver {
	authors := r.Database.Authors()
	var rs []*authorResolver
	for _, a := range authors {
		rs = append(rs, &authorResolver{r.Database, a})
	}
	return rs
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

func (r *authorResolver) Posts() []*postResolver {
	posts := r.database.PostsByAuthor(r.author.Handle)
	var rs []*postResolver
	for _, p := range posts {
		rs = append(rs, &postResolver{r.database, p})
	}
	return rs
}

//=============================================================================
// Posts
//=============================================================================

type postResolver struct {
	database *db.Database
	post     *db.Post
}

func (r *Resolver) Posts() []*postResolver {
	posts := r.Database.Posts()
	var rs []*postResolver
	for _, p := range posts {
		rs = append(rs, &postResolver{r.Database, p})
	}
	return rs
}

func (r *postResolver) ID() graphql.ID {
	return graphql.ID(fmt.Sprintf("%v", r.post.Id))
}

func (r *postResolver) Author() *authorResolver {
	return &authorResolver{r.database, r.database.Author(r.post.Author)}
}

func (r *postResolver) Status() string {
	return r.post.Status
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
