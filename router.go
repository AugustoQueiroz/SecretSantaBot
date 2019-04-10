package main

import (
    // Standard Packages
    "os"
    "log"
    "fmt"
    "strconv"
    "net/http"

    // External Packages
    "github.com/rs/cors"
    "github.com/gorilla/mux"
)

func router() {
    port := os.Getenv("PORT")

    if port == "" {
        log.Fatal("$PORT must be set!")
    }

    repeatStr := os.Getenv("REPEAT")
    repeat, err := strconv.Atoi(repeatStr)
    if err != nil {
        log.Println("$REPEAT not set, using default value")
        repeat = 5
    }

    repeat += 5

    // Create the router
	router := mux.NewRouter()

	// Endpoints and their handlers
	router.HandleFunc("/", isUp).Methods("GET")

	// CORS setting to allow Cross-Origin Requests
	handler := cors.Default().Handler(router)

	// Start router listening and serving
	log.Fatal(http.ListenAndServe(":" + port, handler))
}

func isUp(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Yes")
}
