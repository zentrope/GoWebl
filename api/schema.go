package api

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
	 id: ID!
	 author: Author!
	 status: String!
	 text: String!
	 dateCreated: String!
	 dateUpdated: String!
 }
`
