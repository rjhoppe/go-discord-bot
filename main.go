package main

import (
	"fmt"
	"rjhoppe/go-discord-bot/bot"
	"rjhoppe/go-discord-bot/config"
)

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bot.Start()

	<-make(chan struct{})
	return
}
