package token

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/parthoshuvo/authsvc/table/user"
)

type JWTCustomClaims struct {
	ID  string `json:"id"`
	UID string `json:"uid"`
	jwt.StandardClaims
}

type Cache interface {
	SetRefreshToken(*AuthToken) error
}

type Service struct {
	jwtDef *JWTDef
	cache  Cache
}

func NewService(jwtDef *JWTDef, cache Cache) *Service {
	return &Service{jwtDef, cache}
}

func (svc *Service) NewAuthTokenPair(usr *user.User) (*AuthTokenPair, error) {
	accessToken, err := svc.createAuthToken(usr, svc.jwtDef.AccessToken)
	if err != nil {
		return nil, err
	}
	refreshToken, err := svc.createAuthToken(usr, svc.jwtDef.RefreshToken)
	if err != nil {
		return nil, err
	}
	if err := svc.cache.SetRefreshToken(refreshToken); err != nil {
		return nil, err
	}
	return &AuthTokenPair{AccessToken: accessToken.String(), RefreshToken: refreshToken.String()}, nil
}

func (svc *Service) createAuthToken(usr *user.User, tokenDef *TokenDef) (*AuthToken, error) {
	userID := usr.RowGUID
	uid := uuid.NewString()
	exp := tokenDef.ExpiresAt()
	claims := &JWTCustomClaims{
		ID:  userID,
		UID: uid,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  jwt.TimeFunc().Unix(),
			Subject:   usr.Email.String(),
			ExpiresAt: exp,
		},
	}
	tokenStr, err := svc.signToken(claims, tokenDef.Secret)
	if err != nil {
		return nil, err
	}
	return &AuthToken{tokenStr, exp, userID, uid}, nil
}

func (svc *Service) signToken(claims *JWTCustomClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
