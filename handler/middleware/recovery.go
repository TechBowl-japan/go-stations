package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mileusna/useragent"
)

func Recovery(h http.HandlerFunc) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rp := recover(); rp != nil {
				log.Printf("Recovered from panic: %v", rp)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			log.Println("リカバーできました。")
		}()

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

type contextKey string

// osのcontext　keyを設定
const (
	osKey contextKey = "os"
)

func UserOs(h http.HandlerFunc) http.Handler {
	userAgent := func(w http.ResponseWriter, r *http.Request) {
		//リクエスト時間を取得
		requestTime := time.Now()
		log.Println(requestTime)

		//useragentからOS名を取得
		uaString := useragent.Parse(r.UserAgent())
		//ua := uaString.OS()

		//contextにOSとリクエスト時間を格納
		ctx := r.Context()
		ctx = context.WithValue(ctx, "os", uaString.OS)
		r = r.WithContext(ctx)

		//Handler処理
		h.ServeHTTP(w, r)

		//処理時間(ミリ単位)
		latency := time.Since(requestTime).Milliseconds()

		osInformation, ok := r.Context().Value(osKey).(string)
		if !ok {
			osInformation = "OS不明"
		}

		//構造体に値を格納
		metadata := RecoveryRequest{
			Timestamp: requestTime,
			Latency:   latency,
			Path:      r.URL.Path,
			OS:        osInformation,
		}

		//JSONに変換して出力
		jsonOutput, _ := json.Marshal(metadata)
		fmt.Println(string(jsonOutput))
	}
	return http.HandlerFunc(userAgent)
}
