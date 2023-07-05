# ChatBot-QQ

一个简单的聊天机器人

程序会连接[go-cqhttp](https://github.com/Mrs4s/go-cqhttp)的webSocket进行QQ上的消息接收，将过滤处理之后的消息交给[ChatGLM2-6B](https://github.com/THUDM/ChatGLM2-6B)然后将AI的回复送回go-cqhttp，起到一个类似中间件的作用。

## 配置文件

```json
    "api_cqhttp": "ws://127.0.0.1:5701", //
    "api_chat": "http://127.0.0.1:7000", //
    "mem_static": [// 初始记忆，而且不会因mem_length被丢弃，可以多条
        [
            "你的名字：小樱，你说话风格是简洁可爱活泼的，你是可爱的猫娘，说话结尾要加上“喵~”，因为你过于可爱，有很多人想要和你对话，他们的话像这样：\"人名：说话内容\"",// 用户说的话
            "知道了喵！"// AI说的话
        ]
    ],
    "mem_length": 50,// 最大对话记忆轮数
    "this_qq": 123123123,// go-cqhttp上登录的QQ号码
    "chat_args": {// ChatGLM2-6B的参数
        "max_length ": 2048,
        "top_p": 1.0,
        "temperature": 1.0
    }
```

# 不要用于非法用途
