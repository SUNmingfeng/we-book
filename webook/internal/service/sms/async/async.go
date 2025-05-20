package async

import (
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/service/sms"
	"context"
	"time"
)

type AsyncSMSService struct {
	svc        sms.Service
	repo       repository.AsyncSmsRepository
	windowSize int
	times      []time.Duration
}

func NewAsyncSMSService(svc sms.Service, repo repository.AsyncSmsRepository, winSize int, times []time.Duration) *AsyncSMSService {
	return &AsyncSMSService{
		svc:        svc,
		repo:       repo,
		windowSize: winSize,
		times:      times,
	}
}

func (a AsyncSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	//处理要异步发送的请求
	if a.needAsync() {
		//把请求存入数据库
	}
	return a.svc.Send(ctx, tplId, args, numbers...)
}

func (a AsyncSMSService) needAsync() bool {
	//如果最近连续5次的响应时间都大雨500ms，就判定为服务商崩溃，即需要异步
	//如进入异步的时间已经大于3分钟，则退出异步
}
