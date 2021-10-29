package token

import "time"

type Maker interface {
	//创建新令牌
	CreateToken(username string, duration time.Duration) (string, error)
	//检查令牌
	VerifyToken(token string) (*Payload, error)
}
