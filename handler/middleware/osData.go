package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/mileusna/useragent"
)

type contextKey string //新しいstring型を定義

const osContextKey = contextKey("os") //("os")同じ文字列キーの衝突回避andOS情報
// MiddlewareでHTTPリクエストがサーバーに届く前に共通処理開始
// next http.Handlerは次の処理、HTTPリクエストを処理するためのhttp.Handlerは型を返却
func OSContextMiddleware(next http.Handler) http.Handler {
	//http.HandlerFunc 型に変換
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//UserAgent解析
		ua := useragent.Parse(r.UserAgent())
		osName := ua.OS //windows only
		//WithValueで「新しいContext作成」（親となるContext,新しい値を格納キー,格納する値）
		ctx := context.WithValue(r.Context(), osContextKey, osName)
		//OS情報が追加された新しいリクエストを返却
		//次に実行する処理（Handler)呼び出す
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func AccessLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//スタート時刻
		startTime := time.Now()
		//OS情報をContextから取得
		osName := GetOSFromContext(r.Context()) //65行目以降の処理の結果を返している。
		time.Sleep(10 * time.Millisecond)
		//handler実行
		next.ServeHTTP(w, r)
		//アクセス日時と処理後の処理時間の差分（ミリ秒）
		latency := time.Since(startTime).Milliseconds()
		//type Requestlog structに代入
		logData := model.Requestlog{
			Timestamp: startTime,
			Latency:   latency,
			Path:      r.URL.Path,
			OS:        osName,
		}
		//structからJSONに変換
		logJSON, err := json.Marshal(logData)
		if err != nil {
			fmt.Println("Failed to encode log data:", err)
			return
		}
		//stringにキャストして標準出力に表示
		fmt.Println(string(logJSON))

	})
}

// Contextに格納されたOS情報を取り出す。
func GetOSFromContext(ctx context.Context) string {
	//ctx.ValueでOS情報を取得し、string型に変更
	if osName, ok := ctx.Value(osContextKey).(string); ok {
		return osName
	}
	//格納されていなければunknown
	return "unknown"
}
