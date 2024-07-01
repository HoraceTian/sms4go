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
