// Command verification runs a tiny standalone HTTP service that
// employee-service calls on every employee lookup.
//
// It exists purely so the observability lab has a second real hop: with only
// employee-service running, Tempo's service graph has a single real node (the
// caller shows up as a "virtual" user node because loadgen/curl doesn't
// propagate trace context). Once employee-service calls this service with an
// OTel-instrumented client, the same trace id spans both processes, giving a
// three node graph (user -> employee-service -> verification-service) and a
// waterfall with more than one span to correlate.
package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devpies/employee-service/pkg/trace"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	ctx := context.Background()

	prv, err := trace.NewProvider(trace.ProviderConfig{
		OTLPEndpoint:   getenv("API_TRACE_OTLP_ENDPOINT", "localhost:4318"),
		ServiceName:    "verification-service",
		ServiceVersion: "1.0.0",
		Environment:    "dev",
		Disabled:       getenv("API_TRACE_DISABLED", "false") == "true",
	})
	if err != nil {
		log.Fatalf("trace provider: %v", err)
	}
	defer prv.Close(ctx)

	mux := http.NewServeMux()
	mux.Handle("/verify/", otelhttp.NewHandler(http.HandlerFunc(verifyHandler), "verify"))

	addr := getenv("ADDR", "0.0.0.0:8090")
	srv := &http.Server{Addr: addr, Handler: mux}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("verification-service listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}

// verifyHandler simulates an eligibility check on the employee id in the
// path, with a small random delay so it shows up as real work in the trace.
func verifyHandler(w http.ResponseWriter, r *http.Request) {
	employeeID := r.URL.Path[len("/verify/"):]

	time.Sleep(time.Duration(5+rand.Intn(20)) * time.Millisecond)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"employeeId": employeeID,
		"verified":   true,
	})
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
