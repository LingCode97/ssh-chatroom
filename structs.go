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

// 支持的命令
var commands = map[string]string{
	Help:    "显示帮助信息",
	Color:   "切换颜色",
	Nick:    "修改昵称",
	History: "查看历史消息",
	Users:   "查询在线用户",
}

const (
	Help    string = "help"
	Color   string = "color"
	Nick    string = "nick"
	History string = "history"
	Users   string = "users"
)
