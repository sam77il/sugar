package sugar

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func (s *Sugar) Put(path string, handler HandlerFunction) {
	s.Routes = append(s.Routes, Route{
		Path: path,
		Method: "PUT",
		HandlerFunc: handler,
	})
}

func (s *Sugar) Patch(path string, handler HandlerFunction) {
	s.Routes = append(s.Routes, Route{
		Path: path,
		Method: "PATCH",
		HandlerFunc: handler,
	})
}

func (s *Sugar) Delete(path string, handler HandlerFunction) {
	s.Routes = append(s.Routes, Route{
		Path: path,
		Method: "DELETE",
		HandlerFunc: handler,
	})
}

func (s Sugar) Listen(port uint64) {
	fmt.Println(">> Booting up server <<")
	router := http.NewServeMux()
	routes := make(map[string]map[string]Route)

	for _, route := range s.Routes {
		checkRoute(route, routes)
	}

	for path := range routes {
		if strings.Contains(path, ":") {
			keys := strings.TrimPrefix(path, ":")
			log.Println(keys[len(keys) - 1])
		}

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
	finalPort := strconv.FormatUint(port, 10)
	fmt.Println(">> Successfully booted up on port", finalPort)
	log.Fatal(http.ListenAndServe(":" + finalPort, router))
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