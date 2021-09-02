# reply

`reply` is a Go library that supports developers with standardising the responses sent from their API service(s). It allows users to predefine non-successful messages and their corresponding status code based on errors manifest passed to the `replier`.

## Installation

```sh
go get github.com/ooaklee/reply
```

## Examples

There are several ways you can integrate `reply` into your application. Below, you will find an example of how you can get the most out of this package.

### How to create a `replier`

```go
// (Optional) Create a error manifest, to hold correlating error as string and it's manifest
// item
baseManifest := []reply.ErrorManifest{
                {"example-404-error": reply.ErrorManifestItem{Message: "resource not found", StatusCode: http.StatusNotFound}},
            }

// Create replier to manage the responses going back to consumer(s)
replier := reply.NewReplier(baseManifest)
```

### How to send response(s) 

You can use `reply` for both successful and error based responses.

> `NOTE` - When sending an error response, it is essential to make sure you populate the `replier`'s error manifest with the correct errors. Otherwise, a `500 - Internal Server Error` response will be sent back to the client by default if it cannot match the passed error in the manifest.

#### Making use of error manifest

```go

// ExampleHandler handler to demostrate how to use package for error
// response
func ExampleHandler(w http.ResponseWriter, r *http.Request) {

    // Create error with value corresponding to one of the manifest's entry's key
    exampleErr := errors.New("example-404-error")


    // Pass error to replier's method to return predefined response, else
    // 500
    _ := replier.NewHTTPResponse(&reply.NewResponseRequest{
        Writer: w,
        Error:  exampleErr,
    })
}
```

When the endpoint linked to the handler above is called, you should see the following JSON response.

> `NOTE` - The `baseManifest` was initially declared, and its item represents the response shown below. Although the status code is not shown in the response body, it to has been set accordingly and returned to the consumer.

```JSON

{
    "status": {
        "message": "resource not found"
    }
}
```

#### Sending client successful response


```go

// ExampleGetAllHandler handler to demostrate how to use package for successful 
// response
func ExampleGetAllHandler(w http.ResponseWriter, r *http.Request) {

    // building sample user model 
    type user struct {
        ID int `json:"id"`
        Name string `json:"name"`
    }

    // emulate users pulled from repository
    mockedQueriedUsers := []user{
        {ID: 1, Name: "John Doe"},
        {ID: 2, Name: "Sam Smith"},
    }


    // build and sent default formatted JSON response for consumption
    // by client 
    _ := replier.NewHTTPResponse(&reply.NewResponseRequest{
        Writer: w,
        Data: mockedUsers
        StatusCode: htttp.StatusOK
    })
}
```

When the endpoint linked to the handler above is called, you should see the following JSON response.

> `NOTE` - Unlike the error use case, successful requests expect the `StatusCode` to be defined when creating a successful response. If you do not provide a status code, 200 will be assumed.

```JSON
{
    "data": [
        {"id": 1, "name": "John Doe"},
        {"id": 2, "name": "Sam Smith"}
    ]
}
```

## Response Types

There are currently four core response types supported by `reply`. They are the `Error`, `Token`, `Data` (*Success*) and `Default` response types. Each type has its JSON representation which is defined through a `Transfer Object`.

> NOTE: Unless otherwise stated, the `Transfer Object` assumed will be the [default transfer object (defaultReplyTransferObject)](./model.go)

### Universal Attributes

All core response types share universal attributes, which you can set in addition to their outputs. These include:
- Headers
- Meta
- Status Code

### Error Response Type

The `Error` response notifies the consumer when an error/ unexpected behaviour has occurred on the API. The message and the status code forwarded to the consumer is sourced from the error manifest. In the event the error's string
representation isn't in the manifest; `reply` will return the consumer a "500 - Internal Server Error" response.

<<<<<<< Updated upstream
=======
#### As code

To create an `error` response use the following code snippet:

