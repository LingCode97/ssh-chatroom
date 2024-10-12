package main

import "golang.org/x/crypto/ssh"

type Client struct {
	channel  ssh.Channel
	username string
	msgChan  chan Message
	color    string
}

type Message struct {
	msg        string
	sendClient *Client //发送消息的客户端
}

// ANSI 颜色码
var colors = []string{
	"\033[31m", // 红色
	"\033[32m", // 绿色
	"\033[33m", // 黄色
	"\033[34m", // 蓝色
	"\033[35m", // 紫色
	"\033[36m", // 青色
}
