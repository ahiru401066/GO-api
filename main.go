package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    return
  }
  fmt.Fprintln(w, "Hello go!")
}

func main() {
	http.HandleFunc("/", helloHandler)
  err := http.ListenAndServe(":8080", nil)
  if err!= nil{
    panic(err)
  }
}