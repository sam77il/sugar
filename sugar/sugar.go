package sugar

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
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
	HTTPVersion int
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

type Config struct {
	Postgres bool
	HTTP2 bool
}

func (c *Context) HTML(html string) {
	c.writer.Header().Set("Content-Type", "text/html")
	c.writer.Write([]byte(html))
}

func (c *Context) Page(data any, filenames ...string) {
	if reflect.TypeOf(data).Kind() != reflect.Struct {
		log.Fatal("Only Structs")
	}
	jsBytes, err := os.ReadFile("sugar/frontend/sugar.js")
	if err != nil {
		log.Fatal(err)
	}
	cssBytes, err := os.ReadFile("styles/main.css")
	if err != nil {
		log.Fatal(err)
	}
	resetCssBytes, err := os.ReadFile("sugar/frontend/reset.css")
	if err != nil {
		log.Fatal(err)
	}
	css := "<style>"
	css += string(cssBytes)
	css += "</style>"
	resetCss := "<style>"
	resetCss += string(resetCssBytes)
	resetCss += "</style>"
	script := "<script>"
	script += string(jsBytes)
	script += "</script>"

	wrapped := map[string]any{
		"Data": data,
		"JSLib": template.HTML(script),
		"ResetStyles": template.HTML(resetCss),
		"Styles": template.HTML(css),
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
	return &URL{
		Path: c.Request.URL.Path,
		Query: c.Request.URL.Query(),
	}
}

func New(config Config) *Sugar {
	var httpVersion int
	if config.HTTP2 {
		httpVersion = 2
	} else {
		httpVersion = 1
	}
	if config.Postgres {
		godotenv.Load()
		db, err := pgxpool.New(context.Background(), os.Getenv("SUGAR_POSTGRES"))
		if err != nil {
			log.Fatal(err)
		}

		return &Sugar{
			DB: db,
			HTTPVersion: httpVersion,
		}
	}

	return &Sugar{
		HTTPVersion: httpVersion,
	}
}

func (s *Sugar) Middleware(path string, handler MiddlewareFunction) {
	s.Middlewares = append(s.Middlewares, Middleware{
		Path: path,
		MiddlewareFunc: handler,
	})
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

func (s *Sugar) Listen(port string) {
	router := http.NewServeMux()
	defer s.DB.Close()

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

	if s.HTTPVersion == 2 {
		server := &http.Server{
			Addr:    port, // Standard-TLS-Port wäre 443, aber 8443 ist gut für lokal
			Handler: router,  // Kein h2c nötig! HTTP/2 läuft automatisch über TLS
		}

		log.Println("Starting server on https://localhost" + port + " (with HTTP/2 on TLS)")
		err := server.ListenAndServeTLS("cert.pem", "key.pem")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("Starting server on https://localhost" + port)
		log.Fatal(http.ListenAndServe(port, router))
	}
}