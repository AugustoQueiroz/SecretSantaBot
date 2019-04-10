package main

import (
    // Standard Packages
    "log"
    "time"
    "strings"
    "strconv"

    // Project Packages
    "telegramBot"
)

func StartHandler(message *telegramBot.Message) {
    log.Println("Here")
    if len(message.Body) <= len("/start") {
        return
    }
    // Going to assume that the start message was to join a secret santa
    chatId, err := strconv.Atoi(message.Body[len("/start")+1:])
    if err != nil {
        // The text following the command was not an id, so act as if this is an unrelated command
        return
    }

    joinChannel, exists := activeJoinChannels[chatId]
    if !exists {
        // Act as if this is an unrelated start command
        return
    }

    joinChannel <- message.From

    // Inform user they have joined the santa
    messageBody := "You have joined a secret santa"
    telegramBot.SendMarkdownMessage(messageBody, message.Origin.Id)
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

func NewSantaMessage(body string, chatId int) {
    // Create the inline buttons that will be displayed
    var inlineKeyboard telegramBot.InlineKeyboardMarkup

    var participateButton telegramBot.InlineKeyboardButton
    participateButton.Label = "Participate"
    participateButton.URL = "http://t.me/secretsantainatorbot?start=" + strconv.Itoa(chatId)
    log.Println(participateButton.URL)

    var closeButton telegramBot.InlineKeyboardButton
    closeButton.Label = "Done"
    closeButton.CallbackData = "done:" + strconv.Itoa(chatId)

    inlineKeyboard.Keyboard = [][]telegramBot.InlineKeyboardButton {
        []telegramBot.InlineKeyboardButton { participateButton },
        []telegramBot.InlineKeyboardButton { closeButton },
    }

    // Send the message
    telegramBot.SendMessageWithKeyboard(body, chatId, "Markdown", inlineKeyboard)
}

func OpenSantaHandler(message *telegramBot.Message) {
    log.Println("Create open santa command received")
    if !GuardMessageBelongsToGroup(message) { return } // Guarantee that the message was received in a group, or deal with it otherwise

    // If the message was, indeed, received in a group, handle it

    // First, check if there isn't a secret santa already accepting participants here
    _, exists := activeJoinChannels[message.Origin.Id]
    if exists {
        // Say there is already a santa going on and the quit
        return
    }

    // If one does not exist, create the go routine and channel
    activeJoinChannels[message.Origin.Id] = make(JoinChannel)
    go SantaWriteUp("open", message.Origin, activeJoinChannels[message.Origin.Id])
    timer := time.NewTimer(60*time.Second)
    go func() {
        <-timer.C
        close(activeJoinChannels[message.Origin.Id])
    }()

    // Then send a message with the possible actions
    replyBody := "An *Open Santa* was created by " + message.From.FirstName + ".\n" +
                 "To participate in it, press 'Participate' below.\n" +
                 "When everyone has joined, press 'Close' to create the pairings, which will be sent here."

    NewSantaMessage(replyBody, message.Origin.Id)
}

func CallbackHandler(callback *telegramBot.CallbackQuery) {
    if strings.HasPrefix(callback.Data, "done:") {
        // Done callback recieved
        chatId, err := strconv.Atoi(callback.Data[len("done:"):])
        if err != nil {
            return
        }

        joinChannel, exists := activeJoinChannels[chatId]
        if exists {
            joinChannel <- nil
        }
    }
}
