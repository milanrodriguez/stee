package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/milanrodriguez/stee/stee"
)

func handleMain(core *stee.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path

		if strings.HasPrefix(key, "/") {
			key = key[1:]
		}
		if strings.HasSuffix(key, "/") {
			key = key[:len(key)-1]
		}

		target, ok := core.GetRedirection(key)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Error 404: No redirection found for key \"%s\"", key)
			return
		}
		http.Redirect(w, r, target, http.StatusFound)
	})
}
