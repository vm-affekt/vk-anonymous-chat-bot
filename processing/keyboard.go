package processing

import "vk-anonymous-chat-bot/vkapi"

type ButtonPayload struct {
	CommandName commandName `json:"command"`
}

func outOfChatKeyboard() *vkapi.Keyboard {
	return &vkapi.Keyboard{
		OneTime: true,
		Buttons: [][]vkapi.Button{
			[]vkapi.Button{
				{
					Action: vkapi.Action{
						Type: "text",
						Payload: ButtonPayload{
							CommandName: CommandFind,
						},
						Label: "Искать собеседника",
					},
					Color: vkapi.ColorButtonPrimary,
				}},
		},
	}
}

func onSearchingKeyboard() *vkapi.Keyboard {
	return &vkapi.Keyboard{
		OneTime: false,
		Buttons: [][]vkapi.Button{
			[]vkapi.Button{
				{
					Action: vkapi.Action{
						Type: "text",
						Payload: ButtonPayload{
							CommandName: CommandStop,
						},
						Label: "Остановить",
					},
					Color: vkapi.ColorButtonNegative,
				},
			},
		},
	}
}

func inChatKeyboard() *vkapi.Keyboard {
	return &vkapi.Keyboard{
		OneTime: false,
		Buttons: [][]vkapi.Button{
			[]vkapi.Button{
				{
					Action: vkapi.Action{
						Type: "text",
						Payload: ButtonPayload{
							CommandName: CommandStop,
						},
						Label: "Остановить",
					},
					Color: vkapi.ColorButtonNegative,
				},
				{
					Action: vkapi.Action{
						Type: "text",
						Payload: ButtonPayload{
							CommandName: CommandNext,
						},
						Label: "Следующий собеседник",
					},
					Color: vkapi.ColorButtonDefault,
				}},
		},
	}
}
