package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"webookProgram/webook/internal/service/sms"
)

type SMSService struct {
	svc sms.Service
}

func (s *SMSService) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	var sc Claims
	token, err := jwt.ParseWithClaims(biz, &sc, func(token *jwt.Token) (interface{}, error) {
		return []byte("6uZhFEhonyX0JalbKDkarQMRpzLwuS3N"), nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("token无效")
	}
	tplId := sc.TplId
	return s.Send(ctx, tplId, args, numbers...)
}

type Claims struct {
	jwt.RegisteredClaims
	TplId string
}
