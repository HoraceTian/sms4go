package sms4go

import "fmt"

type Client interface {
	CreateSmsBlender()
}

type smsClient struct {
}

func (smsClient) CreateSmsBlender() {
	fmt.Println("create sms blender")
}
