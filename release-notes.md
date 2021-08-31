# Reply Release Notes

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