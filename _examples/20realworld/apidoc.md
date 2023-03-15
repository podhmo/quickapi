---
title: Conduit API
version: 1.0.0
---

# Conduit API

Conduit API documentation

- [paths](#paths)
- [schemas](#schemas)

## paths

| endpoint | operationId | tags | summary |
| --- | --- | --- | --- |
| `GET /articles` | [main.GetArticles](#maingetarticles-get-articles)  | `Articles` | Get recent articles globally |
| `POST /articles` | [main.CreateArticle](#maincreatearticle-post-articles)  | `Articles` | Create an article |
| `GET /articles/feed` | [main.GetArticlesFeed](#maingetarticlesfeed-get-articlesfeed)  | `Articles` | Get recent articles from users you follow |
| `DELETE /articles/{slug}` | [main.DeleteArticle](#maindeletearticle-delete-articlesslug)  | `Articles` | Delete an article |
| `GET /articles/{slug}` | [main.GetArticle](#maingetarticle-get-articlesslug)  | `Articles` | Get an article |
| `PUT /articles/{slug}` | [main.UpdateArticle](#mainupdatearticle-put-articlesslug)  | `Articles` | Update an article |
| `GET /articles/{slug}/comments` | [main.GetArticleComments](#maingetarticlecomments-get-articlesslugcomments)  | `Comments` | Get comments for an article |
| `POST /articles/{slug}/comments` | [main.CreateArticleComment](#maincreatearticlecomment-post-articlesslugcomments)  | `Comments` | Create a comment for an article |
| `DELETE /articles/{slug}/comments/{id}` | [main.DeleteArticleComment](#maindeletearticlecomment-delete-articlesslugcommentsid)  | `Comments` | Delete a comment for an article |
| `DELETE /articles/{slug}/favorite` | [main.DeleteArticleFavorite](#maindeletearticlefavorite-delete-articlesslugfavorite)  | `Favorites` | Unfavorite an article |
| `POST /articles/{slug}/favorite` | [main.CreateArticleFavorite](#maincreatearticlefavorite-post-articlesslugfavorite)  | `Favorites` | Favorite an article |
| `GET /profiles/{username}` | [main.GetProfileByUsername](#maingetprofilebyusername-get-profilesusername)  | `Profile` | Get a profile |
| `DELETE /profiles/{username}/follow` | [main.UnfollowUserByUsername](#mainunfollowuserbyusername-delete-profilesusernamefollow)  | `Profile` | Unfollow a user |
| `POST /profiles/{username}/follow` | [main.FollowUserByUsername](#mainfollowuserbyusername-post-profilesusernamefollow)  | `Profile` | Follow a user |
| `GET /tags` | [main.GetTags](#maingettags-get-tags)  | `Tags` | Get tags |
| `GET /user` | [main.GetCurrentUser](#maingetcurrentuser-get-user)  | `User and Authentication` | Get current user |
| `PUT /user` | [main.UpdateCurrentUser](#mainupdatecurrentuser-put-user)  | `User and Authentication` | Update current user |
| `POST /users/` | [main.CreateUser](#maincreateuser-post-users)  | `User and Authentication` | Register a new user |
| `POST /users/login` | [main.Login](#mainlogin-post-userslogin)  | `User and Authentication` | Existing user login |


### main.GetArticles `GET /articles`

Get recent articles globally

| name | value |
| --- | --- |
| operationId | main.GetArticles[  <sub>(source)</sub>](https://github.com/podhmo/quickapi/blob/main/_examples/20realworld/main.go#L352) |
| endpoint | `GET /articles` |
| input | Input |
| output | [`GetArticlesOutput[[]Article]`](#article) ｜ [`UnauthorizedError`](#unauthorizederror) ｜ [`GenericError`](#genericerror) ｜ [`ErrorResponse`](#errorresponse) |
| tags | `Articles` |


#### input (application/json)

```go
// GET /articles
type Input struct {
	// Filter by tag
	tag? string `in:"query"`

	// Filter by author (username)
	author? string `in:"query"`

	// Filter by favorites of a user (username)
	favorited? string `in:"query"`

	limit? integer `in:"query"`

	offset? integer `in:"query"`
}
```

#### output (application/json)

```go
// GET /articles (200)
type Output200 struct {	// GetArticlesOutput
	articles []struct {	// Article
		slug string

		title string

		description string

		body string

		tagList []string

		createdAt string `format:"date-time"`

		updatedAt string `format:"date-time"`

		favorited boolean	// default: false

		favoritesCount integer

		author struct {	// Profile
			bio string

			following boolean	// default: false

			image string

			username string
		}
	}

	articlesCount integer
}

// GET /articles (401)
type Output401 struct {	// UnauthorizedError
}

// GET /articles (422)
type Output422 struct {	// GenericError
	errors struct {	// GenericErrorErrors
		body []string
	}
}

// GET /articles (default)
// default error
type OutputDefault struct {	// ErrorResponse
	code integer

	error string

	detail? []string
}
```

#### description

Get most recent articles globally. Use query parameters to filter results. Auth is optional
### main.CreateArticle `POST /articles`

Create an article

| name | value |
| --- | --- |
| operationId | main.CreateArticle[  <sub>(source)</sub>](https://github.com/podhmo/quickapi/blob/main/_examples/20realworld/main.go#L373) |
| endpoint | `POST /articles` |
| input | Input[ [`CreateArticleInput[[[]]]`](#) ] |
| output | [`CreateArticleOutput[Article]`](#article) ｜ [`UnauthorizedError`](#unauthorizederror) ｜ [`GenericError`](#genericerror) ｜ [`ErrorResponse`](#errorresponse) |
| tags | `Articles` |


#### input (application/json)

```go
// POST /articles
type Input struct {
	JSONBody struct {	// CreateArticleInput
		article struct {	// 
			title string

			description string

			body string

			tagList? []string
		}
	}
}
```

#### output (application/json)

```go
// POST /articles (200)
type Output200 struct {	// CreateArticleOutput
	article struct {	// Article
		slug string

		title string

		description string

		body string

		tagList []string

		createdAt string `format:"date-time"`

		updatedAt string `format:"date-time"`

		favorited boolean	// default: false

		favoritesCount integer

		author struct {	// Profile
			bio string

			following boolean	// default: false

			image string

			username string
		}
	}
}

// POST /articles (401)
type Output401 struct {	// UnauthorizedError
}

// POST /articles (422)
type Output422 struct {	// GenericError
	errors struct {	// GenericErrorErrors
		body []string
	}
}

// POST /articles (default)
// default error
type OutputDefault struct {	// ErrorResponse
	code integer

	error string

	detail? []string
}
```

#### description

Create an article. Auth is required
### main.GetArticlesFeed `GET /articles/feed`

Get recent articles from users you follow

| name | value |
| --- | --- |
| operationId | main.GetArticlesFeed[  <sub>(source)</sub>](https://github.com/podhmo/quickapi/blob/main/_examples/20realworld/main.go#L331) |
| endpoint | `GET /articles/feed` |
| input | Input |
| output | [`GetArticlesFeedOutput[[]Article]`](#article) ｜ [`UnauthorizedError`](#unauthorizederror) ｜ [`GenericError`](#genericerror) ｜ [`ErrorResponse`](#errorresponse) |
| tags | `Articles` |


#### input (application/json)

```go
// GET /articles/feed
type Input struct {
	limit? integer `in:"query"`

	offset? integer `in:"query"`
}
```

#### output (application/json)

```go
// GET /articles/feed (200)
type Output200 struct {	// GetArticlesFeedOutput
	articles []struct {	// Article
		slug string

		title string

		description string

		body string

		tagList []string

		createdAt string `format:"date-time"`

		updatedAt string `format:"date-time"`

		favorited boolean	// default: false

		favoritesCount integer

		author struct {	// Profile
			bio string

			following boolean	// default: false

			image string

			username string
		}
	}

	articlesCount integer
}

// GET /articles/feed (401)
type Output401 struct {	// UnauthorizedError
}

// GET /articles/feed (422)
type Output422 struct {	// GenericError
	errors struct {	// GenericErrorErrors
		body []string
	}
}

// GET /articles/feed (default)
// default error
type OutputDefault struct {	// ErrorResponse
	code integer

	error string

	detail? []string
}
```

#### description

Get most recent articles from users you follow. Use query parameters to limit. Auth is required
### main.DeleteArticle `DELETE /articles/{slug}`

Delete an article

| name | value |
| --- | --- |
| operationId | main.DeleteArticle[  <sub>(source)</sub>](https://github.com/podhmo/quickapi/blob/main/_examples/20realworld/main.go#L421) |
| endpoint | `DELETE /articles/{slug}` |
| input | Input |
| output | `<Anonymous>` ｜ [`UnauthorizedError`](#unauthorizederror) ｜ [`GenericError`](#genericerror) ｜ [`ErrorResponse`](#errorresponse) |
| tags | `Articles` |


#### input (application/json)

```go
// DELETE /articles/{slug}
type Input struct {
	// Slug of the article to delete
	slug string `in:"path"`
}
```

#### output (application/json)

```go
// DELETE /articles/{slug} (200)
type Output200 struct {	// 
}

// DELETE /articles/{slug} (401)
type Output401 struct {	// UnauthorizedError
}

// DELETE /articles/{slug} (422)
type Output422 struct {	// GenericError
	errors struct {	// GenericErrorErrors
		body []string
	}
}

// DELETE /articles/{slug} (default)
// default error
type OutputDefault struct {	// ErrorResponse
	code integer

	error string

	detail? []string
}
```

#### description

Delete an article. Auth is required
### main.GetArticle `GET /articles/{slug}`

Get an article

| name | value |
| --- | --- |
| operationId | main.GetArticle[  <sub>(source)</sub>](https://github.com/podhmo/quickapi/blob/main/_examples/20realworld/main.go#L388) |
| endpoint | `GET /articles/{slug}` |
| input | Input |
| output | [`GetArticleOutput[Article]`](#article) ｜ [`GenericError`](#genericerror) ｜ [`ErrorResponse`](#errorresponse) |
| tags | `Articles` |


#### input (application/json)

```go
// GET /articles/{slug}
type Input struct {
	// Slug of the article to get
	slug string `in:"path"`
}
```

#### output (application/json)

```go
// GET /articles/{slug} (200)
type Output200 struct {	// GetArticleOutput
	article struct {	// Article
		slug string

		title string

		description string

		body string

		tagList []string

		createdAt string `format:"date-time"`

		updatedAt string `format:"date-time"`

		favorited boolean	// default: false

		favoritesCount integer

		author struct {	// Profile
			bio string

			following boolean	// default: false

			image string

			username string
		}
	}
}

// GET /articles/{slug} (422)
type Output422 struct {	// GenericError
	errors struct {	// GenericErrorErrors
		body []string
	}
}

// GET /articles/{slug} (default)
// default error
type OutputDefault struct {	// ErrorResponse
	code integer

	error string

	detail? []string
}
```

#### description

Get an article. Auth not required
### main.UpdateArticle `PUT /articles/{slug}`

Update an article

| name | value |
| --- | --- |
| operationId | main.UpdateArticle[  <sub>(source)</sub>](https://github.com/podhmo/quickapi/blob/main/_examples/20realworld/main.go#L410) |
| endpoint | `PUT /articles/{slug}` |
| input | Input[ [`UpdateArticleInput[]`](#) ] |
| output | [`UpdateArticleOutput[Article]`](#article) ｜ [`UnauthorizedError`](#unauthorizederror) ｜ [`GenericError`](#genericerror) ｜ [`ErrorResponse`](#errorresponse) |
| tags | `Articles` |


#### input (application/json)

```go
// PUT /articles/{slug}
type Input struct {
	// Slug of the article to update
	slug string `in:"path"`

	JSONBody struct {	// UpdateArticleInput
		article struct {	// 
			title string

			description string

			body string
		}
	}
}
```

#### output (application/json)

```go
// PUT /articles/{slug} (200)
type Output200 struct {	// UpdateArticleOutput
	article struct {	// Article
		slug string

		title string

		description string

		body string

		tagList []string

		createdAt string `format:"date-time"`

		updatedAt string `format:"date-time"`

		favorited boolean	// default: false

		favoritesCount integer

		author struct {	// Profile
			bio string

			following boolean	// default: false

			image string

			username string
		}
	}
}

// PUT /articles/{slug} (401)
type Output401 struct {	// UnauthorizedError
}

// PUT /articles/{slug} (422)
type Output422 struct {	// GenericError
	errors struct {	// GenericErrorErrors
		body []string
	}
}

// PUT /articles/{slug} (default)
// default error
type OutputDefault struct {	// ErrorResponse
	code integer

	error string

	detail? []string
}
```

#### description

Update an article. Auth is required
### main.GetArticleComments `GET /articles/{slug}/comments`

Get comments for an article

| name | value |
| --- | --- |
| operationId | main.GetArticleComments[  <sub>(source)</sub>](https://github.com/podhmo/quickapi/blob/main/_examples/20realworld/main.go#L436) |
| endpoint | `GET /articles/{slug}/comments` |
| input | Input |
| output | [`GetArticleCommentsOutput[[]Comment]`](#comment) ｜ [`UnauthorizedError`](#unauthorizederror) ｜ [`GenericError`](#genericerror) ｜ [`ErrorResponse`](#errorresponse) |
| tags | `Comments` |


#### input (application/json)

```go
// GET /articles/{slug}/comments
type Input struct {
	// Slug of the article that you want to get comments for
	slug string `in:"path"`
}
```

#### output (application/json)

```go
// GET /articles/{slug}/comments (200)
type Output200 struct {	// GetArticleCommentsOutput
	Comments []struct {	// Comment
		id integer

		createdAt string `format:"date-time"`

		updatedAt string `format:"date-time"`

		body string

		author struct {	// Profile
			bio string

			following boolean	// default: false

			image string

			username string
		}
	}
}

// GET /articles/{slug}/comments (401)
type Output401 struct {	// UnauthorizedError
}

// GET /articles/{slug}/comments (422)
type Output422 struct {	// GenericError
	errors struct {	// GenericErrorErrors
		body []string
	}
}

// GET /articles/{slug}/comments (default)
// default error
type OutputDefault struct {	// ErrorResponse
	code integer

	error string

	detail? []string
}
```

#### description

Get the comments for an article. Auth is optional
### main.CreateArticleComment `POST /articles/{slug}/comments`

Create a comment for an article

| name | value |
| --- | --- |
| operationId | main.CreateArticleComment[  <sub>(source)</sub>](https://github.com/podhmo/quickapi/blob/main/_examples/20realworld/main.go#L456) |
| endpoint | `POST /articles/{slug}/comments` |
| input | Input[ [`CreateArticleCommentInput[]`](#) ] |
| output | [`CreateArticleCommentOutput[Comment]`](#comment) ｜ [`UnauthorizedError`](#unauthorizederror) ｜ [`GenericError`](#genericerror) ｜ [`ErrorResponse`](#errorresponse) |
| tags | `Comments` |


#### input (application/json)

```go
// POST /articles/{slug}/comments
type Input struct {
	// Slug of the article that you want to create a comment
	slug string `in:"path"`

	JSONBody struct {	// CreateArticleCommentInput
		comment struct {	// 
			body string
		}
	}
}
```

#### output (application/json)

```go
// POST /articles/{slug}/comments (200)
type Output200 struct {	// CreateArticleCommentOutput
	Comment struct {	// Comment
		id integer

		createdAt string `format:"date-time"`

		updatedAt string `format:"date-time"`

		body string

		author struct {	// Profile
			bio string

			following boolean	// default: false

			image string

			username string
		}
	}
}

// POST /articles/{slug}/comments (401)
type Output401 struct {	// UnauthorizedError
}

// POST /articles/{slug}/comments (422)
type Output422 struct {	// GenericError
	errors struct {	// GenericErrorErrors
		body []string
	}
}

// POST /articles/{slug}/comments (default)
// default error
type OutputDefault struct {	// ErrorResponse
	code integer

	error string

	detail? []string
}
```

#### description

Create a comment for an article. Auth is required
### main.DeleteArticleComment `DELETE /articles/{slug}/comments/{id}`

Delete a comment for an article

| name | value |
| --- | --- |
| operationId | main.DeleteArticleComment[  <sub>(source)</sub>](https://github.com/podhmo/quickapi/blob/main/_examples/20realworld/main.go#L468) |
| endpoint | `DELETE /articles/{slug}/comments/{id}` |
| input | Input |
| output | `<Anonymous>` ｜ [`UnauthorizedError`](#unauthorizederror) ｜ [`GenericError`](#genericerror) ｜ [`ErrorResponse`](#errorresponse) |
| tags | `Comments` |


#### input (application/json)

```go
// DELETE /articles/{slug}/comments/{id}
type Input struct {
	// Slug of the article that you want to delete a comment
	slug string `in:"path"`

	// ID of the comment you want to delete
	id string `in:"path"`
}
```

#### output (application/json)

```go
// DELETE /articles/{slug}/comments/{id} (200)
type Output200 struct {	// 
}

// DELETE /articles/{slug}/comments/{id} (401)
type Output401 struct {	// UnauthorizedError
}

// DELETE /articles/{slug}/comments/{id} (422)
type Output422 struct {	// GenericError
	errors struct {	// GenericErrorErrors
		body []string
	}
}

// DELETE /articles/{slug}/comments/{id} (default)
// default error
type OutputDefault struct {	// ErrorResponse
	code integer

	error string

	detail? []string
}
```

#### description

Delete a comment for an article. Auth is required.
### main.DeleteArticleFavorite `DELETE /articles/{slug}/favorite`

Unfavorite an article

| name | value |
| --- | --- |
| operationId | main.DeleteArticleFavorite[  <sub>(source)</sub>](https://github.com/podhmo/quickapi/blob/main/_examples/20realworld/main.go#L498) |
| endpoint | `DELETE /articles/{slug}/favorite` |
| input | Input |
| output | [`DeleteArticleFavoriteOutput[Article]`](#article) ｜ [`UnauthorizedError`](#unauthorizederror) ｜ [`GenericError`](#genericerror) ｜ [`ErrorResponse`](#errorresponse) |
| tags | `Favorites` |


#### input (application/json)

```go
// DELETE /articles/{slug}/favorite
type Input struct {
	// Slug of the article that you want to unfavorite
	slug string `in:"path"`
}
```

#### output (application/json)

```go
// DELETE /articles/{slug}/favorite (200)
type Output200 struct {	// DeleteArticleFavoriteOutput
	article struct {	// Article
		slug string

		title string

		description string

		body string

		tagList []string

		createdAt string `format:"date-time"`

		updatedAt string `format:"date-time"`

		favorited boolean	// default: false

		favoritesCount integer

		author struct {	// Profile
			bio string

			following boolean	// default: false

			image string

			username string
		}
	}
}

// DELETE /articles/{slug}/favorite (401)
type Output401 struct {	// UnauthorizedError
}

// DELETE /articles/{slug}/favorite (422)
type Output422 struct {	// GenericError
	errors struct {	// GenericErrorErrors
		body []string
	}
}

// DELETE /articles/{slug}/favorite (default)
// default error
type OutputDefault struct {	// ErrorResponse
	code integer

	error string

	detail? []string
}
```

#### description

Unfavorite an article. Auth is required
### main.CreateArticleFavorite `POST /articles/{slug}/favorite`

Favorite an article

| name | value |
| --- | --- |
| operationId | main.CreateArticleFavorite[  <sub>(source)</sub>](https://github.com/podhmo/quickapi/blob/main/_examples/20realworld/main.go#L483) |
| endpoint | `POST /articles/{slug}/favorite` |
| input | Input |
| output | [`CreateArticleFavoriteOutput[Article]`](#article) ｜ [`UnauthorizedError`](#unauthorizederror) ｜ [`GenericError`](#genericerror) ｜ [`ErrorResponse`](#errorresponse) |
| tags | `Favorites` |


#### input (application/json)

```go
// POST /articles/{slug}/favorite
type Input struct {
	// Slug of the article that you want to favorite
	slug string `in:"path"`
}
```

#### output (application/json)

```go
// POST /articles/{slug}/favorite (200)
type Output200 struct {	// CreateArticleFavoriteOutput
	article struct {	// Article
		slug string

		title string

		description string

		body string

		tagList []string

		createdAt string `format:"date-time"`

		updatedAt string `format:"date-time"`

		favorited boolean	// default: false

		favoritesCount integer

		author struct {	// Profile
			bio string

			following boolean	// default: false

			image string

			username string
		}
	}
}

// POST /articles/{slug}/favorite (401)
type Output401 struct {	// UnauthorizedError
}

// POST /articles/{slug}/favorite (422)
type Output422 struct {	// GenericError
	errors struct {	// GenericErrorErrors
		body []string
	}
}

// POST /articles/{slug}/favorite (default)
// default error
type OutputDefault struct {	// ErrorResponse
	code integer

	error string

	detail? []string
}
```

#### description

Favorite an article. Auth is required
### main.GetProfileByUsername `GET /profiles/{username}`

Get a profile

| name | value |
| --- | --- |
| operationId | main.GetProfileByUsername[  <sub>(source)</sub>](https://github.com/podhmo/quickapi/blob/main/_examples/20realworld/main.go#L278) |
| endpoint | `GET /profiles/{username}` |
| input | Input |
| output | [`GetProfileByUsernameOutput[Profile]`](#profile) ｜ [`UnauthorizedError`](#unauthorizederror) ｜ [`GenericError`](#genericerror) ｜ [`ErrorResponse`](#errorresponse) |
| tags | `Profile` |


#### input (application/json)

```go
// GET /profiles/{username}
type Input struct {
	// username of the profile to get
	username string `in:"path"`
}
```

#### output (application/json)

```go
// GET /profiles/{username} (200)
type Output200 struct {	// GetProfileByUsernameOutput
	profile struct {	// Profile
		bio string

		following boolean	// default: false

		image string

		username string
	}
}

// GET /profiles/{username} (401)
type Output401 struct {	// UnauthorizedError
}

// GET /profiles/{username} (422)
type Output422 struct {	// GenericError
	errors struct {	// GenericErrorErrors
		body []string
	}
}

// GET /profiles/{username} (default)
// default error
type OutputDefault struct {	// ErrorResponse
	code integer

	error string

	detail? []string
}
```

#### description

Get a profile of a user of the system. Auth is optional
### main.UnfollowUserByUsername `DELETE /profiles/{username}/follow`

Unfollow a user

| name | value |
| --- | --- |
| operationId | main.UnfollowUserByUsername[  <sub>(source)</sub>](https://github.com/podhmo/quickapi/blob/main/_examples/20realworld/main.go#L308) |
| endpoint | `DELETE /profiles/{username}/follow` |
| input | Input |
| output | [`UnfollowUserByUsernameOutput[Profile]`](#profile) ｜ [`UnauthorizedError`](#unauthorizederror) ｜ [`GenericError`](#genericerror) ｜ [`ErrorResponse`](#errorresponse) |
| tags | `Profile` |


#### input (application/json)

```go
// DELETE /profiles/{username}/follow
type Input struct {
	// username of the profile you want to unfollow
	username string `in:"path"`
}
```

#### output (application/json)

```go
// DELETE /profiles/{username}/follow (200)
type Output200 struct {	// UnfollowUserByUsernameOutput
	profile struct {	// Profile
		bio string

		following boolean	// default: false

		image string

		username string
	}
}

// DELETE /profiles/{username}/follow (401)
type Output401 struct {	// UnauthorizedError
}

// DELETE /profiles/{username}/follow (422)
type Output422 struct {	// GenericError
	errors struct {	// GenericErrorErrors
		body []string
	}
}

// DELETE /profiles/{username}/follow (default)
// default error
type OutputDefault struct {	// ErrorResponse
	code integer

	error string

	detail? []string
}
```

#### description

Unfollow a user by username
### main.FollowUserByUsername `POST /profiles/{username}/follow`

Follow a user

| name | value |
| --- | --- |
| operationId | main.FollowUserByUsername[  <sub>(source)</sub>](https://github.com/podhmo/quickapi/blob/main/_examples/20realworld/main.go#L293) |
| endpoint | `POST /profiles/{username}/follow` |
| input | Input |
| output | [`FollowUserByUsernameOutput[Profile]`](#profile) ｜ [`UnauthorizedError`](#unauthorizederror) ｜ [`GenericError`](#genericerror) ｜ [`ErrorResponse`](#errorresponse) |
| tags | `Profile` |


#### input (application/json)

```go
// POST /profiles/{username}/follow
type Input struct {
	// username of the profile you want to follow
	username string `in:"path"`
}
```

#### output (application/json)

```go
// POST /profiles/{username}/follow (200)
type Output200 struct {	// FollowUserByUsernameOutput
	profile struct {	// Profile
		bio string

		following boolean	// default: false

		image string

		username string
	}
}

// POST /profiles/{username}/follow (401)
type Output401 struct {	// UnauthorizedError
}

// POST /profiles/{username}/follow (422)
type Output422 struct {	// GenericError
	errors struct {	// GenericErrorErrors
		body []string
	}
}

// POST /profiles/{username}/follow (default)
// default error
type OutputDefault struct {	// ErrorResponse
	code integer

	error string

	detail? []string
}
```

#### description

Follow a user by username
### main.GetTags `GET /tags`

Get tags

| name | value |
| --- | --- |
| operationId | main.GetTags[  <sub>(source)</sub>](https://github.com/podhmo/quickapi/blob/main/_examples/20realworld/main.go#L519) |
| endpoint | `GET /tags` |
| input | Input |
| output | [`GetTagsOutput`](#gettagsoutput) ｜ [`GenericError`](#genericerror) ｜ [`ErrorResponse`](#errorresponse) |
| tags | `Tags` |


#### input (application/json)

```go
// GET /tags
type Input struct {
	// Filter by tag
	tag? string `in:"query"`

	// Filter by author (username)
	author? string `in:"query"`

	// Filter by favorites of a user (username)
	favorited? string `in:"query"`

	limit? integer `in:"query"`

	offset? integer `in:"query"`
}
```

#### output (application/json)

```go
// GET /tags (200)
type Output200 struct {	// GetTagsOutput
	tags string
}

// GET /tags (422)
type Output422 struct {	// GenericError
	errors struct {	// GenericErrorErrors
		body []string
	}
}

// GET /tags (default)
// default error
type OutputDefault struct {	// ErrorResponse
	code integer

	error string

	detail? []string
}
```

#### description

Get tags. Auth not required
### main.GetCurrentUser `GET /user`

Get current user

| name | value |
| --- | --- |
| operationId | main.GetCurrentUser[  <sub>(source)</sub>](https://github.com/podhmo/quickapi/blob/main/_examples/20realworld/main.go#L241) |
| endpoint | `GET /user` |
| input | Input |
| output | [`GetCurrentUserOutput[User]`](#user) ｜ [`UnauthorizedError`](#unauthorizederror) ｜ [`GenericError`](#genericerror) ｜ [`ErrorResponse`](#errorresponse) |
| tags | `User and Authentication` |



#### output (application/json)

```go
// GET /user (200)
type Output200 struct {	// GetCurrentUserOutput
	user struct {	// User
		email string

		token string

		username string

		bio string

		image string
	}
}

// GET /user (401)
type Output401 struct {	// UnauthorizedError
}

// GET /user (422)
type Output422 struct {	// GenericError
	errors struct {	// GenericErrorErrors
		body []string
	}
}

// GET /user (default)
// default error
type OutputDefault struct {	// ErrorResponse
	code integer

	error string

	detail? []string
}
```

#### description

Gets the currently logged-in user
### main.UpdateCurrentUser `PUT /user`

Update current user

| name | value |
| --- | --- |
| operationId | main.UpdateCurrentUser[  <sub>(source)</sub>](https://github.com/podhmo/quickapi/blob/main/_examples/20realworld/main.go#L263) |
| endpoint | `PUT /user` |
| input | Input[ [`UpdateCurrentUserInput[]`](#) ] |
| output | [`UpdateCurrentUserOutput[User]`](#user) ｜ [`UnauthorizedError`](#unauthorizederror) ｜ [`GenericError`](#genericerror) ｜ [`ErrorResponse`](#errorresponse) |
| tags | `User and Authentication` |


#### input (application/json)

```go
// PUT /user
type Input struct {
	JSONBody struct {	// UpdateCurrentUserInput
		user struct {	// 
			email? string

			password? string `format:"password"`

			username? string

			bio? string

			image? string
		}
	}
}
```

#### output (application/json)

```go
// PUT /user (200)
type Output200 struct {	// UpdateCurrentUserOutput
	user struct {	// User
		email string

		token string

		username string

		bio string

		image string
	}
}

// PUT /user (401)
type Output401 struct {	// UnauthorizedError
}

// PUT /user (422)
type Output422 struct {	// GenericError
	errors struct {	// GenericErrorErrors
		body []string
	}
}

// PUT /user (default)
// default error
type OutputDefault struct {	// ErrorResponse
	code integer

	error string

	detail? []string
}
```

#### description

Update user information for current user
### main.CreateUser `POST /users/`

Register a new user

| name | value |
| --- | --- |
| operationId | main.CreateUser[  <sub>(source)</sub>](https://github.com/podhmo/quickapi/blob/main/_examples/20realworld/main.go#L230) |
| endpoint | `POST /users/` |
| input | Input[ [`CreateUserInput[]`](#) ] |
| output | [`CreateUserOutput[User]`](#user) ｜ [`GenericError`](#genericerror) ｜ [`ErrorResponse`](#errorresponse) |
| tags | `User and Authentication` |


#### input (application/json)

```go
// POST /users/
type Input struct {
	JSONBody struct {	// CreateUserInput
		user struct {	// 
			email string

			password string `format:"password"`

			username string
		}
	}
}
```

#### output (application/json)

```go
// POST /users/ (201)
type Output201 struct {	// CreateUserOutput
	user struct {	// User
		email string

		token string

		username string

		bio string

		image string
	}
}

// POST /users/ (422)
type Output422 struct {	// GenericError
	errors struct {	// GenericErrorErrors
		body []string
	}
}

// POST /users/ (default)
// default error
type OutputDefault struct {	// ErrorResponse
	code integer

	error string

	detail? []string
}
```

#### description

Register a new user
### main.Login `POST /users/login`

Existing user login

| name | value |
| --- | --- |
| operationId | main.Login[  <sub>(source)</sub>](https://github.com/podhmo/quickapi/blob/main/_examples/20realworld/main.go#L212) |
| endpoint | `POST /users/login` |
| input | Input[ [`LoginInput[]`](#) ] |
| output | [`LoginOutput[User]`](#user) ｜ [`UnauthorizedError`](#unauthorizederror) ｜ [`GenericError`](#genericerror) ｜ [`ErrorResponse`](#errorresponse) |
| tags | `User and Authentication` |


#### input (application/json)

```go
// POST /users/login
type Input struct {
	JSONBody struct {	// LoginInput
		user struct {	// 
			email string

			password string `format:"password"`
		}
	}
}
```

#### output (application/json)

```go
// POST /users/login (200)
type Output200 struct {	// LoginOutput
	user struct {	// User
		email string

		token string

		username string

		bio string

		image string
	}
}

// POST /users/login (401)
type Output401 struct {	// UnauthorizedError
}

// POST /users/login (422)
type Output422 struct {	// GenericError
	errors struct {	// GenericErrorErrors
		body []string
	}
}

// POST /users/login (default)
// default error
type OutputDefault struct {	// ErrorResponse
	code integer

	error string

	detail? []string
}
```

#### description

Login for existing user





----------------------------------------

## schemas

| name | summary |
| --- | --- |
| [Article](#article) |  |
| [Comment](#comment) |  |
| [ErrorResponse](#errorresponse) | represents a normal error response type |
| [GenericError](#genericerror) | Unexpected error |
| [GenericErrorErrors](#genericerrorerrors) |  |
| [LimitParam](#limitparam) | The numbers of items to return. |
| [OffsetParam](#offsetparam) | The number of items to skip before starting to collect the result set. |
| [Profile](#profile) |  |
| [Time](#time) |  |
| [UnauthorizedError](#unauthorizederror) | Unauthorized |
| [User](#user) |  |



### Article



```go
type Article struct {
	slug string

	title string

	description string

	body string

	tagList []string

	createdAt string `format:"date-time"`

	updatedAt string `format:"date-time"`

	favorited boolean	// default: false

	favoritesCount integer

	author struct {	// Profile
		bio string

		following boolean	// default: false

		image string

		username string
	}
}
```

- [output of main.GetArticles (200) as `GetArticlesOutput[[]Article]`](#maingetarticles-get-articles)
- [output of main.CreateArticle (200) as `CreateArticleOutput[Article]`](#maincreatearticle-post-articles)
- [output of main.GetArticlesFeed (200) as `GetArticlesFeedOutput[[]Article]`](#maingetarticlesfeed-get-articlesfeed)
- [output of main.GetArticle (200) as `GetArticleOutput[Article]`](#maingetarticle-get-articlesslug)
- [output of main.UpdateArticle (200) as `UpdateArticleOutput[Article]`](#mainupdatearticle-put-articlesslug)
- [output of main.DeleteArticleFavorite (200) as `DeleteArticleFavoriteOutput[Article]`](#maindeletearticlefavorite-delete-articlesslugfavorite)
- [output of main.CreateArticleFavorite (200) as `CreateArticleFavoriteOutput[Article]`](#maincreatearticlefavorite-post-articlesslugfavorite)

### Comment



```go
type Comment struct {
	id integer

	createdAt string `format:"date-time"`

	updatedAt string `format:"date-time"`

	body string

	author struct {	// Profile
		bio string

		following boolean	// default: false

		image string

		username string
	}
}
```

- [output of main.GetArticleComments (200) as `GetArticleCommentsOutput[[]Comment]`](#maingetarticlecomments-get-articlesslugcomments)
- [output of main.CreateArticleComment (200) as `CreateArticleCommentOutput[Comment]`](#maincreatearticlecomment-post-articlesslugcomments)

### ErrorResponse

represents a normal error response type

```go
// ErrorResponse represents a normal error response type
type ErrorResponse struct {
	code integer

	error string

	detail? []string
}
```

- [output of main.GetArticles (default) as `ErrorResponse`](#maingetarticles-get-articles)
- [output of main.CreateArticle (default) as `ErrorResponse`](#maincreatearticle-post-articles)
- [output of main.GetArticlesFeed (default) as `ErrorResponse`](#maingetarticlesfeed-get-articlesfeed)
- [output of main.DeleteArticle (default) as `ErrorResponse`](#maindeletearticle-delete-articlesslug)
- [output of main.GetArticle (default) as `ErrorResponse`](#maingetarticle-get-articlesslug)
- [output of main.UpdateArticle (default) as `ErrorResponse`](#mainupdatearticle-put-articlesslug)
- [output of main.GetArticleComments (default) as `ErrorResponse`](#maingetarticlecomments-get-articlesslugcomments)
- [output of main.CreateArticleComment (default) as `ErrorResponse`](#maincreatearticlecomment-post-articlesslugcomments)
- [output of main.DeleteArticleComment (default) as `ErrorResponse`](#maindeletearticlecomment-delete-articlesslugcommentsid)
- [output of main.DeleteArticleFavorite (default) as `ErrorResponse`](#maindeletearticlefavorite-delete-articlesslugfavorite)
- [output of main.CreateArticleFavorite (default) as `ErrorResponse`](#maincreatearticlefavorite-post-articlesslugfavorite)
- [output of main.GetProfileByUsername (default) as `ErrorResponse`](#maingetprofilebyusername-get-profilesusername)
- [output of main.UnfollowUserByUsername (default) as `ErrorResponse`](#mainunfollowuserbyusername-delete-profilesusernamefollow)
- [output of main.FollowUserByUsername (default) as `ErrorResponse`](#mainfollowuserbyusername-post-profilesusernamefollow)
- [output of main.GetTags (default) as `ErrorResponse`](#maingettags-get-tags)
- [output of main.GetCurrentUser (default) as `ErrorResponse`](#maingetcurrentuser-get-user)
- [output of main.UpdateCurrentUser (default) as `ErrorResponse`](#mainupdatecurrentuser-put-user)
- [output of main.CreateUser (default) as `ErrorResponse`](#maincreateuser-post-users)
- [output of main.Login (default) as `ErrorResponse`](#mainlogin-post-userslogin)

### GenericError

Unexpected error

```go
// Unexpected error
type GenericError struct {
	errors struct {	// GenericErrorErrors
		body []string
	}
}
```

- [output of main.GetArticles (422) as `GenericError`](#maingetarticles-get-articles)
- [output of main.CreateArticle (422) as `GenericError`](#maincreatearticle-post-articles)
- [output of main.GetArticlesFeed (422) as `GenericError`](#maingetarticlesfeed-get-articlesfeed)
- [output of main.DeleteArticle (422) as `GenericError`](#maindeletearticle-delete-articlesslug)
- [output of main.GetArticle (422) as `GenericError`](#maingetarticle-get-articlesslug)
- [output of main.UpdateArticle (422) as `GenericError`](#mainupdatearticle-put-articlesslug)
- [output of main.GetArticleComments (422) as `GenericError`](#maingetarticlecomments-get-articlesslugcomments)
- [output of main.CreateArticleComment (422) as `GenericError`](#maincreatearticlecomment-post-articlesslugcomments)
- [output of main.DeleteArticleComment (422) as `GenericError`](#maindeletearticlecomment-delete-articlesslugcommentsid)
- [output of main.DeleteArticleFavorite (422) as `GenericError`](#maindeletearticlefavorite-delete-articlesslugfavorite)
- [output of main.CreateArticleFavorite (422) as `GenericError`](#maincreatearticlefavorite-post-articlesslugfavorite)
- [output of main.GetProfileByUsername (422) as `GenericError`](#maingetprofilebyusername-get-profilesusername)
- [output of main.UnfollowUserByUsername (422) as `GenericError`](#mainunfollowuserbyusername-delete-profilesusernamefollow)
- [output of main.FollowUserByUsername (422) as `GenericError`](#mainfollowuserbyusername-post-profilesusernamefollow)
- [output of main.GetTags (422) as `GenericError`](#maingettags-get-tags)
- [output of main.GetCurrentUser (422) as `GenericError`](#maingetcurrentuser-get-user)
- [output of main.UpdateCurrentUser (422) as `GenericError`](#mainupdatecurrentuser-put-user)
- [output of main.CreateUser (422) as `GenericError`](#maincreateuser-post-users)
- [output of main.Login (422) as `GenericError`](#mainlogin-post-userslogin)

### GenericErrorErrors



```go
type GenericErrorErrors struct {
	body []string
}
```


### LimitParam

The numbers of items to return.

```go
// The numbers of items to return.
type LimitParam integer
// tags: `minimum:"1"`	// default: 20
```


### OffsetParam

The number of items to skip before starting to collect the result set.

```go
// The number of items to skip before starting to collect the result set.
type OffsetParam integer
// tags: `minimum:"0"`
```


### Profile



```go
type Profile struct {
	bio string

	following boolean	// default: false

	image string

	username string
}
```

- [output of main.GetProfileByUsername (200) as `GetProfileByUsernameOutput[Profile]`](#maingetprofilebyusername-get-profilesusername)
- [output of main.UnfollowUserByUsername (200) as `UnfollowUserByUsernameOutput[Profile]`](#mainunfollowuserbyusername-delete-profilesusernamefollow)
- [output of main.FollowUserByUsername (200) as `FollowUserByUsernameOutput[Profile]`](#mainfollowuserbyusername-post-profilesusernamefollow)

### Time



```go
type  string
// tags: `format:"date-time"`
```


### UnauthorizedError

Unauthorized

```go
// Unauthorized
type UnauthorizedError struct {
}
```

- [output of main.GetArticles (401) as `UnauthorizedError`](#maingetarticles-get-articles)
- [output of main.CreateArticle (401) as `UnauthorizedError`](#maincreatearticle-post-articles)
- [output of main.GetArticlesFeed (401) as `UnauthorizedError`](#maingetarticlesfeed-get-articlesfeed)
- [output of main.DeleteArticle (401) as `UnauthorizedError`](#maindeletearticle-delete-articlesslug)
- [output of main.UpdateArticle (401) as `UnauthorizedError`](#mainupdatearticle-put-articlesslug)
- [output of main.GetArticleComments (401) as `UnauthorizedError`](#maingetarticlecomments-get-articlesslugcomments)
- [output of main.CreateArticleComment (401) as `UnauthorizedError`](#maincreatearticlecomment-post-articlesslugcomments)
- [output of main.DeleteArticleComment (401) as `UnauthorizedError`](#maindeletearticlecomment-delete-articlesslugcommentsid)
- [output of main.DeleteArticleFavorite (401) as `UnauthorizedError`](#maindeletearticlefavorite-delete-articlesslugfavorite)
- [output of main.CreateArticleFavorite (401) as `UnauthorizedError`](#maincreatearticlefavorite-post-articlesslugfavorite)
- [output of main.GetProfileByUsername (401) as `UnauthorizedError`](#maingetprofilebyusername-get-profilesusername)
- [output of main.UnfollowUserByUsername (401) as `UnauthorizedError`](#mainunfollowuserbyusername-delete-profilesusernamefollow)
- [output of main.FollowUserByUsername (401) as `UnauthorizedError`](#mainfollowuserbyusername-post-profilesusernamefollow)
- [output of main.GetCurrentUser (401) as `UnauthorizedError`](#maingetcurrentuser-get-user)
- [output of main.UpdateCurrentUser (401) as `UnauthorizedError`](#mainupdatecurrentuser-put-user)
- [output of main.Login (401) as `UnauthorizedError`](#mainlogin-post-userslogin)

### User



```go
type User struct {
	email string

	token string

	username string

	bio string

	image string
}
```

- [output of main.GetCurrentUser (200) as `GetCurrentUserOutput[User]`](#maingetcurrentuser-get-user)
- [output of main.UpdateCurrentUser (200) as `UpdateCurrentUserOutput[User]`](#mainupdatecurrentuser-put-user)
- [output of main.CreateUser (201) as `CreateUserOutput[User]`](#maincreateuser-post-users)
- [output of main.Login (200) as `LoginOutput[User]`](#mainlogin-post-userslogin)