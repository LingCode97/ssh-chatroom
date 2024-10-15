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
	cmds := strings.Split(cmd, " ")
	if len(cmds) == 0 {
		sendError(user, "未知命令！")
		return
	}
	if _, ok := commands[cmds[0]]; !ok {
		sendError(user, "未知命令！")
		return
	}
	switch cmds[0] {
	case Users:
		users(user)
	case Enter:
		enter(user, cmds)
	case Leave:
		leave(user)
	case Rooms:
		rooms(user)
	}
}

func rooms(user *User) {
	var roomsName []string
	for _, v := range availableRooms {
		roomsName = append(roomsName, v.Name)
	}
	send(user, Message{Msg: "[" + strings.Join(roomsName, ",") + "] \n", From: user.Room.Name})
}

func leave(user *User) {
	if user.Room.Name == Default {
		sendError(user, "默认房间无法离开")
		return
	}
	user.Room.Leave(user.Session)
	availableRooms[Default].Enter(user)
}

func enter(user *User, cmds []string) {
	if len(cmds) != 2 {
		sendError(user, "命令错误,进入房间请输入#enter RoomName")
		return
	}
	room := availableRooms[cmds[1]]
	if room == nil {
		sendError(user, "房间不存在,请输入#rooms查看房间列表")
		return
	}
	user.Room.Leave(user.Session)
	room.Enter(user)
}

func users(user *User) {
	var usersName []string
	for _, v := range user.Room.Users {
		usersName = append(usersName, v.Session.User())
	}
	send(user, Message{Msg: "[" + strings.Join(usersName, ",") + "] \n", From: user.Room.Name})
}

func sendError(user *User, err string) {
	msg := Message{
		Msg:  err,
		From: user.Room.Name,
	}
	send(user, msg)
}
