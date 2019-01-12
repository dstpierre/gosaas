<a href="https://buildsaasappingo.com/" title="Build a SaaS app in Go">
	<img src="https://buildsaasappingo.com/public/basaig.jpg" alt="Build a SaaS app in Go" align="right" height="250" />
</a>

# gosaas [![Documentation](https://godoc.org/github.com/dstpierre/gosaas?status.svg)](http://godoc.org/github.com/dstpierre/gosaas) [![Build Status](https://travis-ci.org/dstpierre/gosaas.svg?branch=master)](https://travis-ci.org/dstpierre/gosaas) [![Go Report Card](https://goreportcard.com/badge/github.com/dstpierre/gosaas)](https://goreportcard.com/report/github.com/dstpierre/gosaas)  [![Release](https://img.shields.io/github/release/dstpierre/gosaas.svg?label=Release)](https://github.com/dstpierre/gosaas/releases)

In September 2018 I published a book named [Build a SaaS app in Go](https://buildsaasappingo.com/). This project is the transformation of what the book teaches into a library that can be used to quickly build a web app / SaaS and focusing on your core product instead of common SaaS components.

*This is under development and API might change.*

### Usage quick example

You can create your main package and copy the `docker-compose.yml`. You'll need Redis and MongoDB for the library to work.

```go
package main

import (
	"net/http"
	"github.com/dstpierre/gosaas"
	"github.com/dstpierre/gosaas/model"
)

func main() {
	routes := make(map[string]*gosaas.Route)
	routes["test"] = &gosaas.Route{
		Logger:      true,
		MinimumRole: model.RolePublic,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gosaas.Respond(w, r, http.StatusOK, "hello world!")
		}),
	}

	mux := gosaas.NewServer(routes)
	http.ListenAndServe(":8080", mux)
}
```

## Table of content

* [Installation](#installation)
* [What's included](#what-s-included)
* [Quickstart](#quickstart)
  - [Defining routes](#defining-routes)
	- [How the database is handled](#database)
	- [Responding with JSON or HTML](#responding-to-requests)
	- [Parsing JSON body into types](#json-parsing)
	- [Getting current user and database from the request Context](#context)
* [Status and contributing](#status-and-contributing)
* [Credits](#credits)

## Installation

`go get github.com/dstpierre/gosaas/...`

## What's included

The following aspects are covered by this library:

* Web server capable of serving HTML templates, static files. Also JSON for an API.
* Easy helper functions for parsing and encoding type<->JSON.
* Routing logic in your own code.
* Middlewares: logging, authentication, rate limiting and throttling.
* User authentication and authorization using multiple ways to pass a token and a simple role based security.
* Database agnostic data layer. Currently handling MongoDB and an in-memory provider.
* User management, billing (per account or per user) and webhooks management.
* Simple queue (using Redis) and Pub/Sub for queuing tasks.
* Cron-like scheduling for recurring tasks.

## Quickstart

Here's some quick tips to get you up and running.

### Defining routes

You only need to pass the top level routes that gosaas needs to handle via a `map[string]*gosaas.Route`.

For example, if you have the following routes in your web application:

`/task, /task/mine, /task/done, /ping`

You would pass the following `map` to gosaas's `NewServer` function:

```go
routes := make(map[string]*gosaas.Route)
routes["task"] = &gosaas.Route{
	Logger: true,
	WithDB: true,
	handler: task,
	...
}
routes["ping"] = &gosaas.Route(
	Logger: true,
	Handler: ping,
)
```

Where `task` and `ping` are types that implement `http`'s `ServeHTTP` function, for instance:

```go
type Task struct{}

func (t *Task) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = gosaas.ShiftPath(r.URL.Path)
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
gosaas.Respond(w, r, http.StatusOK, oneTask)
```

**ServePage**: used to return HTML from templates:

```go
gosaas.ServePage(w, r, "template.html", data)
```

### JSON parsing

There a helper function called `gosaas.ParseBody` that handles the JSON decoding into types. This is a typical http handler:

```go
func (t Type) do(w http.ResponseWriter, r *http.Request) {
	var oneTask MyTask
	if err := gosaas.ParseBody(r.Body, &oneTask); err != nil {
		gosaas.Respond(w, r, http.StatusBadRequest, err)
		return
	}
	...
}
```

### Context

You'll most certainly need to get a reference back to the database and the currently 
logged in user. This is done via the request `Context`.

```go
func (t Type) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := ctx.Value(gosaas.ContextDatabase).(*data.DB)
	auth := ctx.Value(ContextAuth).(Auth)

	// you may use the db.Connection in your own data implementation
	tasks := Tasks{DB: db.Connection}
	list, err := tasks.List(auth.AccountID, auth.UserID)
	if err != nil {
		gosaas.Respond(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusOK, list)
}
```

## Status and contributing

I'm currently trying to reach a v1 and planning to use this in production with my next SaaS.

If you'd like to contribute I'd be more than happy to discuss, post an issue and feel free to 
explain what you'd like to add/change/remove.

Here's some aspect that are still a bit rough:

* Redis is **required** and cannot be changed easily, it's also coupled with the `queue` package.
* The controller for managing account/user is not done yet.
* The billing controller will need to be glued.
* The controllers package should be inside a `internal` package.
* Still not sure if the way the data package is written that it is idiomatic / easy to understand.

## Credits

Thanks for the following packages:

* [github.com/globalsign/mgo](https://github.com/globalsign/mgo)
* [github.com/go-redis/redis](https://github.com/go-redis/redis)
* [github.com/robfig/cron](https://github.com/robfig/cron)
* [github.com/satori/go.uuid](https://github.com/satori/go.uuid)
* [github.com/stripe/stripe-go](https://github.com/stripe/stripe-go)