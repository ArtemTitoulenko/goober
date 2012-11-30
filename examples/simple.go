package main

import (
  "../../goober"
  "net/http"
  "io"
)

func hello(w http.ResponseWriter, g *goober.Request) {
  io.WriteString(w, "Hello world\n")
}

func main() {
  var g = goober.New()
  g.Get("/hello/:id", hello)
  g.ErrorPages[404] = "<h1>Not Found</h1>"

  g.ListenAndServe(":8080")
}

