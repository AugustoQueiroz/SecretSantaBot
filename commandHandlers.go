package main

import (
    // Standard Packages
    "log"
    "time"
    "strings"
    "strconv"

    // Project Packages
    "github.com/AugustoQueiroz/telegramBot"
)

func StartHandler(message *telegramBot.Message) {
    if len(message.Body) <= len("/start") {
        return
    }
    // Going to assume that the start message was to join a secret santa
    messageId, err := strconv.Atoi(message.Body[len("/start")+1:])
    if err != nil {
        // The text following the command was not an id, so act as if this is an unrelated command
        return
    }

    JoinSanta(messageId, message.From)
}

func GuardMessageBelongsToGroup(message *telegramBot.Message) bool {
    if message.Origin.Type == "channel" { return false } // If the message is comming from a channel, there is nothing to be done, just leave
    if message.Origin.Type == "private" {
        // If the message is coming from a private conversation, inform about correct usage and leave
        replyBody := "This command is meant to be send in *groups* or *supergroups*, and not on a private conversation.\n\nTo learn more about the usage of this bot send /help"
        telegramBot.SendMarkdownMessage(replyBody, message.Origin.Id)
        return false
    }

    return true
}

func NewSantaMessage(body string, chatId int, messageId int) {
    // Create the inline buttons that will be displayed
    var inlineKeyboard telegramBot.InlineKeyboardMarkup

    var participateButton telegramBot.InlineKeyboardButton
    participateButton.Label = "Participate"
    participateButton.URL = "http://t.me/secretsantainatorbot?start=" + strconv.Itoa(messageId)

    var closeButton telegramBot.InlineKeyboardButton
    closeButton.Label = "Done"
    closeButton.CallbackData = "done:" + strconv.Itoa(messageId)

    inlineKeyboard.Keyboard = [][]telegramBot.InlineKeyboardButton {
        []telegramBot.InlineKeyboardButton { participateButton },
        []telegramBot.InlineKeyboardButton { closeButton },
    }

    // Send the message
    telegramBot.SendMessageWithKeyboard(body, chatId, "Markdown", inlineKeyboard)
}

func OpenSantaHandler(message *telegramBot.Message) {
    if !GuardMessageBelongsToGroup(message) { return } // Guarantee that the message was received in a group, or deal with it otherwise

    // If the message was, indeed, received in a group, handle it

    // First, check if there isn't a secret santa already accepting participants here
    _, exists := activeSantas[message.Id]
    if exists {
        // Say there is already a santa going on and the quit
        return
    }

    // Then send a message with the possible actions
    replyBody := "An *Open Santa* was created by " + message.From.FirstName + ".\n" +
                 "To participate in it, press 'Participate' below.\n" +
                 "When everyone has joined, press 'Close' to create the pairings, which will be sent here."

    NewSantaMessage(replyBody, message.Origin.Id, message.Id)
    participantsMessage := telegramBot.SendMarkdownMessage("#participants: 0", message.Origin.Id)

    // If one does not exist, create the go routine and channel
    activeSantas[message.Id] = &SantaInfo { "open", *participantsMessage.Origin, participantsMessage.Id, []telegramBot.User{} }
    timer := time.NewTimer(60*time.Second)
    go func() {
        <-timer.C
        SantaDone(message.Id)
    }()
}

func SecretSantaHandler(message *telegramBot.Message) {
    if !GuardMessageBelongsToGroup(message) { return } // Guarantee that the message was received in a group, or deal with it otherwise

    // If the message was, indeed, received in a group, hanlde it

    // First, check if there isn't a secret santa already accepting participants here
    _, exists := activeSantas[message.Id]
    if exists {
        // Say there is already a santa going on and then quit
        return
    }

    // Then send a messsage with the possible actions
    replyBody := "A *Secret Santa* was created by *" + message.From.FirstName + "*.\n" +
                 "To participate in it, press 'Participate' below.\n" +
                 "When everyone has joined, press 'Close' to create the pairings, which will be sent here."

    NewSantaMessage(replyBody, message.Origin.Id, message.Id)
    participantsMessage := telegramBot.SendMarkdownMessage("#participants: 0", message.Origin.Id)

    // Then create the active santa
    activeSantas[message.Id] = &SantaInfo { "secret", *participantsMessage.Origin, participantsMessage.Id, []telegramBot.User{} }
    timer := time.NewTimer(60*time.Second)
    go func() {
        <-timer.C
        SantaDone(message.Id)
    }()
}

func CallbackHandler(callback *telegramBot.CallbackQuery) {
    if callback == nil { return }
    if strings.HasPrefix(callback.Data, "done:") {
        // Done callback recieved
        chatId, err := strconv.Atoi(callback.Data[len("done:"):])
        if err != nil {
            return
        }

        _, exists := activeSantas[chatId]
        if !exists {
            return
        }

        SantaDone(chatId)
    }
}
