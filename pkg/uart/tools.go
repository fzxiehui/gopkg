package uart

import (
	"log"
	"runtime"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

/*
 * 获取当前主机上的串口列表
 */
func GetPortsList() []string {

	if runtime.GOOS == "windows" {
		ports, err := serial.GetPortsList()
		if err != nil {
			log.Println("Error listing ports:", err)
			return nil
		}
		return ports
	}

	if runtime.GOOS == "linux" {
		// get usb ports
		usbPorts, err := enumerator.GetDetailedPortsList()
		if err != nil {
			log.Println("Error listing ports:", err)
			return nil
		}
		portNames := make([]string, 0)

		for _, port := range usbPorts {
			portNames = append(portNames, port.Name)
		}
		return portNames
	}

	return nil
}

/*
 * 判断串口是否在列表中
 */
func IsPort(portName string, ports []string) bool {
	for _, p := range ports {
		if p == portName {
			return true
		}
	}
	return false
}
