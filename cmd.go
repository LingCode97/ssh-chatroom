package main

import (
	"github.com/gliderlabs/ssh"
	"strings"
)

func Cmd(s ssh.Session, message string) {
	user := onlineUsers[s]
	//如果是#开头，则为命令
	cmd, _ := strings.CutPrefix(message, "#")
	cmd, _ = strings.CutSuffix(cmd, "\n")
	if _, ok := commands[cmd]; !ok {
		msg := Message{
			Msg:  "未知命令！",
			From: user.Room.Name,
		}
		send(user, msg)
		return
	}
	switch cmd {
	case Users:
		users(user)
	}
}

func users(user *User) {
	msg := strings.Builder{}
	msg.WriteString("在线列表[\u001B[33m")
	for _, v := range onlineUsers {
		msg.WriteString(v.Color + v.Session.User() + ",")
	}
	msgStr := strings.TrimSuffix(msg.String(), ",")
	msgStr += "\u001B[0m] \n"
	send(user, Message{Msg: msgStr, From: user.Room.Name})
}
