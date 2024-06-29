package sms4go

import "fmt"

type Client interface {
	CreateSmsBlender()
	GetBySupplier(string) ISmsBlender
	GetSmsBlender() ISmsBlender
	GetSmsBlenderWithConfigId(string) ISmsBlender
}

type smsClient struct {
	blends         map[string]ISmsBlender
	providerHolder *providerFactoryHolder
}

func (s *smsClient) CreateSmsBlender() {
	fmt.Println("create sms blender")
}

func (s *smsClient) GetBySupplier(supplier string) ISmsBlender {
	return s.blends[supplier]
}

// GetSmsBlender 负载均衡获取 Blender
func (s *smsClient) GetSmsBlender() ISmsBlender {
	return nil
}

// GetSmsBlenderWithConfigId 根据 ConfigId 获取 Blender
func (s *smsClient) GetSmsBlenderWithConfigId(configId string) ISmsBlender {
	return nil
}
