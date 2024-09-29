package uart

import (
	"errors"
	"fmt"
	"log"
	"time"

	"go.bug.st/serial"
)

type uart struct {
	PortName       string      // 串口名字
	BaudRate       int         // 通信波特率
	Port           serial.Port // 通信对象
	SendChannel    chan []byte // 发送缓冲区
	ReceiveChannel chan []byte // 接收缓冲区
	ErrorHandler   func(error) // 错误处理程序
	OpenStatus     bool        // 打开状态
	ClosedChannel  chan bool   // 关闭通道
}

// Config Begin
type Option func(*uart)

// 设置波特率
func WithBaudRate(baudRate int) Option {
	return func(u *uart) {
		u.BaudRate = baudRate
	}
}

// 设置错误处理程序
func WithErrorHandler(handler func(error)) Option {
	return func(u *uart) {
		u.ErrorHandler = handler
	}
}

// Config End
func NewUart(portName string, options ...Option) Uart {
	uart := &uart{
		PortName:       portName,
		BaudRate:       9600,
		SendChannel:    make(chan []byte),
		ReceiveChannel: make(chan []byte),
		ErrorHandler: func(err error) {
			log.Fatal(err)
		},
		ClosedChannel: make(chan bool),
		OpenStatus:    false,
	}
	for _, f := range options {
		f(uart)
	}
	return uart
}

// 获取串口名
func (u *uart) GetName() string {
	return u.PortName
}

// 打开串口
func (u *uart) Open() error {
	mode := &serial.Mode{
		BaudRate: u.BaudRate,
	}
	port, err := serial.Open(u.PortName, mode)
	if err != nil {
		log.Println("Error opening port: ", err)
		return err
	}
	// 串口对像
	u.Port = port
	// 连接状态
	u.OpenStatus = true

	// 开启收发轮训
	go u.receivePoll()
	go u.sendPoll()

	return nil
}

// 发送完成后返回，会阻塞
func (uart *uart) Send(buf []byte) (int, error) {
	if !uart.OpenStatus {
		log.Println("串口未打开, 请使用 Open 函数!")
		return 0, errors.New("串口未打开, 请使用 Open 函数!")
	}
	l, err := uart.Port.Write(buf)
	if err != nil {
		uart.Port.Close()
		binerr := errors.New(fmt.Sprintf("%s -> %s", "Send Error!", err.Error()))
		log.Println(binerr)
		uart.ErrorHandler(binerr)
		return 0, err
	}

	return l, nil
}

// 读取，会阻塞，可以设设置超时时间
func (uart *uart) ReadWithTimeout(timeout time.Duration) ([]byte, error) {
	select {
	case data := <-uart.ReceiveChannel:
		return data, nil
	case <-time.After(timeout):
		return nil, errors.New("read operation timed out")
	}
}

// 读取, 会阻塞
func (uart *uart) Read() []byte {
	return <-uart.ReceiveChannel
}

// 带缓冲区的发送
func (uart *uart) SendWithBuffer(buf []byte) {
	uart.SendChannel <- buf
}

// 内部函数: 发送协程
func (uart *uart) sendPoll() {

	for {
		select {
		case exit := <-uart.ClosedChannel:
			if exit {
				log.Println("sendPoll exit")
				return
			}
		case data := <-uart.SendChannel:
			_, err := uart.Port.Write(data)
			if err != nil {
				binerr := errors.New(fmt.Sprintf("%s -> %s", "sendPoll Error!", err.Error()))
				log.Println(binerr)
				uart.Port.Close()
				// panic(binerr)
				uart.ErrorHandler(binerr)
			}
		}
	}
}

// 内部函数: 接收协程
func (uart *uart) receivePoll() {
	for {
		select {
		case exit := <-uart.ClosedChannel:
			if exit {
				log.Println("readPoll exit")
				return
			}
		default:
			buf := make([]byte, 128)
			n, err := uart.Port.Read(buf)
			if err != nil {
				binerr := errors.New(fmt.Sprintf("%s -> %s", "readPoll Error!", err.Error()))
				uart.Port.Close()
				log.Println(binerr)
				uart.ErrorHandler(binerr)
			}
			// for i := 0; i < n; i++ {
			// log.Debugf("%02X", buf[i])
			// }

			if n > 0 {
				uart.ReceiveChannel <- buf[:n]
			}
		}
	}
}

// 关闭
func (uart *uart) Close() {
	uart.ClosedChannel <- true
	uart.OpenStatus = false
	uart.Port.Close()
}
