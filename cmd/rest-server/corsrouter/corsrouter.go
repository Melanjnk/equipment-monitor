package corsrouter

import (
	"net/http"
	"github.com/gorilla/mux"
)

type CORSRouter struct {
	mux.Router
}

func (cr *CORSRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", 
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
		)
	}
    
	if r.Method != "OPTIONS" {
		// Lets Gorilla work
		cr.Router.ServeHTTP(w, r)
	}
    // Stop here if its Preflighted OPTIONS request
}
