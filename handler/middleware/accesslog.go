package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
)

func AccessLogOutput1(h http.Handler) http.Handler {
	fmt.Println("start AccessLogOutput1")
	fn := func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("access log pre process")
		now := time.Now()
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "timestamp", time.Now())))
		fmt.Println("access log post process")
		var logs model.AccessLog
		logs.OS = "FAIL"
		if nil != r.Context().Value("OS") {
			logs.OS = r.Context().Value("OS").(string)

		}
		fmt.Println(r.Context())
		logs.Path = r.URL.Path
		logs.Timestamp = time.Now()
		fmt.Println("1")
		//startTimes := r.Context().Value("timestamp").(time.Time)
		fmt.Println("2")
		//diff := logs.Timestamp.Sub(startTimes)
		//fmt.Println(diff)
		logs.Latency = time.Since(now).Milliseconds()
		//		logs.Latency = diff.Milliseconds()u
		fmt.Println("accesslog")
		data, err := json.Marshal(logs)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(data))
	}
	return http.HandlerFunc(fn)
}

func AccessLogOutput(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("Start ACcess Log Output")
	var logs model.AccessLog
	logs.OS = "FAIL"
	if nil != r.Context().Value("OS") {
		logs.OS = r.Context().Value("OS").(string)

	}
	fmt.Println(logs)
	logs.Path = r.URL.Path
	logs.Timestamp = time.Now()
	startTimes := r.Context().Value("timestamp").(time.Time)
	diff := logs.Timestamp.Sub(startTimes)
	fmt.Println(diff)
	logs.Latency = diff.Milliseconds()
	fmt.Println("accesslog")
	data, err := json.Marshal(logs)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(data))

	//fmt.Println(json.NewEncoder(w).Encode(&logs))
	return nil
}

func AddStartTime(ctx context.Context) context.Context {
	return context.WithValue(ctx, "timestamp", time.Now())
}

func AccessLogOutput2(next http.HandlerFunc) http.HandlerFunc {
	fmt.Println("start AccessLogOutput2")

	fn := func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("access log pre process")
		now := time.Now()
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "timestamp", time.Now())))
		fmt.Println("access log post process")
		var logs model.AccessLog
		logs.OS = "FAIL"
		if nil != r.Context().Value("OS") {
			logs.OS = r.Context().Value("OS").(string)

		}
		fmt.Println(r.Context())
		logs.Path = r.URL.Path
		logs.Timestamp = time.Now()
		fmt.Println("1")
		//startTimes := r.Context().Value("timestamp").(time.Time)
		fmt.Println("2")
		//diff := logs.Timestamp.Sub(startTimes)
		//fmt.Println(diff)
		logs.Latency = time.Since(now).Milliseconds()
		//		logs.Latency = diff.Milliseconds()u
		fmt.Println("accesslog")
		data, err := json.Marshal(logs)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(data))
	}
	return http.HandlerFunc(fn)
}
