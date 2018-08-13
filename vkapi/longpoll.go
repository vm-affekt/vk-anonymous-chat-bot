package vkapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type LongPollServerData struct {
	Key    string `json:"key"`
	Server string `json:"server"`
	TS     int    `json:"ts"`
}

func (api *VKBotAPI) initLongPollConnection() error {
	req := GetLongPollServerRequest{api.groupID}
	data, err := req.Run(api)
	if err != nil {
		return err
	}
	api.longPoolData = *data
	return nil
}

type LongPollResponse struct {
	TS      int              `json:"ts,string"`
	Updates []LongPollUpdate `json:"updates"`

	FailedCode int `json:"failed,omitempty"`
}

const (
	LPHistoryOld = 1
	LPKeyExpired = 2
	LPInfoLost   = 3
)

type LongPollUpdate struct {
	Type   string          `json:"type"`
	Object json.RawMessage `json:"object"`
}

type MessageObject struct {
	SenderID int64  `json:"from_id"`
	Content  string `json:"text"`
	Payload  string `json:"payload"`
}

func (api *VKBotAPI) processLongPollUpdates(lpActions <-chan LongPollUpdate, messages chan<- MessageObject) {
	for action := range lpActions {
		switch action.Type {
		case "message_new":
			if messages != nil {
				var msg MessageObject
				err := json.Unmarshal(action.Object, &msg)
				if err != nil {
					log.Println("Произошла ошибка при парсинге нового сообщения (message_new) из ответа LongPoll-сервера: " + err.Error())
					continue
				}
				log.Printf("Получено новое сообщение от пользователя id%d. Тело сообщения: %s\n", msg.SenderID, msg.Content)
				messages <- msg
			}

		default:
			log.Printf("Перехвачен необрабатываемый тип (%s) LongPoll-сообщения\n", action.Type)
		}

	}
}

func (api *VKBotAPI) OnLongPoolMessage(messages chan<- MessageObject) error {
	var err error
	if api.longPoolData.Key == "" {
		if err = api.initLongPollConnection(); err != nil {
			return err
		}
	}
	try := 0

	lpActions := make(chan LongPollUpdate, 100) // TODO: Размер буфера вынести в конфиг
	defer close(lpActions)
	go api.processLongPollUpdates(lpActions, messages)

	client := http.Client{}
	const MaxRequestAttempts = 5
	for {
		if try == MaxRequestAttempts {
			log.Fatalf("Количество неуспешных подряд обращений к LongPoll-серверу достигло максимума (%d попыток). Описание последней ошибки: %v\n",
				MaxRequestAttempts, err)
		}
		lpData := &api.longPoolData
		url := fmt.Sprintf("%s?act=a_check&key=%s&ts=%d&wait=%d",
			lpData.Server,
			lpData.Key,
			lpData.TS,
			25, // TODO: Вынести в конфиг
		)
		log.Println("Выполнение запроса к LongPoll-серверу по пути", url)
		// url := fmt.Sprintf(urlFmt, vkAPIEndpoint, method, api.accessToken, vkAPIVersion)
		var resp *http.Response
		resp, err = client.Get(url)
		if err != nil {
			log.Println("Произошла ошибка при попытке выполнить GET-запрос к LongPoll-серверу:", err)
			try++
			continue
		}
		var respBytes []byte
		respBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Произошла ошибка при попытке прочитать тело HTTP-ответа от LongPoll-сервера:", err)
			try++
			continue
		}
		log.Println("Ответ от LongPoll-сервера успешно получен. Тело ответа:", string(respBytes))
		var respModel LongPollResponse
		err = json.Unmarshal(respBytes, &respModel)
		if err != nil {
			log.Println("Произошла ошибка при десериализации JSON-ответа от LongPoll-сервера:", err)
			try++
			continue
		}
		if respModel.FailedCode > 0 {
			log.Println("Ответ от LongPoll-сервера содержит ошибку. Код ошибки:", respModel.FailedCode)
			switch respModel.FailedCode {
			case LPKeyExpired, LPInfoLost:
				log.Println("LongPoll-ключ устарел. Получение нового...")
				if err = api.initLongPollConnection(); err != nil {
					log.Println("Произошла ошибка при обновлении LongPoll-соединения. Причина:", err)
					try++
					continue
				}

			case LPHistoryOld:

				log.Println("История longPoll-событий устарела или была частично утеряна")
				lpData.TS = respModel.TS
				try++
				continue
			}

		}

		for _, upd := range respModel.Updates {
			lpActions <- upd
		}
		try = 0

		lpData.TS = respModel.TS

	}
}
