package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
)

func SendMsg(client *Client) {
	writer := bufio.NewWriter(client.channel)
	for msg := range client.msgChan {
		if msg.sendClient == client {
			continue
		}
		writer.WriteString(msg.msg)
		writer.Flush()
	}
}

func SaveMsg(client *Client, msg string) {
	username := "系统"
	if client != nil {
		username = client.username
	}
	formattedMsg := fmt.Sprintf("[%s]: %s", username, msg)
	Msg := Message{
		msg:        formattedMsg,
		sendClient: client,
	}
	broadcast <- Msg
}

func HandleClient(client *Client) {
	defer client.channel.Close()
	reader := bufio.NewReader(client.channel)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Printf("%s 断开连接", client.username)
				removeClient <- client
				SaveMsg(nil, client.username+"下线了～\n")
				return
			}
			log.Println("客户端消息读取失败:", err)
			return
		}
		SaveMsg(client, message)
	}
}
