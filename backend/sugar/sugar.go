package sugar

import (
	"fmt"
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
	routes := map[string]map[string]Route{
		"/": map[string]Route{
			"GET": Route{

			},
		},
	}
	router := http.NewServeMux()

	for _, route := range s.Routes {
		fmt.Println("for", route.Path)
		checkRoute(route, routes)

		// for i, path := range routes {
			
		// }
		// router.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
		// 	if route.Data[r.Method] == nil {
		// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		// 		return
		// 	}

		// 	if s.Logs {
		// 		log.Printf("%s %s", r.Method, r.URL)
		// 	}

		// 	var bodyBytes []byte
		// 	var err error
		// 	if bodyBytes, err = io.ReadAll(r.Body); err != nil {
		// 		http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	}

		// 	handler := &Handler{
		// 		Response: Response{
		// 			w: w,
		// 		},
		// 		Request: Request{
		// 			Path: r.URL.Path,
		// 			Method: r.Method,
		// 			Body: bodyBytes,
		// 			r: r,
		// 		},
		// 	}

		// 	route.HandlerFunc(handler)
		// })
	}
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("error listening for port %s reason: %s", port, err.Error())
	}
}

func checkRoute(r Route, routes map[string]) {
	if len(routes) > 0 {
		for _, route := range routes {
			route["/"].
		}

		// routes[r.Path] = append(routes[r.Path], Route{
		// 	Path: r.Path,
		// 	Method: r.Method,
		// 	HandlerFunc: r.HandlerFunc,
		// })
	} else {
		// routes[r.Path] = append(routes[r.Path], Route{
		// 	Path: r.Path,
		// 	Method: r.Method,
		// 	HandlerFunc: r.HandlerFunc,
		// })
	}
	// if len(routes) > 0 {
	// 	for i, route := range routes {
	// 		if route.Path == r.Path {
	// 			for _, m := range route.Methods {
	// 				if m == r.Methods[0] {log.Fatal("found a route with identical paths and methods")}
	// 			}
	// 			log.Println("added")
	// 			(*routes)[i].Methods = append((*routes)[i].Methods, r.Methods[0])
	// 		} else {
	// 			*routes = append(*routes, r)
	// 		}
	// 	}
	// } else {
	// 	*routes = append(*routes, r)
	// }
}