package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		server: server,
	}

	//启动user channel的listener

	go user.ListenMessage()
	return user
}

//监听User channel的方法，一旦有消息就发送到客户端
func (user *User) ListenMessage() {
	for {
		msg := <-user.C
		user.conn.Write([]byte(msg + "\n"))
	}
}

func (user *User) Online() {

	user.server.mapLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.mapLock.Unlock()

	user.server.Broadcast(user, "已上线")
}

func (user *User) Offline() {

	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLock.Unlock()

	user.server.Broadcast(user, "下线")
}

func (user *User) SendMessage(msg string) {
	user.conn.Write([]byte(msg))
}
func (user *User) DoMessage(msg string) {//消息核心处理方法
	if msg == "who" {
		user.server.mapLock.Lock()
		for _, user01 := range user.server.OnlineMap {
			onlinemsg := "[" + user01.Addr + "]" + user01.Name + ":" + "在线...\n"
			user.conn.Write([]byte(onlinemsg))
		}
		user.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newname := strings.Split(msg, "|")[1]

		_, ok := user.server.OnlineMap[newname]
		if ok {
			user.SendMessage("当前用户名被使用")
		} else {
			user.server.mapLock.Lock()
			delete(user.server.OnlineMap, user.Name)
			user.server.OnlineMap[newname] = user
			user.server.mapLock.Unlock()
			user.Name = newname
			user.SendMessage("您已更新用户名:" + user.Name + "\n")

		}

	}else if len(msg)>4 && msg[:3] == "to|"{
		remotename := strings.Split(msg,"|")[1]
		if remotename==""{
			user.SendMessage("消息格式不正确")
			return
		}

		remoteuser,ok := user.server.OnlineMap[remotename]
		if !ok{
			user.SendMessage("该用户不存在")
			return
		}

		content := strings.Split(msg,"|")[2]
		if content ==""{
			user.SendMessage("无消息内容,请重发")
			return
		}
		remoteuser.SendMessage(user.Name+"对你说:"+ content)


	} else {
		user.server.Broadcast(user, msg)
	}

}
