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

// osのcontext keyを設定
const (
	osKey contextKey = "os"
)

func UserOs(h http.HandlerFunc) http.Handler {
	userAgent := func(w http.ResponseWriter, r *http.Request) {
		//リクエスト時間を取得
		requestTime := time.Now()
		log.Println(requestTime)

		//UserAgent使ってHTTPリクエストから情報を分析(OS情報を取得するため)
		uaString := r.UserAgent()
		ua := useragent.Parse(uaString)

		//contextにOSとリクエスト時間を格納
		ctx := r.Context()
		ctx = context.WithValue(ctx, "os", ua.OS)
		r = r.WithContext(ctx)

		//Handler処理
		h.ServeHTTP(w, r)

		//処理時間(ミリ単位)
		latency := time.Since(requestTime).Milliseconds()

		//取得したOS情報の型変換を行い、かつ取得できなかったときの判定処理
		osInformation, ok := ctx.Value("os").(string)
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
