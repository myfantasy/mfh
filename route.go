package mfh

import (
	"net/http"
	"sort"
)

// Route is a route manager for get Hande function routing only by url
type Route struct {
	rDir     map[string]func(http.ResponseWriter, *http.Request)
	rMap     map[int]map[string]func(http.ResponseWriter, *http.Request)
	lMap     map[int]struct{}
	lS       []int
	defRoute func(http.ResponseWriter, *http.Request)
}

// AddDirectRouteH Add Exec route
func (r *Route) AddDirectRouteH(uri string, h http.Handler) {
	r.AddDirectRoute(uri, h.ServeHTTP)
}

// AddRouteH Add route
func (r *Route) AddRouteH(uri string, h http.Handler) {
	r.AddRoute(uri, h.ServeHTTP)
}

// AddDefaultRouteH Add default route
func (r *Route) AddDefaultRouteH(h http.Handler) {
	r.AddDefaultRoute(h.ServeHTTP)
}

// AddDirectRoute Add Exec route
func (r *Route) AddDirectRoute(uri string, f func(http.ResponseWriter, *http.Request)) {
	uri = "/" + uri
	if r.rDir == nil {
		r.rDir = make(map[string]func(http.ResponseWriter, *http.Request))
	}
	r.rDir[uri] = f
}

// AddRoute Add route
func (r *Route) AddRoute(uri string, f func(http.ResponseWriter, *http.Request)) {
	uri = "/" + uri
	if r.rMap == nil {
		r.rMap = make(map[int]map[string]func(http.ResponseWriter, *http.Request))
		r.lMap = make(map[int]struct{})
	}

	l := len(uri)
	if _, ok := r.lMap[l]; !ok {
		r.lMap[l] = struct{}{}
		r.rMap[l] = map[string]func(http.ResponseWriter, *http.Request){}
		r.lS = append(r.lS, l)
		sort.Slice(r.lS, func(i, j int) bool { return r.lS[i] > r.lS[j] })
	}

	r.rMap[l][uri] = f
}

// AddDefaultRoute Add default route
func (r *Route) AddDefaultRoute(f func(http.ResponseWriter, *http.Request)) {
	r.defRoute = f
}

// SearchRoute Found Route by url
func (r *Route) SearchRoute(uri string) (fn func(http.ResponseWriter, *http.Request)) {

	if f, ok := r.rDir[uri]; ok {
		return f
	}

	l := len(uri)

	for _, v := range r.lS {
		if v <= l {
			s := uri[:v]
			m := r.rMap[v]
			if f, ok := m[s]; ok {
				return f
			}
		}
	}

	return r.defRoute
}

// ServeHTTP http.Handler interface
func (r *Route) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	f := r.SearchRoute(req.URL.String())
	f(res, req)
}
