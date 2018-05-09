package engine

type key int

const (
	// ContextOriginalPath holds the original request URL
	ContextOriginalPath key = iota
	// ContextRequestStart holds the request start time
	ContextRequestStart
	// ContextDatabase holds a reference to a data.DB database connection and services
	ContextDatabase
	// ContextUserID holds the user ID (this is just for demo)
	ContextUserID
)
