package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// Source:
// - https://github.com/go-chi/chi
// - https://github.com/go-chi/chi/blob/master/_examples/rest/main.go

// Test data structure
type Item struct {
	id string `json:"id"`
}

type contextKey string

var (
	contextKeyID = contextKey("id")
)

// ErrNotFound error not found
var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

// ErrNoError no error
var ErrNoError = &ErrResponse{HTTPStatusCode: 200, StatusText: "OK"}

func main() {

	if _, err := InitDb(); err != nil {
		fmt.Println("FATAL error")
		os.Exit(1)
	}
	r := registerRoutes()
	http.ListenAndServe(":3333", r)
}

func registerRoutes() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(5 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	// RESTy routes for "/api/feature" resource
	r.Route("/api/fetcher", func(r chi.Router) {

		//r.With(paginate).Get("/", listArticles)             // GET /api/fetcher
		r.Get("/", listItems)   // GET /api/fetcher
		r.Post("/", createItem) // POST /api

		// Subrouters:
		r.Route("/{id}", func(r chi.Router) {
			r.Use(itemCtx)
			//r.Use(ArticleCtx)
			r.Delete("/", deleteItem)         // DELETE /api/fetcher/{id}
			r.Get("/history", getHistoryItem) // PUT /api/fetcher/{id}/history
		})
	})
	return r
}

// paginate is a stub, but very possible to implement middleware logic
// to handle the request params for handling a paginated request.
func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// just a stub.. some ideas are to look at URL query params for something like
		// the page number, or the limit, and send a query cursor down the chain
		next.ServeHTTP(w, r)
	})
}

func listItems(w http.ResponseWriter, r *http.Request) {

	render.Render(w, r, ErrNoError)
}

func getHistoryItem(w http.ResponseWriter, r *http.Request) {

	render.Render(w, r, ErrNoError)
}

func itemCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var idItem *Item
		var err error

		if idNum := chi.URLParam(r, "id"); idNum != "" {
			idItem, err = &Item{idNum}, nil
		}
		if err != nil {
			fmt.Fprintf(w, "ctx item: problem 1")
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyID, idItem)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func deleteItem(w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value(contextKeyID).(*Item)

	fmt.Fprintf(w, "delete item: {%s}", id)
	render.Render(w, r, ErrNoError)
}

func createItem(w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value(contextKeyID).(*Item)

	fmt.Fprintf(w, "create item: {%s}", id)
	render.Render(w, r, ErrNoError)
}

func getItem(w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value(contextKeyID).(*Item)

	fmt.Fprintf(w, "Get item: {%s}", id)
	render.Render(w, r, ErrNoError)
}
