package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
)

func DisplayAccessLog(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		timeStamp := time.Now()
		path := r.URL.Path
		os := GetOSFromContext(r.Context())

		h.ServeHTTP(w, r)

		latency := time.Since(timeStamp).Milliseconds()

		log := model.AccessLog{
			Timestamp: timeStamp,
			Latency:   latency,
			Path:      path,
			OS:        os,
		}

		logJson, _ := json.Marshal(log)
		fmt.Println(string(logJson))

	}
	return http.HandlerFunc(fn)
}
