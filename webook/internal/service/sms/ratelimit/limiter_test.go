package ratelimit

import (
	"basic-go/webook/internal/service/sms"
	smsmocks "basic-go/webook/internal/service/sms/mocks"
	"basic-go/webook/pkg/limiter"
	limitmocks "basic-go/webook/pkg/limiter/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRateLimiterSMSService_Send(t *testing.T) {
	testcases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter)
		//不需要输入
		wantErr error
	}{
		{
			name: "不限流",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				smssvc := smsmocks.NewMockService(ctrl)
				l := limitmocks.NewMockLimiter(ctrl)
				smssvc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
				return smssvc, l
			},
			wantErr: nil,
		},
		{
			name: "限流",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				smssvc := smsmocks.NewMockService(ctrl)
				l := limitmocks.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(true, nil)
				return smssvc, l
			},
			wantErr: errLimited,
		},
		{
			name: "限流器错误",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				smssvc := smsmocks.NewMockService(ctrl)
				l := limitmocks.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, errors.New("限流器错误"))
				return smssvc, l
			},
			wantErr: errors.New("限流器错误"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			smsSvc, l := tc.mock(ctrl)
			svc := NewRateLimiterSMSService(l, smsSvc)
			err := svc.Send(context.Background(), "123", []string{"456"}, "19900001111")
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
