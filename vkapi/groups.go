package vkapi

import "encoding/json"

type GetLongPollServerRequest struct {
	GroupID int64
}

func (r *GetLongPollServerRequest) MethodName() string {
	return "groups.getLongPollServer"
}

func (r *GetLongPollServerRequest) Params() (params RequestParams, err error) {
	params = RequestParams{
		"group_id": r.GroupID,
	}
	return
}

func (r *GetLongPollServerRequest) Run(vkAPI *VKBotAPI) (result *LongPollServerData, err error) {
	resp, err := vkAPI.SendAPIRequest(r)
	if err != nil {
		return
	}
	result = new(LongPollServerData)
	err = json.Unmarshal(resp.Response, result)
	return
}
