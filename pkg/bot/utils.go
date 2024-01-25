package bot

import (
	"github.com/SkaceKamen/valetudo-telegram-bot/pkg/valetudo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CurrentStateAttachmentState struct {
	Type     string
	Attached bool
}

type CurrentState struct {
	BatteryStatus       string
	BatteryLevel        int
	Status              string
	Map                 valetudo.RobotStateMap
	Attachments         []CurrentStateAttachmentState
	AttachedAttachments []string
}

func (bot *Bot) getParsedState() (*CurrentState, error) {
	robotState, err := bot.robotApi.GetRobotState()

	if err != nil {
		return nil, err
	}

	return stateObjToData(robotState), nil
}

func (bot *Bot) getRooms() (*[]valetudo.RobotStateMapLayer, error) {
	state, err := bot.robotApi.GetRobotState()

	if err != nil {
		return nil, err
	}

	result := []valetudo.RobotStateMapLayer{}

	for _, layer := range state.Map.Layers {
		if layer.Type == "segment" {
			result = append(result, layer)
		}
	}

	return &result, nil
}

func stateObjToData(state *valetudo.RobotState) *CurrentState {
	result := CurrentState{}

	for _, attribute := range state.Attributes {
		if attribute.Class == "BatteryStateAttribute" {
			if attribute.Flag != nil {
				result.BatteryStatus = *attribute.Flag
			}
			if attribute.Level != nil {
				result.BatteryLevel = *attribute.Level
			}
		}

		if attribute.Class == "StatusStateAttribute" {
			if attribute.Value != nil {
				result.Status = *attribute.Value
			}
		}

		if attribute.Class == "AttachmentStateAttribute" {
			if attribute.Type != nil && attribute.Attached != nil {
				result.Attachments = append(result.Attachments, CurrentStateAttachmentState{
					Type:     *attribute.Type,
					Attached: *attribute.Attached,
				})

				if *attribute.Attached {
					result.AttachedAttachments = append(result.AttachedAttachments, *attribute.Type)
				}
			}
		}
	}

	result.Map = state.Map

	return &result
}

func (bot *Bot) Send(receiverId int64, message string) error {
	_, err := bot.telegramApi.Send(tgbotapi.NewMessage(receiverId, message))

	return err
}
