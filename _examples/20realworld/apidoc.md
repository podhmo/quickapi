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
| `GET /articles` | [main.GetArticles](#maingetarticles-get-articles)  | `main` | Get recent articles globally |
| `POST /articles` | [main.CreateArticle](#maincreatearticle-post-articles)  | `main` | Create an article |
| `GET /articles/feed` | [main.GetArticlesFeed](#maingetarticlesfeed-get-articlesfeed)  | `main` | Get recent articles from users you follow |
| `DELETE /articles/{slug}` | [main.DeleteArticle](#maindeletearticle-delete-articlesslug)  | `main` | Delete an article |
| `GET /articles/{slug}` | [main.GetArticle](#maingetarticle-get-articlesslug)  | `main` | Get an article |
| `PUT /articles/{slug}` | [main.UpdateArticle](#mainupdatearticle-put-articlesslug)  | `main` | Update an article |
| `GET /articles/{slug}/comments` | [main.GetArticleComments](#maingetarticlecomments-get-articlesslugcomments)  | `main` | Get comments for an article |
| `POST /articles/{slug}/comments` | [main.CreateArticleComment](#maincreatearticlecomment-post-articlesslugcomments)  | `main` | Create a comment for an article |
| `DELETE /articles/{slug}/comments/{id}` | [main.DeleteArticleComment](#maindeletearticlecomment-delete-articlesslugcommentsid)  | `main` | Delete a comment for an article |
| `DELETE /articles/{slug}/favorite` | [main.DeleteArticleFavorite](#maindeletearticlefavorite-delete-articlesslugfavorite)  | `main` | Unfavorite an article |
| `POST /articles/{slug}/favorite` | [main.CreateArticleFavorite](#maincreatearticlefavorite-post-articlesslugfavorite)  | `main` | Favorite an article |
| `GET /profiles/{username}` | [main.GetProfileByUsername](#maingetprofilebyusername-get-profilesusername)  | `main` | Get a profile |
| `DELETE /profiles/{username}/follow` | [main.UnfollowUserByUsername](#mainunfollowuserbyusername-delete-profilesusernamefollow)  | `main` | Unfollow a user |
| `POST /profiles/{username}/follow` | [main.FollowUserByUsername](#mainfollowuserbyusername-post-profilesusernamefollow)  | `main` | Follow a user |
| `GET /tags` | [main.GetTags](#maingettags-get-tags)  | `main` | Get tags |
| `GET /user` | [main.GetCurrentUser](#maingetcurrentuser-get-user)  | `main` | Get current user |
| `PUT /user` | [main.UpdateCurrentUser](#mainupdatecurrentuser-put-user)  | `main` | Update current user |
| `POST /users/` | [main.CreateUser](#maincreateuser-post-users)  | `main` |  |
| `POST /users/login` | [main.Login](#mainlogin-post-userslogin)  | `main` | Existing user login |


### main.GetArticles `GET /articles`

Get recent articles globally

| name | value |
| --- | --- |
| operationId | main.GetArticles |
| endpoint | `GET /articles` |
| tags | `main` |




#### output (application/json)

```go
// GET /articles (200)
type Output200 struct {	// 
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
| operationId | main.CreateArticle |
| endpoint | `POST /articles` |
| tags | `main` |




#### output (application/json)

```go
// POST /articles (200)
type Output200 struct {	// 
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
| operationId | main.GetArticlesFeed |
| endpoint | `GET /articles/feed` |
| tags | `main` |




#### output (application/json)

```go
// GET /articles/feed (200)
type Output200 struct {	// 
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
| operationId | main.DeleteArticle |
| endpoint | `DELETE /articles/{slug}` |
| tags | `main` |



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
| operationId | main.GetArticle |
| endpoint | `GET /articles/{slug}` |
| tags | `main` |



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
type Output200 struct {	// 
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
| operationId | main.UpdateArticle |
| endpoint | `PUT /articles/{slug}` |
| tags | `main` |



#### input (application/json)

```go
// PUT /articles/{slug}
type Input struct {
	// Slug of the article to update
	slug string `in:"path"`
}
```

#### output (application/json)

```go
// PUT /articles/{slug} (200)
type Output200 struct {	// 
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
| operationId | main.GetArticleComments |
| endpoint | `GET /articles/{slug}/comments` |
| tags | `main` |



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
type Output200 struct {	// 
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
| operationId | main.CreateArticleComment |
| endpoint | `POST /articles/{slug}/comments` |
| tags | `main` |



#### input (application/json)

```go
// POST /articles/{slug}/comments
type Input struct {
	// Slug of the article that you want to create a comment
	slug string `in:"path"`
}
```

#### output (application/json)

```go
// POST /articles/{slug}/comments (200)
type Output200 struct {	// 
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
| operationId | main.DeleteArticleComment |
| endpoint | `DELETE /articles/{slug}/comments/{id}` |
| tags | `main` |



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
| operationId | main.DeleteArticleFavorite |
| endpoint | `DELETE /articles/{slug}/favorite` |
| tags | `main` |



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
type Output200 struct {	// 
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
| operationId | main.CreateArticleFavorite |
| endpoint | `POST /articles/{slug}/favorite` |
| tags | `main` |



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
type Output200 struct {	// 
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
| operationId | main.GetProfileByUsername |
| endpoint | `GET /profiles/{username}` |
| tags | `main` |



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
type Output200 struct {	// 
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
| operationId | main.UnfollowUserByUsername |
| endpoint | `DELETE /profiles/{username}/follow` |
| tags | `main` |



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
type Output200 struct {	// 
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
| operationId | main.FollowUserByUsername |
| endpoint | `POST /profiles/{username}/follow` |
| tags | `main` |



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
type Output200 struct {	// 
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
| operationId | main.GetTags |
| endpoint | `GET /tags` |
| tags | `main` |




#### output (application/json)

```go
// GET /tags (200)
type Output200 struct {	// 
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
| operationId | main.GetCurrentUser |
| endpoint | `GET /user` |
| tags | `main` |




#### output (application/json)

```go
// GET /user (200)
type Output200 struct {	// 
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
| operationId | main.UpdateCurrentUser |
| endpoint | `PUT /user` |
| tags | `main` |




#### output (application/json)

```go
// PUT /user (200)
type Output200 struct {	// 
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



| name | value |
| --- | --- |
| operationId | main.CreateUser |
| endpoint | `POST /users/` |
| tags | `main` |




#### output (application/json)

```go
// POST /users/ (200)
type Output200 struct {	// 
}

// POST /users/ (default)
// default error
type OutputDefault struct {	// ErrorResponse
	code integer

	error string

	detail? []string
}
```
### main.Login `POST /users/login`

Existing user login

| name | value |
| --- | --- |
| operationId | main.Login |
| endpoint | `POST /users/login` |
| tags | `main` |




#### output (application/json)

```go
// POST /users/login (200)
type Output200 struct {	// 
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
| [ErrorResponse](#errorresponse) | represents a normal error response type |



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