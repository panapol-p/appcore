package appcore_utils

import "github.com/memphisdev/memphis.go"

func NewMessageBroker(configs *Configurations) *memphis.Conn {
	c, err := memphis.Connect(configs.MemphisHost, configs.MemphisUsername, configs.MemphisToken)
	if err != nil {
		panic(err.Error())
	}
	return c
}
