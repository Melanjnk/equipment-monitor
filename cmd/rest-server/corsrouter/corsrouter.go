package corsrouter

import (
	"net/http"
	"github.com/gorilla/mux"
)

type CORSRouter struct {
	mux.Router
}

func (cr *CORSRouter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if origin := request.Header.Get("Origin"); origin != "" {
		writer.Header().Set("Access-Control-Allow-Origin", origin)
		writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", 
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
		)
	}
    
	if request.Method != "OPTIONS" {
		// Lets Gorilla work
		cr.Router.ServeHTTP(writer, request)
	}
    // Stop here if its Preflighted OPTIONS request
}
