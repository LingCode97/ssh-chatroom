package main

import "strings"

func Cmd(client *Client, message string) {
	//如果是#开头，则为命令
	cmd, _ := strings.CutPrefix(message, "#")
	cmd, _ = strings.CutSuffix(cmd, "\n")
	if _, ok := commands[cmd]; !ok {
		SaveMsg(nil, "未知命令！")
		return
	}
	switch cmd {
	case Users:
		users(client)
	}
}

func users(client *Client) {
	msg := strings.Builder{}
	msg.WriteString("在线列表[\u001B[33m")
	for onlineClient := range clients {
		msg.WriteString(onlineClient.color + onlineClient.username + ",")
	}
	msgStr := strings.TrimSuffix(msg.String(), ",")
	msgStr += "\u001B[0m] \n"
	AssignMsg(client, msgStr)
}
