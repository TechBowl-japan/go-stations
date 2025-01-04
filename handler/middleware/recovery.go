package middleware

import (
	"fmt"
	"log"
	"net/http"
)

func Recovery(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// TODO: ここに実装をする
		//関数が終了する直前に実行/defer:システムがダウンする前でも絶対に実行する
		//実行される順番は1defer(),2h.ServeHTTP(w, r),3if err := recover()…
		defer func() {
			//recover()：パニックが発生すればカバー、しなければnilを返す。
			if err := recover(); err != nil {
				// パニックの値をログに出力
				log.Printf("Recovered from panic: %v\n", err)
				//　エラーメッセージの作成
				http.Error(w, fmt.Sprintf("Internal Server Error: %v", err), http.StatusInternalServerError)
			}
		}() //即時実行関数：定義と同時に実行
		//ラップしているハンドラーを実行
		h.ServeHTTP(w, r)
		//下に書くとシャットダウン命令させちゃう　理由：事前に知っていない為、この下は処理されることはない。
		// if err := recover(); err != nil {
		// 	// パニックの値をログに出力
		// 	log.Printf("Recovered from panic: %v\n", err)
		// 	//　エラーメッセージの作成
		// 	http.Error(w, fmt.Sprintf("Internal Server Error: %v", err), http.StatusInternalServerError)
		// }
	}
	return http.HandlerFunc(fn)
}
