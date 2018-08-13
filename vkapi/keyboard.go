package vkapi

import (
	"encoding/json"
)

type ButtonColor string

const (
	ColorButtonPrimary  ButtonColor = "primary"
	ColorButtonDefault  ButtonColor = "default"
	ColorButtonNegative ButtonColor = "negative"
	ColorButtonPositive ButtonColor = "positive"
)

type Keyboard struct {
	OneTime bool       `json:"one_time"`
	Buttons [][]Button `json:"buttons"`
}

type Button struct {
	Action Action      `json:"action"`
	Color  ButtonColor `json:"color"`
}

type Action struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
	Label   string      `json:"label"`
}

func (a *Action) MarshalJSON() ([]byte, error) {
	payloadBytes, err := json.Marshal(a.Payload)
	if err != nil {
		return nil, err
	}

	type alias Action
	return json.Marshal(&struct {
		Payload string `json:"payload"`
		*alias
	}{
		Payload: string(payloadBytes),
		alias:   (*alias)(a),
	})
}
