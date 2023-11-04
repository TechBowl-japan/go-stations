package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
)

func Logging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// 次のハンドラーを呼び出す
		h.ServeHTTP(w, r)

		// リクエストの処理後のログ
		endTime := time.Now()
		latency := endTime.Sub(startTime).Milliseconds()
		os, ok := r.Context().Value(model.OSContextKey("os")).(string)
		if !ok {
			os = "unknown"
		}

		log := model.Log{
			Timestamp: startTime,
			Latency:   latency,
			Path:      r.URL.Path,
			OS:        os,
		}

		jsonOutput, err := json.Marshal(log)
		if err != nil {
			fmt.Printf("Error marshaling JSON: %v\n", err)
			return
		}

		fmt.Println(string(jsonOutput))
	})
}
