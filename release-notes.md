# Reply Release Notes

## [v2.0.0](https://github.com/ooaklee/reply/releases/tag/v2.0.0)
2025-08-03

* Updated module to `github.com/ooaklee/reply/v2`

## [v2.0.0-alpha](https://github.com/ooaklee/reply/releases/tag/v2.0.0-alpha)
2025-07-15

* Updated `ErrorManifest` to use `error` as the key instead of `string`
* Updated `replier.go` to reflect the changes to `ErrorManifest` and use `errors.Is` for error matching
* Updated `README.md` to reflect the changes
* Reflected changes in `examples/example_simple_api.go` and provide examples of how to use it with errors

## [v1.1.0](https://github.com/ooaklee/reply/releases/tag/v1.1.0)
2025-07-13

* Updated `go.mod` to go 1.20
* Updated `NewHTTPErrorResponse` to behave like `NewHTTPMultiErrorResponse` when manifest errors are wrapped, i.e. `errors.Join(errs...)`
* Updated `README.md` to reflect the changes to `NewHTTPErrorResponse` and provide an example of how to use it with wrapped errors
* Added a new test case to verify the changes

## [v1.0.0](https://github.com/ooaklee/reply/releases/tag/v1.0.0)
2021-09-11

* Update top-level response members from `meta`, `status` and `data` **->** `meta`, `errors` and `data`
* Updated `Manifest Error Item` fields
* Updated logic & added new aide `NewHTTPMultiErrorResponse` to support multiple error response objects in response
* Refactored code
* Added logic to support users passing custom *error transfer objects* (`TransferObjectError`), `reply.WithTransferObjectError`
  * Added `RefreshTransferObject` method call to `TransferObjectError`
* Updated Token attributes & methods on base `Transfer Object` to have a generic name to cover cases where API uses different token identifiers
  * `AccessToken` -> `TokenOne` & `RefreshToken` -> `TokenTwo`
* Added logic to set `StatusCode` if not set already. Defaults to `400`.
  
## [v1.0.0-alpha.3](https://github.com/ooaklee/reply/releases/tag/v1.0.0-alpha.3)
2021-09-11

* Refactor Token attributes & methods to have a more general name to make use of cases where API uses different token identifiers less confusing
  * `AccessToken` -> `TokenOne` & `RefreshToken` -> `TokenTwo`
    * Dev can create custom `TransferObject` to set JSON attribute to match their use case
* Added logic to set `StatusCode` if not set already. Defaults to `400`.
* Updated README to better describe the library

## [v1.0.0-alpha.2](https://github.com/ooaklee/reply/releases/tag/v1.0.0-alpha.2)
2021-09-10

* Refactored logic to allow users to pass custom error transfer objects (TransferObjectError), `reply.WithTransferObjectError`
  * Added `RefreshTransferObject` method call to `TransferObjectError`
* Updated `example simple api` to use the new Replier Option and update the `transfer object error`

## [v1.0.0-alpha.1](https://github.com/ooaklee/reply/releases/tag/v1.0.0-alpha.1)
2021-09-10

* Added new aide `NewHTTPMultiErrorResponse` to support multiple error response
* Updated `example simple api` to use the new aide `NewHTTPMultiErrorResponse`
* Added logic to handle/ create multiple error response
* Refactor code to make it more readable with new logic

## [v1.0.0-alpha](https://github.com/ooaklee/reply/releases/tag/v1.0.0-alpha)
2021-09-09

* Update top-level response members from `meta`, `status` and `data` **->** `meta`, `errors` and `data`
* Updated underlying logic to how an error is handled
* Updated `Manifest Error Item` attributes

## [v0.2.0](https://github.com/ooaklee/reply/releases/tag/v0.2.0)
2021-09-04

* Fixed bug in logic for merging error manifests
* Added helper functions (aides) to help users more efficiently use the library

## [v0.2.0-alpha.1](https://github.com/ooaklee/reply/releases/tag/v0.2.0-alpha.1)
2021-09-04

* Fixed bug in logic for merging error manifests

## [v0.2.0-alpha](https://github.com/ooaklee/reply/releases/tag/v0.2.0-alpha)
2021-09-03

* Initial logic for helper functions (`aides`) to help users more efficiently use the library

## [v0.1.1](https://github.com/ooaklee/reply/releases/tag/v0.1.1)
2021-08-31

* Added licensing information
* Fixed typos

## [v0.1.0](https://github.com/ooaklee/reply/releases/tag/v0.1.0)
2021-08-28

* Refactored code
* Utilise error returned in `NewHTTPResponse`
* Added log entry for unfound manifest error item.

## [v0.1.0-alpha.1](https://github.com/ooaklee/reply/releases/tag/v0.1.0-alpha.1)
2021-08-27

* Updated logic for the setting headers.
* Added logic to default to JSON `content-type`

## [v0.1.0-alpha](https://github.com/ooaklee/reply/releases/tag/v0.1.0-alpha)
2021-08-26

* Enforced `response priority`, priority is declared as follows:
  - *Error response*
  - *Token response* (access token/ refresh token)
  - *Data response*
  - *Default response*
> Note - *Default response* occurs if none of the attributes expected to satisfy the other response types is satisfied. The JSON response for the default response returned will be:
```json
{
    "data": "{}"
}
```
* Added ability to send responses based on singular error. 
> NOTE: If the error does not exist in the `error manifest` reply will default to 500 - Internal Server Error