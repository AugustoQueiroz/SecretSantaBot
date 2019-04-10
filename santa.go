package main

import (
    // Standard Packages
    "time"
    "strconv"
    "math/rand"

    // Project Packages
    "telegramBot"
)

type JoinChannel chan *telegramBot.User

var (
    activeJoinChannels map[int]JoinChannel
)

// A function that receives new participants in a specific santa until the channel is closed
// Meant to be run as a goroutine
func SantaWriteUp(santaType string, chat *telegramBot.Chat, joinChannel JoinChannel) {
    var usersParticipating []telegramBot.User

    for {
        user, ok := <-joinChannel

        if !ok || user == nil {
            SantaPairing(santaType, chat, usersParticipating)
            return
        }

        usersParticipating = append(usersParticipating, *user)
    }
}

func SantaPairing(santaType string, chat *telegramBot.Chat, usersParticipating []telegramBot.User) {
    r := rand.New(rand.NewSource(time.Now().Unix()))
    pairs := make([]telegramBot.User, len(usersParticipating))

    // Generate a random permutation and assign the pairing according to it
    permutation := r.Perm(len(pairs))
    for i, randIndex := range permutation {
        pairs[i] = usersParticipating[randIndex]

        if pairs[i] == usersParticipating[i] {
            // If someone got themselves swaps with the person before them
            // If was the first person, change it with the last
            if i == 0 {
                pairs[i] = usersParticipating[permutation[len(pairs)-1]]
                permutation[len(pairs)-1] = 0
            } else {
                temp := pairs[i]
                pairs[i] = pairs[i-1]
                pairs[i-1] = temp
            }
        }
    }

    SharePairing(santaType, chat, usersParticipating, pairs)
}

func SharePairing(santaType string, chat *telegramBot.Chat, usersParticipating []telegramBot.User, pairings []telegramBot.User) {
    if santaType == "open" {
        resultMessageBody := "Secret Santa Pairings:\n"
        for i, user := range usersParticipating {
            santa := pairings[i]
            resultMessageBody += "[" + user.FirstName + "](tg://user?id=" + strconv.Itoa(user.Id) + "): "
            resultMessageBody += "[" + santa.FirstName + "](tg://user?id=" + strconv.Itoa(santa.Id) + ")\n"
        }

        telegramBot.SendMarkdownMessage(resultMessageBody, chat.Id)
    } else if santaType == "secret" {

    }
}
