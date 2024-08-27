package uart

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestUart(t *testing.T) {
	// 列出串口列表
	fmt.Println(GetPortsList())
	// [/dev/ttyS0 /dev/ttyUSB0 /dev/ttyUSB1]

	// 获取串中列表
	ports := GetPortsList()
	fmt.Println(IsPort("/dev/ttyUSB1", ports))
	// true

	// 创建串口
	uart := NewUart("/dev/ttyUSB0", WithBaudRate(9600))

	// 打开串口
	uart.Open()

	// 记得关闭
	defer uart.Close()

	// 收发字符串
	go func(u Uart) {
		time.Sleep(3 * time.Second)
		u.Send([]byte("Hello"))
	}(uart)

	fmt.Println(uart.Read())
	//output: [72 101 108 108 111]

	// 收发结构体
	type User struct {
		Name string
	}
	go func(u Uart) {
		time.Sleep(1 * time.Second)
		user := User{
			Name: "Hello",
		}
		b, _ := json.Marshal(user)
		fmt.Println(string(b))
		//output: {"Name":"Hello"}
		u.Send(b)
	}(uart)

	buf := uart.Read()
	readUser := User{}
	json.Unmarshal(buf, &readUser)
	fmt.Println(readUser)
	//output: {Hello}

}
