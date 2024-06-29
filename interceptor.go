package sms4go

// Interceptor 拦截器规范接口
type Interceptor interface {
	GetOrder() int
}
