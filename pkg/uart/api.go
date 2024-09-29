package uart

import "time"

/*
 * 另外公共接口
 * GetPortsList() []string 获取当前主机上的串口列表
 * IsPort(portName string, ports []string) bool 判断串口是否在列表中
 */

type Uart interface {

	// 获取串口名
	GetName() string

	// 打开串口
	Open() error

	// 发送完成后返回，会阻塞
	Send(buf []byte) (int, error)

	// 读取, 会阻塞
	Read() []byte

	// 读取，会阻塞，可以设设置超时时间
	ReadWithTimeout(timeout time.Duration) ([]byte, error)

	// 带缓冲区的发送
	SendWithBuffer(buf []byte)

	// 关闭连接
	Close()

	// 内部函数: 发送协程
	sendPoll()

	// 内部函数: 接收协程
	receivePoll()
}
