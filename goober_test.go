package goober

import (
  "testing"
  "net/http"
  "io"
)

func TestGoober_New(test *testing.T) {
  g := New()

  if g == nil || g.head == nil || g.ErrorPages == nil {
    test.FailNow()
  }
}

func Test_newRouteTreeNode(test *testing.T) {
  nrtn := newRouteTreeNode()

  if nrtn.children == nil || nrtn.variables == nil {
    test.FailNow()
  }
}

func hello_world_handler(w http.ResponseWriter, r *Request) {
  io.WriteString(w, "Hello, world.")
}

func Test_Goober_Get(test *testing.T) {
  g := New()
  g.Get("/hello", hello_world_handler)

  if g.head["GET"].children["hello"].handler == nil {
    test.FailNow()
  }
}

func Test_Goober_Post(test *testing.T) {
  g := New()
  g.Post("/hello", hello_world_handler)

  if g.head["POST"].children["hello"].handler == nil {
    test.FailNow()
  }
}

func Test_Goober_Put(test *testing.T) {
  g := New()
  g.Put("/hello", hello_world_handler)

  if g.head["PUT"].children["hello"].handler == nil {
    test.FailNow()
  }
}

func Test_Goober_Delete(test *testing.T) {
  g := New()
  g.Delete("/hello", hello_world_handler)

  if g.head["DELETE"].children["hello"].handler == nil {
    test.FailNow()
  }
}

// unit testing code from web.go, they have good stuff
func buildTestRequest(method string, path string, body string, headers map[string][]string, cookies []*http.Cookie) *http.Request {
    host := "127.0.0.1"
    port := "80"
    rawurl := "http://" + host + ":" + port + path
    url_, _ := url.Parse(rawurl)
    proto := "HTTP/1.1"

    if headers == nil {
        headers = map[string][]string{}
    }

    headers["User-Agent"] = []string{"web.go test"}
    if method == "POST" {
        headers["Content-Length"] = []string{fmt.Sprintf("%d", len(body))}
        if headers["Content-Type"] == nil {
            headers["Content-Type"] = []string{"text/plain"}
        }
    }

    req := http.Request{Method: method,
        URL:    url_,
        Proto:  proto,
        Host:   host,
        Header: http.Header(headers),
        Body:   ioutil.NopCloser(bytes.NewBufferString(body)),
    }

    for _, cookie := range cookies {
        req.AddCookie(cookie)
    }
    return &req
}
