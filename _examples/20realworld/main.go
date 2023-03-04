//go:generate go run ./ -gendoc -docfile openapi.json -mdfile apidoc.md
package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"go/token"
	"log"
	"net/http"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/podhmo/quickapi"
	"github.com/podhmo/quickapi/qopenapi/define"
	reflectopenapi "github.com/podhmo/reflect-openapi"
	reflectshape "github.com/podhmo/reflect-shape"
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
		c.GoPositionFunc = func(fset *token.FileSet, f *reflectshape.Func) string {
			// TODO: multiple package
			fpos := fset.Position(f.Pos())
			return fmt.Sprintf("https://github.com/podhmo/quickapi/blob/main/_examples/20realworld/main.go#L%d", fpos.Line)
		}
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
		s.WithMin(1)
	})

	define.Post(bc, "/users/login", Login).Tags(tagUserAndAuthentication).AnotherError(bc, 401, &UnauthorizedError{}, "").AnotherError(bc, 422, &GenericError{}, "")
	define.Post(bc, "/users/", CreateUser).Tags(tagUserAndAuthentication).Status(http.StatusCreated).AnotherError(bc, 422, &GenericError{}, "")
	define.Get(bc, "/user", GetCurrentUser).Tags(tagUserAndAuthentication).AnotherError(bc, 401, &UnauthorizedError{}, "").AnotherError(bc, 422, &GenericError{}, "")
	define.Put(bc, "/user", UpdateCurrentUser).Tags(tagUserAndAuthentication).AnotherError(bc, 401, &UnauthorizedError{}, "").AnotherError(bc, 422, &GenericError{}, "")

	define.Get(bc, "/profiles/{username}", GetProfileByUsername).Tags(tagProfile).AnotherError(bc, 401, &UnauthorizedError{}, "").AnotherError(bc, 422, &GenericError{}, "")
	define.Post(bc, "/profiles/{username}/follow", FollowUserByUsername).Tags(tagProfile).AnotherError(bc, 401, &UnauthorizedError{}, "").AnotherError(bc, 422, &GenericError{}, "")
	define.Delete(bc, "/profiles/{username}/follow", UnfollowUserByUsername).Tags(tagProfile).AnotherError(bc, 401, &UnauthorizedError{}, "").AnotherError(bc, 422, &GenericError{}, "")

	define.Get(bc, "/articles/feed", GetArticlesFeed).Tags(tagArticles).AnotherError(bc, 401, &UnauthorizedError{}, "").AnotherError(bc, 422, &GenericError{}, "")
	define.Get(bc, "/articles", GetArticles).Tags(tagArticles).AnotherError(bc, 401, &UnauthorizedError{}, "").AnotherError(bc, 422, &GenericError{}, "")
	define.Post(bc, "/articles", CreateArticle).Tags(tagArticles).AnotherError(bc, 401, &UnauthorizedError{}, "").AnotherError(bc, 422, &GenericError{}, "")
	define.Get(bc, "/articles/{slug}", GetArticle).Tags(tagArticles).AnotherError(bc, 422, &GenericError{}, "")
	define.Put(bc, "/articles/{slug}", UpdateArticle).Tags(tagArticles).AnotherError(bc, 401, &UnauthorizedError{}, "").AnotherError(bc, 422, &GenericError{}, "")
	define.Delete(bc, "/articles/{slug}", DeleteArticle).Tags(tagArticles).AnotherError(bc, 401, &UnauthorizedError{}, "").AnotherError(bc, 422, &GenericError{}, "")
	define.Get(bc, "/articles/{slug}/comments", GetArticleComments).Tags(tagComments).AnotherError(bc, 401, &UnauthorizedError{}, "").AnotherError(bc, 422, &GenericError{}, "")
	define.Post(bc, "/articles/{slug}/comments", CreateArticleComment).Tags(tagComments).AnotherError(bc, 401, &UnauthorizedError{}, "").AnotherError(bc, 422, &GenericError{}, "")
	define.Delete(bc, "/articles/{slug}/comments/{id}", DeleteArticleComment).Tags(tagComments).AnotherError(bc, 401, &UnauthorizedError{}, "").AnotherError(bc, 422, &GenericError{}, "")
	define.Post(bc, "/articles/{slug}/favorite", CreateArticleFavorite).Tags(tagFavorites).AnotherError(bc, 401, &UnauthorizedError{}, "").AnotherError(bc, 422, &GenericError{}, "")
	define.Delete(bc, "/articles/{slug}/favorite", DeleteArticleFavorite).Tags(tagFavorites).AnotherError(bc, 401, &UnauthorizedError{}, "").AnotherError(bc, 422, &GenericError{}, "")

	define.Get(bc, "/tags", GetTags).Tags(tagTags).AnotherError(bc, 422, &GenericError{}, "")
}

