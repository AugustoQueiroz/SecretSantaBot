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

type SantaInfo struct {
    santaType       string
    chat            telegramBot.Chat
    messageId       int
    participants    []telegramBot.User
}

var (
    activeJoinChannels map[int]JoinChannel
    activeSantas        map[int]*SantaInfo
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

func JoinSanta(messageId int, user *telegramBot.User) {
    santaInfo, exists := activeSantas[messageId]
    if !exists {
        messageBody := "The santa you're trying to join doesn't exist. Perhaps it has expired or someone has closed it"
        telegramBot.SendMarkdownMessage(messageBody, user.Id)
        return
    }

    if NotIn(santaInfo.participants, user) {
        santaInfo.participants = append(santaInfo.participants, *user)

        // Message User
        messageBody := "You have been added as a participant to an " + santaInfo.santaType + " santa"
        telegramBot.SendMarkdownMessage(messageBody, user.Id)

        // Edit message in group
        newMessageBody := "#participants: " + strconv.Itoa(len(santaInfo.participants)) + "\n"
        for _, participant := range santaInfo.participants {
            newMessageBody += "[" + participant.FirstName + "](tg://user?id=" + strconv.Itoa(participant.Id) + ")\n"
        }
        telegramBot.EditMessageText(santaInfo.chat.Id, santaInfo.messageId, newMessageBody, "Markdown")
    } else {
        messageBody := "You have already been added to that " + santaInfo.santaType + " santa."
        telegramBot.SendMarkdownMessage(messageBody, user.Id)
    }
}

func SantaDone(messageId int) {
    santaInfo, exists := activeSantas[messageId]
    if !exists {
        return
    }

    delete(activeSantas, messageId)

    if len(santaInfo.participants) < 2 {
        SantaNotEnoughParticipants(santaInfo, messageId)
        return
    }

    SantaPairing(santaInfo.santaType, &santaInfo.chat, santaInfo.participants)
}

func SantaNotEnoughParticipants(santaInfo *SantaInfo, messageId int) {
    messageBody := "Not enough participants"

    telegramBot.SendMarkdownMessage(messageBody, santaInfo.chat.Id)
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

    SharePairings(santaType, chat, usersParticipating, pairs)
}

func SharePairings(santaType string, chat *telegramBot.Chat, usersParticipating []telegramBot.User, pairs []telegramBot.User) {
    messageBody := "The secret santas have been picked..."
    telegramBot.SendMarkdownMessage(messageBody, chat.Id)

    if santaType == "open" {
        OpenSharing(chat, usersParticipating, pairs)
    } else if santaType == "secret" {
        SecretSharing(usersParticipating, pairs)
    }
}

func OpenSharing(chat *telegramBot.Chat, usersParticipating []telegramBot.User, pairings []telegramBot.User) {
    messageBody := "Open Santa Pairings:\n"
    for i, user := range usersParticipating {
        santa := pairings[i]
        messageBody += "[" + user.FirstName + "](tg://user?id=" + strconv.Itoa(user.Id) + ") -> "
        messageBody += "[" + santa.FirstName + "](tg://user?id=" + strconv.Itoa(santa.Id) + ")\n"
    }

    telegramBot.SendMarkdownMessage(messageBody, chat.Id)
}

func SecretSharing(usersParticipating []telegramBot.User, pairings []telegramBot.User) {
    for i, user := range usersParticipating {
        santa := pairings[i]

        messageBody := "Your secret santa is [" + santa.FirstName + "](tg://user?id=" + strconv.Itoa(santa.Id) + ")"
        telegramBot.SendMarkdownMessage(messageBody, user.Id)
    }
}
