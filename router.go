package main

import (
    // Standard Packages
    "os"
    "fmt"
    "log"
    "strconv"
    "net/http"

    // External Packages
    "github.com/rs/cors"
    "github.com/gorilla/mux"

    // Project Packages
)

func main() {
    // Get the port the bot will be attached to from the environment
    port := os.Getenv("PORT")

    if port == "" {
        log.Fatal("$PORT must be set")
    }

    var err error
    repeatStr := os.Getenv("REPEAT")
    repeat, err := strconv.Atoi(repeatStr)
    if err != nil {
        log.Print("Error converting $REPEAT to an int: %q - Using default", err)
        repeat = 5
    }

    repeat = repeat + 5

    // Create the router
    router := mux.NewRouter()

    // Endpoints and their handlers
    isUpEndpoint := "/"

    router.HandleFunc(isUpEndpoint, isUp).Methods("GET")

    // CORS setting to allow Cross-Origin Requests
    handler := cors.Default().Handler(router)

    // Start router listening and serving
    log.Fatal(http.ListenAndServe(":" + port, handler))
}

func isUp(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "The bot is up and running")
}
