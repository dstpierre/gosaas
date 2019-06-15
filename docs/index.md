---
layout: default
title: A Go library to build SaaS and web app
---

# User documentation for **gosaas** library

> This is a work in progress :). Please check the 
> [GoDoc](https://godoc.org/github.com/dstpierre/gosaas).

### What and why?

In September 2018 I published a book called "[Build a SaaS app in Go](https://buildsaasappingo.com)" after building 
two SaaS with Go and a React front-end. I thought that it might be helpful to extract common pieces into a reusable 
Go library.

If you want to have all the details and support the project, you can buy the book which still contains up-to-date 
information about building a SaaS in Go.

### Table of content

**Main concepts**

* [Configuration file](config.md)
* [Routing](routing.md)
* [Request/Response](req-resp.md)
* [Defining your handlers](handlers.md)

**Built-in middlewares and modules**

* [Database & migration](db.md)
* [Membership](membership.md)
* [Billing & subscription](billing.md)
* [Caching](caching.md)
* [Throttling & rate limit](limits.md)
* [Queue & background tasks](queue-tasks.md)