package ratelimit

import (
	"basic-go/webook/internal/service/sms"
	"basic-go/webook/pkg/limiter"
	"context"
	"errors"
)

var errLimited = errors.New("rate_limit_sms limit exceeded")

type RateLimiterSMSService struct {
	limiter limiter.Limiter
	svc     sms.Service
	key     string
}

func NewRateLimiterSMSService(limiter limiter.Limiter, sms sms.Service) *RateLimiterSMSService {
	return &RateLimiterSMSService{
		limiter: limiter,
		svc:     sms,
		key:     "rate_limit_sms",
	}
}

func (r RateLimiterSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	limited, err := r.limiter.Limit(ctx, r.key)
	if err != nil {
		return err
	}
	if limited {
		return errLimited
	}
	return r.svc.Send(ctx, tplId, args, numbers...)
}
