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
| `GET /articles` | [main.GetArticles](#maingetarticles-get-articles)  | `main` |  |
| `POST /articles` | [main.CreateArticle](#maincreatearticle-post-articles)  | `main` |  |
| `GET /articles/feed` | [main.GetArticlesFeed](#maingetarticlesfeed-get-articlesfeed)  | `main` |  |
| `DELETE /articles/{slug}` | [main.DeleteArticle](#maindeletearticle-delete-articlesslug)  | `main` |  |
| `GET /articles/{slug}` | [main.GetArticle](#maingetarticle-get-articlesslug)  | `main` |  |
| `PUT /articles/{slug}` | [main.UpdateArticle](#mainupdatearticle-put-articlesslug)  | `main` |  |
| `GET /articles/{slug}/comments` | [main.GetArticleComments](#maingetarticlecomments-get-articlesslugcomments)  | `main` |  |
| `POST /articles/{slug}/comments` | [main.CreateArticleComment](#maincreatearticlecomment-post-articlesslugcomments)  | `main` |  |
| `DELETE /articles/{slug}/comments/{id}` | [main.DeleteArticleComment](#maindeletearticlecomment-delete-articlesslugcommentsid)  | `main` |  |
| `DELETE /articles/{slug}/favorite` | [main.DeleteArticleFavorite](#maindeletearticlefavorite-delete-articlesslugfavorite)  | `main` |  |
| `POST /articles/{slug}/favorite` | [main.CreateArticleFavorite](#maincreatearticlefavorite-post-articlesslugfavorite)  | `main` |  |
| `GET /profiles/{username}` | [main.GetProfileByUsername](#maingetprofilebyusername-get-profilesusername)  | `main` |  |
| `DELETE /profiles/{username}/follow` | [main.UnfollowUserByUsername](#mainunfollowuserbyusername-delete-profilesusernamefollow)  | `main` |  |
| `POST /profiles/{username}/follow` | [main.FollowUserByUsername](#mainfollowuserbyusername-post-profilesusernamefollow)  | `main` |  |
| `GET /tags` | [main.GetTags](#maingettags-get-tags)  | `main` |  |
| `GET /user` | [main.GetCurrentUser](#maingetcurrentuser-get-user)  | `main` |  |
| `PUT /user` | [main.UpdateCurrentUser](#mainupdatecurrentuser-put-user)  | `main` |  |
| `POST /users/` | [main.CreateUser](#maincreateuser-post-users)  | `main` |  |
| `POST /users/login` | [main.Login](#mainlogin-post-userslogin)  | `main` | handlers |


### main.GetArticles `GET /articles`



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
### main.CreateArticle `POST /articles`



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
### main.GetArticlesFeed `GET /articles/feed`



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
### main.DeleteArticle `DELETE /articles/{slug}`



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
### main.GetArticle `GET /articles/{slug}`



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
### main.UpdateArticle `PUT /articles/{slug}`



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
### main.GetArticleComments `GET /articles/{slug}/comments`



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
### main.CreateArticleComment `POST /articles/{slug}/comments`



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
### main.DeleteArticleComment `DELETE /articles/{slug}/comments/{id}`



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
### main.DeleteArticleFavorite `DELETE /articles/{slug}/favorite`



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
### main.CreateArticleFavorite `POST /articles/{slug}/favorite`



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
### main.GetProfileByUsername `GET /profiles/{username}`



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
### main.UnfollowUserByUsername `DELETE /profiles/{username}/follow`



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
### main.FollowUserByUsername `POST /profiles/{username}/follow`



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
### main.GetTags `GET /tags`



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
### main.GetCurrentUser `GET /user`



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
### main.UpdateCurrentUser `PUT /user`



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

handlers

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

handlers





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