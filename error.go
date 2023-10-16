package main

import "net/http"

func (c *apiConfig) ErrorHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, http.StatusOK, response{Error: "Internal Server Error"})
}
