# Change history for @indexdata/foliogo

## [0.1.7](https://github.com/indexdata/foliogo/tree/v0.1.7) (2025-03-21)

* `RequestParams` structure now has new optional member `ContentType`. Fixes #20.

## [0.1.6](https://github.com/indexdata/foliogo/tree/v0.1.6) (2024-08-21)

* `NewService` and `NewDefaultSession` both accept optional final arguments, a pre-configured logger. Fixes #16.

## [0.1.5](https://github.com/indexdata/foliogo/tree/v0.1.5) (2023-12-06)

* Add `service.ResumeSession` (requires token to be passed in for each request). Fixes #14.

## [0.1.4](https://github.com/indexdata/foliogo/tree/v0.1.4) (2023-11-23)

* Repair "curl"-category logging to once more include payloads. Fixes #13.

## [0.1.3](https://github.com/indexdata/foliogo/tree/v0.1.3) (2023-11-22)

* Modify "curl"-category logging to include authentication header. Fixes #12.

## [0.1.2](https://github.com/indexdata/foliogo/tree/v0.1.2) (2023-11-22)

* Elements of `RequestParams` structure are now capitalised, for access from other packages.
* New method `session.Fetch0` simply invokes `Fetch` with an empty `RequestParams` structure.

## [0.1.1](https://github.com/indexdata/foliogo/tree/v0.1.1) (2023-11-21)

* Add new session method `GetTenant`.

## [0.1.0](https://github.com/indexdata/foliogo/tree/v0.1.0) (2023-11-19)

* BREAKING: `session.Fetch` returns a byte slice rather than a string-to-any map. Fixes #8.
* Do not specify an empty body to send when making a GET request. Fixes #7.

## [0.0.2](https://github.com/indexdata/foliogo/tree/v0.0.2) (2023-11-17)

* Add new top-level function `NewDefaultSession`. Fixes #6.

## [0.0.1](https://github.com/indexdata/foliogo/tree/v0.0.1) (2023-11-16)

* New module created from scratch, taking inspiration from the Node package [FolioJS](https://github.com/indexdata/foliojs)

