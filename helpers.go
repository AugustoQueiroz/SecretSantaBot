package main

import (
    "telegramBot"
)

func NotIn(group []telegramBot.User, user *telegramBot.User) bool {
    for _, user2 := range group {
        if user.Id == user2.Id {
            return false
        }
    }

    return true
}
