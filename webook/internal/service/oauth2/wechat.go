package oauth2

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lithammer/shortuuid/v4"
	"net/http"
	"net/url"
	"webookProgram/webook/internal/domain"
)

var (
	redirectUri = url.PathEscape("https://xxx.com")
)

type Service interface {
	AuthURL(ctx context.Context) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}
type service struct {
	appId     string
	appSecret string
	client    *http.Client
}

func NewService(appId string, appSecret string) Service {
	return &service{appId: appId, appSecret: appSecret, client: http.DefaultClient}
}
func (s *service) AuthURL(ctx context.Context) (string, error) {
	const urlPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	uid := shortuuid.New()
	return fmt.Sprintf(urlPattern, s.appId, redirectUri, uid), nil
}
func (s *service) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {
	const targetPattern = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	target := fmt.Sprintf(targetPattern, s.appId, s.appSecret, code)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	decoder := json.NewDecoder(resp.Body)
	var res Result
	err = decoder.Decode(&res)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	if res.ErrCode != 0 {
		return domain.WechatInfo{}, fmt.Errorf("微信返回错误响应,错误码:%d，错误信息:%s", res.ErrCode, res.ErrMsg)
	}
	return domain.WechatInfo{
		OpenID:  res.OpenId,
		UnionID: res.UnionId,
	}, nil
}

type Result struct {
	AccessToken  string `json:"access_token,omitempty"`
	Expires      int64  `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	OpenId       string `json:"openid,omitempty"`
	Scope        string `json:"scope,omitempty"`
	UnionId      string `json:"unionid,omitempty"`
	ErrCode      int64  `json:"errcode,omitempty"`
	ErrMsg       string `json:"errmsg,omitempty"`
}
