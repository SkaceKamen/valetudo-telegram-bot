package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/SkaceKamen/valetudo-telegram-bot/pkg/valetudo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	robotApi    *valetudo.ValetudoClient
	telegramApi *tgbotapi.BotAPI
	chatIds     []int64
}

func NewBot(robotApi *valetudo.ValetudoClient, telegramApi *tgbotapi.BotAPI) Bot {
	return Bot{robotApi: robotApi, telegramApi: telegramApi}
}

func (bot *Bot) AddUserId(id int64) {
	bot.chatIds = append(bot.chatIds, id)
}

func (bot *Bot) Start() error {
	err := bot.publishMyCommands()

	if err != nil {
		// TODO: Custom error
		return err
	}

	go func() {
		err = bot.listenToStateChanges()
		if err != nil {
			log.Println(err)
		}
	}()

	err = bot.listenToMessages()
	if err != nil {
		return err
	}

	return nil
}

func (bot *Bot) listenToStateChanges() error {
	lastState, err := bot.getParsedState()

	if err != nil {
		return err
	}

	bot.robotApi.ListenToStateChanges(func(state *valetudo.RobotState) {
		parsed := stateObjToData(state)

		log.Println("Received state, status: ", parsed.Status, " batteryStatus:", parsed.BatteryStatus, " batteryLevel:", parsed.BatteryLevel)

		if lastState.BatteryStatus != parsed.BatteryStatus {
			bot.handleBatteryStatusChange(lastState, parsed)
		}

		if lastState.Status != parsed.Status {
			bot.handleStatusChange(lastState, parsed)
		}

		lastState = parsed
	})

	return nil
}

func (bot *Bot) handleStatusChange(previous *CurrentState, new *CurrentState) {
	for _, user := range bot.chatIds {
		newStatusLabel := localizeRobotStatus(new.Status)
		newStatusIcon := robotStatusEmoji(new.Status)
		statusMessage := newStatusIcon + " " + newStatusLabel

		// Special status transitions that aren't actually a separate statuses
		switch new.Status {
		case "returning":
			if previous.Status == "cleaning" {
				statusMessage = "✅ Cleaning complete, returning home"
			}
		}

		bot.Send(user, statusMessage)
	}
}

func (bot *Bot) handleBatteryStatusChange(previous *CurrentState, new *CurrentState) {
	for _, user := range bot.chatIds {
		statusMessage := ""

		switch new.BatteryStatus {
		case "charging":
			statusMessage = fmt.Sprintf("🪫 Charging battery from %d %%", new.BatteryLevel)
		case "charged":
			statusMessage = "🔋 Battery fully charged"
		}

		if statusMessage != "" {
			bot.Send(user, statusMessage)
		}
	}
}

func (bot *Bot) isAllowedUserId(id int64) bool {
	for _, allowedId := range bot.chatIds {
		if allowedId == id {
			return true
		}
	}

	return false
}

func (bot *Bot) listenToMessages() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.telegramApi.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			if !bot.isAllowedUserId(update.CallbackQuery.Message.Chat.ID) {
				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "You are not allowed to do that")
				if _, err := bot.telegramApi.Request(callback); err != nil {
					log.Println(err)
				}

				continue
			}

			data := strings.Split(update.CallbackQuery.Data, " ")
			switch data[0] {
			case "pause":
				err := bot.robotApi.Pause()
				if err != nil {
					log.Println(err)
					bot.Send(update.CallbackQuery.Message.Chat.ID, "❌ Error pausing robot: "+err.Error())
				} else {
					callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "⏸ Paused")
					if _, err := bot.telegramApi.Request(callback); err != nil {
						log.Println(err)
					}
				}
			case "stop":
				err := bot.robotApi.Stop()
				if err != nil {
					log.Println(err)
					bot.Send(update.CallbackQuery.Message.Chat.ID, "❌ Error stopping robot: "+err.Error())
				} else {
					callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "⏹ Stopped")
					if _, err := bot.telegramApi.Request(callback); err != nil {
						log.Println(err)
					}
				}
			case "home":
				err := bot.robotApi.Home()
				if err != nil {
					log.Println(err)
					bot.Send(update.CallbackQuery.Message.Chat.ID, "❌ Error sending robot home: "+err.Error())
				} else {
					callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "🏠 Going home")
					if _, err := bot.telegramApi.Request(callback); err != nil {
						log.Println(err)
					}
				}
			case "clean":
				if len(data) < 2 {
					bot.handleCleanCommand(update.CallbackQuery.Message.Chat.ID, "")
					callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "You need to pick what")
					if _, err := bot.telegramApi.Request(callback); err != nil {
						log.Println(err)
					}

					continue
				}

				err := bot.robotApi.CleanMapSegments([]string{data[1]}, 1)
				if err != nil {
					bot.Send(update.CallbackQuery.Message.Chat.ID, "❌ Error cleaning room: "+err.Error())
					log.Println(err)
				} else {
					callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "🧹 Cleaning")
					if _, err := bot.telegramApi.Request(callback); err != nil {
						log.Println(err)
					}
				}
			}

			continue
		}

		if update.Message == nil {
			continue
		}

		if !bot.isAllowedUserId(update.Message.From.ID) {
			continue
		}

		// log.Printf("[%d] %s", update.Message.Chat.ID, update.Message.Text)

		if !update.Message.IsCommand() {
			continue
		}

		switch update.Message.Command() {
		case "start":
			bot.chatIds = append(bot.chatIds, update.Message.Chat.ID)
			bot.Send(update.Message.Chat.ID, "✔️ Subscribed to notifications")
		case "status":
			err := bot.handleStatusCommand(update.Message.Chat.ID, update.Message.CommandArguments())
			if err != nil {
				log.Println(err)
				bot.Send(update.Message.Chat.ID, "❌ Error fetching status: "+err.Error())
			}
		case "stop":
			err := bot.robotApi.Stop()
			if err != nil {
				log.Println(err)
				bot.Send(update.Message.Chat.ID, "❌ Error stopping robot: "+err.Error())
			}
		case "home":
			err := bot.robotApi.Home()
			if err != nil {
				log.Println(err)
				bot.Send(update.Message.Chat.ID, "❌ Error sending robot home: "+err.Error())
			}
		case "pause":
			err := bot.robotApi.Pause()
			if err != nil {
				log.Println(err)
				bot.Send(update.Message.Chat.ID, "❌ Error pausing robot: "+err.Error())
			}
		case "clean":
			err := bot.handleCleanCommand(update.Message.Chat.ID, update.Message.CommandArguments())
			if err != nil {
				log.Println(err)
				bot.Send(update.Message.Chat.ID, "❌ Error cleaning: "+err.Error())
			}
		}
	}

	return nil
}
