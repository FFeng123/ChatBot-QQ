package chat

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/apex/log"
)

type ApiSend1 struct {
	Prompt      string      `json:"prompt"`
	History     [][2]string `json:"history"`
	MaxLength   int         `json:"max_length"`
	TopP        float32     `json:"top_p"`
	Temperature float32     `json:"temperature"`
}
type ApiGet1 struct {
	Response string `json:"response"`
	//History  []string `json:"history"`
	Status int    `json:"status"`
	Time   string `json:"time"`
}

type ChatContext struct {
	addr string

	MemMaxLength    int
	API_max_length  int
	API_top_p       float32
	API_temperature float32

	StaticMem [][2]string
	mem       [][2]string
}

func NewContext(Addr string) *ChatContext {
	return &ChatContext{
		addr: Addr,
	}
}

func (c *ChatContext) CleanMem() {
	c.mem = [][2]string(c.StaticMem)
}

func (c *ChatContext) Chat(message string, addMemOnly bool, messageAddMem bool) (replay string) {
	if !addMemOnly {
		c.mem = append(c.mem, c.StaticMem...)
		data, err := json.Marshal(ApiSend1{
			Prompt:  message,
			History: c.mem,
		})
		c.mem = c.mem[:len(c.mem)-len(c.StaticMem)]
		if err != nil {
			log.Error(err.Error())
			return
		}
		req, err := http.Post(c.addr, "application/json", bytes.NewBuffer(data))
		if err != nil {
			log.Error(err.Error())
			return
		}

		data, err = io.ReadAll(req.Body)
		req.Body.Close()

		if err != nil {
			log.Error(err.Error())
			return
		}

		var APIReq ApiGet1

		err = json.Unmarshal(data, &APIReq)
		if err != nil {
			log.Error(err.Error())
			return
		}
		if messageAddMem {
			c.mem = append(c.mem, [2]string{message, APIReq.Response})
		}
		replay = APIReq.Response
	} else {
		if messageAddMem {
			c.mem = append(c.mem, [2]string{message, ""})
		}
	}

	// 处理过长记忆
	if len(c.mem) > c.MemMaxLength {
		c.mem = c.mem[len(c.mem)-c.MemMaxLength:]
		copy(c.StaticMem, c.mem)
	}
	return
}
