package service

import (
	"encoding/json"
	"fmt"
	"gimchat/models"
	"net"
	"net/http"
	"sync"

	"github.com/fatih/set"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

func Chat(c *gin.Context) {
	writer := c.Writer
	request := c.Request
	// 校验token 等合法性
	token := c.GetHeader("token")
	user, error := models.FindUserByToken(token)
	if error != nil {
		return
	}
	userId := user.Identity

	isValida := true // 校验token结果
	conn, err := (&websocket.Upgrader{
		// token校验
		CheckOrigin: func(r *http.Request) bool {
			return isValida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Print(err)
		return
	}
	// 获取连接
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	// 用户关系  userid和node绑定  并加锁处理
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()
	// 完成发送逻辑
	go sendProc(node)
	// 完成接收逻辑
	go recvProc(node)

	sendMsg(userId, []byte("你好,欢迎来到ginchat"))
}

// 发送给客户端消息
func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("[ws] sendProc >>>:", string(data))
		}
	}
}

// 接收客户端消息
func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		broadMsg(data)
		fmt.Println("[ws] recvProc >>>:", string(data))
	}
}

var udpsendChan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpsendChan <- data
}

func init() {
	go udpSendProc()
	go udpRecvProc()
}

// udp数据发送协程
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP: net.IPv4(255, 255, 255, 255), //全网段广播地址
		// IP:   net.IPv4(192, 168, 2, 255), //单网段广播地址（本机ip为192.168.2.26时）
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()

	for {
		select {
		case data := <-udpsendChan:
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("[ws] udpSendProc >>>:", string(data))
		}
	}
}

// udp数据接收协程
func udpRecvProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero, //代表本机所有网卡地址
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
	}
	defer con.Close()
	for {
		fmt.Println("[ws] udpRecvProc 。。。。。。")
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(buf[0:n])
	}
}

// 后端调度逻辑处理
func dispatch(data []byte) {
	fmt.Println("[ws] dispatch >>>:", string(data))
	msg := models.Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1:
		sendMsg(msg.TargetUserID, data)
	}
}

func sendMsg(userId int64, msg []byte) {
	fmt.Println("[ws] sendMsg >>>:", userId, msg)
	rwLocker.RLock()
	node, ok := clientMap[userId]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}
