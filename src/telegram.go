package main

import (
	"fmt"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func sendTelegramMessage(message string) (err error) {
	// Send message to Telegram user

	var telegramBotToken string
	var telegramChatID string

	config, err := ConfigParser()

	if err != nil {
		log.Error("Error reading config file: ", err)
		log.Warnf("Unable to send Telegram message. Message: %s", message)
		return fmt.Errorf("error reading config file: %v", err)
	}

	if !config.Integrations.Telegram.Enabled {
		log.Warn("Telegram integration was called even though it is not enabled. Not attempting to send message.")
		return fmt.Errorf("telegram integration is not enabled")
	}

	// set vars
	telegramBotToken = config.Integrations.Telegram.TelegramBotToken
	telegramChatID = config.Integrations.Telegram.TelegramChatID

	// create bot
	bot, err := tgbot.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Error("Error creating bot: ", err)
		return fmt.Errorf("error creating bot: %v", err)
	}
	bot.Debug = false // just to make sure we are not using dev mode.

	// create message
	log.Infof("Telegram integration is authorised on username: %s", bot.Self.UserName)

	msg := tgbot.NewMessageToChannel(telegramChatID, message)
	msg.ParseMode = "markdown"

	// send message
	_, err = bot.Send(msg)
	if err != nil {
		log.Error("Error sending message: ", err)
		return fmt.Errorf("error sending message: %v", err)
	}

	log.Infof("Message sent to Telegram chat ID: %s from username: %s", telegramChatID, bot.Self.UserName)
	return nil
}
