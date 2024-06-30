package infra

// BaseConfig Sms 基础配置
type BaseConfig struct {
	Factory         string
	AccessKeyId     string
	SDKAppId        string
	AccessKeySecret string
	Signature       string
	TemplateId      string
	Weight          int
	ConfigId        string
	RetryInterval   int
	MaxRetries      int
	Maximum         int
}

func (b *BaseConfig) GetConfigId() string {
	return b.ConfigId
}
func (b *BaseConfig) GetSupplier() string {
	return b.Factory
}
func (b *BaseConfig) SetRetryInterval(retryInterval int) {
	if retryInterval <= 0 {
		panic("retryInterval must be greater than 0 seconds")
	}
	b.RetryInterval = retryInterval
}
func (b *BaseConfig) SetMaxRetries(maxRetries int) {
	if maxRetries < 0 {
		panic("maxRetries cannot be less than 0")
	}
	b.MaxRetries = maxRetries
}
func (b *BaseConfig) SetMaximum(maximum int) {
	if maximum < 0 {
		panic("maximum cannot be less than 0")
	}
	b.Maximum = maximum
}
