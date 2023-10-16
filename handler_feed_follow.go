package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/squashd/blog-aggregator/internal/database"
	"net/http"
	"time"
)

func (c *apiConfig) handleCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	feedFollow, err := c.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    params.FeedID,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		return
	}

	respondWithJSON(w, http.StatusCreated, feedFollow)
}

func (c *apiConfig) handleDeleteFeedFollow(w http.ResponseWriter, r *http.Request) {
	feedFollowId := chi.URLParam(r, "feedFollowId")
	id, err := uuid.Parse(feedFollowId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Not a valid feed follow id")
		return
	}
	err = c.DB.DeleteFeedFollow(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete feed follow")
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}

func (c *apiConfig) handleGetAllFeedsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	feeds, err := c.DB.GetFeedFollowsForUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get feeds")
		return
	}

	respondWithJSON(w, http.StatusOK, feeds)
}
