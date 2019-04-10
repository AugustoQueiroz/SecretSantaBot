package main

import (
    "telegramBot"
)

func main() {
    telegramBot.HandleFunc("/start", StartHandler)
    telegramBot.HandleFunc("/createopensanta", OpenSantaHandler)

    telegramBot.CallbackHandler = CallbackHandler

    activeJoinChannels = make(map[int]JoinChannel)

    telegramBot.Poller()
}
