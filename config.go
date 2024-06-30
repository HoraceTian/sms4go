package sms4go

import "time"

// SmsConfig Sms4Go 配置内容
type SmsConfig struct {
	IsPrint     bool        // 是否打印日志
	HttpTimeout int         // http 超时时间
	PoolConfig  *PoolConfig //协程池配置
}

// SupplierConfig 供应商配置规范
type SupplierConfig interface {
	// GetConfigId 用于获取配置对象的唯一标识符
	GetConfigId() string

	// GetSupplier 用于获取配置对象所属的供应商名称或标识符
	GetSupplier() string
}

// PoolConfig 协程池配置
type PoolConfig struct {
	ExpiryDuration   time.Duration // 过期时间
	MaxBlockingTasks int           // 最大阻塞任务数
	RoutinePoolSize  int           // 协程池大小
}
