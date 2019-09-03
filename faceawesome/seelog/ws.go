package seelog

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
)

/**
client for websocket
*/
type client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
}

/**
client for manager
*/
type clientManager struct {
	clients    map[*client]bool
	broadcast  chan []byte
	register   chan *client
	unregister chan *client
}

var manager = clientManager{
	//广播
	broadcast:  make(chan []byte),
	register:   make(chan *client),
	unregister: make(chan *client),
	clients:    make(map[*client]bool),
}

func (manager *clientManager) start() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("[seelog] error:%+v", err)
		}
	}()

	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				//关闭管道
				close(conn.send)
				//关闭websocket
				conn.socket.Close()
				//从clients中删除conn 管道
				delete(manager.clients, conn)
			}
		case message := <-manager.broadcast:
			for conn := range manager.clients {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(manager.clients, conn)
				}
			}
		}
	}
}

/**
socket数据写入
*/
func (c *client) write() {
	for msg := range c.send {
		if _, err := c.socket.Write(msg); err != nil {
			fmt.Println("write msg failed. ", err)
			break
		}
	}
	log.Println("web socket closed")
}
