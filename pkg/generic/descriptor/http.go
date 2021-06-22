package descriptor

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// Cookies ...
type Cookies map[string]string

// HTTPRequest ...
type HTTPRequest struct {
	Header  http.Header
	Query   url.Values
	Cookies Cookies
	Method  string
	Host    string
	Path    string
	Params  *Params // path params
	RawBody []byte
	Body    map[string]interface{}
}

// HTTPResponse ...
type HTTPResponse struct {
	Header     http.Header
	StatusCode int32
	Body       map[string]interface{}
}

// NewHTTPResponse ...
func NewHTTPResponse() *HTTPResponse {
	return &HTTPResponse{
		Header: http.Header{},
		Body:   map[string]interface{}{},
	}
}

// Write to ResponseWriter
func (resp *HTTPResponse) Write(w http.ResponseWriter) error {
	w.WriteHeader(int(resp.StatusCode))
	for k := range resp.Header {
		w.Header().Set(k, resp.Header.Get(k))
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(resp.Body)
}

// Param in request path
type Param struct {
	Key   string
	Value string
}

// Params and recyclable
type Params struct {
	params   []Param
	recycle  func(*Params)
	recycled bool
}

// Recycle the Params
func (ps *Params) Recycle() {
	if ps.recycled {
		return
	}
	ps.recycled = true
	ps.recycle(ps)
}

// ByName search Param by given name
func (ps *Params) ByName(name string) string {
	for _, p := range ps.params {
		if p.Key == name {
			return p.Value
		}
	}
	return ""
}
