package vkapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/CossackPyra/pyraconv"
)

const (
	vkAPIEndpoint = "https://api.vk.com/method"
	vkAPIVersion  = "5.80"
)

type VKBotAPI struct {
	accessToken string
	groupID     int64
	httpClient  http.Client

	longPoolData LongPollServerData
}

func NewVKBotAPI(accessToken string, groupID int64) *VKBotAPI {
	//preparedURL := fmt.Sprintf("%s/%s?access_token=%s&v=%s", vkAPIEndpoint, method, AccessToken, VKAPIVersion)
	return &VKBotAPI{
		accessToken: accessToken,
		groupID:     groupID,
	}
}

type VKAPIResponse struct {
	Response json.RawMessage `json:"response,omitempty"`
}

func (api *VKBotAPI) SendAPIRequest(apiReq VKAPIRequest) (*VKAPIResponse, error) {
	params, err := apiReq.Params()
	if err != nil {
		return nil, fmt.Errorf("Невозможно получить параметры этого API-запроса: %s", err.Error())
	}
	return api.SendAPIRequestByParams(apiReq.MethodName(), params)
}

func (api *VKBotAPI) SendAPIRequestByParams(method string, params map[string]interface{}) (*VKAPIResponse, error) {
	url := fmt.Sprintf("%s/%s?access_token=%s&v=%s", vkAPIEndpoint, method, api.accessToken, vkAPIVersion)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if params != nil {
		reqValues := req.URL.Query()
		for param, value := range params {
			reqValues.Add(param, pyraconv.ToString(value))
		}
		req.URL.RawQuery = reqValues.Encode()
	}
	log.Printf("Отправка GET-запроса на %v ...\n", req.URL.String())

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	resultBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("VK API сервер вернул ошибочный HTTP-код: %d", resp.StatusCode)
	}

	log.Println("Тело запроса к VK API: ", string(resultBytes))

	var result VKAPIResponse
	err = json.Unmarshal(resultBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, err

}
