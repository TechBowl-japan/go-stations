package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
)

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
	fmt.Println(json.Marshal(logs))
	data, err := json.Marshal(logs)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(data)) // {"age":30,"name":"Tanaka"}

	//fmt.Println(json.NewEncoder(w).Encode(&logs))
	return nil
}

func AddStartTime(ctx context.Context) context.Context {
	return context.WithValue(ctx, "timestamp", time.Now())
}
