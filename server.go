package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

func main() {
	server := NewServer("127.0.0.1", 8888)
	server.Start()

}

type Server struct {
	Ip   string
	Port int

	//在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//消息广播
	Message chan string
}

//server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}
func (server *Server) ListenMessager() {
	for {
		msg := <-server.Message

		server.mapLock.Lock()
		for _, client := range server.OnlineMap {
			client.C <- msg
		}
		server.mapLock.Unlock()
	}
}

func (server *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	server.Message <- sendMsg
}

func (server *Server) Handler(conn net.Conn) {
	//当前链接的业务
	fmt.Println("链接建立成功")
	user := NewUser(conn, server)

	//用户上线
	user.Online()

	isLive := make(chan bool)

	//接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf) ////////////////////////////////////////////////////////保持用户畅通
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			msg := string(buf[:n-1])
			user.DoMessage(msg) //domessage方法内有执行conn.wirte方法但没有再触发conn.read是由于在一个for循环内

			isLive <- true
		}

	}()
	for {
		select {
		case <-isLive:
			//重置定时器

		case <-time.After(time.Second * 100):
			user.SendMessage("你被踢了")
			close(user.C)
			conn.Close()
			return
		}
	}

}

//启动服务器接口
func (server *Server) Start() {
	//socker listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return

	}
	//close listen socket
	defer listener.Close()

	//启动服务端监听器
	go server.ListenMessager() ///////////////////////////////////////////////////////保持服务器ke达
	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err", err)
			continue

		}

		//do handle
		go server.Handler(conn)
	}

}
