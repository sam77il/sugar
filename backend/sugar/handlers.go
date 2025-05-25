package sugar

import (
	"encoding/json"
	"net/http"
)


type Handler struct {
	Request
	Response
}

type Request struct {
	Body []byte
	Path    string
	Method string
	r *http.Request
}

type Response struct {
	w http.ResponseWriter
}

func (h Response) JSON(data any) {
	h.w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(h.w).Encode(data); err != nil {
		http.Error(h.w, err.Error(), http.StatusInternalServerError)
	}
}