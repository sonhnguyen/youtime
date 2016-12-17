package main

import (
	"encoding/json"
	"net/http"
	"time"
	"youtime"
)

func (a *App) GetVideoByLinkHandler() HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {

		queryValues := req.URL.Query()
		site := queryValues.Get("site")
		link := queryValues.Get("link")

		video, err := youtime.GetVideoByLink(site, link, a.mongodb)
		if err != nil {
			a.logr.Log("error when return json %s", err)
			return newAPIError(404, "error when return json %s", err)
		}

		err = json.NewEncoder(w).Encode(video)
		if err != nil {
			a.logr.Log("error when return json %s", err)
			return newAPIError(404, "error when return json %s", err)
		}
		return nil
	}
}

func (a *App) GetVideoByIdHandler() HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {
		params := GetParamsObj(req)
		id := params.ByName("id")

		video, err := youtime.GetVideoById(id, a.mongodb)
		if err != nil {
			a.logr.Log("error when return json %s", err)
			return newAPIError(404, "error when return json %s", err)
		}

		err = json.NewEncoder(w).Encode(video)
		if err != nil {
			a.logr.Log("error when return json %s", err)
			return newAPIError(404, "error when return json %s", err)
		}

		return nil
	}
}

func (a *App) PostCommentByIdHandler() HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {
		var comment youtime.Comment

		err := json.NewDecoder(req.Body).Decode(&comment)
		comment.DateCreated = time.Now().UTC()
		if err != nil {
			a.logr.Log("error decode param: %s", err)
			return newAPIError(400, "error param: %s", err)
		}

		params := GetParamsObj(req)
		id := params.ByName("id")

		err = youtime.PostCommentById(id, comment, a.mongodb)
		if err != nil {
			a.logr.Log("error when return json %s", err)
			return newAPIError(404, "error when return json %s", err)
		}

		return nil
	}
}
