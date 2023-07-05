package cq

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/apex/log"
	"golang.org/x/net/websocket"
)

type MessageEvent struct {
	QQCode  int64
	QQName  string
	Message string
}

type ReqMessageFunc func(string)
type MessageEventListenerFunc func(MessageEvent, ReqMessageFunc)
type cqMessageData struct {
	QQ   int64  `json:"qq"`
	Text string `json:"text"`
}
type cqMessage struct {
	Type string        `json:"type"`
	Data cqMessageData `json:"data"`
}
type cqEvent struct {
	PoostType   string `json:"post_type"`
	MessageType string `json:"message_type"`

	UserID  int64 `json:"user_id"`
	GroupID int64 `json:"group_id"`
	Sender  struct {
		Nickname string `json:"nickname"`
		Card     string `json:"card"`
	} `json:"sender"`

	Message interface{} `json:"message"`
}

type cqSend struct {
	Action string      `json:"action"`
	Params interface{} `json:"params"`
	//UserID  int64                  `json:"user_id"`
	//Message []cqMessage            `json:"message"`
}

type ParamsUserMessage struct {
	UserID  int64       `json:"user_id"`
	Message []cqMessage `json:"message"`
}
type ParamsGroupMessage struct {
	GroupID int64       `json:"group_id"`
	Message []cqMessage `json:"message"`
}

var (
	MessageEventListener MessageEventListenerFunc = func(me MessageEvent, rmf ReqMessageFunc) {}
)

// 使字符串变得适合发送给程序,并返回是否需要处理
func toNorString(m interface{}) (s string, ok bool) {
	switch m := m.(type) {
	case string:
		s = m
		news := strings.ReplaceAll(s, fmt.Sprintf("[CQ:at,qq=%d]", ThisQQ), "")
		ok = len(news) != len(s)
	case []cqMessage:
		for _, v := range m {
			switch v.Type {
			case "text":
				s += v.Data.Text
			case "at":
				if v.Data.QQ == ThisQQ {
					ok = true
				}
			}
		}
	}

	return
}

func listen(ws *websocket.Conn) {
	defer ws.Close()
	var data []byte = make([]byte, 2048)
	for {
		l, err := ws.Read(data)
		if err != nil {
			log.Warn("CQHttp的连接已断开")
			break
		}
		var Evt cqEvent
		json.Unmarshal(data[:l], &Evt)

		switch Evt.PoostType {
		case "message":
			messagestr, ok := toNorString(Evt.Message)
			if messagestr == "" {
				continue
			}
			var is_group bool = Evt.MessageType == "group"

			if ok || !is_group {
				var name string = Evt.Sender.Card
				if name == "" {
					name = Evt.Sender.Nickname
				}
				if name == "" {
					name = fmt.Sprint(Evt.UserID)
				}
				log.Info(fmt.Sprintf("收到%s的消息：%s", name, messagestr))
				MessageEventListener(MessageEvent{
					QQCode:  Evt.UserID,
					QQName:  name,
					Message: messagestr,
				}, func(s string) {
					log.Info(fmt.Sprintf("给予%s回复：%s", name, s))
					sendObj := cqSend{}
					if is_group {
						sendObj.Action = "send_group_msg"
						sendObj.Params = ParamsGroupMessage{
							GroupID: Evt.GroupID,
							Message: []cqMessage{
								{
									Type: "at",
									Data: cqMessageData{
										QQ: Evt.UserID,
									},
								},
								{
									Type: "text",
									Data: cqMessageData{
										Text: s,
									},
								},
							},
						}

					} else {
						sendObj.Action = "send_private_msg"
						sendObj.Params = ParamsUserMessage{
							UserID: Evt.UserID,
							Message: []cqMessage{
								{
									Type: "text",
									Data: cqMessageData{
										Text: s,
									},
								},
							},
						}
					}

					senddata, err := json.Marshal(sendObj)
					if err != nil {
						log.Error(err.Error())
					}
					_, err = ws.Write(senddata)
					if err != nil {
						ReConnect()
						ws.Write(senddata)
					}
				})
			}
		}
	}
	ReConnect()
}
