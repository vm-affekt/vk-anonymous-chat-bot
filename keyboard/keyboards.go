package keyboard

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"vkApi/vkapi"
)

type keyboardFileName string

const (
	keyboardDefsDir = "keyboard_defs"

	inChatKeyboardFile keyboardFileName = "inchat.json"
)

var (
	inChatKeyboard *vkapi.Keyboard
)

func loadKeyboardFromFile(fileName keyboardFileName) (*vkapi.Keyboard, error) {
	path := filepath.Join(keyboardDefsDir, string(fileName))
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var keyboard vkapi.Keyboard
	if err := json.Unmarshal(b, &keyboard); err != nil {
		return nil, err
	}
	return &keyboard, nil
}

func InChatKeyboard() *vkapi.Keyboard {
	// var err error
	// if inChatKeyboard == nil {
	// 	if inChatKeyboard, err = loadKeyboardFromFile(inChatKeyboardFile); err != nil {
	// 		log.Fatalf("An error occured while loading inChat keyboard: %v", err)
	// 	}
	// }
	// return inChatKeyboard

	// return &vkapi.Keyboard{
	// 	OneTime: false,
	// 	Buttons: [][]vkapi.Button{
	// 		[]vkapi.Button{{
	// 			Action: vkapi.Action{
	// 				Type: "text",
	// 				Payload: ButtonPayload{
	// 					CommandName: CommandFind,
	// 				},
	// 				Label: "Следующий собеседник",
	// 			},
	// 			Color: vkapi.ColorButtonPrimary,
	// 		}},
	// 	},
	return nil
}
