package rest

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/omriza/go-msvc-base/pkg/adapters"
)

type restAdapter struct {
	server *http.Server
	wait   time.Duration
	state  adapters.DrivingAdapterState
}

func (ra *restAdapter) Describe() adapters.DrivingAdapterDescription {
	return adapters.DrivingAdapterDescription{
		Title: "REST Adapter",
		State: ra.state,
	}
}

func (ra *restAdapter) Start() {
	// Run server in a goroutine so that it doesn't block.
	go func() {
		if err := ra.server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	ra.state = adapters.Started
	log.Println("Server is running...")

	c := make(chan os.Signal, 1)
	// Graceful shutdown when quit via SIGINT (Ctrl+C)
	signal.Notify(c, os.Interrupt)

	// Block until SIGINT is received, then stopping
	<-c
	ra.Stop()
}

func (ra *restAdapter) Stop() {
	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), ra.wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	ra.server.Shutdown(ctx)

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation
	log.Println("Shutting down")
	ra.state = adapters.Stopped
	os.Exit(0)
}

func (ra *restAdapter) GetCommand() adapters.DrivingAdapterCmd {
	return adapters.DrivingAdapterCmd{}
}

// NewRestAdapter defines a factory method for rest adapters
func NewRestAdapter(config adapters.DrivingAdapterConfig) (adapters.DrivingAdapter, error) {
	ra := restAdapter{
		state: adapters.Initialized,
	}

	flag.DurationVar(&ra.wait, "graceful-timeout", time.Second*2, "the duration for which the server gracefully wait for existing connections to finish")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/health", HealthHandler)

	ra.server = &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return &ra, nil
}

// Handlers

// HomeHandler handles the "/" (home) route
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Welcome home")
	log.Println("Request:", r)

	w.WriteHeader(http.StatusOK)

}

// HealthHandler handles the "/health" route
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Health called")
	log.Println("Request:", r)

	w.WriteHeader(http.StatusOK)

}

// Register ... kinda crappy, mb replace this with injection + adding it to the interface
func Register() {
	adapters.RegisterDrivingAdapter("rest", NewRestAdapter)
}
