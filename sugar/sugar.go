package sugar

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type HandlerFunction func(*Context)
type MiddlewareFunction func(http.HandlerFunc) http.HandlerFunc

type Route struct {
	Path        string
	Method      string
	HandlerFunc HandlerFunction
}

type Middleware struct {
	Path string
	MiddlewareFunc MiddlewareFunction
}

type Sugar struct {
	Routes []Route
	Middlewares []Middleware
	NotFoundRoute HandlerFunction
}

type Context struct {
	request *http.Request
	writer http.ResponseWriter
}

func (c *Context) HTML(html string) {
	c.writer.Header().Set("Content-Type", "text/html")
	c.writer.Write([]byte(html))
}

func (c *Context) Redirect(url string) {
	http.Redirect(c.writer, c.request, url, http.StatusFound)
}

func (c *Context) JSON(content any) {
	c.writer.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(c.writer)
	enc.Encode(content)
}

func New(middlewares []Middleware) *Sugar {
	return &Sugar{
		Middlewares: middlewares,
	}
}

func (s *Sugar) Get(path string, handler HandlerFunction) {
	s.Routes = append(s.Routes, Route{
		Path: path,
		Method: http.MethodGet,
		HandlerFunc: handler,
	})
}

func (s *Sugar) Listen(port int) error {
	router := http.NewServeMux()
	for _, route := range s.Routes {
		handler := func(w http.ResponseWriter, r *http.Request) {
			if r.Method != route.Method {
				http.NotFound(w, r)
				return
			}
			ctx := Context{
				request: r,
				writer: w,
			}
			route.HandlerFunc(&ctx)
		}

		for _, handlerrr := range s.Middlewares {
			if handlerrr.Path == "*" {
				log.Println("Middleware added for path: " + route.Path)
				handler = handlerrr.MiddlewareFunc(handler)
			} else if handlerrr.Path == route.Path {
				log.Println("Middleware added for path: " + route.Path)
				handler = handlerrr.MiddlewareFunc(handler)
			}
		}

		router.HandleFunc(route.Path, handler)
	}

	err := http.ListenAndServe(":" + strconv.Itoa(port), router)
	return err
}