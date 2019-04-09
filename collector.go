package main

import (
    "os"
    "log"
    "fmt"
    "strconv"
    "net/http"

    // External Packages
    "github.com/rs/cors"
    "github.com/gorilla/mux"

    // Project Packages
    "telegramBot"
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

    telegramBot.DeleteWebhook()

    // Create the router
    router := mux.NewRouter()

    // Endpoints
    isUpEndpoint := "/"

    // Handlers for endpoints
    router.HandleFunc(isUpEndpoint, isUp).Methods("GET")

    // CORS setting to allow Cross-Origin Requests
    handler := cors.Default().Handler(router)

    // Set the bot's webhook to it's endpoint

    // Start router listening and serving
    log.Fatal(http.ListenAndServe(":" + port, handler))
}

func isUp(w http.ResponseWriter, r *http.Request) {
    updates := telegramBot.GetUpdates()

    if len(updates) > 0 {
        update := updates[len(updates)-1]
        messageBody := update.Message.Body
        entities := update.Message.Entities

        if len(entities) > 0 {
            entity := entities[0]
            fmt.Fprintf(w, entity.Type)
            fmt.Fprintf(w, strconv.Itoa(entity.Offset))
            fmt.Fprintf(w, strconv.Itoa(entity.Length))
        }

        fmt.Fprintf(w, messageBody)
    } else {
        fmt.Fprintf(w, "Hello") 
    }
}
