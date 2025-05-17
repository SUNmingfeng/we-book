package failover

import (
	"basic-go/webook/internal/service/sms"
	smsmocks "basic-go/webook/internal/service/sms/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestFailoverSMSService_Send(t *testing.T) {
	testcases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) []sms.Service
		wantErr error
	}{
		{
			name: "一次发送成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0}
			},
			wantErr: nil,
		},
		{
			name: "服务商全部失败",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error1"))
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error2"))
				svc2 := smsmocks.NewMockService(ctrl)
				svc2.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error3"))
				return []sms.Service{svc0, svc1, svc2}
			},
			wantErr: errors.New("所有服务商均发送失败"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewFailoverSMSService(tc.mock(ctrl))
			err := svc.Send(context.Background(), "123", []string{"123"}, []string{"19900002222"}...)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
