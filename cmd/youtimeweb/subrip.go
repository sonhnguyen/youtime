package main

import (
	"fmt"
	"net/http"
	"time"
	"youtime"
)

func (a *App) GetSubtitleByIDHandler() HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {

		params := GetParamsObj(req)
		id := params.ByName("id")

		video, err := youtime.GetSubtitleByID(id, a.mongodb)
		if err != nil {
			a.logr.Log("error when return json %s", err)
			return newAPIError(404, "error when return json %s", err)
		}
		srtFile := []youtime.Subtitle{}
		for i, v := range video.Comment {
			var start time.Duration
			var end time.Duration
			start = time.Duration(v.Time) * time.Millisecond
			if i+1 < len(video.Comment) {
				end = time.Duration(video.Comment[i+1].Time) * time.Millisecond
			} else {
				end = start + 2000*time.Millisecond
			}
			if end-start > 5000*time.Millisecond {
				end = start + 5000*time.Millisecond
			}
			fmt.Println("%v %v %v", i, start, end)
			srtFile = append(srtFile, youtime.Subtitle{Number: i + 1, Text: v.Content, Start: start, End: end})
		}

		//copy the relevant headers. If you want to preserve the downloaded file name, extract it with go's url parser.
		w.Header().Set("Content-Disposition", "attachment; filename=subtitle.srt")

		err = youtime.WriteSubs(w, srtFile)
		if err != nil {
			a.logr.Log("error when making srt %s", err)
			return newAPIError(404, "error when making srt %s", err)
		}

		return nil
	}
}
