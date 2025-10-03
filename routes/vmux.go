package routes

import (
	"container/list"
	"net/http"
)

type VMuxMiddleware interface {
	Run(http.ResponseWriter, *http.Request, func(http.ResponseWriter, *http.Request))
}

// declare a type for shortcut. The func returns the same func but wrappee around with the middleware
type MiddlewareFunc func(http.ResponseWriter, *http.Request, func(http.ResponseWriter, *http.Request))

type VMux struct {
	http.ServeMux
	middlewareFuncs list.List
}

func (vmux *VMux) Add(middlewareFunc MiddlewareFunc) {
	vmux.middlewareFuncs.PushBack(middlewareFunc)
}

func (vmux *VMux) nextMiddleware(elem *list.Element) func(w http.ResponseWriter, r *http.Request) {
	if elem != nil {
		return func(w http.ResponseWriter, r *http.Request) {
			//This recurses and returns funcs from the list.
			elem.Value.(MiddlewareFunc)(w, r, vmux.nextMiddleware(elem.Next()))
		}
	}
	return vmux.ServeMux.ServeHTTP
}

func (vmux *VMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	vmux.nextMiddleware(vmux.middlewareFuncs.Front())(w, req)
}
