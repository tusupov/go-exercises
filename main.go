package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tusupov/go-exercises/handle"
	"github.com/tusupov/go-exercises/middleware"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	port = kingpin.Flag("port", "Listening port").Short('p').Default("8080").Int()
)

func init() {
	kingpin.Parse()
}

func main() {

	r := mux.NewRouter()

	// Handle function
	h := handle.New()
	r.HandleFunc("/account", h.AccountOpen).Methods(http.MethodPut)
	r.HandleFunc("/account", h.AccountBalance).Methods(http.MethodGet)
	r.HandleFunc("/account", h.AccountDeposit).Methods(http.MethodPost)
	r.HandleFunc("/account", h.AccountClose).Methods(http.MethodDelete)

	// Middleware
	middleware.SetLogger(os.Stderr)
	r.Use(middleware.Panic, middleware.AccessLog)

	// Start server
	log.Printf("Listening port [%d] ...", *port)
	if err := http.ListenAndServe(":"+strconv.Itoa(*port), r); err != nil {
		log.Fatal(err)
	}

}
