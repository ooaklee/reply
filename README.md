# reply

`reply` is a Go library that supports developers with shaping and standardising the responses sent from their API service(s). It also allows users to predefine non-successful error objects, giving a granularity down to `title`, `description`, `status code`, and many more through the error manifest(s) passed to the `Replier`.

## Table of Contents

- [Installation](#installation)
- [Getting Started](#getting-started)
  - [How to create a `Replier`](#how-to-create-a-replier)
  - [More about the `ErrorManifest`](#more-about-the-errormanifest)
    - [Deeper look into `ErrorManifestItem`](#deeper-look-into-errormanifestitem)
- [How to send a response(s)](#how-to-send-a-responses)
  - [Making use of error manifest](#making-use-of-error-manifest)
  - [Sending  "successful responses"](#sending--successful-responses)
- [Transfer Objects](#transfer-objects)
  - [Base Transfer Object (`TransferObject`)](#base-transfer-object-transferobject)
  - [Error Transfer Object (`TransferObjectError`)](#error-transfer-object-transferobjecterror)
- [Response Types](#response-types)
  - [Universal Attributes](#universal-attributes)
  - [Error Response Type](#error-response-type)
    - [As code (including aide example)](#as-code-including-aide-example)
    - [JSON Representation](#json-representation)
      - [With `Meta`](#with-meta)
  - [Token Response Type](#token-response-type)
    - [As code (including aide example)](#as-code-including-aide-example-1)
    - [JSON Representation](#json-representation-1)
      - [With `Meta`](#with-meta-1)
  - [Data Response Type](#data-response-type)
    - [As code (including aide example)](#as-code-including-aide-example-2)
    - [JSON Representation](#json-representation-2)
      - [With `Meta`](#with-meta-2)
  - [Default (Blank) Response Type](#default-blank-response-type)
    - [As code (including aide example)](#as-code-including-aide-example-3)
    - [JSON Representation](#json-representation-3)
      - [With `Meta`](#with-meta-3)
- [Copyright](#copyright)

---


## Installation

```sh
  go get github.com/ooaklee/reply
```

## Getting Started

There are several ways you can integrate `reply` into your application. Below, you will find an example of how you can get the most out of this package.

### How to create a `Replier`

When creating a `Replier`, you only have to pass a `reply.ErrorManifest` collection. The collection can be empty or contain as many entries as you'd like.

Just remember, when creating an `Error Response` (Multi or Single), the passed manifest will be used.

```go
// Have a definition of errors you want to use in your application
var (
	// errExample404 is an example 404 error
	errExample404 = errors.New("example-404-error")
	// errExampleDobValidation is an example dob validation error
	errExampleDobValidation = errors.New("example-dob-validation-error")
	// errExampleNameValidation is an example name validation error
	errExampleNameValidation = errors.New("example-name-validation-error")
	// errExampleMissing is an example missing error
	errExampleMissing = errors.New("example-missing-error")
	// errExampleEmailValidation is an example email validation error
	errExampleEmailValidation = errors.New("example-email-validation-error")
)

// (Optional) Create an error manifest to hold correlating errors and their manifest item.
// These can be also be sourced from relevant packages to populate
//
// See how we have to reply.ErrorManifests, on with mulitple
// items and the other with just one.
baseManifest := []reply.ErrorManifest{
    {
      errExample404: reply.ErrorManifestItem{Title: "resource not found", StatusCode: http.StatusNotFound},
      errExampleNameValidation: reply.ErrorManifestItem{Title: "Validation Error", Detail: "The name provided does not meet validation requirements", StatusCode: http.StatusBadRequest, About: "www.example.com/reply/validation/1011", Code: "1011"},
    },
    {errExampleDobValidation: reply.ErrorManifestItem{Title: "Validation Error", Detail: "Check your DoB, and try again.", Code: "100YT", StatusCode: http.StatusBadRequest}},
    // example using error from another package
    // {somepackage.ErrNotFound: reply.ErrorManifestItem{Title: "Not Found", Detail: "The requested resource was not found", Code: "404", StatusCode: http.StatusNotFound}},
  }

// Create Replier to manage the responses going back to consumer(s)
replier := reply.NewReplier(baseManifest)
```

> NOTE - By default, if an `Error Manifest Item` does not have a `StatusCode` set, `reply` will default to `400 (Bad Request)`.

### More about the `ErrorManifest` 

The `ErrorManifest` contains an error and its corresponding `ErrorManifestItem`.

The `ErrorManifestItem` is used to explicitly define the attributes to include in your response's error object. Like previously mentioned, the `ErrorManifestItem` should be created against the explicit error used in your code or taken from a package you're using.

#### Deeper look into `ErrorManifestItem`

`ErrorManifestItems` will come in various sizes depending on how much information you what to make visible to your consumer.

> It is essential to evaluate the exposure level of your API continuously. Is it something that will be used external to your team/ business, and thus minimal information should be given?

The key attributes of the `ErrorManifestItem` are:

- **Title** (`string`): Summary of the error being returned. Try keeping it short and sweet.

- **Detail** (`string`): Gives a more descriptive outline of the error, something with more context.
 
- **StatusCode** (`int`): The [HTTP Status Code](https://httpstatuses.com) associated with respective error. If it's a `5XX` error, it will be the sole error object returned in a `multi error response` scenario.
 
- **About** (`string`): The URL to a page that gives more context about the error

- **Code** (`string`): The `internal code` (application or business) that's used to identify the error

- **Meta** (`interface{}`): Any additional meta-information you may want to pass with your error object

Assuming an `ErrorManifest` containing the following entry was passed to a `Replier`,

```go
// manifest item references error defined in previous example
{errExampleNameValidation: reply.ErrorManifestItem{Title: "Validation Error", Detail: "The name provided does not meet validation requirements", StatusCode: http.StatusBadRequest, About: "www.example.com/reply/validation/1011", Code: "1011"}}
``` 

And its respective error was passed when creating a new error response (`NewHTTPErrorResponse`). `reply` would return the following JSON response:

```json
{
  "errors": [
    {
      "title": "Validation Error",
      "detail": "The name provided does not meet validation requirements",
      "about": "www.example.com/reply/validation/1011",
      "status": "400",
      "code": "1011"
    }
  ]
}
```

If instead, multiple errors were passed to the `NewHTTPMultiErrorResponse` method or errors keys were wrapped and passed to `NewHTTPErrorResponse`, and all had an entry in the `ErrorManifest`, `reply` would return a response similar to the following JSON response:


```json
{
  "errors": [
    {
      "title": "Validation Error",
      "detail": "The name provided does not meet validation requirements",
      "about": "www.example.com/reply/validation/1011",
      "status": "400",
      "code": "1011"
    },
    {
      "title": "Validation Error",
      "detail": "The email provided does not meet validation requirements",
      "status": "400"
    }
  ]
}
```

> NOTE - Not all attributes in the `ErrorManifestItem` have to be specified. By default, if a `StatusCode` is not provided in the item `400` would be set. 

> NOTE - You can create your own custom error json response shape, by using the `reply.WithTransferObjectError` option when creating your replier. Check the [**example simple api ** implementation (`replierWithCustomTransitionObjs`)](examples/example_simple_api.go) for a working example.

## How to send a response(s) 

At the core, you can use `reply` two send both successful and error responses.

When sending an error response, it is essential to make sure you populate the `Error Manifest` passed to the `Replier` with the correct errors source from your code or packages. Otherwise, a `500 - Internal Server Error` response will be sent back to the client by default if it cannot match the passed error in the manifest.

> When matching `errors.Is` is used, the error in the manifest should be the same error as the one passed.

Having expected `ErrorManifest` entries are especially important for `Multi Error` responses. One unmatched error will return a single `500 - Internal Server Error` instead of the array of passed error responses.

### Making use of error manifest

There are currently **3** Replier methods that make use of the `Error Manifest`. These methods are `NewHTTPResponse`, `NewHTTPMultiErrorResponse` and `NewHTTPErrorResponse`.

> NOTE - `NewHTTPResponse` is the base of both the `NewHTTPMultiErrorResponse` and `NewHTTPErrorResponse` aides.
>
> NOTE - You can get `NewHTTPErrorResponse` to behave like `NewHTTPMultiErrorResponse` by wrapping your error keys with `errors.Join(errs...)`

Below you will find an example using `NewHTTPResponse`. However, for simplicity, it's recommended you use one of the error aides. The error [aide implementation is outlined **HERE**](#error-response-type).

```go

// ExampleHandler handler to demostrate how to use package for error
// response
func ExampleHandler(w http.ResponseWriter, r *http.Request) {

  // Pass error to Replier's method to return predefined response, else
  // 500
  _ = replier.NewHTTPResponse(&reply.NewResponseRequest{
    Writer: w,
    // errExample404 defined in previous example
    Error:  errExample404,
  })
}
```

When the endpoint linked to the handler above is called, you should see the following JSON response.

```JSON
{
  "errors": [
    {
      "title": "resource not found",
      "status": "404"
    }
  ]
}
```

> `NOTE` - The `baseManifest` was initially declared, and its item represents the response shown below. The status code is both shown in the response body as a string, and it is also set accordingly. 

### Sending  "successful responses"

The **3** Replier methods that can send "successful responses" are `NewHTTPResponse`, `NewHTTPBlankResponse` and `NewHTTPDataResponse`.

> NOTE - `NewHTTPResponse` is the base of both the `NewHTTPBlankResponse` and `NewHTTPDataResponse` aides.

Below you will find an example using `NewHTTPResponse`, however for simplicity, it's recommended you use either of the follow aide implementations:
 - [Blank Response Aide](#default-blank-response-type)
 - [Data Response Aide](#data-response-type) 

```go

// ExampleGetAllHandler handler to demostrate how to use package for successful 
// response
func ExampleGetAllHandler(w http.ResponseWriter, r *http.Request) {

  // building sample user model
  type user struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
  }

  // emulate users pulled from repository
  mockedQueriedUsers := []user{
    {ID: 1, Name: "John Doe"},
    {ID: 2, Name: "Sam Smith"},
  }

  // build and sent default formatted JSON response for consumption
  // by client
  _ = replier.NewHTTPResponse(&reply.NewResponseRequest{
    Writer:     w,
    Data:       mockedUsers,
    StatusCode: htttp.StatusOK,
  })
}
```

When the endpoint linked to the handler above is called, you should see the following JSON response.

```JSON
{
  "data": [
    {
      "id": 1,
      "name": "John Doe"
    },
    {
      "id": 2,
      "name": "Sam Smith"
    }
  ]
}
```

> `NOTE` - Unlike the error use case, successful requests expect the `StatusCode` to be defined when creating a successful response. If you do not provide a status code, 200 will be assumed.
> 
> It is recommend to use use either the [Blank Response Aide](#default-blank-response-type) or [Data Response Aide](#data-response-type) based on your desired ouput

## Transfer Objects

`Transfer objects` are used to define the shape of various elements within the overall response. In particular, they are used for the `base response object` and the `individual error response object`.

If desired, users can create their own `transfer object` for the `base` and `individual error` response objects with additional logic.

### Base Transfer Object (`TransferObject`)

The `Transfer Object` used for the `base response object` **must** satisfy the following interface:

```go
// TransferObject outlines expected methods of a transfer object
type TransferObject interface {
  SetHeaders(headers map[string]string)
  SetStatusCode(code int)
  SetMeta(meta map[string]interface{})
  SetTokenOne(token string)
  SetTokenTwo(token string)
  GetWriter() http.ResponseWriter
  GetStatusCode() int
  SetWriter(writer http.ResponseWriter)
  SetStatus(transferObjectStatus *TransferObjectStatus)
  RefreshTransferObject() TransferObject
  SetData(data interface{})
}
```

The interface uses relatively self-explanatory method names. Still, if you want to see an example of how one might create your own `transfer object`, you can find the `default transfer object` used by `reply` [here (defaultReplyTransferObject)](./model.go). 

Once your `transfer object` has been created and is valid, you can overwrite the default `transfer object` in your newly created version by using the following code when declaring your `Replier`:

```go
// some implementation of your desired transfer object
var customTransferObject reply.TransferObject

customTransferObject = &foo{}


// create a Replier, overwriting the default transfer object
replier := reply.NewReplier([]reply.ErrorManifest{}, reply.WithTransferObject(customTransferObject))

// use the new Replier as you otherwise would
```

> *NOTE:* you can also pass in your custom transfer object with `&foo{}`, for example:
>
> `replier := reply.NewReplier([]reply.ErrorManifest{}, reply.WithTransferObject(&foo{}))`


> For a live example on how you can use a custom `transfer object`, please look at the [`simple API examples` in this repo](./examples/example_simple_api.go). You are looking out for the `fooReplyTransferObject` implementation.

### Error Transfer Object (`TransferObjectError`)

The `Transfer Object` used for the `individual error response object` **must** satisfy the following interface:

```go
// TransferObjectError outlines expected methods of a transfer object error
type TransferObjectError interface {
  SetTitle(title string)
  GetTitle() string
  SetDetail(detail string)
  GetDetail() string
  SetAbout(about string)
  GetAbout() string
  SetStatusCode(status int)
  GetStatusCode() string
  SetCode(code string)
  GetCode() string
  SetMeta(meta interface{})
  GetMeta() interface{}
  RefreshTransferObject() TransferObjectError
}
```

The interface uses relatively self-explanatory method names. Look at the [deeper dive of the Error Manifest Item](#deeper-look-into-errormanifestitem) to better understand how these methods are used. 

Suppose you want to see an example of how one might create your own `transfer object error`. In that case, you can find the `default transfer object error` used by `reply` [here (defaultReplyTransferObjectError)](./model.go). 

You can overwrite the default `transfer object error` used by your `Replier` by using the following code when declaring your `Replier`:

```go
// some implementation of your desired transfer object
var customTransferObjectError reply.TransferObjectError

customTransferObjectError = &bar{}


// create a Replier, overwriting the default transfer object error (error transfer object)
replier := reply.NewReplier([]reply.ErrorManifest{}, reply.WithTransferObjectError(customTransferObjectError))

// use the new Replier as you otherwise would
```

> *NOTE:* you can also pass in your custom transfer object with `&bar{}`, for example:
>
> `replier := reply.NewReplier([]reply.ErrorManifest{}, reply.customTransferObjectError(&bar{}))`


> For a live example on how you can use a custom `transfer object error` in combination with a custom `transfer object`, please look at the [`simple API examples` in this repo](./examples/example_simple_api.go). You are looking out for the `replierWithCustomTransitionObjs` implementation.
>
> You can set your custom `transfer object error` individually, as shown above.


## Response Types

There are currently four core response types supported by `reply`. They are the `Error`, `Token`, `Data` and `Default` response types. Each type has its JSON representation defined through a [`Transfer Object`](#transfer-objects).

> NOTE: Unless otherwise stated, the `Transfer Objects` assumed will be the [default transfer object (defaultReplyTransferObject)](./model.go) and [default transfer object error (defaultReplyTransferObjectError)](./model.go).

### Universal Attributes

All core response types share universal attributes, which you can set in addition to their outputs. These include:
- Headers
- Meta
- Status Code

> NOTE - `Status Code` is set at different levels dependant on `response type`. For example, the `error response type` is handled in the [`ErrorManifest`](#more-about-the-errormanifest).

### Error Response Type

The `Error` response notifies the consumer when an error/ unexpected behaviour has occurred on the API. There are **2** types of `Error Response Types`, `NewHTTPErrorResponse` and `NewHTTPMultiErrorResponse`.

> Where `NewHTTPMultiErrorResponse` explicitly expects a slice of errors, `NewHTTPErrorResponse` can also return multi errors response by wrapping the manifest errors with `errors.Join(errs...)`.

The error response object forwarded to the consumer is sourced from the [error manifest](#more-about-the-errormanifest). In the event the error's string
representation isn't in the manifest; `reply` will return the consumer a "500 - Internal Server Error" response.

#### As code (including aide example)

To create an individual `error` response use the following code snippet:

```go
// create error manifest using errors defined in application (see previous example for error definitions)
baseManifest := []reply.ErrorManifest{
  {errExample404: reply.ErrorManifestItem{Title: "resource not found", StatusCode: http.StatusNotFound},
    errExampleNameValidation: reply.ErrorManifestItem{Title: "Validation Error", Detail: "The name provided does not meet validation requirements", StatusCode: http.StatusBadRequest, About: "www.example.com/reply/validation/1011", Code: "1011"},
  },
  errExampleDobValidation: reply.ErrorManifestItem{Title: "Validation Error", Detail: "Check your DoB, and try again.", Code: "100YT", StatusCode: http.StatusBadRequest},
}

// create Replier based on error manifest
replier := reply.NewReplier(baseManifest)

func ExampleHandler(w http.ResponseWriter, r *http.Request) {

  _ = replier.NewHTTPResponse(&reply.NewResponseRequest{
    Writer: w,
    Error:  errExample404,
  })
}
```

If you wanted to send a `multi error response`, you could use the following, assuming the same `Replier` from above is being used:

```go
func ExampleHandler(w http.ResponseWriter, r *http.Request) {

  // errors returned
  exampleErrs := []errors{
    errExampleNameValidation,
    errExampleDobValidation,
  }

  _ = replier.NewHTTPResponse(&reply.NewResponseRequest{
    Writer: w,
    Errors: exampleErrs,
  })
}

```

You can also send a `multi error response` with a single error by wrapping the error keys with `errors.Join(errs...)`:

```go
func ExampleHandler(w http.ResponseWriter, r *http.Request) {

  // errors returned
  exampleWrappedErr := errors.Join(errExampleNameValidation, errExampleEmailValidation)

  _ = replier.NewHTTPResponse(&reply.NewResponseRequest{
    Writer: w,
    Error: exampleWrappedErr,
  })
}

```

For readability and simplicity, you can use the `HTTP error response aides`. You can find code snippets using these aides below:

- `Individual Error`

```go
// inside of the request handler
_ = replier.NewHTTPErrorResponse(w, exampleErr)
```

- `Multi Error`

```go
// inside of the request handler
_ = replier.NewHTTPMultiErrorResponse(w, exampleErrs)
```

You can also add additional `headers` and `meta data` to the response by using the optional `WithHeaders` and/ or `WithMeta` response attributes respectively. For example:

```go
_ = replier.NewHTTPErrorResponse(w, exampleErr, reply.WithMeta(map[string]interface{}{
    "example": "meta in error reponse",
  }))
```

**OR**

```go
_ = replier.NewHTTPMultiErrorResponse(w, exampleErrs, reply.WithMeta(map[string]interface{}{
    "example": "meta in error reponse",
  }))
```

#### JSON Representation

`Error` responses are returned with the format. The following responses are based on the examples above, so your response content will vary.

- `Individual Error`

```JSON
{
  "errors": [
    {
      "title": "resource not found",
      "status": "404"
    }
  ]
}
```

- `Multi Error`

```json
{
  "errors": [
    {
      "title": "Validation Error",
      "detail": "The name provided does not meet validation requirements",
      "about": "www.example.com/reply/validation/1011",
      "status": "400",
      "code": "1011"
    },
    {
      "title": "Validation Error",
      "detail": "The email provided does not meetvalidation requirements",
      "status": "400"
    }
  ]
}
```

##### With `Meta`

When a `meta` is also declared, the response will have the following format. It can be as big or small as needed.

- `Individual Error`

```JSON
{
  "errors": [
    {
      "title": "resource not found",
      "status": "404"
    }
  ],
  "meta": {
    "example": "meta in error reponse"
  }
}
```

- `Multi Error`

```json
{
  "errors": [
    {
      "title": "Validation Error",
      "detail": "The name provided does not meet validation requirements",
      "about": "www.example.com/reply/validation/1011",
      "status": "400",
      "code": "1011"
    },
    {
      "title": "Validation Error",
      "detail": "The email provided does not meetvalidation requirements",
      "status": "400"
    }
  ],
  "meta": {
    "example": "meta in error reponse"
  }
}
```

### Token Response Type

The `token` response sends the consumer tokens. Currently, it is limited to **2** tokens, and with the default `Transfer Object`, `TokenOne` represents `access_token`, and `TokenTwo` represents `refresh_token`. However, if you use other ID/ JSON attributes to describe your tokens for your API, you can [**create a `Custom Transfer Object`**](#transfer-objects).

Again, when using the default `Transfer Object`, the supported `TokenOne` and `TokenTwo` represent `acccess_token` and `refresh_token`, respectively. If either is passed in the response request, `reply` will default to this response type.

#### As code (including aide example)

To create a `token` response use the following code snippet:

```go
replier := reply.NewReplier([]reply.ErrorManifest{})

func ExampleHandler(w http.ResponseWriter, r *http.Request) {

  // do something to get tokens

  _ = replier.NewHTTPResponse(&reply.NewResponseRequest{
    Writer:     w,
    TokenOne:   "08a0a043-b532-4cea-8117-364739f2d994",
    TokenTwo:   "08b29914-09a8-4a4a-8aa5-b1ffaff266e6",
    StatusCode: 200,
  })
}
```

For readability and simplicity, you can use the `HTTP token response aide`. You can find a code snippet using this aide below:

```go
// inside of the request handler
_ = replier.NewHTTPTokenResponse(w, 200, "08a0a043-b532-4cea-8117-364739f2d994", "08b29914-09a8-4a4a-8aa5-b1ffaff266e6")
```

You can also add additional `headers` and `meta data` to the response by using the optional `WithHeaders` and/ or `WithMeta` response attributes respectively. For example:

```go
_ = replier.NewHTTPTokenResponse(w, 200, "08a0a043-b532-4cea-8117-364739f2d994", "08b29914-09a8-4a4a-8aa5-b1ffaff266e6", reply.WithMeta(map[string]interface{}{
 "example": "meta in token reponse",
}))
```

> NOTE: If you only want to return one token, pass an empty string, i.e. `""`. Although, you must give at least one token string. 

#### JSON Representation

`Error` responses are returned with the format.

```JSON
{
  "access_token": "08a0a043-b532-4cea-8117-364739f2d994",
  "refresh_token": "08b29914-09a8-4a4a-8aa5-b1ffaff266e6"
}
```

##### With `Meta`

When a `meta` is also declared, the response will have the following format. It can be as big or small as needed.

```JSON
{
  "access_token": "08a0a043-b532-4cea-8117-364739f2d994",
  "refresh_token": "08b29914-09a8-4a4a-8aa5-b1ffaff266e6",
  "meta": {
    "example": "meta in token reponse"
  }
}
```

### Data Response Type

The `data` response can be seen as a *successful* response. It parses the passed struct into its JSON representation and passes it to the consumer in the JSON response. The JSON response below will represent a response if the data passed was a user struct with the:
- `id` 1
- `name` john doe
- `dob` 1/1/1970

#### As code (including aide example)

To create a `data` response use the following code snippet:

```go
type user struct {
  id   int    `json:"id"`
  name string `json:"name"`
  dob  string `json:"dob"`
}

replier := reply.NewReplier([]reply.ErrorManifest{})

func ExampleHandler(w http.ResponseWriter, r *http.Request) {

  u := user{
    id:   1,
    name: "john doe",
    dob:  "1/1/1970",
  }

  _ = replier.NewHTTPResponse(&reply.NewResponseRequest{
    Writer:     w,
    Data:       u,
    StatusCode: 201,
  })
}
```


For readability and simplicity, you can use the `HTTP data (successful) response aide`. You can find a code snippet using this aide below:

```go
// inside of the request handler
_ = replier.NewHTTPDataResponse(w, 201, u)
```

You can also add additional `headers` and `meta data` to the response by using the optional `WithHeaders` and/ or `WithMeta` response attributes respectively. For example:

```go
_ = replier.NewHTTPDataResponse(w, 201, u, reply.WithMeta(map[string]interface{}{
    "example": "meta in data reponse",
  }))
```

#### JSON Representation

`Data` responses are returned with the format.

```JSON
{
  "data": {
    "id": 1,
    "name": "john doe",
    "dob": "1/1/1970"
  }
}
```

##### With `Meta`

When a `meta` is also declared, the response will have the following format. It can be as big or small as needed.

```JSON
{
  "data": {
    "id": 1,
    "name": "john doe",
    "dob": "1/1/1970"
  },
  "meta": {
    "example": "meta in data reponse"
  }
}
```

### Default (Blank) Response Type

The `default` (blank) response returns `"{}"` with a status code of `200` if no `error`, `tokens`, `data` and `status code` is passed. If desired, another `status code` can be specified with `default` responses.

#### As code (including aide example)

To create a `default` response use the following code snippet:

```go
replier := reply.NewReplier([]reply.ErrorManifest{})

func ExampleHandler(w http.ResponseWriter, r *http.Request) {

  _ = replier.NewHTTPResponse(&reply.NewResponseRequest{
    Writer:     w,
    StatusCode: 200,
  })
}
```

For readability and simplicity, you can use the `HTTP default (blank) response aide`. You can find a code snippet using this aide below:

```go
// inside of the request handler
_ = replier.NewHTTPBlankResponse(w, 200)
```

You can also add additional `headers` and `meta data` to the response by using the optional `WithHeaders` and/ or `WithMeta` response attributes respectively. For example:

```go
_ = replier.NewHTTPBlankResponse(w, 200, reply.WithMeta(map[string]interface{}{
    "example": "meta in default reponse",
  }))
```


#### JSON Representation

`Default` responses are returned with the format.

```JSON
{
  "data": "{}"
}
```

##### With `Meta`

When a `meta` is also declared, the response will have the following format. It can be as big or small as needed.

```JSON
{
  "data": "{}",
  "meta": {
    "example": "meta in default reponse"
  }
}
```

## Copyright

Copyright (C) 2021 by Leon Silcott <leon@boasi.io>.

reply library released under MIT License.
See [LICENSE](https://github.com/ooaklee/reply/blob/master/LICENSE) for details.