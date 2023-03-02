//go:generate go run ./ -gendoc -docfile openapi.json -mdfile apidoc.md
package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/podhmo/quickapi"
	"github.com/podhmo/quickapi/qopenapi/define"
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
	bc := define.MustBuildContext(doc, router)

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
	define.Post(bc, "/users/login", Login)
	define.Post(bc, "/users/", CreateUser)
	define.Get(bc, "/user", GetCurrentUser)
	define.Put(bc, "/user", UpdateCurrentUser)

	define.Get(bc, "/profiles/{username}", GetProfileByUsername)
	define.Post(bc, "/profiles/{username}/follow", FollowUserByUsername)
	define.Delete(bc, "/profiles/{username}/follow", UnfollowUserByUsername)

	define.Get(bc, "/articles/feed", GetArticlesFeed)
	define.Get(bc, "/articles", GetArticles)
	define.Post(bc, "/articles", CreateArticle)
	define.Get(bc, "/articles/{slug}", GetArticle)
	define.Put(bc, "/articles/{slug}", UpdateArticle)
	define.Delete(bc, "/articles/{slug}", DeleteArticle)
	define.Get(bc, "/articles/{slug}/comments", GetArticleComments)
	define.Post(bc, "/articles/{slug}/comments", CreateArticleComment)
	define.Delete(bc, "/articles/{slug}/comments/{id}", DeleteArticleComment)
	define.Post(bc, "/articles/{slug}/favorite", CreateArticleFavorite)
	define.Delete(bc, "/articles/{slug}/favorite", DeleteArticleFavorite)

	define.Get(bc, "/tags", GetTags)
}

// handlers
func Login(ctx context.Context, input struct{}) (output struct{}, err error) {
	return struct{}{}, nil
}
func CreateUser(ctx context.Context, input struct{}) (output struct{}, err error) {
	return struct{}{}, nil
}
func GetCurrentUser(ctx context.Context, input struct{}) (output struct{}, err error) {
	return struct{}{}, nil
}
func UpdateCurrentUser(ctx context.Context, input struct{}) (output struct{}, err error) {
	return struct{}{}, nil
}

func GetProfileByUsername(ctx context.Context, input struct{}) (output struct{}, err error) {
	return struct{}{}, nil
}
func FollowUserByUsername(ctx context.Context, input struct{}) (output struct{}, err error) {
	return struct{}{}, nil
}
func UnfollowUserByUsername(ctx context.Context, input struct{}) (output struct{}, err error) {
	return struct{}{}, nil
}

func GetArticlesFeed(ctx context.Context, input struct{}) (output struct{}, err error) {
	return struct{}{}, nil
}
func GetArticles(ctx context.Context, input struct{}) (output struct{}, err error) {
	return struct{}{}, nil
}
func CreateArticle(ctx context.Context, input struct{}) (output struct{}, err error) {
	return struct{}{}, nil
}
func GetArticle(ctx context.Context, input struct{}) (output struct{}, err error) {
	return struct{}{}, nil
}
func UpdateArticle(ctx context.Context, input struct{}) (output struct{}, err error) {
	return struct{}{}, nil
}
func DeleteArticle(ctx context.Context, input struct{}) (output struct{}, err error) {
	return struct{}{}, nil
}
func GetArticleComments(ctx context.Context, input struct{}) (output struct{}, err error) {
	return struct{}{}, nil
}
func CreateArticleComment(ctx context.Context, input struct{}) (output struct{}, err error) {
	return struct{}{}, nil
}
func DeleteArticleComment(ctx context.Context, input struct{}) (output struct{}, err error) {
	return struct{}{}, nil
}
func CreateArticleFavorite(ctx context.Context, input struct{}) (output struct{}, err error) {
	return struct{}{}, nil
}
func DeleteArticleFavorite(ctx context.Context, input struct{}) (output struct{}, err error) {
	return struct{}{}, nil
}
func GetTags(ctx context.Context, input struct{}) (output struct{}, err error) {
	return struct{}{}, nil
}
