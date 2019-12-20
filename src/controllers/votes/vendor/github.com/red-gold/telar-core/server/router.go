package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	cf "github.com/red-gold/telar-core/config"
)

type ServerRouter struct {
	router *httprouter.Router
}

var router *httprouter.Router

func NewServerRouter() *ServerRouter {
	router = httprouter.New()
	return &ServerRouter{
		router: router,
	}
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (r *ServerRouter) GET(path string, handle Handle, protected RouteProtection) {
	r.router.GET(path, Req(handle, protected))
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func (r *ServerRouter) HEAD(path string, handle Handle, protected RouteProtection) {
	r.router.HEAD(path, Req(handle, protected))
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func (r *ServerRouter) OPTIONS(path string, handle Handle, protected RouteProtection) {
	r.router.OPTIONS(path, Req(handle, protected))
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (r *ServerRouter) POST(path string, handle Handle, protected RouteProtection) {
	r.router.POST(path, Req(handle, protected))
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (r *ServerRouter) PUT(path string, handle Handle, protected RouteProtection) {
	r.router.PUT(path, Req(handle, protected))
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (r *ServerRouter) PATCH(path string, handle Handle, protected RouteProtection) {
	r.router.PATCH(path, Req(handle, protected))
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (r *ServerRouter) DELETE(path string, handle Handle, protected RouteProtection) {
	r.router.DELETE(path, Req(handle, protected))
}

// GETWR with main http request parameters is a shortcut for router.HandleWR("GET", path, handle)
func (r *ServerRouter) GETWR(path string, handle HandleWR, protected RouteProtection) {
	r.router.GET(path, ReqWR(handle, protected))
}

// HEADWR is a shortcut for router.HandleWR("HEAD", path, handle)
func (r *ServerRouter) HEADWR(path string, handle HandleWR, protected RouteProtection) {
	r.router.HEAD(path, ReqWR(handle, protected))
}

// OPTIONSWR is a shortcut for router.HandleWR("OPTIONS", path, handle)
func (r *ServerRouter) OPTIONSWR(path string, handle HandleWR, protected RouteProtection) {
	r.router.OPTIONS(path, ReqWR(handle, protected))
}

// POSTWR is a shortcut for router.HandleWR("POST", path, handle)
func (r *ServerRouter) POSTWR(path string, handle HandleWR, protected RouteProtection) {
	r.router.POST(path, ReqWR(handle, protected))
}

// PUTWR is a shortcut for router.HandleWR("PUT", path, handle)
func (r *ServerRouter) PUTWR(path string, handle HandleWR, protected RouteProtection) {
	r.router.PUT(path, ReqWR(handle, protected))
}

// PATCHWR is a shortcut for router.HandleWR("PATCH", path, handle)
func (r *ServerRouter) PATCHWR(path string, handle HandleWR, protected RouteProtection) {
	r.router.PATCH(path, ReqWR(handle, protected))
}

// DELETEWR is a shortcut for router.HandleWR("DELETE", path, handle)
func (r *ServerRouter) DELETEWR(path string, handle HandleWR, protected RouteProtection) {
	r.router.DELETE(path, ReqWR(handle, protected))
}

// POSTFILE is a shortcut for router.HandleWR("POST", path, handle)
func (r *ServerRouter) POSTFILE(path string, handle HandleWR, protected RouteProtection) {
	r.router.POST(path, ReqFileWR(handle, protected))
}

// ServeHTTP makes the router implement the http.Handler interface.
func (r *ServerRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	config := cf.AppConfig
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Content-Type", "application/json")
	// if origin := req.Header.Get("Origin"); origin != "" {
	w.Header().Set("Access-Control-Allow-Origin", *config.Origin)
	// }
	w.Header().Set("Access-Control-Allow-Headers", "'X-Requested-With, X-HTTP-Method-Override, Accept, Content-Type,access-control-allow-origin, access-control-allow-headers")
	r.router.ServeHTTP(w, req)
}
