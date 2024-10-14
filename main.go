package main

import (
	"fmt"
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"math/rand/v2"
)

var (
	onlineUsers    map[ssh.Session]*User
	availableRooms map[string]*Room
)

func main() {
	onlineUsers = make(map[ssh.Session]*User)
	availableRooms = make(map[string]*Room)
	availableRooms[Default] = &Room{Name: Default}
	availableRooms[Tech] = &Room{Name: Tech}
	availableRooms[Music] = &Room{Name: Music}

	ssh.Handle(func(s ssh.Session) {
		//监听断开连接
		go func(s ssh.Session) {
			select {
			case <-s.Context().Done():
				if _, ok := onlineUsers[s]; ok {
					onlineUsers[s].Room.Leave(s)
				}
			}
		}(s)
		chat(s)
	})

	log.Println("开始监听端口 2222...")
	log.Fatal(ssh.ListenAndServe(":2222", nil, ssh.HostKeyFile("key.txt")))
}

func chat(s ssh.Session) {
	term := terminal.NewTerminal(s, fmt.Sprintf("%s > ", s.User()))
	if _, ok := onlineUsers[s]; !ok {
		//初次加入聊天室,分配默认房间
		d := availableRooms[Default]
		user := &User{Session: s, Terminal: term, Room: d}
		//随机分配颜色
		index := rand.IntN(len(Colors))
		user.Color = Colors[index]
		onlineUsers[s] = user
		d.Enter(user)
	}
	for {
		line, err := term.ReadLine()
		if err != nil {
			break
		}
		if len(line) > 0 && onlineUsers[s] != nil {
			if line[0] == '#' {
				Cmd(s, line)
			} else {
				onlineUsers[s].Room.SendMessage(s.User(), line)
			}
		}
	}
}
