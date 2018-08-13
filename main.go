package main

import (
	"fmt"
	"log"
	"time"
	"vk-anonymous-chat-bot/processing"
)

func msgProcessor(messages <-chan string) {
	for m := range messages {
		fmt.Println("msgProcessor: " + m)
		time.Sleep(3 * time.Second)
	}
	fmt.Println("msgProcessor() END!")
}

func main() {
	group := processing.NewVKChatGroup("836f3876149cabf78dd584bc9732736e9197ec21d0d8a2cc3b168d607c07f41a8d2ffb93db674d4b65412", 82448081) // 1 8
	err := group.SendMessageTo(30284936)
	if err != nil {
		log.Println("An error occured: ", err.Error())
	}

	go func() {
		err := group.Start()
		if err != nil {
			fmt.Println("on long polling listener error occured:", err.Error())
		}
	}()

	fmt.Scanln()
}
