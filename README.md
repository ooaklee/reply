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
// Create a error manifest, to hold correlating error as string and it's manifest
// item
baseManifest := []reply.ErrorManifest{
                {"example-404-error": reply.ErrorManifestItem{Message: "resource not found", StatusCode: http.StatusNotFound}},
            }

// Create replier to manage the responses going back to consumer(s)
replier := reply.NewReplier(baseManifest)
```

### How to send response(s) 

You can use `reply` for both successful and error based responses.

> `NOTE` - When sending an error response, it is essential to make sure you populate the `replier`'s error manifest with the correct errors. Otherwise, a `500 - Internal Server Error` response will be sent back to the client by default if it cannot match the passed error with on in the manifest.

#### Making use of error manifest

```go

// ExampleHandler handler to demostrate how to use package for error
// response
func ExampleHandler(w http.ResponseWriter, r *http.Request) {

    // Create error with value corresponding to one of the manifest's entry's key
    exampleErr := errors.New("example-404-error")


    // Pass error to replier's method to return predefined response, else
    // 500
    replier.NewHTTPResponse(&reply.NewResponseRequest{
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
    replier.NewHTTPResponse(&reply.NewResponseRequest{
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