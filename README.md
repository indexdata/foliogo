# FolioGo - FOLIO client library for Go

Copyright (C) 2023 Index Data Aps.

This software is distributed under the terms of the Apache License, Version 2.0. See the file "[LICENSE](LICENSE)" for more information.

<!-- md2toc -l 2 README.md -->
* [Overview](#overview)
* [API](#api)
* [Environment](#environment)
* [Logging](#logging)
* [Author](#author)


## Overview

FolioGo is a simple Go module to allow the creation of code that manipulate instances of [the FOLIO library services platform](https://www.folio.org/). For example, [a very simple program](bin/folio-list-users.go) to list the first 20 usernames, with asterisks next to the active ones, might read as follows:
```
package main

import "fmt"
import "github.com/indexdata/foliogo"

func main() {
	service := foliogo.NewService("https://folio-snapshot-okapi.dev.folio.org")
	session, _ := service.Login("diku", "user-basic-view", "user-basic-view")
	body, _ := session.Fetch("users?limit=20", foliogo.RequestParams{})
	fmt.Println(body)
}
```

This module is a port of [the FolioJS Node package](https://github.com/indexdata/foliojs) which does the same thing for JavaScript. We need the Go version to become part of [the FOLIO module mod-reporting](https://github.com/indexdata/mod-reporting), but it will likely have plenty of other uses.


## API

The API is described in a separate document, [The FolioGo API](doc/api.md).


## Environment

The behaviour of the FolioGo library can be modified by the values of the following environment variables:

* `LOGGING_CATEGORIES` or `LOGCAT` -- see [Logging](#logging) below.
* `FOLIOGO_SESSION_TIMEOUT` (or `FOLIOJS_SESSION_TIMEOUT`) -- if defined, the number of seconds after which a new session cookie will be requested. (If not defined, the default is to request a new cookie after half of the lifetime of the old one, which is typically about ten minutes.)


## Logging

This library uses the tiny but beautiful [`catlogger`](https://github.com/MikeTaylor/catlogger) library to provide optional logging. This is configured at run-time by the `LOGGING_CATEGORIES` or `LOGCAT` environment variable, which is set to a comma-separated list of categories such as `op,curl,status`. Messages in all the listed categories are logged.

Apart from categories used by `log` invocations in application code, the following categories are used by the libarary itself:
* `service`: log when a new service is created.
* `session`: log when a new session is created.
* `op`: whenever a high-level Okapi operation is about to be executed, its name and parameters are logged.
* `auth`: emits messages when authenticating or re-authenticating a session.
* `curl`: whenever an HTTP request is made, the equivalent `curl` command is logged. This can be useful for reproducing bugs.
* `status`: whenever an HTTP response is received, its HTTP status and content-type are logged. The combination of `op,status` is useful for tracing what a program is doing.
* `response`: whenever an HTTP response is received, its content is logged.

## Author

Mike Taylor (mike@indexdata.com),
for [Index Data Aps](https://www.indexdata.com/).



