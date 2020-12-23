package http

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/milanrodriguez/stee/internal/stee"
)

type apiHandler struct {
	http.Handler
	enable    bool
	prefix    string
	simpleAPI struct {
		enable bool
	}
}

func handleAPI(core *stee.Core, prefix string, simpleAPIEnabled bool) http.Handler {
	router := httprouter.New()

	// General API route(s), not yet implemented
	router.GET(prefix+"/", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, "Error 501: API not implemented in this version. Try to upgrade Stee.")
	})

	// Simple API routes
	if simpleAPIEnabled {
		router.GET(prefix+"/simple/add/:base64target", handleSimpleAdd(core))
		router.GET(prefix+"/simple/add/:base64target/:key", handleSimpleAdd(core))
		router.GET(prefix+"/simple/get/:key", handleSimpleGet(core))
		router.GET(prefix+"/simple/del/:key", handleSimpleDel(core))
	}

	return router
}
