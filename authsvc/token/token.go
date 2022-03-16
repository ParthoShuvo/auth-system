package token

import (
	"time"
)

const tokenTypeBearer = "bearer"

type ExpireTime int

func (et ExpireTime) duration() time.Duration {
	return time.Minute * time.Duration(et)
}

type TokenDef struct {
	Secret string
	Exp    ExpireTime
}

func (td TokenDef) ExpiresAt() int64 {
	return time.Now().Add(td.Exp.duration()).Unix()
}

type JWTDef struct {
	AccessToken  *TokenDef
	RefreshToken *TokenDef
}

type AuthToken struct {
	tokenStr string
	exp      int64
	id       string
	uid      string
}

func (t *AuthToken) UserID() string {
	return t.id
}

func (t *AuthToken) UUID() string {
	return t.uid
}

func (t *AuthToken) String() string {
	return t.tokenStr
}

func (t *AuthToken) Expires() time.Duration {
	diff := t.exp - time.Now().Unix()
	return (time.Second * time.Duration(diff))
}

func (t *AuthToken) expiresInSeconds() time.Duration {
	return time.Duration(t.Expires().Seconds())
}

type AuthTokenPair struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	TokenType    string        `json:"token_type"`
	Expires      time.Duration `json:"expires"`
}
