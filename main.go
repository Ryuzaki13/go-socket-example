package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type Data struct {
	Type    int
	User    string
	Message string
	Raw     json.RawMessage
}

type Client struct {
	Socket      *net.TCPConn
	IsConnected bool
}

// Список подключенных клиентов
//var clientList []*Client
var clientList map[string]*Client

func main() {
	clientList = make(map[string]*Client)

	// 1. Создать слушатель TCP
	addr, e := net.ResolveTCPAddr("tcp", "127.0.0.1:8888")
	if e != nil {
		fmt.Println("ERROR RESOLVE TCP ADDRESS")
		return
	}

	// 2. Запустить слушатель
	listener, e := net.ListenTCP("tcp", addr)
	if e != nil {
		fmt.Println("ERROR LISTEN TCP ADDRESS")
		return
	}

	var n int
	nickName := make([]byte, 40)
	for {
		fmt.Println("Ожидание подключений...")
		// 3. Ожидание подключения клиента
		clientSocket, e := listener.AcceptTCP()
		if e != nil {
			fmt.Println("ERROR CLIENT ACCEPTING")
			return
		}

		go func() {
			for {
				n, e = clientSocket.Read(nickName)
				if e != nil {
					fmt.Println("ERROR READ CLIENT NICKNAME")
					return
				}

				nickNameString := string(nickName[0:n])

				_, ok := clientList[nickNameString]
				if ok {
					_, _ = clientSocket.Write([]byte("1"))
					continue
				} else {
					_, _ = clientSocket.Write([]byte("0"))
				}

				fmt.Println("Подключен новый клиент", nickNameString, clientSocket.RemoteAddr().String())

				client := &Client{
					Socket:      clientSocket,
					IsConnected: true,
				}

				// 4. Отправить работу с подключенным
				//	  клиентов в отдельный поток
				go ClientWorker(client)

				// 5. Добавить в общий список подключенных
				//    клиентов подключенного клиента
				clientList[nickNameString] = client

				i := 0
				clientNames := make([]string, len(clientList))
				for nick := range clientList {
					clientNames[i] = nick
					i++
				}

				var data Data
				data.Type = 3
				data.Raw, e = json.Marshal(clientNames)
				if e != nil {
					fmt.Println(e)
					return
				}

				bytes, e := json.Marshal(data)
				if e != nil {
					fmt.Println(e)
					return
				}

				for _, client := range clientList {
					_, _ = client.Socket.Write(bytes)
				}

				return
			}
		}()
	}

}

func ClientWorker(client *Client) {
	buffer := make([]byte, 8192)
	var data Data
	for client.IsConnected {
		n, e := client.Socket.Read(buffer)
		if e != nil {
			client.IsConnected = false
			break
		}
		e = json.Unmarshal(buffer[:n], &data)
		if e != nil {
			fmt.Println(e)
			return
		}

		switch data.Type {
		case 2:
			{
				user, ok := clientList[data.User]
				fmt.Println(user)
				if ok {
					_, _ = user.Socket.Write(buffer[:n])
				}
			}

		case 4:
			{
				for _, user := range clientList {
					if client.Socket != user.Socket {
						_, _ = user.Socket.Write(buffer[:n])
					}
				}
			}
		}

	}
}
