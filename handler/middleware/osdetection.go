package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mileusna/useragent"
)

// type contextKey string

// const tokenContextKey contextKey = "OS"

func DetectOS(h http.Handler) http.Handler {
	fmt.Println("start OS detection")
	fn := func(w http.ResponseWriter, r *http.Request) {
		ua := useragent.Parse(r.UserAgent())
		ctx := SetOSVersion(r.Context(), ua.OS)
		//ctx = context.WithValue(ctx, "timestamp", time.Now())
		//ctx1 := AddStartTime(ctx)
		fmt.Println("OS is")
		fmt.Println(ctx)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func DetectOS1(h http.HandlerFunc) http.HandlerFunc {
	fmt.Println("start OS detection")
	fn := func(w http.ResponseWriter, r *http.Request) {
		ua := useragent.Parse(r.UserAgent())
		ctx := SetOSVersion(r.Context(), ua.OS)
		//ctx = context.WithValue(ctx, "timestamp", time.Now())
		//ctx1 := AddStartTime(ctx)
		fmt.Println("OS is")
		fmt.Println(ctx)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func SetOSVersion(ctx context.Context, t string) context.Context {
	fmt.Println("SetOSVersion function starts")
	return context.WithValue(ctx, "OS", t)
}

func GetOSVersion(ctx context.Context) (string, error) {
	v := ctx.Value("OS")
	token, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("token not found")
	}

	return token, nil
}
