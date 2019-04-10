package main

import (
    "log"
    "telegramBot"
)

func main() {
    telegramBot.HandleFunc("/start", StartHandler)
    telegramBot.HandleFunc("/createopensanta", OpenSantaHandler)
    telegramBot.HandleFunc("/createsecretsanta", SecretSantaHandler)

    telegramBot.CallbackHandler = CallbackHandler

    activeSantas = make(map[int]*SantaInfo)

    log.Println("Starting poller")
    telegramBot.Poller()
}
