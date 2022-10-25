package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {

    http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello World")
    })

        // ポート = 8686
    if err := http.ListenAndServe(":8686", nil); err != nil {
        log.Fatal("ListenAndServe failed.", err)
    }
}