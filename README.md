<a href="https://buildsaasappingo.com/" title="Build a SaaS app in Go">
	<img src="https://buildsaasappingo.com/public/basaig.jpg" alt="Build a SaaS app in Go" align="right" />
</a>

# gosaas [![Documentation](https://godoc.org/github.com/dstpierre/gosaas?status.svg)](http://godoc.org/github.com/dstpierre/gosaas) [![Build Status](https://travis-ci.org/dstpierre/gosaas.svg?branch=master)](https://travis-ci.org/dstpierre/gosaas) [![Go Report Card](https://goreportcard.com/badge/github.com/dstpierre/gosaas)](https://goreportcard.com/report/github.com/dstpierre/gosaas)  [![Coverage Status](https://coveralls.io/repos/github/dstpierre/gosaas/badge.svg?branch=master)](https://coveralls.io/github/dstpierre/gosaas?branch=master) [![GitHub issues](https://img.shields.io/github/issues/dstpierre/gosaas.svg)](https://github.com/dstpierre/gosaas/issues) [![license](https://img.shields.io/github/license/dstpierre/gosaas.svg?maxAge=2592000)](https://github.com/dstpierre/gosaas/LICENSE) [![Release](https://img.shields.io/github/release/dstpierre/gosaas.svg?label=Release)](https://github.com/dstpierre/gosaas/releases)

In September 2018 I published a book named [Build a SaaS app in Go](https://buildsaasappingo.com/). This project is the transformation of what the book teaches into a library that can be used to quickly build a web app / SaaS and focusing on your core product instead of common SaaS components.

*This is under development.*

### Usage quick example

```go
package main

import (
	"net/http"
	"github.com/dstpierre/gosaas"
	"github.com/dstpierre/gosaas/engine"
	"github.com/dstpierre/gosaas/data/model"
)

func main() {
	routes := make(map[string]*engine.Route)
	routes["test"] = &engine.Route{
		Logger:      true,
		MinimumRole: model.RolePublic,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			engine.Respond(w, r, http.StatusOK, "hello world!")
		}),
	}

	mux := gosaas.NewServer(routes)
	http.ListenAndServe(":8080", mux)
}
```

## Table of content

- [Installation](#installation)
- [Quickstart](#quickstart)
  - [Defining routes](#defining-routes)
	- [How the database is handled](#database)
	- [Responding with JSON or HTML](#responding-to-requests)
	- [Parsing JSON body into types](#json-parsing)

## Installation

`go get github.com/dstpierre/gosaas/...`

## Quickstart

Here's some quick tips to get you up and running.

### Defining routes

You only need to pass the top level routes that gosaas needs to handle via a `map[string]*engine.Route`.

For example, if you have the following routes in your web application:

`/task, /task/mine, /task/done, /ping`

You would pass the following `map` to gosaas's `NewServer` function:

```go
routes := make(map[string]*engine.Route)
routes["task"] = &engine.Route{
	Logger: true,
	WithDB: true,
	handler: task,
	...
}
routes["ping"] = &engine.Route(
	Logger: true,
	Handler: ping,
)
```

Where `task` and `ping` are types that implement `http`'s `ServeHTTP` function, for instance:

```go
type Task struct{}

func (t *Task) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = engine.ShiftPath(r.URL.Path)
	if head =="/" {
		t.list(w, r)
	} else if head == "mine" {
		t.mine(w, r)
	}
	...
}
```

You may define `Task` in its own package or inside your `main` package.

### Database

The `data` package is provider agnostic. For now, only a MongoDB and an in-memory implementation are available (Postgres will be implemented next). There's a `data.DB` type that abstract away the database proprietary code.

Before calling `http.ListenAndServe` you have to initialize the `DB` field of the `Server` type:

```go
db := &data.DB{}

if err := db.Open(*dn, *ds); err != nil {
	log.Fatal("unable to connect to the database:", err)
}

mux.DB = db
```

Where `*dn` and `*ds` are flags containing "mongo" and "localhost" respectively which are the driver name and the datasource connection string.

### Responding to requests

The `engine` package exposes two useful functions:

**Respond**: used to return JSON:

```go
engine.Respond(w, r, http.StatusOK, oneTask)
```

**ServePage**: used to return HTML from templates:

```go
engine.ServePage(w, r, "template.html", data)
```

### JSON parsing

There a helper function called `engine.ParseBody` that handles the JSON decoding into types. This is a typical http handler:

```go
func (t Type) do(w http.ResponseWriter, r *http.Request) {
	var oneTask MyTask
	if err := engine.ParseBody(r.Body, &oneTask); err != nil {
		engine.Respond(w, r, http.StatusBadRequest, err)
		return
	}
	...
}
```

