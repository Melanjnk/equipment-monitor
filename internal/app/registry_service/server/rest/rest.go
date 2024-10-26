package rest

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/Melanjnk/equipment-monitor/cmd/rest-server/corsrouter"
)

type RestServer struct {
	http.Server
}

func (rs *RestServer) start(goRoutine func()) {
	go goRoutine()
 
	// Wait for an interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
 
	// Attempt a graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	rs.Shutdown(ctx)
}

func (rs *RestServer) StartHTTP(port string, router *corsrouter.CORSRouter) {
	rs.Addr = port
	rs.Handler = router
	go func() {
        if err := rs.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
            log.Fatalf("HTTP server error: %v", err)
        }
        log.Println("Stopped serving new connections.")
    }()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan

    shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10 * time.Second)
    defer shutdownRelease()

    if err := rs.Shutdown(shutdownCtx); err != nil {
        log.Fatalf("HTTP shutdown error: %v", err)
    }
    log.Println("Graceful shutdown complete.")
}
