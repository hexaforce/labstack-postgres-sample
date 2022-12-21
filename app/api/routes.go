package api

import (
	"encoding/json"
	"go-microservice-example/pkg/db/models"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-pg/pg/v10"
)

// start api with the pgdb and return a chi router
func StartAPI(pgdb *pg.DB) *chi.Mux {
	//get the router
	r := chi.NewRouter()
	//add middleware
	//in this case we will store our DB to use it later
	r.Use(middleware.Logger, middleware.WithValue("DB", pgdb))

	//routes for our service
	r.Route("/comments", func(r chi.Router) {
		r.Post("/", createComment)
		r.Get("/", getComments)
		r.Get("/{commentID}", getCommentByID)
		r.Put("/{commentID}", updateCommentByID)
		r.Delete("/{commentID}", deleteCommentByID)
	})

	//test route to make sure everything works
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("up and running"))
	})

	return r
}

//-- UTILS --

func handleErr(w http.ResponseWriter, err error) {
	res := &CommentResponse{
		Success: false,
		Error:   err.Error(),
		Comment: nil,
	}
	err = json.NewEncoder(w).Encode(res)
	//if there's an error with encoding handle it
	if err != nil {
		log.Printf("error sending response %v\n", err)
	}
	//return a bad request and exist the function
	w.WriteHeader(http.StatusBadRequest)
}

func handleDBFromContextErr(w http.ResponseWriter) {
	res := &CommentResponse{
		Success: false,
		Error:   "could not get the DB from context",
		Comment: nil,
	}
	err := json.NewEncoder(w).Encode(res)
	//if there's an error with encoding handle it
	if err != nil {
		log.Printf("error sending response %v\n", err)
	}
	//return a bad request and exist the function
	w.WriteHeader(http.StatusBadRequest)
}

func succCommentResponse(comment *models.Comment, w http.ResponseWriter) {
	//return successful response
	res := &CommentResponse{
		Success: true,
		Error:   "",
		Comment: comment,
	}
	//send the encoded response to responsewriter
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("error encoding comment: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//send a 200 response
	w.WriteHeader(http.StatusOK)
}

// -- handle routes
