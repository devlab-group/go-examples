package main

import (
  "fmt"
  "net/http"
)

const (
  port = ":8080"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Hello, world!")
  fmt.Print("Request!")
}

func init() {
  fmt.Printf("Started server at http://localhost%v.\n", port)
  http.HandleFunc("/", HelloWorld)
  http.ListenAndServe(port, nil)
}

func main() {}
