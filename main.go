package main

import (
	"fmt"
	"gopkg/pkg/config"
)

func main() {
	fmt.Println("gopkg")
	cfg := config.New("test", "/home/hello/netconfig.yaml", nil)
	cfg.Config()
	cfg.ReadConfigFromFile()
	fmt.Println(cfg.Config().Get("network"))
	// fmt.Println(cfg.Config().GetString("network"))
}
