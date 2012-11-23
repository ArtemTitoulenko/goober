package goober

import  (
  "net/http"
  "strings"
  "io"
  "time"
  "fmt"
)

// Main goober struct. Abides the handler interface.
type Goober struct {
  head map[string]*routeTreeNode
  ErrorPages map[int]string
}

// Goober handlers, for simplicity, are just functions with a given
// signature.
type Handler func(http.ResponseWriter, *Request)

// We use this a few places, so we can give it a type as well.
type RouteMap map[string]*routeTreeNode

// Our parse tree structure for routes
type routeTreeNode struct {
  // Handler if a node is a terminal
  handler Handler
  // Statis children
  children RouteMap
  // Dynamic/variable children
  variables RouteMap
}

// Augment http.Request with URLParams that will be grabbed
// from the request in the form of /:variables/
type Request struct {
  http.Request
  URLParams map[string]string
}

// A quick initializer for routeTreeNodes
func newRouteTreeNode() (node *routeTreeNode) {
  node = &routeTreeNode{
    children: make(RouteMap),
    variables: make(RouteMap),
  }

  return
}

// Initialize our Goober object
func New() (* Goober) {
  var head = make(RouteMap)
  head = head
  head["GET"] = newRouteTreeNode()
  head["HEAD"] = newRouteTreeNode()
  head["POST"] = newRouteTreeNode()
  head["PUT"] = newRouteTreeNode()
  head["DELETE"] = newRouteTreeNode()

  g := &Goober{
    head: head,
    ErrorPages: make(map[int]string),
  }

  return g
}

// Simple helper to allow us to trim leading and trailing /'s
func isSlash(s rune) (bool) {
  return s == '/'
}

// Adds a handler to our route tree
func (g *Goober) AddHandler(method string, route string, handler Handler) (err int){
  err = 0
  route = strings.TrimFunc(route, isSlash)
  var parts = strings.Split(route, "/")

  // Iterate through the bits of our path and add to the tree
  var cur = g.head[method]
  for i := range parts {
    var part = parts[i]

    // No // empty paths
    if (len(part) == 0) {
      err = 1
      return
    }

    // Check for variables
    if strings.HasPrefix(part, ":") {
      // dynamic
      if (cur.variables[part] != nil) {
        cur = cur.variables[part]
      } else {
        cur.variables[part] = newRouteTreeNode()
        cur = cur.variables[part]
      }
    } else {
      // static
      if (cur.children[part] != nil) {
        cur = cur.children[part]
      } else {
        cur.children[part] = newRouteTreeNode()
        cur = cur.children[part]
      }
    }
  }

  // add handler
  cur.handler = handler
  return
}

// Wrapper functions for common types of request
func (g *Goober) Get(route string, handler Handler) (int) {
  return g.AddHandler("GET", route, handler)
  return g.AddHandler("HEAD", route, handler)
}

func (g *Goober) Post(route string, handler Handler) (int) {
  return g.AddHandler("POST", route, handler)
}

func (g *Goober) Put(route string, handler Handler) (int) {
  return g.AddHandler("PUT", route, handler)
}

func (g *Goober) Delete(route string, handler Handler) (int) {
  return g.AddHandler("DELETE", route, handler)
}

func walkTree(node *routeTreeNode, parts []string, r *Request) (handler Handler, err int) {
  err = 0
  handler = nil

  if len(parts) == 0 {
    // if we've reached a terminal state, return handler
    handler = node.handler
  } else {
    // else, look for it
    var part = parts[0]

    if node.children[part] != nil {
      // check static routes first, they have priority
      return walkTree(node.children[part], parts[1:], r)
    } else {
      for k, v := range node.variables {
        // check all dynamic routes, taking first match
        handler, err = walkTree(v, parts[1:], r)
        if err == 0 {
          // goofy recursive way to build up params
          r.URLParams[k] = part
          return
        }
      }

      // if we don't find any dynamic matches, there was an error
      err = -1
    }
  }

  return
}

// Given a request, find the appropriate handler
func (g *Goober) GetHandler(r *Request) (handler Handler, err int) {
  var path = strings.TrimFunc(r.URL.Path, isSlash)
  var parts = strings.Split(path, "/")
  return walkTree(g.head[r.Method], parts, r)
}

// A simple function to handle error pages for us
func (g *Goober) errorHandler(w http.ResponseWriter, r *Request, code int) {
  if page, ok := g.ErrorPages[code]; ok {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(code)
    io.WriteString(w, page)
  }
}

// Routes requests
func (g *Goober) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  var startTime = time.Now()
  defer func() {
    fmt.Printf("[%s] %s - took %s\n", r.Method, r.URL.Path, time.Since(startTime))
  }()

  // create augmented request object
  var request = &Request{
    Request: *r,
    URLParams: make(map[string]string),
  }

  // get the handler for the request
  var f, err = g.GetHandler(request)
  if err == 0 && f != nil {
    // user response. pad with content-type.
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    f(w, request)
  } else {
    g.errorHandler(w, request, 404)
  }

}

