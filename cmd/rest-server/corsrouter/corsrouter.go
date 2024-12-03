package corsrouter

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

type CORSRouter struct {
	mux.Router
}

func (corsrouter *CORSRouter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	origin := request.Header.Get("Origin")
	if request.Method == http.MethodOptions {
		if origin != "" {
			log.Printf("CORS preflight request from origin: %s", origin)
			corsrouter.setCORSHeaders(writer, origin)
		} else {
			log.Println("OPTIONS request without Origin header, skipping CORS handling.")
		}
		// Stop here for preflight OPTIONS request with 200 status
		writer.WriteHeader(http.StatusOK)
		return
	}

	// Set CORS headers for other requests if the Origin header is present
	if origin != "" {
		log.Printf("CORS request from origin: %s, method: %s", origin, request.Method)
		corsrouter.setCORSHeaders(writer, origin)
	}

	// Delegate to Gorilla Router for non-OPTIONS requests
	corsrouter.Router.ServeHTTP(writer, request)
}

func (corsrouter *CORSRouter) setCORSHeaders(writer http.ResponseWriter, origin string) {
	log.Printf("Setting CORS headers for origin: %s", origin)
	writer.Header().Set("Access-Control-Allow-Origin", origin)
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
	)
}
