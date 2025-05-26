package sugar

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Sugar struct {
	Config
	Routes []Route
}

type Route struct {
	Path string
	Method string
	HandlerFunc HandlerFunction
}

type Config struct {
	AppName string
	Logs bool
}

type HandlerFunction func(*Controller)

func New(config Config) *Sugar {
	sugar := &Sugar{
		Config: config,
	}

	return sugar
}

func (s *Sugar) Get(path string, handler HandlerFunction) {
	s.Routes = append(s.Routes, Route{
		Path: path,
		Method: "GET",
		HandlerFunc: handler,
	})
}

func (s *Sugar) Post(path string, handler HandlerFunction) {
	s.Routes = append(s.Routes, Route{
		Path: path,
		Method: "POST",
		HandlerFunc: handler,
	})
}

func (s Sugar) Listen(port string) {
	router := http.NewServeMux()
	routes := make(map[string]map[string]Route)

	for _, route := range s.Routes {
		fmt.Println("for", route.Path)
		checkRoute(route, routes)
	}

	for path,_ := range routes {
		router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			var route Route
			var ok bool
			if route, ok = routes[path][r.Method]; !ok {
				http.Error(w, "not allowed http method", http.StatusMethodNotAllowed)
				return
			}
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "could not read body", http.StatusInternalServerError)
				return
			}

			controller := Controller{
				Request: Request{
					r: r,
					Body: bodyBytes,
					Path: r.URL.Path,
					Method: r.Method,
				},
				Response: Response{
					w: w,
				},
			}

			route.HandlerFunc(&controller)
		})
	}

	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("error listening for port %s reason: %s", port, err.Error())
	}
}

func checkRoute(r Route, routes map[string]map[string]Route) {
	_, ok := routes[r.Path]
	if ok {
		_, ok2 := routes[r.Path][r.Method]
		if ok2 {
			log.Fatal("already existing path and route detected")
			return
		}
	} else {
		routes[r.Path] = make(map[string]Route)
	}
	
	routes[r.Path][r.Method] = Route{
		Path: r.Path,
		Method: r.Method,
		HandlerFunc: r.HandlerFunc,
	}
}