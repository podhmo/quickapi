//go:generate go run ./ -gendoc -docfile openapi.json -mdfile apidoc.md
package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/podhmo/quickapi"
	"github.com/podhmo/quickapi/qopenapi/define"
	reflectopenapi "github.com/podhmo/reflect-openapi"
)

//go:embed openapi.json
var openapiDocData []byte

//go:embed apidoc.md
var mdDocData []byte

var options struct {
	gendoc  bool
	docfile string
	mdfile  string
	port    int
}

func main() {
	flag.BoolVar(&options.gendoc, "gendoc", false, "generate openapi doc")
	flag.IntVar(&options.port, "port", 8080, "port")
	flag.StringVar(&options.docfile, "docfile", "", "file name of openapi doc. if this value is empty output to stdout.")
	flag.StringVar(&options.mdfile, "mdfile", "", "")
	flag.Parse()
	if err := run(); err != nil {
		log.Fatalf("!! %+v", err)
	}

}

func run() error {
	ctx := context.Background()

	// TODO: contact: name: Realworld, url: https://realworld.how
	// ToDO: licence: mit

	doc := define.Doc().
		Title("Conduit API").
		Description("Conduit API documentation").
		Version("1.0.0").
		Server(fmt.Sprintf("http://localhost:%d", options.port), "local development")

	if !options.gendoc {
		doc = doc.LoadFromData(openapiDocData)
	}

	router := quickapi.DefaultRouter()
	bc, err := define.NewBuildContext(doc, router, func(c *reflectopenapi.Config) {
		c.EnableAutoTag = false
	})
	if err != nil {
		return fmt.Errorf("build context: %w", err)
	}

	mount(bc)

	if options.gendoc {
		if err := bc.EmitDoc(ctx, options.docfile); err != nil {
			return err
		}

		if options.mdfile != "" {
			if err := bc.EmitMDDoc(ctx, options.mdfile); err != nil {
				return err
			}
		}
		return nil
	}

	handler, err := bc.BuildHandler(ctx)
	if err != nil {
		return err
	}
	dochandler, err := bc.BuildDocHandler(ctx, "/_doc", mdDocData)
	if err != nil {
		return err
	}
	bc.Router().Mount("/_doc", dochandler)

	if err := quickapi.NewServer(fmt.Sprintf(":%d", options.port), handler, 5*time.Second).ListenAndServe(ctx); err != nil {
		log.Printf("[Error] !! %+v", err)
	}
	return nil
}

func mount(bc *define.BuildContext) {
	tagArticles := "Articles"
	tagComments := "Comments"
	tagFavorites := "Favorites"
	tagProfile := "Profile"
	tagTags := "Tags"
	tagUserAndAuthentication := "User and Authentication"

	define.Type(bc, OffsetParam(0)).After(func(s *openapi3.Schema) {
		s.WithMin(0)
	})
	define.Type(bc, LimitParam(20)).After(func(s *openapi3.Schema) {
		// TODO: default value by ↑
		s.WithMin(1).WithDefault(20)
	})

	define.Post(bc, "/users/login", Login).Tags(tagUserAndAuthentication)
	define.Post(bc, "/users/", CreateUser).Tags(tagUserAndAuthentication)
	define.Get(bc, "/user", GetCurrentUser).Tags(tagUserAndAuthentication)
	define.Put(bc, "/user", UpdateCurrentUser).Tags(tagUserAndAuthentication)

	define.Get(bc, "/profiles/{username}", GetProfileByUsername).Tags(tagProfile)
	define.Post(bc, "/profiles/{username}/follow", FollowUserByUsername).Tags(tagProfile)
	define.Delete(bc, "/profiles/{username}/follow", UnfollowUserByUsername).Tags(tagProfile)

	define.Get(bc, "/articles/feed", GetArticlesFeed).Tags(tagArticles)
	define.Get(bc, "/articles", GetArticles).Tags(tagArticles)
	define.Post(bc, "/articles", CreateArticle).Tags(tagArticles)
	define.Get(bc, "/articles/{slug}", GetArticle).Tags(tagArticles)
	define.Put(bc, "/articles/{slug}", UpdateArticle).Tags(tagArticles)
	define.Delete(bc, "/articles/{slug}", DeleteArticle).Tags(tagArticles)
	define.Get(bc, "/articles/{slug}/comments", GetArticleComments).Tags(tagComments)
	define.Post(bc, "/articles/{slug}/comments", CreateArticleComment).Tags(tagComments)
	define.Delete(bc, "/articles/{slug}/comments/{id}", DeleteArticleComment).Tags(tagComments)
	define.Post(bc, "/articles/{slug}/favorite", CreateArticleFavorite).Tags(tagFavorites)
	define.Delete(bc, "/articles/{slug}/favorite", DeleteArticleFavorite).Tags(tagFavorites)

	define.Get(bc, "/tags", GetTags).Tags(tagTags)
}

// handlers

// Existing user login
//
// Login for existing user
func Login(ctx context.Context, input struct{}) (output struct{}, err error) {
	return struct{}{}, nil
}

// Register a new user
func CreateUser(ctx context.Context, input struct{}) (output struct{}, err error) {
	return struct{}{}, nil
}

// Get current user
//
// Gets the currently logged-in user
func GetCurrentUser(ctx context.Context, input struct{}) (output struct{}, err error) {
	return struct{}{}, nil
}

