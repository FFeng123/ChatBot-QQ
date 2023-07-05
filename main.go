package main

import (
	"bot/chat"
	"bot/cq"
	"fmt"
	"os"
	"os/signal"

	"github.com/apex/log"
)

var chatContext *chat.ChatContext

func main() {
	log.Info("ChatBot启动···")

	ReadConfig()

	log.Info("连接到CQHttp···")

	cq.ThisQQ = Config.ThisQQ
	chatContext = chat.NewContext(Config.ChatAPI)
	chatContext.API_max_length = Config.ChatAPIArgs.API_max_length
	chatContext.API_temperature = Config.ChatAPIArgs.API_temperature
	chatContext.API_top_p = Config.ChatAPIArgs.API_top_p
	chatContext.MemMaxLength = Config.MemLength
	chatContext.StaticMem = Config.StaticMem

	cq.Connect(Config.CQHttpAddr)

	cq.MessageEventListener = MessageEvt

	log.Info("已启动")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func MessageEvt(m cq.MessageEvent, req cq.ReqMessageFunc) {
	rep := chatContext.Chat(fmt.Sprintf("%s：%s", m.QQName, m.Message), false, true)
	if rep == "" {
		rep = "ごめんなさい、桜今はちょっと...少々お待ちください。"
	}

	req(rep)
}
