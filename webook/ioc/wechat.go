package ioc

import (
	"basic-go/webook/internal/service/oauth2/wechat"
)

func InitWechatService() wechat.Service {
	//appID, ok := os.LookupEnv("APP_ID")
	//if !ok {
	//	panic(errors.New("找不到环境变量APP_ID"))
	//}
	return wechat.NewService("123456", "")
}