// Update current user
//
// Update user information for current user
func UpdateCurrentUser(ctx context.Context, input struct{}) (output struct{}, err error) {
	return struct{}{}, nil
}

type GetProfileByUsernameInput struct {
	Username string `in:"path" path:"username"` // username of the profile to get
}

// Get a profile
//
// Get a profile of a user of the system. Auth is optional
func GetProfileByUsername(ctx context.Context, input GetProfileByUsernameInput) (output struct{}, err error) {
	return struct{}{}, nil
}

type FollowUserByUsernameInput struct {
	Username string `in:"path" path:"username"` // username of the profile you want to follow
}

// Follow a user
//
// Follow a user by username
func FollowUserByUsername(ctx context.Context, input FollowUserByUsernameInput) (output struct{}, err error) {
	return struct{}{}, nil
}

type UnfollowUserByUsernameInput struct {
	Username string `in:"path" path:"username"` // username of the profile you want to unfollow
}

// Unfollow a user
//
// Unfollow a user by username
func UnfollowUserByUsername(ctx context.Context, input UnfollowUserByUsernameInput) (output struct{}, err error) {
	return struct{}{}, nil
}

// The numbers of items to return.
type LimitParam int

// The number of items to skip before starting to collect the result set.
type OffsetParam int

type GetArticlesFeedInput struct {
	Limit  LimitParam  `in:"query" query:"limit"`
	Offset OffsetParam `in:"query" query:"offset"`
}

// Get recent articles from users you follow
//
// Get most recent articles from users you follow. Use query parameters to limit. Auth is required
func GetArticlesFeed(ctx context.Context, input GetArticlesFeedInput) (output struct{}, err error) {
	return struct{}{}, nil
}

// Get recent articles globally
//
// Get most recent articles globally. Use query parameters to filter results. Auth is optional
func GetArticles(ctx context.Context, input struct{}) (output struct{}, err error) {
	return struct{}{}, nil
}

// Create an article
//
// Create an article. Auth is required
func CreateArticle(ctx context.Context, input struct{}) (output struct{}, err error) {
	return struct{}{}, nil
}

type GetArticleInput struct {
	Slug string `in:"path" path:"slug"` // Slug of the article to get
}

// Get an article
//
// Get an article. Auth not required
func GetArticle(ctx context.Context, input GetArticleInput) (output struct{}, err error) {
	return struct{}{}, nil
}

type UpdateArticleInput struct {
	Slug string `in:"path" path:"slug"` // Slug of the article to update
}

// Update an article
//
// Update an article. Auth is required
func UpdateArticle(ctx context.Context, input UpdateArticleInput) (output struct{}, err error) {
	return struct{}{}, nil
}

type DeleteArticleInput struct {
	Slug string `in:"path" path:"slug"` // Slug of the article to delete
}

// Delete an article
//
// Delete an article. Auth is required
func DeleteArticle(ctx context.Context, input DeleteArticleInput) (output struct{}, err error) {
	return struct{}{}, nil
}

type GetArticleCommentsInput struct {
	Slug string `in:"path" path:"slug"` // Slug of the article that you want to get comments for
}

// Get comments for an article
//
// Get the comments for an article. Auth is optional
func GetArticleComments(ctx context.Context, input GetArticleCommentsInput) (output struct{}, err error) {
	return struct{}{}, nil
}

type CreateArticleCommentInput struct {
	Slug string `in:"path" path:"slug"` // Slug of the article that you want to create a comment
}

// Create a comment for an article
//
// Create a comment for an article. Auth is required
func CreateArticleComment(ctx context.Context, input CreateArticleCommentInput) (output struct{}, err error) {
	return struct{}{}, nil
}

type DeleteArticleCommentInput struct {
	Slug string `in:"path" path:"slug"` // Slug of the article that you want to delete a comment
	ID   string `in:"path" path:"id"`   // ID of the comment you want to delete
}

// Delete a comment for an article
//
// Delete a comment for an article. Auth is required.
func DeleteArticleComment(ctx context.Context, input DeleteArticleCommentInput) (output struct{}, err error) {
	return struct{}{}, nil
}

type CreateArticleFavoriteInput struct {
	Slug string `in:"path" path:"slug"` // Slug of the article that you want to favorite
}

// Favorite an article
//
// Favorite an article. Auth is required
func CreateArticleFavorite(ctx context.Context, input CreateArticleFavoriteInput) (output struct{}, err error) {
	return struct{}{}, nil
}

type DeleteArticleFavoriteInput struct {
	Slug string `in:"path" path:"slug"` // Slug of the article that you want to unfavorite
}

// Unfavorite an article
//
// Unfavorite an article. Auth is required
func DeleteArticleFavorite(ctx context.Context, input DeleteArticleFavoriteInput) (output struct{}, err error) {
	return struct{}{}, nil
}

type GetTagsInput struct {
	Query     string `in:"query" query:"tag"`       // Filter by tag
	Author    string `in:"query" query:"author"`    // Filter by author (username)
	Favorited string `in:"query" query:"favorited"` // Filter by favorites of a user (username)

	Limit  LimitParam  `in:"query" query:"limit"`
	Offset OffsetParam `in:"query" query:"offset"`
}

// Get tags
//
// Get tags. Auth not required
func GetTags(ctx context.Context, input GetTagsInput) (output struct{}, err error) {
	return struct{}{}, nil
}
