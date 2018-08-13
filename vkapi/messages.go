package vkapi

import (
	"encoding/json"
)

type SendMessageRequest struct {
	UserID   int64
	Content  string
	Keyboard *Keyboard
}

func (r *SendMessageRequest) MethodName() string {
	return "messages.send"
}

func (r *SendMessageRequest) Params() (params RequestParams, err error) {
	params = RequestParams{
		"peer_id": r.UserID,
		"message": r.Content,
	}
	if r.Keyboard != nil {
		var b []byte
		b, err = json.Marshal(r.Keyboard)
		if err != nil {
			return nil, err
		}
		params["keyboard"] = string(b)
	}
	return
}

func (r *SendMessageRequest) Run(vkAPI *VKBotAPI) (int64, error) {
	resp, err := vkAPI.SendAPIRequest(r)
	if err != nil {
		return 0, err
	}
	var result int64
	err = json.Unmarshal(resp.Response, &result)
	return result, err

}
