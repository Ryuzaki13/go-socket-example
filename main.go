package main

import (
	"encoding/json"
	"fmt"
	"net"
	"unsafe"
)

type Data struct {
	A int
	S string
}

type Client struct {
	Socket      *net.TCPConn
	IsConnected bool
}

// Список подключенных клиентов
var clientList []*Client

func main() {
	clientList = make([]*Client, 0)

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

	for {
		fmt.Println("Ожидание подключений...")
		// 3. Ожидание подключения клиента
		clientSocket, e := listener.AcceptTCP()
		if e != nil {
			fmt.Println("ERROR CLIENT ACCEPTING")
			return
		}

		fmt.Println("Подключен новый клиени", clientSocket.RemoteAddr().String())

		client := &Client{
			Socket:      clientSocket,
			IsConnected: true,
		}

		// 4. Отправить работу с подключенным
		//	  клиентов в отдельный поток
		go ClientWorker(client)

		// 5. Добавить в общий список подключенных
		//    клиентов подключенного клиента
		clientList = append(clientList, client)
	}

}

func ClientWorker(client *Client) {

	dataLen := make([]byte, 8)

	for client.IsConnected {

		n, e := client.Socket.Read(dataLen)
		if e != nil {
			client.IsConnected = false
			break
		}

		buffer := make([]byte, ByteArrayToInt(dataLen))

		n, e = client.Socket.Read(buffer)
		if e != nil {
			client.IsConnected = false
			break
		}

		fmt.Println(client.Socket.RemoteAddr().String(), "принято", n+8, "байт")

		fmt.Println(string(buffer))

		var data Data
		e = json.Unmarshal(buffer, &data)
		if e != nil {
			fmt.Println(e)
			return
		}

		fmt.Println(data)
	}

}

func ByteArrayToInt(arr []byte) int64 {
	val := int64(0)
	size := len(arr)
	for i := 0; i < size; i++ {
		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&val)) + uintptr(i))) = arr[i]
	}
	return val
}
