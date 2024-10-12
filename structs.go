package main

import "golang.org/x/crypto/ssh"

type Client struct {
	channel  ssh.Channel
	username string
	msgChan  chan Message
}

type Message struct {
	msg        string
	sendClient *Client //发送消息的客户端
}
