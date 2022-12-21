package endpoint

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func createComment(w http.ResponseWriter, r *http.Request) {
	//get the request body and decode it
	req := &CommentRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	//if there's an error with decoding the information
	//send a response with an error
	if err != nil {
		handleErr(w, err)
		return
	}
	//get the db from context
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	//if we can't get the db let's handle the error
	//and send an adequate response
	if !ok {
		handleDBFromContextErr(w)
		return
	}
	//if we can get the db then
	comment, err := models.CreateComment(pgdb, &models.Comment{
		Comment: req.Comment,
		UserID:  req.UserID,
	})
	if err != nil {
		handleErr(w, err)
		return
	}
	//everything is good
	//let's return a positive response
	succCommentResponse(comment, w)
}

func getComments(w http.ResponseWriter, r *http.Request) {
	//get db from ctx
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		handleDBFromContextErr(w)
		return
	}
	//call models package to access the database and return the comments
	comments, err := models.GetComments(pgdb)
	if err != nil {
		handleErr(w, err)
		return
	}
	//positive response
	res := &CommentsResponse{
		Success:  true,
		Error:    "",
		Comments: comments,
	}
	//encode the positive response to json and send it back
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("error encoding comments: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getCommentByID(w http.ResponseWriter, r *http.Request) {
	//get the id from the URL parameter
	//alternatively you could use a URL query
	commentID := chi.URLParam(r, "commentID")

	//get the db from ctx
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		handleDBFromContextErr(w)
		return
	}

	//get the comment from the DB
	comment, err := models.GetComment(pgdb, commentID)
	if err != nil {
		handleErr(w, err)
		return
	}

	//if the retrieval from the db was successful send the data
	succCommentResponse(comment, w)
}

func updateCommentByID(w http.ResponseWriter, r *http.Request) {
	//get the data from the request
	req := &CommentRequest{}
	//decode the data
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		handleErr(w, err)
		return
	}
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		handleDBFromContextErr(w)
		return
	}
	//get the commentID to know what comment to modify
	commentID := chi.URLParam(r, "commentID")
	//we get a string but we need to send an int so we convert it
	intCommentID, err := strconv.ParseInt(commentID, 10, 64)
	if err != nil {
		handleErr(w, err)
		return
	}

	//update the comment
	comment, err := models.UpdateComment(pgdb, &models.Comment{
		ID:      intCommentID,
		Comment: req.Comment,
		UserID:  req.UserID,
	})
	if err != nil {
		handleErr(w, err)
	}
	succCommentResponse(comment, w)
}

func deleteCommentByID(w http.ResponseWriter, r *http.Request) {
	//parse in the req body
	req := &CommentRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		handleErr(w, err)
		return
	}

	//get the db from ctx
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		handleDBFromContextErr(w)
		return
	}

	//get the commentID
	commentID := chi.URLParam(r, "commentID")
	intCommentID, err := strconv.ParseInt(commentID, 10, 64)
	if err != nil {
		handleErr(w, err)
		return
	}

	//delete comment
	err = models.DeleteComment(pgdb, intCommentID)
	if err != nil {
		handleErr(w, err)
	}

	//send successful response
	succCommentResponse(nil, w)
}
