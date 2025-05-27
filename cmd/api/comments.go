package main

import "net/http"

type CommentPayload struct {
	UserID  int    `json:"user_id"`
	PostID  int    `json:"post_id"`
	Content string `json:"content"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFormCtx(r)

	var payload CommentPayload

	payload.PostID = post.ID

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.Comments.Create(r.Context(), payload.UserID, payload.PostID, payload.Content); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, 201, "comment created successfully")
}
