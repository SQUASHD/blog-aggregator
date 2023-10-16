package main

import (
	"github.com/squashd/blog-aggregator/internal/database"
	"net/http"
)

func (c *apiConfig) handleGetPosts(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := c.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get posts")
		return
	}
	respondWithJSON(w, http.StatusOK, posts)
}
