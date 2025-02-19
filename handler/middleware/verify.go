package middleware

//認証機能を作る
//アイパスはセキュリティ観点から環境変数を使う
//ID が BASIC_AUTH_USER_ID
//Password が BASIC_AUTH_PASSWORD
//制限したいHandlerの前にこの機能を設定する(todoの前に)
//アイパスが正しい場合は通過、正しくない時はHTTP Status Code を 401にし、処理終了する

import (
	"context"
	"fmt"
	"log"

	//"net/http"
	//"os"
	//"strings"

	//"github.com/joho/godotenv"
	"github.com/aws/aws-sdk-go-v2/aws"
	//"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

//func authentication(h http.HandlerFunc) http.Handler {
//basicauth := func(r *Request) basicAuth() (username, password string, ok bool) {

//err := godotenv.Load(".env")
//if err != nil {
//fmt.Println("Error loading .env file")
//}

//id := os.Getenv("BASIC_AUTH_USER_ID")
//pass := os.Getenv("BASIC_AUTH_PASSWORD")

//fmt.Println(id)
//fmt.Println(pass)
//}
//return http.HandleFunc(authentication)
//}

// Use this code snippet in your app.
// If you need more information about configurations or implementing the sample code, visit the AWS docs:
// https://aws.github.io/aws-sdk-go-v2/docs/getting-started/

func main() {
	secretName := "Authentication"
	region := "ap-northeast-1"

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		log.Fatal(err.Error())
	}

	// Decrypts secret using the associated KMS key.
	var secretString string = *result.SecretString


	// Your code goes here.
}




https://zoom.us/j/94144706773?pwd=Fpl2MVVrYLzGC2ZYPkpbiPtQ2IMh0c.1

ミーティング ID: 941 4470 6773
パスコード: 770665