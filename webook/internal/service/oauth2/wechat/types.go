package wechat

import (
	"basic-go/webook/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Service interface {
	AuthURL(context context.Context, state string) (string, error)
	VerifyCode(context context.Context, code string) (domain.WechatInfo, error)
}

var redirectURL = url.PathEscape("https://meoying.com/oauth2/wechat/callback")

type service struct {
	appID     string
	appSecret string
	client    *http.Client
}

func NewService(appID string, appSecret string) Service {
	return &service{
		appID:     appID,
		appSecret: appSecret,
		client:    http.DefaultClient,
	}
}

func (s service) AuthURL(context context.Context, state string) (string, error) {
	const authURLPattern = "`https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect`"
	return fmt.Sprintf(authURLPattern, s.appID, redirectURL, state), nil
}

func (s service) VerifyCode(context context.Context, code string) (domain.WechatInfo, error) {
	accesstokenURL := fmt.Sprintf(`https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code`,
		s.appID, s.appSecret, code)
	req, err := http.NewRequestWithContext(context, http.MethodGet, accesstokenURL, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	httpResp, err := s.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	var res Result
	err = json.NewDecoder(httpResp.Body).Decode(&res)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	if res.ErrCode != 0 {
		return domain.WechatInfo{}, fmt.Errorf("wechat verify code error code:%d, error msg:%s", res.ErrCode, res.ErrMsg)
	}
	return domain.WechatInfo{
		OpenId:  res.Openid,
		UnionId: res.UnionId,
	}, nil
}

type Result struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionId      string `json:"unionid"`
	//错误返回
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
