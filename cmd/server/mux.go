package main

import (
	"net/http"
	"strings"
	"sync"
)

var paramMap = sync.Map{}

func storeParams(r *http.Request, m map[string]string) {
	paramMap.Store(r, m)
}

func LoadParams(r *http.Request) map[string]string {
	v, ok := paramMap.Load(r)
	if !ok {
		return nil
	}

	return v.(map[string]string)
}

func delParams(r *http.Request) {
	paramMap.Delete(r)
}

// Mux is an http.Handler for setting routes
type Mux struct {
	routes   []route
	endpoint http.HandlerFunc
}

type route struct {
	path    []string
	handler http.HandlerFunc
}

func (mux *Mux) HandleFunc(path string, handlerFunc http.HandlerFunc) {
	if path[0] != '/' {
		panic("path must start with /")
	}

	arr := strings.Split(path, "/")
	for _, v := range arr[1:] {
		if v == "" {
			panic("null string cannot be used as directory name")
		}

		if v[0] == ':' && len(v) < 2 {
			panic("null string cannot be used as parameter name")
		}
	}

	rt := route{
		path:    arr,
		handler: handlerFunc,
	}
	mux.routes = append(mux.routes, rt)
}

func (mux *Mux) Handle(path string, handler http.Handler) {
	mux.HandleFunc(path, handler.ServeHTTP)
}

func matchPath(ref, in []string) (map[string]string, bool) {
	l := len(ref)
	if len(in) < l {
		return nil, false
	}

	m := make(map[string]string)

	for i := 0; i < l; i++ {
		if ref[i][0] == ':' {
			m[string(ref[i][1:])] = in[i]
		} else if ref[i] != in[i] {
			return nil, false
		}
	}

	return m, true
}

func (mux Mux) serveRoute(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")

	var rt route
	var params map[string]string
	var ok bool

	// remove file extension
	path[len(path)-1] = strings.SplitN(path[len(path)-1], ".", 2)[0]

	for _, rt = range mux.routes {
		params, ok = matchPath(rt.path[1:], path[1:])
		if ok {
			break
		}
	}

	if !ok {
		http.Error(w, "path not found "+r.URL.Path, http.StatusNotFound)
		return
	}

	storeParams(r, params)
	rt.handler(w, r)
	delParams(r)

}

type MiddlewareFunc func(next http.HandlerFunc) (this http.HandlerFunc)

func (mux *Mux) Middleware(fn MiddlewareFunc) {

	if mux.endpoint == nil {
		mux.endpoint = fn(mux.serveRoute)
	}

	mux.endpoint = fn(mux.serveRoute)
}

func (mux Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if mux.endpoint == nil {
		mux.serveRoute(w, r)
	} else {
		mux.endpoint(w, r)
	}
}
