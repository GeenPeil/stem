package cors

import "net/http"

type wildcardCORSHandler struct {
	h http.Handler
}

// WrapWildcardCORSHandler returns given handler wrapped with wildcard cors functionality.
func WrapWildcardCORSHandler(h http.Handler) http.Handler {
	return &wildcardCORSHandler{h}
}

func (c *wildcardCORSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))

	if r.Method == "OPTIONS" {
		return
	}

	c.h.ServeHTTP(w, r)
}
