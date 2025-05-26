package sugar

import (
	"encoding/json"
	"net/http"
	"time"
)

type SameSite int

type Cookie struct {
		Name   string
	Value  string
	Quoted bool // indicates whether the Value was originally quoted

	Path       string    // optional
	Domain     string    // optional
	Expires    time.Time // optional
	RawExpires string    // for reading cookies only

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge      int
	Secure      bool
	HttpOnly    bool
	SameSite    SameSite
	Partitioned bool
	Raw         string
	Unparsed    []string // Raw text of unparsed attribute-value pairs
}

type Controller struct {
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

func (r Response) Error(message string, code int) {
	http.Error(r.w, message, code)
}

func (r Response) JSON(data any) {
	r.w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(r.w).Encode(data); err != nil {
		http.Error(r.w, err.Error(), http.StatusInternalServerError)
	}
}

func (h Response) SetCookie(cookie *Cookie) {
	c := http.Cookie{
		Name: cookie.Name,
		Value: cookie.Value,
		Quoted: cookie.Quoted,
		Path: cookie.Path,
		Domain: cookie.Domain,
		Expires: cookie.Expires,
		RawExpires: cookie.RawExpires,
		MaxAge: cookie.MaxAge,
		Secure: cookie.Secure,
		HttpOnly:cookie.HttpOnly ,
		SameSite: http.SameSite(cookie.SameSite),
		Partitioned: cookie.Partitioned,
		Raw: cookie.Raw,
		Unparsed: cookie.Unparsed,
	}
	http.SetCookie(h.w, &c)
}

func (r Request) Cookie(cookieVal string) (*http.Cookie, error) {
	cookie, err := r.r.Cookie(cookieVal)

	return cookie, err
}