package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"net"
	"sync"
)

type Client struct {
	channel  ssh.Channel
	username string
	msgChan  chan string
}

var (
	clients      = make(map[*Client]bool)
	broadcast    = make(chan string)
	addClient    = make(chan *Client)
	removeClient = make(chan *Client)
	mu           sync.Mutex
)

func handleClient(client *Client) {
	defer client.channel.Close()
	reader := bufio.NewReader(client.channel)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Printf("%s 断开连接", client.username)
				removeClient <- client
				return
			}
			log.Println("客户端消息读取失败:", err)
			return
		}

		formattedMsg := fmt.Sprintf("[%s]: %s", client.username, message)
		broadcast <- formattedMsg
	}
}

func broadcaster() {
	for {
		select {
		case msg := <-broadcast:
			mu.Lock()
			for client := range clients {
				client.msgChan <- msg
			}
			mu.Unlock()
		case client := <-addClient:
			mu.Lock()
			clients[client] = true
			mu.Unlock()
		case client := <-removeClient:
			mu.Lock()
			if _, ok := clients[client]; ok {
				delete(clients, client)
				close(client.msgChan)
			}
			mu.Unlock()
		}
	}
}

func sshHandler(conn net.Conn, config *ssh.ServerConfig) {
	sshConn, chans, reqs, err := ssh.NewServerConn(conn, config)
	if err != nil {
		log.Println("连接失败:", err)
		return
	}
	defer sshConn.Close()

	log.Printf("新的SSH连接： %s (%s)", sshConn.RemoteAddr(), sshConn.ClientVersion())

	go ssh.DiscardRequests(reqs)

	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Println("无法接受channel:", err)
			return
		}

		go func(in <-chan *ssh.Request) {
			for req := range in {
				switch req.Type {
				case "shell":
					req.Reply(true, nil)
				case "window-change":
					req.Reply(true, nil)
				default:
					req.Reply(false, nil)
				}
			}
		}(requests)

		client := &Client{
			channel:  channel,
			username: sshConn.User(),
			msgChan:  make(chan string),
		}

		addClient <- client
		go handleClient(client)

		go func() {
			writer := bufio.NewWriter(channel)
			for msg := range client.msgChan {
				writer.WriteString(msg)
				writer.Flush()
			}
		}()
	}
}

func LoadHostKey() (ssh.Signer, error) {
	key, err := ioutil.ReadFile("key.txt")
	if err != nil {
		return nil, err
	}
	return ssh.ParsePrivateKey(key)
}

func main() {
	go broadcaster()

	signer, err := LoadHostKey()
	if err != nil {
		log.Fatalf("私钥加载失败: %v", err)
	}

	config := &ssh.ServerConfig{
		NoClientAuth: true,
	}
	config.AddHostKey(signer)

	listener, err := net.Listen("tcp", "0.0.0.0:2222")
	if err != nil {
		log.Fatalf("监听连接失败: %v", err)
	}
	defer listener.Close()

	log.Println("开始监听 0.0.0.0:2222...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("建立连接失败: %v", err)
			continue
		}
		go sshHandler(conn, config)
	}
}
