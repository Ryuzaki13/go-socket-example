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

func main() {

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

	// TEST Only
	var data Data
	data.A = 342
	data.S = "Привет Мир!"

	// Преобразовать структуру в байты в формате JSON
	bytes, e := json.Marshal(data)
	if e != nil {
		fmt.Println(e)
		return
	}

	arr := append(IntToByteArray(int64(len(bytes))), bytes...)

	n, e := socket.Write(arr)
	if e != nil {
		fmt.Println("ERROR SEND MESSAGE TO SERVER")
		return
	}

	fmt.Println("отправлено", n, "байт")

	_ = socket.Close()
}

func IntToByteArray(num int64) []byte {
	size := int(unsafe.Sizeof(num))
	arr := make([]byte, size)
	for i := 0; i < size; i++ {
		byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
		arr[i] = byt
	}
	return arr
}