```go
// create error manifest
baseManifest := []reply.ErrorManifest{
                {"example-404-error": reply.ErrorManifestItem{Message: "resource not found", StatusCode: http.StatusNotFound}},
            }

// create replier based on error manifest
replier := reply.NewReplier(baseManifest)

func ExampleHandler(w http.ResponseWriter, r *http.Request) {

    // error returned
    exampleErr := errors.New("example-404-error")

    _ := replier.NewHTTPResponse(&reply.NewResponseRequest{
        Writer: w,
        Error:  exampleErr,
    })
}
```

>>>>>>> Stashed changes
#### JSON Representation

`Error` responses are returned with the format.

```JSON
{
    "status": {
        "message": "resource not found"
    }
}
```

##### With `Meta`

When a `meta` is also declared, the response will have the following format. It can be as big or small as needed.

```JSON
{
    "meta": {
        "example": "meta in error reponse"
    },
    "status": {
        "message": "resource not found"
    }
}
```

### Token Response Type

The `token` response sends the consumer tokens; currently, the supported token types are `acccess_token` and `refresh_token`. If either is passed in the response request, `reply` will default to this response type.

<<<<<<< Updated upstream
=======
#### As code

To create a `token` response use the following code snippet:

```go
replier := reply.NewReplier([]reply.ErrorManifest{})

func ExampleHandler(w http.ResponseWriter, r *http.Request) {

    // do something to get tokens

    _ := replier.NewHTTPResponse(&reply.NewResponseRequest{
        Writer: w,
        AccessToken: "08a0a043-b532-4cea-8117-364739f2d994",
        RefreshToken: "08b29914-09a8-4a4a-8aa5-b1ffaff266e6",
        StatusCode: 200,
    })
}
```

>>>>>>> Stashed changes
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
    "meta": {
        "example": "meta in token reponse"
    },
    "access_token": "08a0a043-b532-4cea-8117-364739f2d994",
    "refresh_token": "08b29914-09a8-4a4a-8aa5-b1ffaff266e6"
}
```

### Data Response Type

The `data` response can be seen as a *successful* response. It parses the passed struct into its JSON representation and passes it to the consumer in the JSON response. The JSON response below will represent a response if the data passed was a user struct with the:
- `id` 1
- `name` john doe
<<<<<<< Updated upstream
- `dob` 1/1/1970 
=======
- `dob` 1/1/1970

#### As code

To create a `data` response use the following code snippet:

```go
type user struct {
    id int `json:"id"`
    name string `json:"name"`
    dob string `json:"dob"`
}

replier := reply.NewReplier([]reply.ErrorManifest{})

func ExampleHandler(w http.ResponseWriter, r *http.Request) {

    u := user{
        id: 1,
        name: "john doe",
        dob: "1/1/1970",
    }

    _ := replier.NewHTTPResponse(&reply.NewResponseRequest{
        Writer: w,
        Data: u,
        StatusCode: 201,
    })
}
```
>>>>>>> Stashed changes

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
    "meta": {
        "example": "meta in data reponse"
    },
     "data": {
        "id": 1,
        "name": "john doe",
        "dob": "1/1/1970"
    }
}
```

### Default Response Type

The `default` response returns `"{}"` with a status code of `200` if no `error`, `tokens`, `data` and `status code` is passed. If desired, another `status code` can be specified with `default` responses.

<<<<<<< Updated upstream
=======
#### As code

To create a `default` response use the following code snippet:

```go
replier := reply.NewReplier([]reply.ErrorManifest{})

func ExampleHandler(w http.ResponseWriter, r *http.Request) {

    _ := replier.NewHTTPResponse(&reply.NewResponseRequest{
        Writer: w,
    })
}
```

>>>>>>> Stashed changes
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
    "meta": {
        "example": "meta in default reponse"
    },
     "data": "{}"
}
```

## Copyright

Copyright (C) 2021 by Leon Silcott <leon@boasi.io>.

reply library released under MIT License.
See [LICENSE](https://github.com/ooaklee/reply/blob/master/LICENSE) for details.