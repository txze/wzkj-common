package jiguang

import (
	"fmt"
	"log"
)

// ExampleUsage 展示如何使用极光一键登录组件
func ExampleUsage() {
	// 注意：这是一个示例，实际使用时需要先初始化客户端
	// 1. 在应用启动时初始化
	// err := Init()
	// if err != nil {
	//     log.Fatal(err)
	// }

	// 2. 验证登录Token
	// phoneInfo, err := VerifyLoginToken("your_login_token", "your_ex_id")
	// if err != nil {
	//     log.Printf("验证失败: %v", err)
	//     return
	// }

	// 3. 处理结果
	// if phoneInfo.Success {
	//     fmt.Printf("手机号: %s\n", phoneInfo.Phone)
	//     fmt.Printf("自定义ID: %s\n", phoneInfo.ExID)
	// } else {
	//     fmt.Printf("验证失败: %s\n", phoneInfo.Error)
	// }
}

// ExampleCustomClient 展示如何使用自定义客户端
func ExampleCustomClient() {
	// 创建配置
	config := &Config{
		AppKey:       "your_app_key",
		MasterSecret: "your_master_secret",
		PrivateKey:   "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----",
		APIURL:       "https://api.verification.jpush.cn/v1/web/loginTokenVerify",
	}

	// 创建客户端
	client, err := NewClient(config)
	if err != nil {
		log.Printf("创建客户端失败: %v", err)
		return
	}

	// 使用客户端验证登录Token
	phoneInfo, err := client.VerifyLoginToken("your_login_token", "your_ex_id")
	if err != nil {
		log.Printf("验证失败: %v", err)
		return
	}

	// 处理结果
	if phoneInfo.Success {
		fmt.Printf("验证成功，手机号: %s\n", phoneInfo.Phone)
	} else {
		fmt.Printf("验证失败: %s\n", phoneInfo.Error)
	}
}

// ExampleWithErrorHandling 展示完整的错误处理
func ExampleWithErrorHandling() {
	// 初始化客户端
	err := Init()
	if err != nil {
		log.Fatalf("初始化极光客户端失败: %v", err)
	}

	// 验证登录Token
	loginToken := "STsid0000001542695429579Ob28vB7b0cYTI9w0GGZrv8ujUu05qZvw"
	exID := "user_123"

	phoneInfo, err := VerifyLoginToken(loginToken, exID)
	if err != nil {
		log.Printf("验证失败: %v", err)
		return
	}

	// 检查验证结果
	if !phoneInfo.Success {
		log.Printf("验证失败: %s", phoneInfo.Error)
		return
	}

	// 验证成功，处理手机号
	fmt.Printf("验证成功！\n")
	fmt.Printf("手机号: %s\n", phoneInfo.Phone)
	fmt.Printf("自定义ID: %s\n", phoneInfo.ExID)

	// 在这里可以进行后续的业务处理
	// 比如：保存用户信息、创建用户账号等
}