// components

type User struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}

type Profile struct {
	Bio       string `json:"bio"`
	Following bool   `json:"following"`
	Image     string `json:"image"`
	Username  string `json:"username"`
}

type Article struct {
	Slug           string    `json:"slug"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Body           string    `json:"body"`
	TagList        []string  `json:"tagList"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Favorited      bool      `json:"favorited"`
	FavoritesCount int       `json:"favoritesCount"`
	Author         Profile   `json:"author"`
}

type Comment struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Body      string    `json:"body"`
	Author    Profile   `json:"author"`
}

// Unauthorized
type UnauthorizedError struct {
}

// Unexpected error
type GenericError struct {
	Errors GenericErrorErrors `json:"errors"`
}
type GenericErrorErrors struct {
	Body []string `json:"body"`
}

// handlers

type LoginInput struct {
	// Credentials to use
	User struct {
		Email    string `json:"email"`
		Password string `json:"password" openapi-override:"{'format': 'password'}"`
	} `json:"user"`
}

type LoginOutput struct {
	User User `json:"user"`
}

// Existing user login
//
// Login for existing user
func Login(ctx context.Context, input LoginInput) (output LoginOutput, err error) {
	return LoginOutput{}, nil
}

type CreateUserInput struct {
	// Details of the new user to register
	User struct {
		Email    string `json:"email"`
		Password string `json:"password" openapi-override:"{'format': 'password'}"`
		Username string `json:"username"`
	} `json:"user"`
}

type CreateUserOutput struct {
	User User `json:"user"`
}

// Register a new user
func CreateUser(ctx context.Context, input CreateUserInput) (output CreateUserOutput, err error) {
	return CreateUserOutput{}, nil
}

type GetCurrentUserOutput struct {
	User User `json:"user"`
}

// Get current user
//
// Gets the currently logged-in user
func GetCurrentUser(ctx context.Context, input struct{}) (output GetCurrentUserOutput, err error) {
	return GetCurrentUserOutput{}, nil
}

type UpdateCurrentUserInput struct {
	// User details to update. At least **one** field is required
	User struct {
		Email    string `json:"email,omitempty"`
		Password string `json:"password,omitempty" openapi-override:"{'format': 'password'}"`
		Username string `json:"username,omitempty"`
		Bio      string `json:"bio,omitempty"`
		Image    string `json:"image,omitempty"`
	} `json:"user"`
}

type UpdateCurrentUserOutput struct {
	User User `json:"user"`
}

// Update current user
//
// Update user information for current user
func UpdateCurrentUser(ctx context.Context, input UpdateCurrentUserInput) (output UpdateCurrentUserOutput, err error) {
	return UpdateCurrentUserOutput{}, nil
}

type GetProfileByUsernameInput struct {
	Username string `in:"path" path:"username"` // username of the profile to get
}

type GetProfileByUsernameOutput struct {
	Profile Profile `json:"profile"`
}

// Get a profile
//
// Get a profile of a user of the system. Auth is optional
func GetProfileByUsername(ctx context.Context, input GetProfileByUsernameInput) (output GetProfileByUsernameOutput, err error) {
	return GetProfileByUsernameOutput{}, nil
}

type FollowUserByUsernameInput struct {
	Username string `in:"path" path:"username"` // username of the profile you want to follow
}

type FollowUserByUsernameOutput struct {
	Profile Profile `json:"profile"`
}

// Follow a user
//
// Follow a user by username
func FollowUserByUsername(ctx context.Context, input FollowUserByUsernameInput) (output FollowUserByUsernameOutput, err error) {
	return FollowUserByUsernameOutput{}, nil
}

type UnfollowUserByUsernameInput struct {
	Username string `in:"path" path:"username"` // username of the profile you want to unfollow
}

type UnfollowUserByUsernameOutput struct {
	Profile Profile `json:"profile"`
}

// Unfollow a user
//
// Unfollow a user by username
func UnfollowUserByUsername(ctx context.Context, input UnfollowUserByUsernameInput) (output UnfollowUserByUsernameOutput, err error) {
	return UnfollowUserByUsernameOutput{}, nil
}

// The numbers of items to return.
type LimitParam int

// The number of items to skip before starting to collect the result set.
type OffsetParam int

type GetArticlesFeedInput struct {
	Limit  LimitParam  `in:"query" query:"limit"`
	Offset OffsetParam `in:"query" query:"offset"`
}

type GetArticlesFeedOutput struct {
	Articles      []Article `json:"articles"`
	ArticlesCount int       `json:"articlesCount"`
}

