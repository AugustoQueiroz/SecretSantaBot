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
    "telegramBot"
)

var (
    domain = "https://secretsanta5000.herokuapp.com/"
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
        log.Print("Error converting $REPEAT to an int: ", err)
        repeat = 5
    }

    repeat = repeat + 5

    // Create the router
    router := mux.NewRouter()

    // Endpoints
    isUpEndpoint := "/"
    telegramEndpoint := "/telegram/"

    // Handlers for endpoints
    router.HandleFunc(isUpEndpoint, isUp).Methods("GET")
    router.HandleFunc(telegramEndpoint + "/{token}/", telegramBot.HandleUpdate).Methods("POST")

    // CORS setting to allow Cross-Origin Requests
    handler := cors.Default().Handler(router)

    // Set the bot's webhook to it's endpoint
    success := telegramBot.SetWebhook(domain + telegramEndpoint)
    if !success {
        log.Fatal("Could not set telegram webhook")
    }
    log.Println("Telegram webhook set")

    // Start router listening and serving
    log.Fatal(http.ListenAndServe(":" + port, handler))
}

func isUp(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "The bot is up and running")
}
