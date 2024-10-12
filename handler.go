package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

func SendMsg(client *Client) {
	writer := bufio.NewWriter(client.channel)
	for msg := range client.msgChan {
		if client == msg.sendClient {
			continue
		}
		writer.WriteString(msg.msg)
		writer.Flush()
	}
}

// SaveMsg 将消息保存到全局广播变量中，异步发送
func SaveMsg(client *Client, msg string) {
	username := "系统"
	color := "\u001B[33m"
	if client != nil {
		username = client.username
		color = client.color
	}
	formattedMsg := fmt.Sprintf("%s[%s]: %s\033[0m", color, username, msg)
	Msg := Message{
		msg:        formattedMsg,
		sendClient: client,
	}
	broadcast <- Msg
}

// AssignMsg 给指定客户端发送消息
func AssignMsg(client *Client, msg string) {
	client.msgChan <- Message{msg: msg}
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
		//判断是命令还是消息，#开头是命令
		if strings.HasPrefix(message, "#") {
			Cmd(client, message)
		} else {
			SaveMsg(client, message)
		}
	}
}

func GetNextColor(colorIndex int) string {
	// 循环分配颜色
	color := colors[colorIndex%len(colors)]
	colorIndex++
	return color
}
