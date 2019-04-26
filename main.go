package main

import (
    "os"
    "log"
    "github.com/AugustoQueiroz/telegramBot"
)

var (
    bot     telegramBot.Bot
)

func main() {
    botToken := os.Getenv("BOT_TOKEN")
    if botToken == "" {
        log.Fatal("$BOT_TOKEN must be set")
    }
    bot = telegramBot.NewBot(botToken)

    bot.HandleFunc("/start", StartHandler)
    bot.HandleFunc("/createopensanta", OpenSantaHandler)
    bot.HandleFunc("/createsecretsanta", SecretSantaHandler)

    bot.CallbackHandler = CallbackHandler

    activeSantas = make(map[int]*SantaInfo)

    log.Println("Creating router")
    //go router()

    log.Println("Starting poller")
    bot.Poller()
}
