package sms4go

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"net/http"
)

// Client 接口定义了客户端的行为
type Client interface {
	// CreateSmsBlender 用于创建一个新的 Blender
	CreateSmsBlender()

	// GetBySupplier 根据供应商的名称获取一个 Blender
	GetBySupplier(supplierName string) ISmsBlender

	// GetSmsBlender 获取默认的 Blender
	GetSmsBlender() ISmsBlender

	// GetSmsBlenderWithConfigId 根据配置 ID 获取一个 Blender
	GetSmsBlenderWithConfigId(configId string) ISmsBlender
}

// Sms 客户端实现
type smsClient struct {
	blends         map[string]ISmsBlender
	providerHolder *providerFactoryHolder
	httpClient     *http.Client
	routinePool    *ants.Pool
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
