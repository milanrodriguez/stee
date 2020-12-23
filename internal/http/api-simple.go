package http

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/milanrodriguez/stee/internal/stee"
)

func handleSimpleAdd(core *stee.Core) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		targetBytes, err := base64.URLEncoding.DecodeString(ps.ByName("base64target"))
		if err != nil {
			// If URLEncoding doesn't, we use RawURLEncoding (omits padding characters)
			targetBytes, err = base64.RawURLEncoding.DecodeString(ps.ByName("base64target"))
		}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error: %s", err.Error())
			return
		}
		target := string(targetBytes)

		key := ps.ByName("key")
		if key != "" {
			err = core.AddRedirectionWithKey(key, target)
		} else {
			key, err = core.AddRedirectionWithoutKey(target)
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %s", err.Error())
			return
		}
		fmt.Fprintf(w, "✔️ Added redirection: '%s' -> %s", key, target)
	}
}

func handleSimpleGet(core *stee.Core) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		key := ps.ByName("key")
		target, err := core.GetRedirection(key)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Impossible to read key: %s", err)
			return
		}
		fmt.Fprintf(w, "Key \"%s\" is pointing to \"%s\"", key, target)
	}
}

func handleSimpleDel(core *stee.Core) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		key := ps.ByName("key")
		err := core.DeleteRedirection(key)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %s", err.Error())
			return
		}
		fmt.Fprintf(w, "Deleted key \"%s\"", key)
	}
}
