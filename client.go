package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       99,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("new.Dial error:", err)
		return nil

	}
	client.conn = conn
	return client

}

func (client *Client) DealResponse() {
	//永久阻塞监听
	io.Copy(os.Stdout, client.conn)
}
func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>请输入合法数字<<<")
		return false

	}

}

func (client *Client) UpdateName() bool {

	fmt.Println(">>>请输入用户名:")
	fmt.Scanln(&client.Name)
	sendmsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendmsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false

	}
	return true
}
func (client *Client) PublicChat() {
	var chatmsg string
	fmt.Println(">>>请输入聊天内容,exit退出")
	fmt.Scanln(&chatmsg)

	for chatmsg != "exit" {
		if len(chatmsg) != 0 {
			sendmsg := chatmsg + "\n"
			_, err := client.conn.Write([]byte(sendmsg))
			if err != nil {
				fmt.Println("conn Write err:", err)
				break

			}
		}
		chatmsg = ""
		fmt.Println(">>>公聊模式,exit退出")
		fmt.Scanln(&chatmsg)

	}

}
func (client *Client) SelectUsers() {
	sendmsg := "who\n"
	_, err := client.conn.Write([]byte(sendmsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return

	}
}
func (client *Client) PrivateChat() {
	var remotename string
	var chatmsg string
	client.SelectUsers()
	fmt.Println(">>>请输入聊天对象,exit退出")
	fmt.Scanln(&remotename)

	for remotename != "exit" {
		fmt.Println("请输入消息内容,exit退出:")
		fmt.Scanln(&chatmsg)

		for chatmsg != "exit" {
			if len(chatmsg) != 0 {
				sendmsg := "to|" + remotename + "|" + chatmsg + "\n\n"
				_, err := client.conn.Write([]byte(sendmsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					break

				}
			}
			chatmsg = ""
			fmt.Println(">>>请输入消息内容,exit退出")
			fmt.Scanln(&chatmsg)
		}
		client.SelectUsers()
		fmt.Println(">>>请输入聊天对象,exit退出")
		fmt.Scanln(&remotename)
	}

}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		switch client.flag {
		case 1:
			client.PublicChat()

			break
		case 2:
			client.PrivateChat()

			break
		case 3:
			client.UpdateName()

			break
		}
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器ip")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口")
}

func main() {
	//命令行解析
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>链接服务器失败...")
		return

	}
	//处理server回执消息
	go client.DealResponse()
	fmt.Println(">>>链接服务器成功...")
	//启动客户端业务
	client.Run()

}
