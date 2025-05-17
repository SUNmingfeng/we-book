package ioc

import (
	"basic-go/webook/internal/service/sms"
	"basic-go/webook/internal/service/sms/localsms"
	"basic-go/webook/internal/service/sms/tencent"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111" // 引入sms
	"os"
)

func InitSMSService() sms.Service {
	//return ratelimit.NewRateLimiterSMSService(limiter.NewRedisSlidingWindowsLimiter(), localsms.NewService())
	return localsms.NewService()
	//return initTencentSMSService()
}

func initTencentSMSService() sms.Service {
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		panic("SMS_SECRET_ID is not set")
	}
	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
	if !ok {
		panic("SMS_SECRET_KEY is not set")
	}
	c, err := tencentSMS.NewClient(
		common.NewCredential(secretId, secretKey),
		"ap-nanjing",
		profile.NewClientProfile(),
	)
	if err != nil {
		panic(err)
	}
	return tencent.NewService(c, "1400984222", "兀一天津有限公司")
}
