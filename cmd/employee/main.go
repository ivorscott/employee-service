package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	r := mux.NewRouter()
	r.Methods("GET").Path("/employees/{employee-id}").HandlerFunc(EmployeesHandler)
	http.Handle("/", r)

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Println("starting server...")
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-shutdown:
		log.Println("main : start shutdown", sig)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			log.Println("shutdown error", err)

		}
		log.Println("shutting down")
		os.Exit(0)
	}
}

type Employee struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Job       string `json:"job"`
}

func EmployeesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	e := map[int]Employee{
		0: {
			ID:        0,
			FirstName: "Alan",
			LastName:  "Watts",
			Job:       "Philosopher",
		},
		1: {
			ID:        1,
			FirstName: "John",
			LastName:  "Locke",
			Job:       "Philosopher",
		},
	}

	id, err := strconv.Atoi(vars["employee-id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "bad request"}`))
		return
	}

	if employee, ok := e[id]; ok {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(employee)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
	}
}
