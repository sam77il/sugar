package sugar

import (
	"io"
	"log"
	"net/http"
)

type Sugar struct {
	Config
	Routes []Route
}

type Route struct {
	Method      string
	Path        string
	HandlerFunc HandlerFunction
}

type Config struct {
	AppName string
	Logs bool
}

type HandlerFunction func(*Handler)

func New(config Config) *Sugar {
	sugar := &Sugar{
		Config: config,
	}

	return sugar
}

func (s *Sugar) Get(path string, handler HandlerFunction) {
	s.Routes = append(s.Routes, Route{
		Method:      "GET",
		Path:        path,
		HandlerFunc: handler,
	})
}

func (s *Sugar) Post(path string, handler HandlerFunction) {
	s.Routes = append(s.Routes, Route{
		Method:      "POST",
		Path:        path,
		HandlerFunc: handler,
	})
}

func (s Sugar) Listen(port string) {
	router := http.NewServeMux()

	for _, route := range s.Routes {
		router.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != route.Method {
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}

			if s.Logs {
				log.Printf("%s %s", r.Method, r.URL)
			}

			var bodyBytes []byte
			var err error
			if bodyBytes, err = io.ReadAll(r.Body); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			handler := &Handler{
				Response: Response{
					w: w,
				},
				Request: Request{
					Path: r.URL.Path,
					Method: r.Method,
					Body: bodyBytes,
					r: r,
				},
			}

			

			route.HandlerFunc(handler)
		})
	}
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("error listening for port %s reason: %s", port, err.Error())
	}
}