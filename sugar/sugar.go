package sugar

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
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
	DB *pgxpool.Pool
	Routes []Route
	Middlewares []Middleware
	NotFoundRoute HandlerFunction
}

type Context struct {
	DB *pgxpool.Pool
	Request *http.Request
	writer http.ResponseWriter
}

type URL struct {
	Path string
	Query map[string][]string
}

func (c *Context) HTML(html string) {
	c.writer.Header().Set("Content-Type", "text/html")
	c.writer.Write([]byte(html))
}

func (c *Context) Page(data any, filenames ...string) {
	if reflect.TypeOf(data).Kind() != reflect.Struct {
		log.Fatal("Only Structs")
	}
	jsBytes, err := os.ReadFile("sugar/sugar.js")
	if err != nil {
		log.Fatal(err)
	}
	script := "<script>"
	script += string(jsBytes)
	script += "</script>"

	wrapped := map[string]any{
		"Data": data,
		"JSLibrary": template.HTML(script),
	}
	withSuffix := make([]string, len(filenames))
	for i, file := range filenames {
		if !strings.HasSuffix(file, ".sugar") {
			file += ".sugar"
		}
		withSuffix[i] = file
	}
	tmpl := template.Must(template.ParseFiles(withSuffix...))

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "page", wrapped)
	if err != nil {
		log.Fatal(err)
	}
	c.HTML(buf.String())
}

func (c *Context) Redirect(url string) {
	http.Redirect(c.writer, c.Request, url, http.StatusFound)
}

func (c *Context) JSON(content any) {
	c.writer.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(c.writer)
	enc.Encode(content)
}

func (c *Context) NotFound() {
	http.NotFound(c.writer, c.Request)
}

func (c *Context) Form() url.Values {
	err := c.Request.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	return c.Request.Form
}

func (c *Context) Get(url string) (string, http.Header, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body), resp.Header, nil
}

func (c *Context) URL() *URL {
	// data := struct{
	// 	Path string
	// 	Query map[string][]string
	// }{
	// 	Path: c.request.URL.Path,
	// 	Query: c.request.URL.Query(),
	// }

	return &URL{
		Path: c.Request.URL.Path,
		Query: c.Request.URL.Query(),
	}
}

func New(db *pgxpool.Pool, middlewares []Middleware) *Sugar {
	return &Sugar{
		DB: db,
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

func (s *Sugar) Post(path string, handler HandlerFunction) {
	s.Routes = append(s.Routes, Route{
		Path: path,
		Method: http.MethodPost,
		HandlerFunc: handler,
	})
}

func (s *Sugar) Listen(port int) {
	router := http.NewServeMux()
	for _, route := range s.Routes {
		handler := func(w http.ResponseWriter, r *http.Request) {
			if r.Method != route.Method {
				http.NotFound(w, r)
				return
			}

			ctx := Context{
				DB: s.DB,
				Request: r,
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

	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), router))
}