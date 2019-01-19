package gosaas

import (
	"encoding/json"
	"io"
	"net/http"
)

// Respond return an strruct with specific status as JSON.
//
// If data is an error it will be wrapped in a generic JSON object:
//
// 	{
// 		"status": 401,
// 		"error": "the result of data.Error()"
// 	}
//
// Example usage:
//
// 	func handler(w http.ResponseWriter, r *http.Request) {
// 		task := Task{ID: 123, Name: "My Task", Done: false}
// 		gosaas.Respond(w, r, http.StatusOK, task)
// 	}
func Respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) error {
	// change error into a real JSON serializable object
	if e, ok := data.(error); ok {
		var tmp = new(struct {
			Status string `json:"status"`
			Error  string `json:"error"`
		})
		tmp.Status = "error"
		tmp.Error = e.Error()
		data = tmp
	}

	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	// write the request ID
	reqID, ok := r.Context().Value(ContextRequestID).(string)
	if ok {
		w.Header().Set("X-Request-ID", reqID)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	logRequest(r, status)

	return nil
}

// ParseBody parses the request JSON body into a struct.ParseBody
//
// Example usage:
//
// 	func handler(w http.ResponseWriter, r *http.Request) {
// 		var task Task
// 		if err := gosaas.ParseBody(r.Body, &task); err != nil {
// 			gosaas.Respond(w, r, http.StatusBadRequest, err)
// 			return
// 		}
// 	}
func ParseBody(body io.ReadCloser, result interface{}) error {
	decoder := json.NewDecoder(body)
	return decoder.Decode(result)
}
