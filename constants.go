package sms4go

import "time"

const (
	// Tencent 腾讯云短信服务
	Tencent = "tencent"
)

const (
	ParamSeparate = "&"
)

const (
	DefaultRoutinePoolSize  = 10
	DefaultMaxBlockingTasks = 10
	DefaultExpiryDuration   = 2 * time.Minute
)
