package main

//認証機能を作る
//アイパスはセキュリティ観点から環境変数を使う
//ID が BASIC_AUTH_USER_ID
//Password が BASIC_AUTH_PASSWORD
//制限したいHandlerの前にこの機能を設定する(todoの前に)
//アイパスが正しい場合は通過、正しくない時はHTTP Status Code を 401にし、処理終了する

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func authentication() {
	err := godotenv.Load(".evn")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	id := os.Getenv("BASIC_AUTH_USER_ID")
	pass := os.Getenv("BASIC_AUTH_PASSWORD")

	fmt.Println(id)
	fmt.Println(pass)
}
