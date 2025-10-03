package routes

import "net/http"

type Headers struct {
}

func (headers *Headers) ServeHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Test", "test44")
}
