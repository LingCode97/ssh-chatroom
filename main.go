package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"golang.org/x/crypto/ssh"
)

// 定义一个全局聊天室管理器
type ChatRoom struct {
	mu      sync.Mutex
	clients map[string]io.Writer
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		clients: make(map[string]io.Writer),
	}
}

func (c *ChatRoom) Join(clientID string, writer io.Writer) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.clients[clientID] = writer
}

func (c *ChatRoom) Leave(clientID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.clients, clientID)
}

func (c *ChatRoom) Broadcast(sender, msg string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for id, client := range c.clients {
		if id != sender { // 不要广播给发送者自己
			fmt.Fprintln(client, msg)
		}
	}
}

var chatRoom = NewChatRoom()

// 生成 SSH 私钥
func generatePrivateKey() (ssh.Signer, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	privateKeyBytes := pem.EncodeToMemory(privateKeyPEM)

	// 将生成的私钥转换为 SSH 格式
	signer, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		return nil, err
	}
	return signer, nil
}

// 简单的 SSH 服务端配置
func handleSSHConnection(conn net.Conn) {
	config := &ssh.ServerConfig{
		NoClientAuth: true, // 不验证客户端
	}

	private, err := generatePrivateKey()
	if err != nil {
		log.Fatal("生成私钥失败:", err)
	}
	fmt.Println(private)

	config.AddHostKey(private)

	// 协商 SSH 会话
	sshConn, chans, reqs, err := ssh.NewServerConn(conn, config)
	if err != nil {
		log.Println("SSH 握手失败:", err)
		return
	}
	defer sshConn.Close()

	go ssh.DiscardRequests(reqs)

	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "未知的通道类型")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Println("无法接受通道:", err)
			continue
		}
		defer channel.Close()

		clientID := sshConn.RemoteAddr().String()
		chatRoom.Join(clientID, channel)
		defer chatRoom.Leave(clientID)

		go func() {
			for req := range requests {
				if req.Type == "shell" {
					req.Reply(true, nil)
				}
			}
		}()

		// 处理客户端输入
		go func() {
			buf := make([]byte, 1024)
			for {
				n, err := channel.Read(buf)
				if err != nil {
					break
				}
				msg := string(buf[:n])
				chatRoom.Broadcast(clientID, msg)
			}
		}()
	}
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:2222")
	if err != nil {
		log.Fatal("无法监听端口:", err)
	}
	defer listener.Close()

	log.Println("SSH 聊天室服务器启动，等待连接...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("无法接受连接:", err)
			continue
		}

		go handleSSHConnection(conn)
	}
}
