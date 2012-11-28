package main

import (
  "../../goober"
  "net/http"
  "io"
  "log"
  "strings"
)

func helloWorld(w http.ResponseWriter, r *goober.Request) {
  io.WriteString(w, "Hello, world.\n")
}

func helloAnyone(w http.ResponseWriter, r *goober.Request) {
  io.WriteString(w, "Hello, " + r.URLParams[":name"] + ".<br>\n")
  for k, v := range r.URL.Query() {
    io.WriteString(w, k + ": " + strings.Join(v, ", ") + "<br>\n")
  }
}

func static(w http.ResponseWriter, r *goober.Request) {
  var fileName = "./" + r.URLParams["*"]
  http.ServeFile(w, &r.Request, fileName)
}

func helloPreFunc (w http.ResponseWriter, r *goober.Request) (error) {
  var err error = nil
  if r.URLParams[":name"] != "Cole" {
    return &SomeError{}
  }
  return err
}

func helloPostFunc (w http.ResponseWriter, r *goober.Request) (error) {
  log.Println("Successfully served hello :name request.")
  return nil
}

type SomeError struct {}
func (e *SomeError) Error() string {
  return "There was an error."
}

func main() {
  var g = goober.New()
  g.Get("/hello", helloWorld)
  g.Get("/hello/:name", helloAnyone).
    AddPreFunc(helloPreFunc).
    AddPostFunc(helloPostFunc)
  g.Get("/assets/*", static)
  g.ErrorPages[404] = "<h1>Not found.</h1>"

  err := g.ListenAndServe(":8080")

  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}

