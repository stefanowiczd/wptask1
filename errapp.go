package main

import (
	"net/http"

	"github.com/go-chi/render"
)

//var m map[string]int = map[int]string{200: "", 400: "Invalid request"

// ErrResponse is a type to represent API errors and corresponding message
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// Render prepares an error
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrInvalidRequest is a plaveholder to keep detailed info about error returned by the application
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

// ErrRender is a plaveholder to keep detailed info about error returned by the application
func ErrRender(err error, errCode int, errStatus string) render.Renderer {

	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: errCode,
		StatusText:     errStatus,
		ErrorText:      err.Error(),
	}
}
