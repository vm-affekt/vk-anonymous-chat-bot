package processing

import (
	"errors"
	"log"
	"math/rand"
	"time"
	"vk-anonymous-chat-bot/vkapi"
)

var (
	ErrUserOutOfChat  = errors.New("user is not found in chat")
	ErrUserOutOfQueue = errors.New("user is not found in queue")
)

type ChatProc struct {
	chats []Chat
	queue []Member

	outcomingMsg chan<- vkapi.SendMessageRequest
}

func (cp *ChatProc) removeMemberFromQueueByIdx(index int) {
	cp.queue = append(cp.queue[:index], cp.queue[index+1:]...)
}

func (cp *ChatProc) removeMemberFromQueue(userID int64) error {
	for i, m := range cp.queue {
		if m.UserID == userID {
			cp.removeMemberFromQueueByIdx(i)
			return nil
		}
	}
	return ErrUserOutOfQueue
}

func (cp *ChatProc) createChat(firstMember, secondMember Member) {
	chat := Chat{
		FirstMember:  firstMember,
		SecondMember: secondMember,
	}
	cp.chats = append(cp.chats, chat)
	msg := "Мы нашли вам собеседника! Общайтесь ;-)"
	keyboard := inChatKeyboard()
	cp.outcomingMsg <- vkapi.SendMessageRequest{
		UserID:   firstMember.UserID,
		Content:  msg,
		Keyboard: keyboard,
	}
	cp.outcomingMsg <- vkapi.SendMessageRequest{
		UserID:   secondMember.UserID,
		Content:  msg,
		Keyboard: keyboard,
	}
}

func (cp *ChatProc) removeChat(index int) {
	c := &cp.chats[index]
	cp.chats = append(cp.chats[:index], cp.chats[index+1:]...)
	log.Printf("Чат с индексом %d удален. (первый участник: id%d ; второй участник: id%d)\n", index, c.FirstMember.UserID, c.SecondMember.UserID)
}

func (cp *ChatProc) DistributeByChatsWorker(tickerCh <-chan time.Time) {
	for _ = range tickerCh {
		mCount := len(cp.queue)
		chatsCount := int(mCount / 2)
		if chatsCount == 0 {
			continue
		}
		log.Printf("Начало распределения участников по чатам. Всего должно быть %d чатов (в очереди %d участников)\n", chatsCount, mCount)
		for i := 0; i < chatsCount; i++ {
			rand.Seed(time.Now().Unix())
			fIdx := rand.Intn(len(cp.queue))
			firstMember := cp.queue[fIdx]
			cp.removeMemberFromQueueByIdx(fIdx)
			sIdx := rand.Intn(len(cp.queue))
			secondMember := cp.queue[sIdx]
			cp.removeMemberFromQueueByIdx(sIdx)
			cp.createChat(firstMember, secondMember)
			log.Printf("Чат создан! Первый участник id:%d [индекс:%d], второй участник id:%d [индекс:%d (после удаления)]\n",
				firstMember.UserID, fIdx, secondMember.UserID, sIdx)
		}
	}
}

func (cp *ChatProc) PushMemberToQueue(member Member) {
	cp.queue = append(cp.queue, member)
	log.Printf("Пользователь id%d добавлен в очередь поиска! Размер очереди на данный момент: %d\n", member.UserID, len(cp.queue))
}

func (cp *ChatProc) checkUserInQueue(userID int64) bool {
	for _, m := range cp.queue {
		if m.UserID == userID {
			return true
		}
	}
	return false
}

// findChatWithUser ищет в чатах пользователя с идентификатором userID и в случае успеха возвращает ссылку на чат, в котром был найден.
// Вторым результатом возвращается индекс чата в слайсе
func (cp *ChatProc) findChatWithUser(userID int64) (*Chat, int) {
	for i, c := range cp.chats {
		if c.FirstMember.UserID == userID || c.SecondMember.UserID == userID {
			return &c, i
		}
	}
	return nil, -1
}

func (cp *ChatProc) ProcessSimpleMessage(msg *vkapi.MessageObject) error {
	senderID := msg.SenderID
	chat, _ := cp.findChatWithUser(senderID)
	if chat == nil {
		return ErrUserOutOfChat
	}
	return chat.ProcessUserMessage(msg, cp.outcomingMsg)
}

type Member struct {
	UserID int64
}

type Chat struct {
	FirstMember  Member
	SecondMember Member

	CreatedDate time.Time
}

func (c *Chat) getOpponentBy(userID int64) *Member {
	switch userID {
	case c.FirstMember.UserID:
		return &c.SecondMember
	case c.SecondMember.UserID:
		return &c.FirstMember
	default:
		return nil
	}
}

func (c *Chat) ProcessUserMessage(msg *vkapi.MessageObject, outMsgCh chan<- vkapi.SendMessageRequest) error {
	opponent := c.getOpponentBy(msg.SenderID)
	if opponent == nil {
		return errors.New("opponent not found for this chat")
	}
	outMsgCh <- vkapi.SendMessageRequest{
		UserID:  opponent.UserID,
		Content: msg.Content,
	}
	log.Printf("Сообщение отправлено пользователю с id:%d от пользователя с id:%d. Текст сообщения: %s\n",
		opponent.UserID, msg.SenderID, msg.Content)
	return nil
}
