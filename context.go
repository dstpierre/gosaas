package gosaas

type key int

const (
	// ContextOriginalPath holds the original requested URL.
	ContextOriginalPath key = iota
	// ContextRequestStart holds the request start time.
	ContextRequestStart
	// ContextDatabase holds a reference to a data.DB database connection and services.
	ContextDatabase
	// ContextAuth holds the authenticated user account id and user id.
	ContextAuth
	// ContextMinimumRole holds the minimum role to access this resource.
	ContextMinimumRole
	// ContextRequestID unique ID for the request.
	ContextRequestID
	// ContextRequestDump holds the request data dump.
	ContextRequestDump
	// ContextLanguage holds the request language.
	ContextLanguage
	// ContextContentIsJSON indicates if the request Content-Type is application/json
	ContextContentIsJSON
)
