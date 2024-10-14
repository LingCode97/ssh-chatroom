package main

import (
	"fmt"
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"log"
)

type ChatService interface {
	Enter(sess ssh.Session, term *terminal.Terminal)
	SendMessage(from, message string)
	Leave(sess ssh.Session)
}

func (r *Room) Enter(u *User) {
	log.Println(fmt.Sprintf("用户%s加入聊天室", u.Session.User()))
	r.Users = append(r.Users, u)
	entryMsg := Message{From: r.Name, Msg: "Welcome to my room!"}
	send(u, entryMsg)
	//自动回溯最新10条消息
	sendCount := 0
	for i := len(r.History) - 1; i >= 0; i-- {
		if sendCount >= 10 {
			break
		}
		send(u, r.History[i])
		sendCount++
	}
}

func (r *Room) Leave(sess ssh.Session) {
	log.Println(fmt.Sprintf("用户%s离开聊天室", sess.User()))
	r.Users = removeByUsername(r.Users, sess.User())
}

func (r *Room) SendMessage(from, message string) {
	log.Println(fmt.Sprintf("用户%s发送消息：%s", from, message))
	messageObj := Message{From: from, Msg: message}
	r.History = append(r.History, messageObj)
	for _, u := range r.Users {
		if (u.Session.User()) != from {
			send(u, messageObj)
		}
	}
}

func removeByUsername(s []*User, n string) []*User {
	var index int
	for i, u := range s {
		if u.Session.User() == n {
			index = i
			break
		}
	}
	return append(s[:index], s[index+1:]...)
}

func send(u *User, m Message) {
	raw := u.Color + m.From + "> " + m.Msg + "\033[0m" + "\n"
	_, _ = u.Terminal.Write([]byte(raw))
}
