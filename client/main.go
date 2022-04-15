package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type Data struct {
	// 1 - запрос на получение новости
	// 2 - отправить сообщение
	// 3 - получить список клиентов
	Type int
	// Если Type=2, то указать ник пользователя
	User string
	Who  string
	// Если Type=2, то указать сообщение пользователю
	Message string

	Raw json.RawMessage
}

var n int
var nickName string
var data Data
var clientList []string

func main() {
	clientList = make([]string, 0)

	addr, e := net.ResolveTCPAddr("tcp", "127.0.0.1:8888")
	if e != nil {
		fmt.Println("ERROR RESOLVE TCP ADDRESS")
		return
	}

	socket, e := net.DialTCP("tcp", nil, addr)
	if e != nil {
		fmt.Println("ERROR CONNECT TO SERVER")
		return
	}

	for {
		fmt.Println("Введите ник:")
		_, _ = fmt.Scanln(&nickName)

		if len(nickName) > 20 {
			fmt.Println("Максимальная длина ника 20 букв")
			continue
		}

		n, e = socket.Write([]byte(nickName))
		if e != nil {
			fmt.Println("ERROR SEND CLIENT NICKNAME")
			continue
		}

		data := make([]byte, 2)
		n, e = socket.Read(data)
		if e != nil {
			fmt.Println("ERROR READ RESPONSE")
			continue
		}
		data = data[0:n]
		if string(data) == "1" {
			continue
		} else {
			break
		}
	}

	go func() {
		var recv Data
		clientData := make([]byte, 4096)
		for {
			l, e := socket.Read(clientData)
			if e != nil {
				fmt.Println(e)
				continue
			}
			e = json.Unmarshal(clientData[:l], &recv)
			if e != nil {
				fmt.Println(e)
				continue
			}

			switch recv.Type {
			case 2:
				fmt.Println(recv.Message)
			case 3:
				e = json.Unmarshal(recv.Raw, &clientList)
				if e != nil {
					fmt.Println(e)
					continue
				}
				fmt.Println("Client list now updated")
			}

		}
	}()

	var command int

	for {

		fmt.Println("Выберите команду:")
		fmt.Println("1 - получить последнюю новость")
		fmt.Println("2 - отправить сообщение")
		fmt.Println("0 - завершить программу")

		_, _ = fmt.Scanln(&command)

		switch command {
		case 0:
			fmt.Println("До встречи!")
			return
		case 1:
			//
		case 2:

			fmt.Println("Список клиентов")
			for i, client := range clientList {
				fmt.Println(i+1, "-", client)
			}

			fmt.Println("Выберите пользователя")
			var id int
			_, _ = fmt.Scanln(&id)
			id--
			if id < 0 || id > len(clientList) {
				fmt.Println("Неверный номер пользователя")
				break
			}

			data.Type = 2
			data.User = clientList[id]

			fmt.Print("Введите сообщение: ")
			_, _ = fmt.Scanln(&data.Message)

			bytes, e := json.Marshal(data)
			if e != nil {
				fmt.Println(e)
				break
			}

			_, e = socket.Write(bytes)
			if e != nil {
				fmt.Println(e)
				_ = socket.Close()
				return
			}

		default:
			fmt.Println("Неверная команда")
		}
	}
}