// Get recent articles from users you follow
//
// Get most recent articles from users you follow. Use query parameters to limit. Auth is required
func GetArticlesFeed(ctx context.Context, input GetArticlesFeedInput) (output GetArticlesFeedOutput, err error) {
	return GetArticlesFeedOutput{}, nil
}

type GetArticlesInput struct {
	Tag       string `in:"query" query:"tag"`       // Filter by tag
	Author    string `in:"query" query:"author"`    // Filter by author (username)
	Favorited string `in:"query" query:"favorited"` // Filter by favorites of a user (username)

	Limit  LimitParam  `in:"query" query:"limit"`
	Offset OffsetParam `in:"query" query:"offset"`
}

type GetArticlesOutput struct {
	Articles      []Article `json:"articles"`
	ArticlesCount int       `json:"articlesCount"`
}

// Get recent articles globally
//
// Get most recent articles globally. Use query parameters to filter results. Auth is optional
func GetArticles(ctx context.Context, input GetArticlesInput) (output GetArticlesOutput, err error) {
	return GetArticlesOutput{}, nil
}

type CreateArticleInput struct {
	// Article to create
	Article struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Body        string   `json:"body"`
		TagList     []string `json:"tagList,omitempty"`
	} `json:"article"`
}

type CreateArticleOutput struct {
	Article Article `json:"article"`
}

// Create an article
//
// Create an article. Auth is required
func CreateArticle(ctx context.Context, input CreateArticleInput) (output CreateArticleOutput, err error) {
	return CreateArticleOutput{}, nil
}

type GetArticleInput struct {
	Slug string `in:"path" path:"slug"` // Slug of the article to get
}

type GetArticleOutput struct {
	Article Article `json:"article"`
}

// Get an article
//
// Get an article. Auth not required
func GetArticle(ctx context.Context, input GetArticleInput) (output GetArticleOutput, err error) {
	return GetArticleOutput{}, nil
}

type UpdateArticleInput struct {
	Slug string `in:"path" path:"slug"` // Slug of the article to update

	// Article to update
	Article struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Body        string `json:"body"`
	} `json:"article"`
}

type UpdateArticleOutput struct {
	Article Article `json:"article"`
}

// Update an article
//
// Update an article. Auth is required
func UpdateArticle(ctx context.Context, input UpdateArticleInput) (output UpdateArticleOutput, err error) {
	return UpdateArticleOutput{}, nil
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

type GetArticleCommentsOutput struct {
	Comments []Comment `json:"Comments"`
}

// Get comments for an article
//
// Get the comments for an article. Auth is optional
func GetArticleComments(ctx context.Context, input GetArticleCommentsInput) (output GetArticleCommentsOutput, err error) {
	return GetArticleCommentsOutput{}, nil
}

type CreateArticleCommentInput struct {
	Slug string `in:"path" path:"slug"` // Slug of the article that you want to create a comment

	// Comment you want to create
	Comment struct {
		Body string `json:"body"`
	} `json:"comment"`
}

type CreateArticleCommentOutput struct {
	Comment Comment `json:"Comment"`
}

// Create a comment for an article
//
// Create a comment for an article. Auth is required
func CreateArticleComment(ctx context.Context, input CreateArticleCommentInput) (output CreateArticleCommentOutput, err error) {
	return CreateArticleCommentOutput{}, nil
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

type CreateArticleFavoriteOutput struct {
	Article Article `json:"article"`
}

// Favorite an article
//
// Favorite an article. Auth is required
func CreateArticleFavorite(ctx context.Context, input CreateArticleFavoriteInput) (output CreateArticleFavoriteOutput, err error) {
	return CreateArticleFavoriteOutput{}, nil
}

type DeleteArticleFavoriteInput struct {
	Slug string `in:"path" path:"slug"` // Slug of the article that you want to unfavorite
}

type DeleteArticleFavoriteOutput struct {
	Article Article `json:"article"`
}

// Unfavorite an article
//
// Unfavorite an article. Auth is required
func DeleteArticleFavorite(ctx context.Context, input DeleteArticleFavoriteInput) (output DeleteArticleFavoriteOutput, err error) {
	return DeleteArticleFavoriteOutput{}, nil
}

type GetTagsInput struct {
	Query     string `in:"query" query:"tag"`       // Filter by tag
	Author    string `in:"query" query:"author"`    // Filter by author (username)
	Favorited string `in:"query" query:"favorited"` // Filter by favorites of a user (username)

	Limit  LimitParam  `in:"query" query:"limit"`
	Offset OffsetParam `in:"query" query:"offset"`
}

// Tags
type GetTagsOutput struct {
	Tags string `json:"tags"`
}

// Get tags
//
// Get tags. Auth not required
func GetTags(ctx context.Context, input GetTagsInput) (output GetTagsOutput, err error) {
	return GetTagsOutput{}, nil
}
