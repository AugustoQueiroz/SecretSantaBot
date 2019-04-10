package main

import (
    "log"
    "github.com/AugustoQueiroz/telegramBot"
)

func main() {
    telegramBot.HandleFunc("/start", StartHandler)
    telegramBot.HandleFunc("/createopensanta", OpenSantaHandler)
    telegramBot.HandleFunc("/createsecretsanta", SecretSantaHandler)

    telegramBot.CallbackHandler = CallbackHandler

    activeSantas = make(map[int]*SantaInfo)

    log.Println("Creating router")
    go router()

    log.Println("Starting poller")
    telegramBot.Poller()
}
