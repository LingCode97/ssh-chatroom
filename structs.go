package main

import (
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type Message struct {
	Msg  string
	From string
}

type User struct {
	Session  ssh.Session
	Terminal *terminal.Terminal
	Room     *Room
	Color    string
}

type Room struct {
	Name    string
	History []Message
	Users   []*User
}

// ANSI 颜色码
var Colors = []string{
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
	Rooms:   "查询房间列表",
	Enter:   "进入房间",
	Leave:   "离开房间",
}

const (
	Help    string = "help"
	Color   string = "color"
	Nick    string = "nick"
	History string = "history"
	Users   string = "users"
	Rooms   string = "rooms"
	Enter   string = "enter"
	Leave   string = "leave"
)

// 频道名称
const (
	Default string = "default"
	Tech    string = "tech"
	Music   string = "music"
)
