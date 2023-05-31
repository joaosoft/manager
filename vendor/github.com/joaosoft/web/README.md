# web
[![Build Status](https://travis-ci.org/joaosoft/web.svg?branch=master)](https://travis-ci.org/joaosoft/web) | [![codecov](https://codecov.io/gh/joaosoft/web/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/web) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/web)](https://goreportcard.com/report/github.com/joaosoft/web) | [![GoDoc](https://godoc.org/github.com/joaosoft/web?status.svg)](https://godoc.org/github.com/joaosoft/web)

A simple and fast web server and client.

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for 
* Common http methods
* Single/Multiple File Upload
* Single/Multiple File Download
* Middlewares
* Filters

## With support for methods
* HEAD
* GET
* POST
* PUT
* DELETE
* PATCH
* COPY
* CONNECT
* OPTIONS
* TRACE
* LINK
* UNLINK
* PURGE
* LOCK
* UNLOCK
* PROPFIND
* VIEW
* ANY (used only on filters)

## With authentication types
* basic
* jwt

## With attachment modes
* [default] zip files when returns more then one file 
  - on client WithClientAttachmentMode(web.MultiAttachmentModeZip)
  - on server WithServerAttachmentMode(web.MultiAttachmentModeZip)
* [experimental] returns attachmentes splited by a boundary defined on header Content-Type 
  - on client WithClientAttachmentMode(web.MultiAttachmentModeBoundary)
  - on server WithServerAttachmentMode(web.MultiAttachmentModeBoundary)


>### Go
```
go get github.com/joaosoft/web
```

## Usage 
This examples are available in the project at [web/examples](https://github.com/joaosoft/web/tree/master/examples)

### Server
```go
func main() {
	// create a new server
	w, err := web.NewServer(web.WithServerMultiAttachmentMode(web.MultiAttachmentModeBoundary))
	if err != nil {
		panic(err)
	}

	claims := jwt.Claims{"name": "joao", "age": 30}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte("bananas"), nil
	}

	checkFunc := func(c jwt.Claims) (bool, error) {
		if claims["name"] == c["name"].(string) &&
			claims["age"] == int(c["age"].(float64)) {
			return true, nil
		}
		return false, fmt.Errorf("invalid jwt session token")
	}

	// add filters
	w.AddMiddlewares(MyMiddlewareOne(), MyMiddlewareTwo())
	w.AddFilter("/hello/:name", web.PositionBefore, MyFilterOne(), web.MethodPost)
	w.AddFilter("/hello/:name/upload", web.PositionBefore, MyFilterTwo(), web.MethodPost)
	w.AddFilter("*", web.PositionBefore, MyFilterThree(), web.MethodPost)

	w.AddFilter("/hello/:name", web.PositionBefore, MyFilterTwo(), web.MethodGet)

	w.AddFilter("/form-data", web.PositionAfter, MyFilterThree(), web.MethodAny)

	// add routes
	w.AddRoutes(
		web.NewRoute(web.MethodOptions, "*", HandlerHelloForOptions, web.MiddlewareOptions()),
		web.NewRoute(web.MethodGet, "/auth-basic", HandlerForGet, web.MiddlewareCheckAuthBasic("joao", "ribeiro")),
		web.NewRoute(web.MethodGet, "/auth-jwt", HandlerForGet, web.MiddlewareCheckAuthJwt(keyFunc, checkFunc)),
		web.NewRoute(web.MethodHead, "/hello/:name", HandlerHelloForHead),
		web.NewRoute(web.MethodGet, "/hello/:name", HandlerHelloForGet, MyMiddlewareThree()),
		web.NewRoute(web.MethodPost, "/hello/:name", HandlerHelloForPost),
		web.NewRoute(web.MethodPut, "/hello/:name", HandlerHelloForPut),
		web.NewRoute(web.MethodDelete, "/hello/:name", HandlerHelloForDelete),
		web.NewRoute(web.MethodPatch, "/hello/:name", HandlerHelloForPatch),
		web.NewRoute(web.MethodCopy, "/hello/:name", HandlerHelloForCopy),
		web.NewRoute(web.MethodConnect, "/hello/:name", HandlerHelloForConnect),
		web.NewRoute(web.MethodOptions, "/hello/:name", HandlerHelloForOptions, web.MiddlewareOptions()),
		web.NewRoute(web.MethodTrace, "/hello/:name", HandlerHelloForTrace),
		web.NewRoute(web.MethodLink, "/hello/:name", HandlerHelloForLink),
		web.NewRoute(web.MethodUnlink, "/hello/:name", HandlerHelloForUnlink),
		web.NewRoute(web.MethodPurge, "/hello/:name", HandlerHelloForPurge),
		web.NewRoute(web.MethodLock, "/hello/:name", HandlerHelloForLock),
		web.NewRoute(web.MethodUnlock, "/hello/:name", HandlerHelloForUnlock),
		web.NewRoute(web.MethodPropFind, "/hello/:name", HandlerHelloForPropFind),
		web.NewRoute(web.MethodView, "/hello/:name", HandlerHelloForView),
		web.NewRoute(web.MethodGet, "/hello/:name/download", HandlerHelloForDownloadFiles),
		web.NewRoute(web.MethodGet, "/hello/:name/download/one", HandlerHelloForDownloadOneFile),
		web.NewRoute(web.MethodPost, "/hello/:name/upload", HandlerHelloForUploadFiles),
		web.NewRoute(web.MethodGet, "/form-data", HandlerFormDataForGet),
		web.NewRoute(web.MethodGet, "/url-form-data", HandlerUrlFormDataForGet),
	)

	w.AddNamespace("/p").AddRoutes(
		web.NewRoute(web.MethodGet, "/hello/:name", HandlerHelloForGet, MyMiddlewareFour()),
	)

	// start the server
	if err := w.Start(); err != nil {
		panic(err)
	}
}

func MyFilterOne() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx *web.Context) error {
			fmt.Println("HELLO I'M THE FILTER ONE")
			return next(ctx)
		}
	}
}

func MyFilterTwo() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx *web.Context) error {
			fmt.Println("HELLO I'M THE FILTER TWO")
			return next(ctx)
		}
	}
}

func MyFilterThree() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx *web.Context) error {
			fmt.Println("HELLO I'M THE FILTER THREE")
			return next(ctx)
		}
	}
}

func MyMiddlewareOne() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx *web.Context) error {
			fmt.Println("HELLO I'M THE MIDDLEWARE ONE")
			return next(ctx)
		}
	}
}

func MyMiddlewareTwo() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx *web.Context) error {
			fmt.Println("HELLO I'M THE MIDDLEWARE TWO")
			return next(ctx)
		}
	}
}

func MyMiddlewareThree() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx *web.Context) error {
			fmt.Println("HELLO I'M THE MIDDLEWARE THREE")
			return next(ctx)
		}
	}
}

func MyMiddlewareFour() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx *web.Context) error {
			fmt.Println("HELLO I'M THE MIDDLEWARE FOUR")
			return next(ctx)
		}
	}
}

func HandlerForGet(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR GET")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \"guest\" }"),
	)
}

func HandlerHelloForHead(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR HEAD")

	return ctx.Response.NoContent(web.StatusOK)
}

func HandlerHelloForGet(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR GET")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForPost(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR POST")

	data := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{}
	ctx.Request.Bind(&data)
	fmt.Printf("DATA: %+v", data)

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForPut(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR PUT")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForDelete(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR DELETE")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForPatch(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR PATCH")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForCopy(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR COPY")

	return ctx.Response.NoContent(web.StatusOK)
}

func HandlerHelloForConnect(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR CONNECT")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForOptions(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR OPTIONS")

	return ctx.Response.NoContent(web.StatusNoContent)
}

func HandlerHelloForTrace(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR TRACE")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForLink(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR LINK")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForUnlink(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR UNLINK")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForPurge(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR PURGE")

	return ctx.Response.NoContent(web.StatusOK)
}

func HandlerHelloForLock(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR LOCK")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForUnlock(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR UNLOCK")

	return ctx.Response.NoContent(web.StatusOK)
}
func HandlerHelloForPropFind(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR PROPFIND")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForView(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR VIEW")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForDownloadOneFile(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR DOWNLOAD ONE FILE")

	dir, _ := os.Getwd()
	body, _ := web.ReadFile(fmt.Sprintf("%s%s", dir, "/examples/data/a.text"), nil)
	ctx.Response.Attachment("text_a.text", body)

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForDownloadFiles(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR DOWNLOAD FILES")

	dir, _ := os.Getwd()
	body, _ := web.ReadFile(fmt.Sprintf("%s%s", dir, "/examples/data/a.text"), nil)
	ctx.Response.Attachment("text_a.text", body)

	body, _ = web.ReadFile(fmt.Sprintf("%s%s", dir, "/examples/data/b.text"), nil)
	ctx.Response.Attachment("text_b.text", body)

	body, _ = web.ReadFile(fmt.Sprintf("%s%s", dir, "/examples/data/c.text"), nil)
	ctx.Response.Inline("text_c.text", body)

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForUploadFiles(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR UPLOAD FILES")

	fmt.Printf("\nAttachments: %+v\n", ctx.Request.FormData)
	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerFormDataForGet(ctx *web.Context) error {
	fmt.Println("HANDLING FORM DATA FOR GET")

	fmt.Printf("\nreceived")
	fmt.Printf("\nvar_one: %s", ctx.Request.GetFormDataString("var_one"))
	fmt.Printf("\nvar_two: %s", ctx.Request.GetFormDataString("var_two"))

	ctx.Response.SetFormData("var_one", "one")
	ctx.Response.SetFormData("var_two", "2")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \"form-data\" }"),
	)
}

func HandlerUrlFormDataForGet(ctx *web.Context) error {
	fmt.Println("HANDLING URL FORM DATA FOR GET")

	fmt.Printf("\nreceived")
	fmt.Printf("\nvar_one: %s", ctx.Request.GetFormDataString("var_one"))
	fmt.Printf("\nvar_two: %s", ctx.Request.GetFormDataString("var_two"))

	ctx.Response.SetFormData("var_one", "one")
	ctx.Response.SetFormData("var_two", "2")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \"form-data\" }"),
	)
}
```

### Client
```go
func main() {
	// create a new client
	c, err := web.NewClient(web.WithClientMultiAttachmentMode(web.MultiAttachmentModeBoundary))
	if err != nil {
		panic(err)
	}

	requestGet(c)
	requestPost(c)

	requestGetBoundary(c)

	requestAuthBasic(c)
	requestAuthJwt(c)

	requestOptionsOK(c)
	requestOptionsNotFound(c)

	bindFormData(c)
	bindUrlFormData(c)
}

func requestGet(c *web.Client) {
	request, err := c.NewRequest(web.MethodGet, "localhost:9001/hello/joao?a=1&b=2&c=1,2,3", web.ContentTypeApplicationJSON, nil)
	if err != nil {
		panic(err)
	}

	response, err := request.Send()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", response)
}

func requestPost(c *web.Client) {
	request, err := c.NewRequest(web.MethodPost, "localhost:9001/hello/joao?a=1&b=2&c=1,2,3", web.ContentTypeApplicationJSON, nil)
	if err != nil {
		panic(err)
	}

	data := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "joao",
		Age:  30,
	}

	bytes, _ := json.Marshal(data)

	response, err := request.WithBody(bytes).Send()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", string(response.Body))
}

func requestGetBoundary(c *web.Client) {
	request, err := c.NewRequest(web.MethodGet, "localhost:9001/hello/joao/download", web.ContentTypeApplicationJSON, nil)
	if err != nil {
		panic(err)
	}

	response, err := request.Send()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", response)
}

func requestOptionsOK(c *web.Client) {
	request, err := c.NewRequest(web.MethodOptions, "localhost:9001/auth-basic", web.ContentTypeApplicationJSON, nil)
	if err != nil {
		panic(err)
	}

	_, err = request.WithAuthBasic("joao", "ribeiro")
	if err != nil {
		panic(err)
	}

	request.SetHeader(web.HeaderAccessControlRequestMethod, []string{string(web.MethodGet)})
	response, err := request.Send()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n\n%d: %s\n\n", response.Status, string(response.Body))
}

func requestOptionsNotFound(c *web.Client) {
	request, err := c.NewRequest(web.MethodOptions, "localhost:9001/auth-basic-invalid", web.ContentTypeApplicationJSON, nil)
	if err != nil {
		panic(err)
	}

	_, err = request.WithAuthBasic("joao", "ribeiro")
	if err != nil {
		panic(err)
	}

	request.SetHeader(web.HeaderAccessControlRequestMethod, []string{string(web.MethodGet)})
	response, err := request.Send()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n\n%d: %s\n\n", response.Status, string(response.Body))
}

func requestAuthBasic(c *web.Client) {
	request, err := c.NewRequest(web.MethodGet, "localhost:9001/auth-basic", web.ContentTypeApplicationJSON, nil)
	if err != nil {
		panic(err)
	}

	_, err = request.WithAuthBasic("joao", "ribeiro")
	if err != nil {
		panic(err)
	}

	response, err := request.Send()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n\n%d: %s\n\n", response.Status, string(response.Body))
}

func requestAuthJwt(c *web.Client) {
	request, err := c.NewRequest(web.MethodGet, "localhost:9001/auth-jwt", web.ContentTypeApplicationJSON, nil)
	if err != nil {
		panic(err)
	}

	_, err = request.WithAuthJwt(jwt.Claims{"name": "joao", "age": 30}, "bananas")
	if err != nil {
		panic(err)
	}

	response, err := request.Send()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n\n%d: %s\n\n", response.Status, string(response.Body))
}

func bindFormData(c *web.Client) {
	request, err := c.NewRequest(web.MethodGet, "localhost:9001/form-data", web.ContentTypeMultipartFormData, nil)
	if err != nil {
		panic(err)
	}

	request.SetFormData("var_one", "one")
	request.SetFormData("var_two", "2")

	response, err := request.Send()
	if err != nil {
		panic(err)
	}

	formData := struct {
		VarOne string `json:"var_one"`
		VarTwo int    `json:"var_two"`
	}{}

	if err := response.BindFormData(&formData); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("\nvar_one: %s", response.GetFormDataString("var_one"))
	fmt.Printf("\nvar_two: %s", response.GetFormDataString("var_two"))

	fmt.Printf("\n\nFORM DATA: %+v\n", formData)
}

func bindUrlFormData(c *web.Client) {
	request, err := c.NewRequest(web.MethodGet, "localhost:9001/url-form-data", web.ContentTypeApplicationForm, nil)
	if err != nil {
		panic(err)
	}

	request.SetFormData("var_one", "one")
	request.SetFormData("var_two", "2")

	response, err := request.Send()
	if err != nil {
		panic(err)
	}

	formData := struct {
		VarOne string `json:"var_one"`
		VarTwo int    `json:"var_two"`
	}{}

	if err := response.BindFormData(&formData); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("\nvar_one: %s", response.GetFormDataString("var_one"))
	fmt.Printf("\nvar_two: %s", response.GetFormDataString("var_two"))

	fmt.Printf("\n\nURL FORM DATA: %+v\n", formData)
}
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
