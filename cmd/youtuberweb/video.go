package main

import (
	"encoding/json"
	"net/http"
	"youtuber"
)

func (a *App) GetYoutubeHandler() HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {
		hello := youtuber.Youtube{String: "hello"}

		err := json.NewEncoder(w).Encode(hello)
		if err != nil {
			a.logr.Log("error when return json %s", err)
			return newAPIError(404, "error when return json %s", err)
		}
		return nil
	}
}
func (a *App) PostYoutubeHandler() HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {
		return nil
	}
}
