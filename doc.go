/*
Package gosaas contains helper functions, middlewares, user management and billing functionalities commonly used in
typical Software as a Service web application.

The primary goal of this library is to handle repetitive components letting you focus on the core part of your project.

You use the NewServer function to get a working server MUX. You need to pass the top level routes to the NewServer
function to get the initial routing working.

For instance if your web application handles the following routes:

	/task
	/task/mine
	/task/done
	/ping
	/ping/stat

You only pass the "task" and "ping" routes to the server. Anything after the
top-level will be handled by your code. You will be interested in ShiftPath,
Respond, ParseBody and ServePage functions to get started.

The most important aspect of a route is the Handler field which corresponds to the code to execute. The Handler is a
standard http.Handler meaning that your code will need to implement the ServeHTTP function.

The remaining fields for a route control if specific middlewares are part of the request life-cycle or not. For
instance, the Logger flag will output request information to stdout when enabled.
*/
package gosaas
