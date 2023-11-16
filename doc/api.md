# The FolioGo API

Copyright (C) 2023 Index Data Aps.

<!-- md2toc -l 2 api.md -->
* [Introduction](#introduction)
* [Public types](#public-types)
* [Top level function](#top-level-function)
    * [foliogo.NewService(url)](#foliogonewserviceurl)
* [class `foliogo.Service`](#class-foliogoservice)
    * [service.Log(category string, args ...string)](#servicelogcategory-string-args-string)
    * [service.Login(tenant string, username string, password string)](#servicelogintenant-string-username-string-password-string)
* [class `foliogo.Session`](#class-foliogosession)
    * [session.Log(category string, args ...string)](#sessionlogcategory-string-args-string)
    * [session.Fetch(path string, params RequestParams)](#sessionfetchpath-string-params-requestparams)
* [Differences from FolioJS](#differences-from-foliojs)



## Introduction

FolioGo is a simple Go library to allow the creation of scripts that manipulate instances of [the FOLIO library services platform](https://www.folio.org/).

The API provides four types and a single exported function: that function creates an object with methods as described below.



## Public types

* `foliogo.Service` represents a service, as discussed below. Its structure is private, and service should only be accessed via its public API.
* `foliogo.Session` represents a session, as discussed below. Its structure is private, and service should only be accessed via its public API.
* `foliogo.RequestParams` represents a set of optional parameters that can be passed to `session.Fetch`, as discussed below.
* `foliogo.Hash` is simply an alias for `map[string]interface{}`, a mapping of strings to aribitrary data objects. It is the return type of `session.Fetch`, as discussed below.



## Top level function

A single top-level function is provided:


### foliogo.NewService(url)

Creates and returns a new `foliogo.Service` object associated with the specified Okapi URL. It is possible for a program to use multiple FOLIO services. See below for details of the `foliogo.Service` class.



## class `foliogo.Service`

Service objects are not created directly by client code, but by the `foliogo.NewService` factory function.

A service object is not associated with any particular tenant: for that, you need a session.

The following methods exist:


### service.Log(category string, args ...string)

Emits a log message in the specified category: see [the top-level `README.md`](../README.md#logging) for details.


### service.Login(tenant string, username string, password string)

Creates and a new `foliogo.Service` object, representing a session in the specified tenant of the service, logged in with the specified credentials (username and password). The session object retains the authentication token, and re-uses it for subsequent operations.

Returns a service object and an error indication. The latter is non-`nil` if an error has occurred.



## class `foliogo.Session`

Session objects are not created directly by client code, but by the `service.Login` factory function. Each session is permanently associated with a particular service, and permanently pertains to a particular tenant on that service. It is possible for a program to use sessions on the same or different FOLIO services.

The following public methods exist:


### session.Log(category string, args ...string)

Emits a log message in the specified category: see [the top-level `README.md`](../README.md#logging) for details.


### session.Fetch(path string, params RequestParams)

Performs an HTTP operation on the session, using an API much like that of [JavaScript `fetch`](https://developer.mozilla.org/en-US/docs/Web/API/fetch). The `path` is interpreted relative to the URL of the service that the session was created for, and should not begin with a slash (`/`). The `params` object can contain any subset of the following parameters:

* `body` (`string`) -- if provided, this content is sent to the HTTP service as the body of a POST or PUT.
* `json` (`interface{}`) -- if provided, this is serialised into a string and sent as though it had been provided as the `body`.
* `method` (`string`) -- specifies which HTTP method to use (GET, PUT, POST, DELETE, etc.). If this is not explicitly specified, and content is provided (as `body` or `json`) then it defaults to POST, otherwise to GET.

FOLIO authentication cookies are automatically included, along with the `X-Okapi-Tenant` header.

If content was provided as a `json` parmeter, then the `Content-type: application/json` header is added.

The value returned from a successful call is the parsed JSON of the response, expressed as a `foliogo.Hash` object, and a `nil` error object. If an error occurs, a non-`nil` error is returned.



## Differences from FolioJS

This library is based in part on [FolioJS](https://github.com/indexdata/foliojs), an analogous library for JavaScript/Node. Because it is written in Go, there are significant differences in how similar functionality is expressed:

* There is no single top-level object, just a top-level function.
* The `session.Fetch` and function is synchronous: concurrency can be implemented at the appliction level using goroutines.
* Because there is no re-authentication background thread, sessions need not be (and cannot be) closed.
* No exceptions are thrown. Non-2xx HTTP responses are returned as regular errors and must be checked for an handled by the caller.
* The data returned from `session.Fetch` is more cumbersome to handle than JavaScript's nice in-memory representation of JSON.

Because this library was was written to fulfil a specific purpose (use in[ the FOLIO module mod-reporting](https://github.com/indexdata/mod-reporting)), it lacks some high-level facilities provided by its ancestor:

* No default session setup from environment variables.
* No support for old-style (non-expiring) FOLIO authentication.
* No way to resume an existing session by providing a token.
* No access to tokens generated by the login process.
* No application-level methods (`postModule`, `addPermsToUser`, etx.)
* No support for Node module descriptors.



