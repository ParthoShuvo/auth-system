package token

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/parthoshuvo/authsvc/table/user"
)

type JWTCustomClaims struct {
	ID  string `json:"id"`
	UID string `json:"uid"`
	jwt.StandardClaims
}

func (claims *JWTCustomClaims) Subject() string {
	return claims.StandardClaims.Subject
}

type Cache interface {
	SetRefreshToken(*AuthToken) error
	GetRefreshToken(string) (string, error)
	RevokeRefreshToken(string) error
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
	return &AuthTokenPair{
		AccessToken:  accessToken.String(),
		RefreshToken: refreshToken.String(),
		TokenType:    tokenTypeBearer,
		Expires:      accessToken.expiresInSeconds(),
	}, nil
}

func (svc *Service) VerifyAccessToken(tokenStr string) (*JWTCustomClaims, error) {
	return svc.parseToken(tokenStr, svc.jwtDef.AccessToken.Secret)
}

func (svc *Service) VerifyRefreshToken(tokenStr string) (*JWTCustomClaims, error) {
	claims, err := svc.parseToken(tokenStr, svc.jwtDef.RefreshToken.Secret)
	if err != nil {
		return nil, err
	}
	tokenID, err := svc.cache.GetRefreshToken(claims.ID)
	if err != nil {
		return nil, err
	}
	if tokenID == "" || tokenID != claims.UID {
		return nil, errors.New("refresh token is invalid or expired")
	}
	return claims, nil
}

func (svc *Service) RevokeRefreshToken(tokenStr string) error {
	claims, err := svc.parseToken(tokenStr, svc.jwtDef.RefreshToken.Secret)
	if err != nil {
		return err
	}
	return svc.cache.RevokeRefreshToken(claims.ID)
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

func (svc *Service) parseToken(tokenStr, secret string) (*JWTCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: [%v]", token.Header["alg"])
			}
			return []byte(secret), nil
		})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token is not valid")
	}
	if claims, ok := token.Claims.(*JWTCustomClaims); ok {
		return claims, nil
	}
	return nil, errors.New("token parsing error")
}

func (svc *Service) signToken(claims *JWTCustomClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
