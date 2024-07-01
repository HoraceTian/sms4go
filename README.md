# sms4go
> 作者：Horace

## 简介
仿照 Sms4j 实现的 golang 版本短信发送服务，未来将会支持多个短信厂商，目前已支持腾讯云短信，未来规划支持以下厂商：
- 阿里云
- 七牛云
- 天翼云
- 待补充

## 如何使用
### 引入依赖

### 简单案例
```golang
package main

import (
	"fmt"
	"sms4go"
	"sms4go/supplier/tencent"
)

func main() {
	// 1.  提供 SmsConfig
	smsConfig := &sms4go.SmsConfig{
		IsPrint: true,
	}

	// 2. 提供拦截器

	// 3. 提供厂商配置
	configMap := make(map[string]sms4go.SupplierConfig)
	configMap["tencent"] = &tencent.Config{
		BaseConfig: sms4go.BaseConfig{
			AccessKeyId:     "你的腾讯云 AccessKeyId",
			AccessKeySecret: "你的腾讯云 AccessKeySecret",
			Signature:       "你的 Signature",
			SDKAppId:        "你的 SDKAppId",
			TemplateId:      "你的 模板Id",
		},
	}

	// 4. 初始化厂商工厂
	client := sms4go.NewSmsClient(sms4go.WithConfigMap(configMap),
		sms4go.WithInterceptors(),
		sms4go.WithSmsConfig(smsConfig),
		sms4go.WithProviderFactories(tencent.NewFactory()),
	)
	client.CreateSmsBlender()
	blender := client.GetBySupplier("tencent")
	message := blender.SendMessage("手机号", "参数")
	fmt.Println(fmt.Sprintf("结果是: %v", message))
}
```