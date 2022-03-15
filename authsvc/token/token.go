package token

import (
	"time"
)

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

func (t *AuthToken) Duration() time.Duration {
	diff := t.exp - time.Now().Unix()
	return time.Second * time.Duration(diff)
}

type AuthTokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
