package processing

import (
	"log"
	"vk-anonymous-chat-bot/vkapi"
)

type commandName string

const (
	CommandFind commandName = "find"
	CommandStop commandName = "stop"
	CommandNext commandName = "next"
)

func (grp *VKChatGroup) onFindCommand(msg *vkapi.MessageObject) {
	userID := msg.SenderID
	if grp.chatProc.checkUserInQueue(userID) {
		log.Printf("[WARNING] Пользователь id%v уже находится в очереди, но пришел запрос на добавление его в очередь поиска!", userID)
		return
	}
	grp.chatProc.PushMemberToQueue(Member{userID})
	outMsg := vkapi.SendMessageRequest{
		UserID:   userID,
		Content:  "Вы добавлены в очередь поиска собеседника! Ждите :-)",
		Keyboard: onSearchingKeyboard(),
	}

	grp.outcomingMessages <- outMsg
}

func (grp *VKChatGroup) onStopCommand(msg *vkapi.MessageObject) {
	userID := msg.SenderID
	chat, chatIdx := grp.chatProc.findChatWithUser(userID)
	keyboard := outOfChatKeyboard()
	if chat != nil {
		opponentID := chat.getOpponentBy(userID).UserID
		grp.chatProc.removeChat(chatIdx)
		grp.outcomingMessages <- vkapi.SendMessageRequest{
			UserID:   userID,
			Content:  "Вы прервали общение с этим собеседником!",
			Keyboard: keyboard,
		}
		grp.outcomingMessages <- vkapi.SendMessageRequest{
			UserID:   opponentID,
			Content:  "Собеседник остановил беседу с вами!",
			Keyboard: keyboard,
		}
	} else {
		if err := grp.chatProc.removeMemberFromQueue(userID); err != nil {
			log.Printf("Произошла ошибка при удалении пользователя id%d из очереди: %v", userID, err)
		}
		grp.outcomingMessages <- vkapi.SendMessageRequest{
			UserID:   userID,
			Content:  "Вы остановили поиск собеседника!",
			Keyboard: keyboard,
		}

	}

}

func (grp *VKChatGroup) processCommand(msg *vkapi.MessageObject, command commandName) {
	switch command {
	case CommandFind:
		grp.onFindCommand(msg)
	case CommandStop:
		grp.onStopCommand(msg)
	}
}
